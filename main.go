package main

import (
	"github.com/gin-gonic/gin"
	"simpleWebUtils/router/code"
	"simpleWebUtils/router/ip"
	"simpleWebUtils/router/minecraft"
	"simpleWebUtils/router/ua"
)

func main() {

	r := gin.Default()
	r.GET("/", readme)

	ua.Init(r)
	ip.Init(r)
	code.Init(r)
	minecraft.Init(r)

	err := r.Run(":4399")
	if err != nil {
		return
	}
}
