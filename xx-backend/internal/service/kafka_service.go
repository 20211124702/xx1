package service

import (
	"context"
	"log"

	"xx-backend/pkg/kafka"
)

type KafkaService struct {
	client *kafka.KafkaClient
}

func NewKafkaService(client *kafka.KafkaClient) *KafkaService {
	return &KafkaService{
		client: client,
	}
}

// LogUserLogin 记录用户登录事件
func (ks *KafkaService) LogUserLogin(userID uint, username string, ip string) error {
	data := map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"ip":       ip,
		"action":   "login",
	}

	return ks.client.SendUserEvent("user_login", data)
}

// LogUserLogout 记录用户登出事件
func (ks *KafkaService) LogUserLogout(userID uint, username string) error {
	data := map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"action":   "logout",
	}

	return ks.client.SendUserEvent("user_logout", data)
}

// LogUserRegister 记录用户注册事件
func (ks *KafkaService) LogUserRegister(userID uint, username string, email string) error {
	data := map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"email":    email,
		"action":   "register",
	}

	return ks.client.SendUserEvent("user_register", data)
}

// LogUserUpdate 记录用户更新事件
func (ks *KafkaService) LogUserUpdate(userID uint, username string, fields map[string]interface{}) error {
	data := map[string]interface{}{
		"user_id":  userID,
		"username": username,
		"action":   "update",
		"fields":   fields,
	}

	return ks.client.SendUserEvent("user_update", data)
}

// LogSystemError 记录系统错误
func (ks *KafkaService) LogSystemError(service string, error string, details map[string]interface{}) error {
	data := map[string]interface{}{
		"service": service,
		"error":   error,
		"details": details,
		"level":   "error",
	}

	return ks.client.SendSystemLog("system_error", data)
}

// LogSystemInfo 记录系统信息
func (ks *KafkaService) LogSystemInfo(service string, message string, details map[string]interface{}) error {
	data := map[string]interface{}{
		"service": service,
		"message": message,
		"details": details,
		"level":   "info",
	}

	return ks.client.SendSystemLog("system_info", data)
}

// StartUserEventConsumer 启动用户事件消费者
func (ks *KafkaService) StartUserEventConsumer(ctx context.Context) {
	go func() {
		err := ks.client.ConsumeUserEvents(ctx, func(message kafka.Message) error {
			log.Printf("Received user event: %s", message.Type)

			// 根据事件类型处理
			switch message.Type {
			case "user_login":
				return ks.handleUserLogin(message)
			case "user_logout":
				return ks.handleUserLogout(message)
			case "user_register":
				return ks.handleUserRegister(message)
			case "user_update":
				return ks.handleUserUpdate(message)
			default:
				log.Printf("Unknown user event type: %s", message.Type)
				return nil
			}
		})

		if err != nil {
			log.Printf("Error consuming user events: %v", err)
		}
	}()
}

// StartSystemLogConsumer 启动系统日志消费者
func (ks *KafkaService) StartSystemLogConsumer(ctx context.Context) {
	go func() {
		err := ks.client.ConsumeSystemLogs(ctx, func(message kafka.Message) error {
			log.Printf("Received system log: %s", message.Type)

			// 根据日志类型处理
			switch message.Type {
			case "system_error":
				return ks.handleSystemError(message)
			case "system_info":
				return ks.handleSystemInfo(message)
			default:
				log.Printf("Unknown system log type: %s", message.Type)
				return nil
			}
		})

		if err != nil {
			log.Printf("Error consuming system logs: %v", err)
		}
	}()
}

// 处理用户登录事件
func (ks *KafkaService) handleUserLogin(message kafka.Message) error {
	// 这里可以添加具体的业务逻辑
	// 例如：更新用户最后登录时间、记录登录统计等
	log.Printf("Handling user login event: %+v", message.Data)
	return nil
}

// 处理用户登出事件
func (ks *KafkaService) handleUserLogout(message kafka.Message) error {
	// 这里可以添加具体的业务逻辑
	// 例如：清理用户会话、记录在线时长等
	log.Printf("Handling user logout event: %+v", message.Data)
	return nil
}

// 处理用户注册事件
func (ks *KafkaService) handleUserRegister(message kafka.Message) error {
	// 这里可以添加具体的业务逻辑
	// 例如：发送欢迎邮件、初始化用户配置等
	log.Printf("Handling user register event: %+v", message.Data)
	return nil
}

// 处理用户更新事件
func (ks *KafkaService) handleUserUpdate(message kafka.Message) error {
	// 这里可以添加具体的业务逻辑
	// 例如：同步用户信息到其他服务、记录变更历史等
	log.Printf("Handling user update event: %+v", message.Data)
	return nil
}

// 处理系统错误
func (ks *KafkaService) handleSystemError(message kafka.Message) error {
	// 这里可以添加具体的业务逻辑
	// 例如：发送告警通知、记录错误统计等
	log.Printf("Handling system error: %+v", message.Data)
	return nil
}

// 处理系统信息
func (ks *KafkaService) handleSystemInfo(message kafka.Message) error {
	// 这里可以添加具体的业务逻辑
	// 例如：记录系统状态、更新监控指标等
	log.Printf("Handling system info: %+v", message.Data)
	return nil
}

// Close 关闭Kafka服务
func (ks *KafkaService) Close() error {
	return ks.client.Close()
}
