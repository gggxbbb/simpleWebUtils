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
	tepl, err := template.New("ip").Parse(ipHTML)
	if err != nil {
		c.String(500, "failed to parse template")
		return
	}
	ip := c.ClientIP()
	data := analyzeIP(ip)
	err = tepl.Execute(c.Writer, data)
	if err != nil {
		c.String(500, "failed to execute template")
		return
	}
}
