package ua

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"html/template"
)

//go:embed ua.html
var uaHTML string

func html(c *gin.Context) {
	generateHTML(c)
}

func generateHTML(c *gin.Context) {
	tepl, _ := template.New("code").Parse(uaHTML)
	data := useragent.Parse(c.GetHeader("User-Agent"))
	err := tepl.Execute(c.Writer, data)
	if err != nil {
		return
	}
}
