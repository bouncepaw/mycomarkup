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
		return tag.NewClosed("section", map[string]string{
			"class": "transclusion transclusion_failed transclusion_not-exists",
		}, "",
			tag.NewClosed("p", map[string]string{}, "Transclusion depth limit")).String()
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
	if xcl.HasError() {
		return genhtml.MapTransclusionErrorToTag(xcl).String()
	}

	// V3
	// V4 This part is awful, bloody hell. Move to the parser module
	// Now, to real transclusion:
	rawText, binaryHtml, err := globals.HyphaAccess(xcl.Target)
	if err != nil {
		xcl.TransclusionError.Reason = blocks.TransclusionErrorNotExists
		return genhtml.MapTransclusionErrorToTag(xcl).String()
	}
	xclVisistor, result := transclusionVisitor(xcl)
	xclctx, _ := mycocontext.ContextFromStringInput(xcl.Target, rawText) // FIXME: it will bite us one day UPDATE: is it the day? I don't feel the bite.
	_ = BlockTree(xclctx, xclVisistor)                                   // Call for side-effects
	xclText := generateHTML(ctx.WithIncrementedRecursionLevel(), result(), counter.UnusableCopy())

	if xcl.Selector == blocks.SelectorAttachment || xcl.Selector == blocks.SelectorFull || xcl.Selector == blocks.SelectorOverview {
		xclText = binaryHtml + xclText
	}

	return tag.NewClosed("section", map[string]string{
		"id": xcl.ID(counter),
		"class": "transclusion transclusion_ok transclusion_" + util.TernaryConditionString(
			xcl.Blend,
			"blend",
			"stand-out",
		),
	}, "",
		tag.NewClosed("a", map[string]string{
			"class": "transclusion__link",
			"href":  "/hypha/" + xcl.Target,
		}, xcl.Target),
		tag.NewClosed("div", map[string]string{
			"class": "transclusion__content",
		}, xclText),
	).String()
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
