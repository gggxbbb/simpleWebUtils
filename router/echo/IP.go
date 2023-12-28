package echo

import "github.com/gin-gonic/gin"

func echoIP(c *gin.Context) {
	ip := c.ClientIP()
	c.Data(200, "text/plain", []byte(ip))
}
