package ip

import "github.com/gin-gonic/gin"

func echo(c *gin.Context) {
	ip := c.ClientIP()
	c.String(200, ip)
}
