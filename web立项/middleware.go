package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 验证 JWT Token
func Authcheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "缺少 Token"})
			c.Abort()
			return
		}

		// 去掉 "Bearer " 前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Token 无效"})
			c.Abort()
			return
		}

		// 把用户信息存入上下文，后面的函数可以取出来用
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// 验证是否是管理员
func Admincheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"msg": "权限不足，需要管理员角色"})
			c.Abort()
			return
		}
		c.Next()
	}
}
