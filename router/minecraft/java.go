package minecraft

import (
	"github.com/gin-gonic/gin"
	"simpleWebUtils/utils"
)

type minecraftJavaRemote struct {
	Server string `json:"server"`
	Port   string `json:"port"`
}

// motd-java
func utilsMinecraftJava(ctx *gin.Context) {
	var server minecraftJavaRemote

	if ctx.Request.Method == "GET" {
		server = minecraftJavaRemote{
			Server: ctx.Param("server"),
			Port:   ctx.Param("port"),
		}
	} else if ctx.Request.Method == "POST" {
		err := ctx.Bind(&server)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error":       "cannot parse request",
				"description": utils.LocalAddressCleaner(err.Error()),
			})
			return
		}
	}

	data, err := QueryJavaMOTD(server.Server, server.Port)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":       "cannot query java motd",
			"description": utils.LocalAddressCleaner(err.Error()),
		})
		return
	}

	ctx.JSON(200, data)
}
