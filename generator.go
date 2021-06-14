package mycomarkup

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/globals"
)

const maxRecursionLevel = 3

func generateHTML(ast []blocks.Block, recursionLevel int, counter *blocks.IDCounter) (html string) {
	if recursionLevel > maxRecursionLevel {
		return "Transclusion depth limit"
	}
	for _, line := range ast {
		switch v := line.(type) {
		case blocks.List:
			var ret string
			for _, item := range v.Items {
				ret += fmt.Sprintf(markerToTemplate(item.Marker), generateHTML(item.Contents, recursionLevel, counter.UnusableCopy()))
			}
			html += fmt.Sprintf(listToTemplate(v), idAttribute(v, counter), ret)
		case blocks.Table:
			t := v
			var ret string
			if t.Caption != "" {
				ret = fmt.Sprintf("<caption>%s</caption>", t.Caption)
			}
			ret += "<tbody>\n"
			for _, tr := range t.Rows {
				ret += "<tr>"
				for _, tc := range tr.Cells {
					ret += fmt.Sprintf(
						"\n\t<%[1]s>%[2]s</%[1]s>",
						tc.TagName(),
						generateHTML(tc.Contents, recursionLevel, counter.UnusableCopy()),
					)
				}
				ret += "</tr>\n"
			}
			html += fmt.Sprintf(`
<table%s>%s</tbody></table>`, idAttribute(t, counter), ret)
		case blocks.Transclusion:
			html += transclusionToHTML(v, recursionLevel, counter.UnusableCopy())
		case blocks.Formatted, blocks.Paragraph, blocks.Img, blocks.HorizontalLine, blocks.LaunchPad, blocks.Heading, blocks.CodeBlock, blocks.Quote:
			html += BlockToHTML(v, counter)
		default:
			html += "<b class='error'>Unknown element.</b>"
		}
	}
	return html
}

func transclusionToHTML(xcl blocks.Transclusion, recursionLevel int, counter *blocks.IDCounter) string {
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
		_ = fmt.Sprintf(messageBase, "failed",
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
	}

	rawText, binaryHtml, err := globals.HyphaAccess(xcl.Target)
	if err != nil {
		return fmt.Sprintf(messageNotExists, xcl.Target)
	}
	ctx, _ := mycocontext.ContextFromStringInput(xcl.Target, rawText)
	ast := BlockTree(ctx)
	xclText := generateHTML(ast, recursionLevel+1, counter)
	return fmt.Sprintf(messageOK, xcl.Target, binaryHtml+xclText)
}

func listToTemplate(list blocks.List) string {
	switch list.Marker {
	case blocks.MarkerOrdered:
		return `
<ol%s>%s</ol>`
	default:
		return `
<ul%s>%s</ul>`
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
