package mycomarkup

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/genhtml"
	"github.com/bouncepaw/mycomarkup/v3/genhtml/tag"
	"github.com/bouncepaw/mycomarkup/v3/globals"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"github.com/bouncepaw/mycomarkup/v3/util"
)

const maxRecursionLevel = 3

// V3 Kinda hard to get rid of that
func generateHTML(ctx mycocontext.Context, ast []blocks.Block, counter *blocks.IDCounter) (html string) {
	if ctx.RecursionLevel() > maxRecursionLevel {
		return "Transclusion depth limit"
	}
	for _, line := range ast {
		switch v := line.(type) {
		case blocks.Quote:
			html += tag.NewClosed(
				"blockquote",
				map[string]string{"id": v.ID(counter)},
				generateHTML(ctx, v.Contents(), counter.UnusableCopy()),
			).String()
		case blocks.List:
			var ret string
			for _, item := range v.Items {
				ret += fmt.Sprintf(
					markerToTemplate(item.Marker),
					generateHTML(ctx, item.Contents, counter.UnusableCopy()),
				)
			}
			html += fmt.Sprintf(listToTemplate(v), idAttribute(v, counter), ret)
		case blocks.Table:
			var ret string
			if v.Caption() != "" {
				ret = fmt.Sprintf("<caption>%s</caption>", v.Caption())
			}
			ret += "<tbody>\n"
			for _, tr := range v.Rows() {
				ret += "<tr>"
				for _, tc := range tr.Cells() {
					ret += fmt.Sprintf(
						"\n\t<%[1]s%[3]s>%[2]s</%[1]s>",
						util.TernaryConditionString(tc.IsHeaderCell(), "th", "td"),
						generateHTML(ctx, tc.Contents(), counter.UnusableCopy()),
						util.TernaryConditionString(
							tc.Colspan() <= 1,
							"",
							fmt.Sprintf(` colspan="%d"`, tc.Colspan()),
						),
					)
				}
				ret += "</tr>\n"
			}
			html += fmt.Sprintf(`
<table%s>%s</tbody></table>`, idAttribute(v, counter), ret)
		case blocks.Transclusion:
			html += transclusionToHTML(ctx, v, counter.UnusableCopy())
		default:
			html += genhtml.BlockToTag(ctx, v, counter).String()
		}
	}
	return html
}

func transclusionToHTML(ctx mycocontext.Context, xcl blocks.Transclusion, counter *blocks.IDCounter) string {
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

	// Nothing will match if there is no error:
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

	// V4 This part is awful
	// Now, to real transclusion:
	rawText, binaryHtml, err := globals.HyphaAccess(xcl.Target)
	if err != nil {
		return fmt.Sprintf(messageNotExists, xcl.Target)
	}
	xclVisistor, result := transclusionVisitor(xcl)
	xclctx, _ := mycocontext.ContextFromStringInput(xcl.Target, rawText) // FIXME: it will bite us one day UPDATE: is it the day? I don't feel the bite.
	_ = BlockTree(xclctx, xclVisistor)
	xclText := generateHTML(ctx.WithIncrementedRecursionLevel(), result(), counter.UnusableCopy())

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

func idAttribute(b blocks.Block, counter *blocks.IDCounter) string {
	switch id := b.ID(counter); {
	case !counter.ShouldUseResults(), id == "":
		return ""
	default:
		return fmt.Sprintf(` id="%s"`, id)
	}
}
