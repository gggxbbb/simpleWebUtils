package ua

import (
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
)

func analyze(c *gin.Context) {
	ua := c.GetHeader("User-Agent")
	uaData := useragent.Parse(ua)
	c.JSON(200, uaData)
}
