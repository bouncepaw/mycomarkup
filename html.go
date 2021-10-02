package mycomarkup

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/genhtml"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"github.com/bouncepaw/mycomarkup/v3/parser"
	"github.com/bouncepaw/mycomarkup/v3/util"
)

// BlockToHTML turns the given block into HTML. It supports only a subset of Mycomarkup.
func BlockToHTML(ctx mycocontext.Context, block blocks.Block, counter *blocks.IDCounter) string {
	switch b := block.(type) {
	case blocks.Formatted, blocks.HorizontalLine, blocks.Paragraph, blocks.RocketLink, blocks.LaunchPad, blocks.CodeBlock, blocks.Heading, blocks.ImgEntry:
		return genhtml.BlockToTag(ctx, b, counter).String()

	case blocks.Img:
		return imgToHTML(ctx, b, counter)

	}

	fmt.Printf("%q\n", block)
	return "<b>UNKNOWN ELEMENT</b>"
}

func idAttribute(b blocks.Block, counter *blocks.IDCounter) string {
	switch id := b.ID(counter); {
	case !counter.ShouldUseResults(), id == "":
		return ""
	default:
		return fmt.Sprintf(` id="%s"`, id)
	}
}

func imgEntryToHTML(ctx mycocontext.Context, entry blocks.ImgEntry, counter *blocks.IDCounter) string {
	var ret string
	if entry.Target.IsBlueLink() {
		ret += fmt.Sprintf(
			`<a href="%s"><img src="%s"%s%s></a>`,
			entry.Target.Href(),
			entry.Target.ImgSrc(),
			util.TernaryConditionString(
				entry.GetWidth() == "",
				"",
				` width="`+entry.GetWidth()+`"`,
			),
			util.TernaryConditionString(
				entry.GetHeight() == "",
				"",
				` height="`+entry.GetHeight()+`"`,
			),
		)
	} else {
		ret += fmt.Sprintf(
			`<a class="%s" href="%s">Hypha <i>%s</i> does not exist</a>`,
			entry.Target.Classes(),
			entry.Target.Href(),
			entry.Target.TargetHypha())
	}
	return fmt.Sprintf(
		`<figure class="img-gallery__entry">
	%s
	<figcaption>%s</figcaption>
</figure>
`,
		ret,
		BlockToHTML(ctx, parser.MakeFormatted(entry.Description, entry.HyphaName), counter),
	)
}

func imgToHTML(ctx mycocontext.Context, img blocks.Img, counter *blocks.IDCounter) string {
	img.MarkExistenceOfSrcLinks()
	var ret string
	for _, entry := range img.Entries {
		ret += BlockToHTML(ctx, entry, counter)
	}
	return fmt.Sprintf(`<section class="img-gallery %s"%s>
%s</section>`,
		func() string {
			if img.HasOneImage() {
				return "img-gallery_one-image"
			}
			return "img-gallery_many-images"
		}(),
		idAttribute(img, counter),
		ret)
}
