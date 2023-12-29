package ua

import "github.com/gin-gonic/gin"

func echo(c *gin.Context) {
	ua := c.GetHeader("User-Agent")
	c.String(200, ua)
}
