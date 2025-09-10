package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go-api-scaffold/internal/config"
	"go-api-scaffold/internal/handler"
	"go-api-scaffold/internal/middleware"
	"go-api-scaffold/internal/repository"
	"go-api-scaffold/internal/service"
	"go-api-scaffold/pkg/database"
	"go-api-scaffold/pkg/logger"
)

// App 应用结构体
type App struct {
	config *config.Config
	server *http.Server
	logger logger.Logger
}

// New 创建新的应用实例
func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
	}
}

// Run 启动应用
func (a *App) Run() error {
	// 初始化日志
	a.logger = logger.New(a.config.Log.Level)

	// 初始化数据库
	db, err := database.New(a.config.Database)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 初始化仓储层
	repos := repository.New(db)

	// 初始化服务层
	services := service.New(repos, a.logger)

	// 初始化处理器
	handlers := handler.New(services, a.logger)

	// 设置Gin模式
	if a.config.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.New()

	// 注册中间件
	router.Use(middleware.Logger(a.logger))
	router.Use(middleware.Recovery(a.logger))
	router.Use(middleware.CORS())

	// 注册路由
	handlers.RegisterRoutes(router)

	// 创建HTTP服务器
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.Server.Port),
		Handler: router,
	}

	// 启动服务器
	go func() {
		a.logger.Info(fmt.Sprintf("Server starting on port %d", a.config.Server.Port))
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error(fmt.Sprintf("Server failed to start: %v", err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Server shutting down...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	a.logger.Info("Server exited")
	return nil
}