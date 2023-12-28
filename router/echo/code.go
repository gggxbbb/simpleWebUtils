package echo

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
