package handler

import (
	"net/http"

	"xx-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SendTestMessage 发送测试消息到Kafka
func SendTestMessage(kafkaService *service.KafkaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if kafkaService == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"code":    503,
				"message": "Kafka服务未初始化",
			})
			return
		}

		var req struct {
			Message string `json:"message" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		// 发送系统信息到Kafka
		err := kafkaService.LogSystemInfo("api", req.Message, map[string]interface{}{
			"endpoint": "/api/kafka/test",
			"method":   "POST",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "发送消息失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "消息发送成功",
		})
	}
}

// GetKafkaStatus 获取Kafka状态
func GetKafkaStatus(kafkaService *service.KafkaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "disabled"
		if kafkaService != nil {
			status = "enabled"
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data": gin.H{
				"status": status,
			},
		})
	}
}
