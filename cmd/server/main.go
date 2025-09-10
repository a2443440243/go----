package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"go-api-scaffold/internal/config"
	"go-api-scaffold/internal/handler"
	"go-api-scaffold/internal/middleware"
	"go-api-scaffold/internal/model"
	"go-api-scaffold/internal/repository"
	"go-api-scaffold/internal/router"
	"go-api-scaffold/internal/service"
	"go-api-scaffold/pkg/database"
	"go-api-scaffold/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	appLogger := logger.New(cfg.Log.Level)

	// 连接数据库
	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	appLogger.Info("Database connected successfully")

	// 自动迁移数据库表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	appLogger.Info("Database migration completed")

	// 初始化仓储层
	userRepo := repository.NewUserRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo, appLogger)

	// 初始化处理器
	healthHandler := handler.NewHealthHandler(appLogger)
	userHandler := handler.NewUserHandler(userService, appLogger)

	// 初始化路由
	r := router.New()

	// 注册中间件
	r.Use(middleware.Logger(appLogger))
	r.Use(middleware.Recovery(appLogger))
	r.Use(middleware.CORS())

	// 注册路由
	v1 := r.Group("/api/v1")
	{
		// 健康检查
		health := v1.Group("/health")
		{
			health.GET("/check", healthHandler.Check)
			health.GET("/ready", healthHandler.Ready)
			health.GET("/live", healthHandler.Live)
		}

		// 用户相关接口
		users := v1.Group("/users")
		{
			users.POST("/", userHandler.Create)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.GET("/", userHandler.List)
		}

		// 认证相关接口
		v1.POST("/login", userHandler.Login)

		// 示例API
		v1.GET("/ping", PingHandler)
		v1.GET("/version", VersionHandler(cfg))
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 启动服务器
	go func() {
		log.Printf("Server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// PingHandler Ping测试接口
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// VersionHandler 获取版本信息
func VersionHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    cfg.App.Name,
			"version": cfg.App.Version,
			"env":     cfg.App.Environment,
		})
	}
}
