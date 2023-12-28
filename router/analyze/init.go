package analyze

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	r.GET("/analyze/ua", analyzeUA)
}
