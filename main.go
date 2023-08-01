package main

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {

	r := gin.Default()

	r.GET("/code/:code", func(c *gin.Context) {
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
	})

	err := r.Run()
	if err != nil {
		return
	}
}
