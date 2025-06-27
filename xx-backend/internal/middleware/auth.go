package middleware

import (
	"net/http"
	"strings"

	"xx-backend/internal/model"
	"xx-backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证token",
			})
			c.Abort()
			return
		}

		// 检查token格式
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token格式错误",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// 验证token（这里需要从context中获取authService）
		authService, exists := c.Get("auth_service")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "认证服务未初始化",
			})
			c.Abort()
			return
		}

		userID, err := authService.(*service.AuthService).ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "token无效或已过期",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 获取用户名
		username := "unknown"
		if db, exists := c.Get("db"); exists {
			var user model.User
			if err := db.(*gorm.DB).First(&user, userID).Error; err == nil {
				username = user.Username
			}
		}

		// 将用户ID和用户名存储到context中
		c.Set("user_id", userID)
		c.Set("username", username)
		c.Next()
	}
}
