package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go-api-scaffold/pkg/logger"
	"go-api-scaffold/pkg/response"
)

// Recovery 错误恢复中间件
func Recovery(log logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.WithFields(map[string]interface{}{
				"error":      err,
				"stack":      string(debug.Stack()),
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}).Error("Panic recovered")
		}

		if err, ok := recovered.(error); ok {
			log.WithFields(map[string]interface{}{
				"error":      err.Error(),
				"stack":      string(debug.Stack()),
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"client_ip":  c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}).Error("Panic recovered")
		}

		log.WithFields(map[string]interface{}{
			"error":      fmt.Sprintf("%v", recovered),
			"stack":      string(debug.Stack()),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Error("Panic recovered")

		// 返回500错误
		response.ServerError(c, "Internal server error")
		c.Abort()
	})
}