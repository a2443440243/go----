package handler

import (
	"github.com/gin-gonic/gin"
	"go-api-scaffold/internal/service"
	"go-api-scaffold/pkg/logger"
)

// Handler 处理器集合
type Handler struct {
	User   *UserHandler
	Health *HealthHandler
	logger logger.Logger
}

// New 创建处理器实例
func New(service *service.Service, logger logger.Logger) *Handler {
	return &Handler{
		User:   NewUserHandler(service.User, logger),
		Health: NewHealthHandler(logger),
		logger: logger,
	}
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	// 健康检查
	health := router.Group("/health")
	{
		health.GET("/", h.Health.Check)
		health.GET("/ready", h.Health.Ready)
		health.GET("/live", h.Health.Live)
	}

	// API路由组
	api := router.Group("/api/v1")
	{
		// 用户相关路由
		users := api.Group("/users")
		{
			users.POST("/", h.User.Create)
			users.GET("/:id", h.User.GetByID)
			users.PUT("/:id", h.User.Update)
			users.DELETE("/:id", h.User.Delete)
			users.GET("/", h.User.List)
		}

		// 认证相关路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.User.Login)
		}
	}
}
