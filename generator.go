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
		return tag.NewClosed("section").
			WithAttrs(map[string]string{
				"class": "transclusion transclusion_failed transclusion_not-exists",
			}).
			WithChildren(tag.NewClosed("p").
				WithContentsStrings("Transclusion depth limit")).
			String()
	}
	for _, line := range ast {
		switch v := line.(type) {
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

	return tag.NewClosed("section").
		WithAttrs(map[string]string{
			"id": xcl.ID(counter),
			"class": "transclusion transclusion_ok transclusion_" + util.TernaryConditionString(
				xcl.Blend,
				"blend",
				"stand-out",
			),
		}).
		WithChildren(
			tag.NewClosed("a").
				WithContentsStrings(xcl.Target).
				WithAttrs(map[string]string{
					"class": "transclusion__link",
					"href":  "/hypha/" + xcl.Target,
				}),
			tag.NewClosed("div").
				WithContentsStrings(xclText).
				WithAttrs(map[string]string{
					"class": "transclusion__content",
				}),
		).
		String()
}

func idAttribute(b blocks.Block, counter *blocks.IDCounter) string {
	switch id := b.ID(counter); {
	case !counter.ShouldUseResults(), id == "":
		return ""
	default:
		return fmt.Sprintf(` id="%s"`, id)
	}
}
