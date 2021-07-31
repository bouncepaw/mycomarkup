package mycomarkup

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/globals"
	"github.com/bouncepaw/mycomarkup/mycocontext"
)

const maxRecursionLevel = 3

func generateHTML(ast []blocks.Block, recursionLevel int, counter *blocks.IDCounter) (html string) {
	if recursionLevel > maxRecursionLevel {
		return "Transclusion depth limit"
	}
	for _, line := range ast {
		switch v := line.(type) {
		case blocks.Quote:
			html += fmt.Sprintf(
				"\n<blockquote%s>%s\n</blockquote>",
				idAttribute(v, counter.UnusableCopy()),
				generateHTML(v.Contents(), recursionLevel, counter.UnusableCopy()),
			)
		case blocks.List:
			var ret string
			for _, item := range v.Items {
				ret += fmt.Sprintf(markerToTemplate(item.Marker), generateHTML(item.Contents, recursionLevel, counter.UnusableCopy()))
			}
			html += fmt.Sprintf(listToTemplate(v), idAttribute(v, counter), ret)
		case blocks.Table:
			var ret string
			if v.Caption != "" {
				ret = fmt.Sprintf("<caption>%s</caption>", v.Caption)
			}
			ret += "<tbody>\n"
			for _, tr := range v.Rows {
				ret += "<tr>"
				for _, tc := range tr.Cells {
					ret += fmt.Sprintf(
						"\n\t<%[1]s%[3]s>%[2]s</%[1]s>",
						tc.TagName(),
						generateHTML(tc.Contents, recursionLevel, counter.UnusableCopy()),
						tc.ColspanAttribute(),
					)
				}
				ret += "</tr>\n"
			}
			html += fmt.Sprintf(`
<table%s>%s</tbody></table>`, idAttribute(v, counter), ret)
		case blocks.Transclusion:
			html += transclusionToHTML(v, recursionLevel, counter.UnusableCopy())
		case blocks.Formatted, blocks.Paragraph, blocks.Img, blocks.HorizontalLine, blocks.LaunchPad, blocks.Heading, blocks.CodeBlock:
			html += BlockToHTML(v, counter)
		default:
			html += "<v class='error'>Unknown element.</v>"
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
	<p class="error">Cannot transclude hypha <a class="wikilink wikilink_new" href="/hypha/%[1]s">%[1]s</a> because it does not exist</p>
</section>`
		messageOK = `<section class="transclusion transclusion_ok%[3]s">
	<a class="transclusion__link" href="/hypha/%[1]s">%[1]s</a>
	<div class="transclusion__content">%[2]s</div>
</section>`
	)

	// Nonthing will match if there is no error:
	switch xcl.TransclusionError.Reason {
	case blocks.TransclusionErrorNotExists:
		return fmt.Sprintf(messageNotExists, xcl.Target)
	case blocks.TransclusionErrorNoTarget:
		return messageNoTarget
	case blocks.TransclusionInTerminal:
		return messageCLI
	case blocks.TransclusionErrorOldSyntax:
		return messageOldSyntax
	}

	// Now, to real transclusion:
	rawText, binaryHtml, err := globals.HyphaAccess(xcl.Target)
	if err != nil {
		return fmt.Sprintf(messageNotExists, xcl.Target)
	}
	xclVisistor, result := transclusionVisitor(xcl)
	ctx, _ := mycocontext.ContextFromStringInput(xcl.Target, rawText) // FIXME: it will bite us one day
	_ = BlockTree(ctx, xclVisistor)                                   // TODO: inject transclusion visitors here
	xclText := generateHTML(result(), recursionLevel+1, counter.UnusableCopy())

	if xcl.Selector == blocks.SelectorAttachment || xcl.Selector == blocks.SelectorFull || xcl.Selector == blocks.SelectorOverview {
		xclText = binaryHtml + xclText
	}
	return fmt.Sprintf(
		messageOK,
		xcl.Target,
		xclText,
		func() string {
			if xcl.Blend {
				return " transclusion_blend"
			}
			return " transclusion_stand-out"
		}(),
	)
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
