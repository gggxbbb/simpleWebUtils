package echo

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	r.GET("/echo/code/:code", echoCode)
	r.GET("/echo/ua", echoUA)
	r.GET("/echo/ip", echoIP)
}
