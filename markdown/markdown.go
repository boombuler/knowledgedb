package markdown

import (
	"html/template"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

const enabled_md_extensions = blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS

func ToPlaintext(content string) string {
	return string(blackfriday.Markdown([]byte(content), new(PlaintextRenderer), enabled_md_extensions))
}

func ToHtml(content string) template.HTML {
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS

	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	html := blackfriday.Markdown([]byte(content), renderer, enabled_md_extensions)
	return template.HTML(bluemonday.UGCPolicy().SanitizeBytes(html))
}
