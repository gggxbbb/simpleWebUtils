package echo

import "github.com/gin-gonic/gin"

func echoUA(c *gin.Context) {
	ua := c.GetHeader("User-Agent")
	c.Data(200, "text/plain", []byte(ua))
}
