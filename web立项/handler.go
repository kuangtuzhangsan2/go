package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var jwtKey = []byte("my_secret_key") // 生产环境请放在环境变量里

// Claims 定义 JWT 里的数据
type Claims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// 注册接口
func Register(c *gin.Context, db *gorm.DB) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	// 设置密码
	if err := input.SetPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "设置失败"})
		return
	}

	// 默认给普通用户角色
	input.Role = "user"

	// 存入数据库
	if err := db.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "用户名已存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}

// 登录接口
func Login(c *gin.Context, db *gorm.DB, rdb *redis.Client) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	var user User
	// 查找用户
	db.Where("username = ?", input.Username).First(&user)

	// 校验密码
	if user.ID == 0 || !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "账号或密码错误，如果遗忘，可以使用口令登录(/forget)"})
		return
	}

	// 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserId:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})
	tokenString, _ := token.SignedString(jwtKey)

	// 将用户信息存入 Redis 缓存
	rdb.Set(c.Request.Context(), "user_info_"+string(rune(user.ID)), user.Username, time.Hour*24)

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "msg": "登录成功"})
}

// 忘记密码
func forget(c *gin.Context, db *gorm.DB, rdb *redis.Client) {
	var input struct {
		Username    string `json:"username"`
		Passproblem string `json:"passproblem"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	var user User
	// 查找用户
	db.Where("username = ?", input.Username).First(&user)

	// 校验口令
	if user.ID == 0 || !user.CheckPassproblem(input.Passproblem) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "口令回答错误"})
		return
	}

	// 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserId:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	})
	tokenString, _ := token.SignedString(jwtKey)

	// 将用户信息存入 Redis 缓存
	rdb.Set(c.Request.Context(), "user_info_"+string(rune(user.ID)), user.Username, time.Hour*24)

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "msg": "登录成功"})
}

// 修改密码
func change(c *gin.Context, db *gorm.DB) {
	var input struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		Newpassword string `json:"newpassword"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	var user User
	// 查找用户
	db.Where("username = ?", input.Username).First(&user)

	// 校验密码
	if user.ID == 0 || !user.CheckPassword(input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "账号或密码错误，如果遗忘，可以使用口令登录(/forget)"})
		return
	}

	// 修改密码

	user.Password = input.Newpassword
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"msg": "修改成功"})
}

// 获取用户列表
func GetUsers(c *gin.Context, db *gorm.DB) {
	var users []User
	db.Find(&users)
	c.JSON(http.StatusOK, users)
}
