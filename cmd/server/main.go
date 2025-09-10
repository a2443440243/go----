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
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	"go-api-scaffold/internal/config"
	"go-api-scaffold/internal/handler"
	"go-api-scaffold/internal/middleware"
	"go-api-scaffold/internal/router"
	"go-api-scaffold/pkg/logger"

	_ "go-api-scaffold/docs" // swagger docs
)

// @title Go API Scaffold
// @version 1.0
// @description 一个基于Go语言的RESTful API脚手架框架
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	appLogger := logger.New(cfg.Log.Level)

	// 暂时跳过数据库连接，直接启动服务器
	log.Println("Starting server without database connection for demo purposes...")

	// 初始化处理器
	healthHandler := handler.NewHealthHandler(appLogger)

	// 初始化路由
	r := router.New()

	// 注册中间件
	r.Use(middleware.Logger(appLogger))
	r.Use(middleware.Recovery(appLogger))
	r.Use(middleware.CORS())

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

// PingHandler godoc
// @Summary Ping测试
// @Description 简单的ping测试接口
// @Tags 示例API
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "成功响应"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// VersionHandler godoc
// @Summary 获取版本信息
// @Description 获取应用程序版本信息
// @Tags 示例API
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "版本信息"
// @Router /version [get]
func VersionHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    cfg.App.Name,
			"version": cfg.App.Version,
			"env":     cfg.App.Environment,
		})
	}
}