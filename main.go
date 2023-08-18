package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/code/:code", echoCode)

	err := r.Run()
	if err != nil {
		return
	}
}
