package ip

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"html/template"
)

//go:embed ip.html
var ipHTML string

func html(c *gin.Context) {
	generateHTML(c)
}

func generateHTML(c *gin.Context) {
	tepl, _ := template.New("ip").Parse(ipHTML)
	ip := c.ClientIP()
	data := analyzeIP(ip)
	err := tepl.Execute(c.Writer, data)
	if err != nil {
		return
	}
}
