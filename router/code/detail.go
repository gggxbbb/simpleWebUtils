package code

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func codeDetail(c *gin.Context) {
	code := c.Param("code")
	// parse code to int
	intCode, err := strconv.Atoi(code)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid code",
		})
		return
	}
	if code_detail[intCode] == "" {
		c.JSON(400, gin.H{
			"error": "invalid code",
			"see":   "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status",
		})
		return
	}
	c.JSON(intCode, gin.H{
		"code":        intCode,
		"description": code_detail[intCode],
		"see":         "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/" + code,
		"cat":         "https://http.cat/" + code,
	})

}
