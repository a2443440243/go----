package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go-api-scaffold/pkg/logger"
)

// Logger 日志中间件
func Logger(log logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 记录请求日志
		log.WithFields(map[string]interface{}{
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
			"status":      param.StatusCode,
			"latency":     param.Latency,
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"user_agent":  param.Request.UserAgent(),
			"body_size":   param.BodySize,
			"error":       param.ErrorMessage,
		}).Info("HTTP Request")

		return ""
	})
}