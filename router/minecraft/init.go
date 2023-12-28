package minecraft

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	r.GET("/minecraft/bedrock/:server", utilsMinecraftBedrock)
	r.GET("/minecraft/bedrock/:server/:port", utilsMinecraftBedrock)
	r.POST("/minecraft/bedrock", utilsMinecraftBedrock)
}
