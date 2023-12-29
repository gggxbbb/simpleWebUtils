package code

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	r.GET("/code/:code", echo)
	r.GET("/code/detail/:code", codeDetail)
}
