package ip

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	r.GET("/ip", echo)
	r.GET("/ip/analyze", analyze)
	r.GET("/ip/html", html)
}
