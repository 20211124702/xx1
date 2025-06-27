package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"xx-backend/config"
	"xx-backend/internal/handler"
	"xx-backend/internal/middleware"
	"xx-backend/internal/model"
	"xx-backend/internal/service"
	"xx-backend/pkg/database"
	"xx-backend/pkg/kafka"
	"xx-backend/pkg/redis"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	db := database.InitMySQL(cfg.MySQL)

	// 自动迁移数据库表
	err := db.AutoMigrate(&model.User{}, &model.Role{}, &model.Menu{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化Redis连接
	redisClient := redis.InitRedis(cfg.Redis)

	// 初始化Kafka客户端
	kafkaConfig := &kafka.Config{
		Brokers:         strings.Split(cfg.Kafka.Brokers, ","),
		TopicUserEvents: cfg.Kafka.TopicUserEvents,
		TopicSystemLogs: cfg.Kafka.TopicSystemLogs,
	}

	kafkaClient, err := kafka.NewKafkaClient(kafkaConfig)
	if err != nil {
		log.Printf("Failed to initialize Kafka client: %v", err)
		// 不中断程序启动，Kafka是可选的
	} else {
		// 创建Kafka主题
		if err := kafkaClient.CreateTopics(); err != nil {
			log.Printf("Failed to create Kafka topics: %v", err)
		}
		log.Println("Kafka client initialized successfully")
	}

	// 初始化Kafka服务
	var kafkaService *service.KafkaService
	if kafkaClient != nil {
		kafkaService = service.NewKafkaService(kafkaClient)

		// 启动Kafka消费者
		ctx := context.Background()
		kafkaService.StartUserEventConsumer(ctx)
		kafkaService.StartSystemLogConsumer(ctx)
		log.Println("Kafka consumers started")
	}

	// 初始化服务层
	userService := service.NewUserService(db, redisClient, kafkaService)
	authService := service.NewAuthService(db, redisClient, kafkaService)

	// 初始化gRPC服务器
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// 启动gRPC服务器（在goroutine中）
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 设置Gin模式
	if cfg.App.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 将服务注入到context中
	r.Use(func(c *gin.Context) {
		c.Set("auth_service", authService)
		c.Set("db", db)
		c.Next()
	})

	// 路由组
	api := r.Group("/api")
	{
		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.Login(authService))
			auth.POST("/logout", middleware.AuthMiddleware(), handler.Logout(authService))
			auth.GET("/profile", middleware.AuthMiddleware(), handler.GetProfile(userService))
			auth.POST("/register", handler.Register(userService))
		}

		// 用户管理路由
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("", handler.GetUsers(userService))
			users.GET("/:id", handler.GetUser(userService))
			users.POST("", handler.CreateUser(userService))
			users.PUT("/:id", handler.UpdateUser(userService))
			users.DELETE("/:id", handler.DeleteUser(userService))
		}

		// 角色管理路由
		roles := api.Group("/roles")
		roles.Use(middleware.AuthMiddleware())
		{
			roles.GET("", handler.GetRoles(userService))
			roles.POST("", handler.CreateRole(userService))
			roles.PUT("/:id", handler.UpdateRole(userService))
			roles.DELETE("/:id", handler.DeleteRole(userService))
		}

		// 菜单管理路由
		menus := api.Group("/menus")
		menus.Use(middleware.AuthMiddleware())
		{
			menus.GET("", handler.GetMenus(userService))
			menus.POST("", handler.CreateMenu(userService))
			menus.PUT("/:id", handler.UpdateMenu(userService))
			menus.DELETE("/:id", handler.DeleteMenu(userService))
		}

		// Kafka管理路由
		kafka := api.Group("/kafka")
		{
			kafka.GET("/status", handler.GetKafkaStatus(kafkaService))
			kafka.POST("/test", middleware.AuthMiddleware(), handler.SendTestMessage(kafkaService))
		}
	}

	// 启动HTTP服务器
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server is running on port %s", cfg.App.Port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// 关闭gRPC服务器
	grpcServer.GracefulStop()

	// 关闭Kafka连接
	if kafkaService != nil {
		if err := kafkaService.Close(); err != nil {
			log.Printf("Error closing Kafka service: %v", err)
		}
	}

	log.Println("Server exiting")
}
