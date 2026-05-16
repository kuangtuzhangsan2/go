package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var rdb *redis.Client

func main() {
	// 连接 MySQL
	dsn := "root:cgs197905@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("MySQL 连接失败:", err)
	}
	// 自动迁移表结构
	db.AutoMigrate(&User{})

	// 连接 Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis 连接失败:", err)
	}

	// 初始化 Gin 路由
	r := gin.Default()

	// 公开接口
	r.POST("/register", func(c *gin.Context) { Register(c, db) })
	r.POST("/login", func(c *gin.Context) { Login(c, db, rdb) })
	r.POST("/forget", func(c *gin.Context) { forget(c, db, rdb) })
	r.POST("/change", func(c *gin.Context) { change(c, db) })

	//受保护的接口 (需要登录)
	api := r.Group("/api")
	api.Use(Authcheck())
	{
		api.GET("/users", func(c *gin.Context) { GetUsers(c, db) })

		// 管理员专用接口 (需要 admin)
		admin := api.Group("/admin")
		admin.Use(Admincheck())
		{
			admin.DELETE("/delete-all", func(c *gin.Context) {
				db.Where("1=1").Delete(&User{})
				c.JSON(200, gin.H{"msg": "所有用户已删除 (管理员特权)"})
			})
		}
	}

	//启动服务
	log.Println("服务启动在 http://localhost:8080")
	r.Run(":8080")
}
