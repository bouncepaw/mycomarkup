package doc

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/generator"
	"github.com/bouncepaw/mycomarkup/globals"
	"strings"
)

const maxRecursionLevel = 3

func Parse(ast []Token, recursionLevel int) (html string) {
	if recursionLevel > maxRecursionLevel {
		return "Transclusion depth limit"
	}
	for _, line := range ast {
		switch v := line.Value.(type) {
		case blocks.Transclusion:
			html += transclusionToHTML(v, recursionLevel)
		case blocks.Img, blocks.HorizontalLine, blocks.LaunchPad, blocks.Heading, blocks.Table, blocks.CodeBlock:
			html += generator.BlockToHTML(v)
		case *blocks.List:
			html += v.RenderAsHtml()
		case string:
			html += v
		default:
			html += "<b class='error'>Unknown element.</b>"
		}
	}
	return html
}

func transclusionToHTML(xcl blocks.Transclusion, recursionLevel int) string {
	var (
		messageBase = `<section class="transclusion transclusion_%s">
	%s
</section>`
		messageCLI = fmt.Sprintf(messageBase, "failed",
			`<p>Transclusion is not supported in documents generated using Mycomarkup CLI</p>`)
		messageNoTarget = fmt.Sprintf(messageBase, "failed",
			`<p>Transclusion target not specified</p>`)
		messageOldSyntax = fmt.Sprintf(messageBase, "failed",
			`<p>This transclusion is using the old syntax. Please update it to the new one</p>`)
		messageGenericError = fmt.Sprintf(messageBase, "failed",
			`<p>An error occured while transcluding</p>`)
		messageNotExists = `<section class="transclusion transclusion_failed">
	<p class="error">Hypha <a class="wikilink wikilink_new" href="/hypha/%[1]s">%[1]s</a> does not exist</p>
</section>`
		messageOK = `<section class="transclusion transclusion_ok">
	<a class="transclusion__link" href="/hypha/%[1]s">%[1]s</a>
	<div class="transclusion__content">%[2]s</div>
</section>`
	)

	switch {
	case globals.UseBatch:
		return messageCLI
	case xcl.Target == "":
		return messageNoTarget
	case strings.Contains(xcl.Target, ":"):
		return messageOldSyntax
	case xcl.Selector == blocks.TransclusionError:
		return messageGenericError
	}

	rawText, binaryHtml, err := globals.HyphaAccess(xcl.Target)
	if err != nil {
		return fmt.Sprintf(messageNotExists, xcl.Target)
	}
	md := Doc(xcl.Target, rawText)
	xclText := Parse(md.Lex(recursionLevel+1), recursionLevel+1)
	return fmt.Sprintf(messageOK, xcl.Target, binaryHtml+xclText)
}
