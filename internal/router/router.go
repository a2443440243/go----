package router

import (
	"github.com/gin-gonic/gin"
)

// New 创建新的路由器
func New() *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	
	// 创建路由器
	r := gin.New()
	
	return r
}