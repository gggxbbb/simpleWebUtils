package code

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"strconv"
)

//go:embed code.html
var codeHTML string

func html(c *gin.Context) {
	code := c.Param("code")
	codeInt, err := strconv.Atoi(code)
	if err != nil {
		generateHTML(400, "invalid code", c)
		return
	}
	if code_detail[codeInt] == "" {
		generateHTML(400, "invalid code", c)
		return
	}
	generateHTML(codeInt, code_detail[codeInt], c)

}

type codeDetailS struct {
	Code    int
	Message string
}

func generateHTML(code int, detail string, c *gin.Context) {
	tepl, _ := template.New("code").Parse(codeHTML)
	data := codeDetailS{
		Code:    code,
		Message: detail,
	}
	err := tepl.Execute(c.Writer, data)
	if err != nil {
		return
	}
}
