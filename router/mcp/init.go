package mcp

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	r.GET("/mcp", index)
	r.POST("/mcp", handleJSONRPC)
}
