package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func echoCode(c *gin.Context) {
	code := c.Param("code")
	// parse code to int
	intCode, err := strconv.Atoi(code)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid code",
		})
		return
	}
	c.Data(intCode, "", nil)
}

func echoUA(c *gin.Context) {
	ua := c.GetHeader("User-Agent")
	c.Data(200, "text/plain", []byte(ua))
}
