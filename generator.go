package mycomarkup

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/globals"
)

const maxRecursionLevel = 3

func generateHTML(ast []blocks.Block, recursionLevel int) (html string) {
	if recursionLevel > maxRecursionLevel {
		return "Transclusion depth limit"
	}
	for _, line := range ast {
		switch v := line.(type) {
		case blocks.List:
			var ret string
			for _, item := range v.Items {
				ret += fmt.Sprintf(markerToTemplate(item.Marker), generateHTML(item.Contents, recursionLevel))
			}
			html += fmt.Sprintf(listToTemplate(v), ret)
		case blocks.Transclusion:
			html += transclusionToHTML(v, recursionLevel)
		case blocks.Formatted, blocks.Paragraph, blocks.Img, blocks.HorizontalLine, blocks.LaunchPad, blocks.Heading, blocks.Table, blocks.TableRow, blocks.CodeBlock, blocks.Quote:
			html += BlockToHTML(v, &blocks.IDCounter{})
		default:
			html += "<b class='error'>Unknown element.</b>"
		}
	}
	return html
}

func transclusionToHTML(xcl blocks.Transclusion, recursionLevel int) string {
	var (
		messageBase = `
<section class="transclusion transclusion_%s">
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
	case globals.CalledInShell:
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
	ctx, _ := mycocontext.ContextFromStringInput(xcl.Target, rawText)
	ast := BlockTree(ctx)
	xclText := generateHTML(ast, recursionLevel+1)
	return fmt.Sprintf(messageOK, xcl.Target, binaryHtml+xclText)
}

func listToTemplate(list blocks.List) string {
	switch list.Marker {
	case blocks.MarkerOrdered:
		return `
<ol>%s</ol>`
	default:
		return `
<ul>%s</ul>`
	}
}

func markerToTemplate(m blocks.ListMarker) string {
	switch m {
	case blocks.MarkerUnordered:
		return `
	<li class="item_unordered">%s</li>`
	case blocks.MarkerOrdered:
		return `
	<li class="item_ordered">%s</li>`
	case blocks.MarkerTodoDone:
		return `
	<li class="item_todo item_todo-done"><input type="checkbox" disabled checked>%s</li>`
	case blocks.MarkerTodo:
		return `
	<li class="item_todo"><input type="checkbox" disabled>%s</li>`
	}
	panic("unreachable")
}
