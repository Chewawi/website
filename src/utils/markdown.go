package utils

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// MarkdownToHTML converts Markdown content to HTML
func MarkdownToHTML(md string) string {
	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	
	// Parse markdown to AST
	doc := p.Parse([]byte(md))
	
	// Create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	
	// Render AST to HTML
	return string(markdown.Render(doc, renderer))
}