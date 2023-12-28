package main

import (
	"github.com/gin-gonic/gin"
	"simpleWebUtils/router/analyze"
	"simpleWebUtils/router/echo"
	"simpleWebUtils/router/minecraft"
)

func main() {

	r := gin.Default()
	r.GET("/", readme)

	echo.Init(r)
	analyze.Init(r)
	minecraft.Init(r)

	err := r.Run(":4399")
	if err != nil {
		return
	}
}
