package main

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
)
import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed README.md
var readmeMD string

//go:embed readmeTemplate.html
var readmeTemplate string

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func readme(c *gin.Context) {
	h := mdToHTML([]byte(readmeMD))
	t := fmt.Sprintf(readmeTemplate, string(h))
	c.Data(200, "text/html", []byte(t))
}
