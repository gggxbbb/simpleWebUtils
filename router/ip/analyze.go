package ip

import "github.com/gin-gonic/gin"

func analyze(c *gin.Context) {
	ip := c.ClientIP()
	data := analyzeIP(ip)
	c.JSON(200, data)
}
