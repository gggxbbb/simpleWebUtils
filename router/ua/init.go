package ua

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	r.GET("/ua", echo)
	r.GET("/ua/analyze", analyze)
}
