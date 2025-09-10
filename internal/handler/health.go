package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go-api-scaffold/pkg/logger"
	"go-api-scaffold/pkg/response"
)

// HealthHandler 健康检查处理器
type HealthHandler struct {
	logger logger.Logger
}

// NewHealthHandler 创建健康检查处理器实例
func NewHealthHandler(logger logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
}

var startTime = time.Now()

// Check godoc
// @Summary 基础健康检查
// @Description 检查服务基本健康状态
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse "健康状态信息"
// @Router /health/check [get]
func (h *HealthHandler) Check(c *gin.Context) {
	uptime := time.Since(startTime)

	health := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Uptime:    uptime.String(),
	}

	response.Success(c, health)
}

// Ready godoc
// @Summary 就绪检查
// @Description 检查服务是否准备好接收请求
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "就绪状态"
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	// 这里可以添加依赖服务的检查，如数据库连接等
	// 目前简单返回就绪状态
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"timestamp": time.Now(),
	})
}

// Live godoc
// @Summary 存活检查
// @Description 检查服务是否正在运行
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "存活状态"
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	// 存活检查，确认服务正在运行
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"timestamp": time.Now(),
	})
}