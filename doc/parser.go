package doc

import (
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/generator"
)

const maxRecursionLevel = 3

func Parse(ast []Line, from, to int, recursionLevel int) (html string) {
	if recursionLevel > maxRecursionLevel {
		return "Transclusion depth limit"
	}
	for _, line := range ast {
		if line.Id >= from && (line.Id <= to || to == 0) || line.Id == -1 {
			switch v := line.Contents.(type) {
			case Transclusion:
				html += Transclude(v, recursionLevel)
			case blocks.Img:
				html += generator.BlockToHTML(v)
			case blocks.Table:
				html += v.AsHtml()
			case *blocks.List:
				html += v.RenderAsHtml()
			case string:
				html += v
			default:
				html += "<b class='error'>Unknown element.</b>"
			}
		}
	}
	return html
}
