package main

import (
	_ "embed"
	"github.com/gin-gonic/gin"
)

//go:embed readme.html
var readmeHTML string

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html", []byte(readmeHTML))
	})

	r.GET("/echo/code/:code", echoCode)
	r.GET("/echo/ua", echoUA)

	r.GET("/analyze/ua", analyzeUA)

	err := r.Run()
	if err != nil {
		return
	}
}
