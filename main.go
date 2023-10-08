package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/", readme)

	r.GET("/echo/code/:code", echoCode)
	r.GET("/echo/ua", echoUA)

	r.GET("/analyze/ua", analyzeUA)

	err := r.Run(":4399")
	if err != nil {
		return
	}
}
