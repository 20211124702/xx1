package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"xx-backend/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
	db           *gorm.DB
	redis        *redis.Client
	kafkaService *KafkaService
}

func NewAuthService(db *gorm.DB, redis *redis.Client, kafkaService *KafkaService) *AuthService {
	return &AuthService{
		db:           db,
		redis:        redis,
		kafkaService: kafkaService,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (s *AuthService) Login(req *LoginRequest, c *gin.Context) (*LoginResponse, error) {
	var user model.User

	// 查询用户
	if err := s.db.Preload("Role").Where("username = ?", req.Username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 验证密码
	if !s.checkPassword(req.Password, user.Password) {
		return nil, fmt.Errorf("密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, fmt.Errorf("用户已被禁用")
	}

	// 生成JWT token
	token, err := s.generateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	// 将token存储到Redis
	ctx := context.Background()
	err = s.redis.Set(ctx, fmt.Sprintf("token:%d", user.ID), token, 24*time.Hour).Err()
	if err != nil {
		return nil, err
	}

	// 记录登录事件到Kafka
	clientIP := s.getClientIP(c)
	if s.kafkaService != nil {
		if err := s.kafkaService.LogUserLogin(user.ID, user.Username, clientIP); err != nil {
			// 记录Kafka错误但不影响登录流程
			fmt.Printf("Failed to log user login to Kafka: %v\n", err)
		}
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Logout(userID int, username string) error {
	ctx := context.Background()
	err := s.redis.Del(ctx, fmt.Sprintf("token:%d", userID)).Err()
	if err != nil {
		return err
	}

	// 记录登出事件到Kafka
	if s.kafkaService != nil {
		if err := s.kafkaService.LogUserLogout(uint(userID), username); err != nil {
			// 记录Kafka错误但不影响登出流程
			fmt.Printf("Failed to log user logout to Kafka: %v\n", err)
		}
	}

	return nil
}

func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int(claims["user_id"].(float64))

		// 检查Redis中是否存在token
		ctx := context.Background()
		exists, err := s.redis.Exists(ctx, fmt.Sprintf("token:%d", userID)).Result()
		if err != nil || exists == 0 {
			return 0, fmt.Errorf("token已过期")
		}

		return userID, nil
	}

	return 0, fmt.Errorf("无效的token")
}

func (s *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key"))
}

func (s *AuthService) checkPassword(inputPassword, storedPassword string) bool {
	// 简单的MD5密码验证，实际项目中应该使用bcrypt等更安全的方式
	hash := md5.Sum([]byte(inputPassword))
	return hex.EncodeToString(hash[:]) == storedPassword
}

func (s *AuthService) getClientIP(c *gin.Context) string {
	// 尝试从各种头部获取真实IP
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return ip
	}
	return c.ClientIP()
}
