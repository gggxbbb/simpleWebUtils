package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/", readme)

	r.GET("/echo/code/:code", echoCode)
	r.GET("/echo/ua", echoUA)
	r.GET("/echo/ip", echoIP)

	r.GET("/analyze/ua", analyzeUA)

	r.GET("utils/minecraft/bedrock/:server", utilsMinecraftBedrock)
	r.GET("/utils/minecraft/bedrock/:server/:port", utilsMinecraftBedrock)

	err := r.Run(":4399")
	if err != nil {
		return
	}
}
