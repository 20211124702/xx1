package handler

import (
	"fmt"
	"net/http"

	"xx-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func Login(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req service.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}
		fmt.Println(&req)
		resp, err := authService.Login(&req, c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "登录失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登录成功",
			"data":    resp,
		})
	}
}

func Logout(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			return
		}

		// 获取用户名用于Kafka记录
		username, exists := c.Get("username")
		if !exists {
			username = "unknown"
		}

		if err := authService.Logout(userID.(int), username.(string)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "登出失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登出成功",
		})
	}
}

func GetProfile(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
			})
			return
		}

		user, err := userService.GetProfile(userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取用户信息失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    user,
		})
	}
}

func Register(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Email    string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": "参数错误", "error": err.Error()})
			return
		}
		if err := userService.Register(req.Username, req.Password, req.Email); err != nil {
			c.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"code": 200, "message": "注册成功"})
	}
}
