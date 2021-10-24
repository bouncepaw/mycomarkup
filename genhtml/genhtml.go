// Package genhtml provides utilities for converting Mycomarkup blocks into HTML documents. As of now, some parts of HTML generation are in other parts of the library, WIP.
package genhtml

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v3/util"
	"github.com/bouncepaw/mycomarkup/v3/util/lines"
	"html"
	"strings"

	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/genhtml/tag"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"github.com/bouncepaw/mycomarkup/v3/parser"
)

// This package shall not depend on anything other than blocks, links, mycocontext, util, tag.

// BlockToTag turns the given Block into a Tag depending on the Context and IDCounter.
func BlockToTag(ctx mycocontext.Context, block blocks.Block, counter *blocks.IDCounter) tag.Tag {
	var attrs = map[string]string{}
	if counter.ShouldUseResults() {
		attrs["id"] = block.ID(counter)
	}

	switch block := block.(type) {
	case blocks.Formatted:
		var (
			contents string
			tagState = blocks.CleanStyleState()
		)
		for i, line := range block.Lines {
			if i > 0 {
				contents += `<br>`
			}

			for _, span := range line {
				switch s := span.(type) {
				case blocks.SpanTableEntry:
					contents += blocks.TagFromState(s.Kind(), tagState)

				case blocks.InlineLink:
					contents += tag.NewClosed("a").
						WithAttrs(map[string]string{
							"href":  s.Href(),
							"class": s.Classes(),
						}).
						WithContentsStrings(s.Display()).
						String()

				case blocks.InlineText:
					contents += html.EscapeString(s.Contents)

				default:
					panic("unknown span")
				}
			}

			for stt, open := range tagState { // Close the unclosed
				if open {
					contents += blocks.TagFromState(stt, tagState)
				}
			}
		}
		// TODO: fumu fumu
		return tag.NewWrapper().WithContentsStrings(contents)

	case blocks.Paragraph:
		return tag.NewClosed("p").
			WithAttrs(attrs).
			WithChildren(BlockToTag(ctx, block.Formatted, counter))

	case blocks.Heading:
		return tag.NewClosed(fmt.Sprintf("h%d", block.Level())).
			WithAttrs(attrs).
			WithChildren(
				BlockToTag(ctx, block.Contents(), counter),
				tag.NewClosed("a").
					WithAttrs(map[string]string{
						"href":  "#" + attrs["id"],
						"class": "heading__link",
					}))

	case blocks.RocketLink:
		return tag.NewClosed("li").
			WithAttrs(map[string]string{
				"class": "launchpad__entry",
			}).
			WithChildren(
				tag.NewClosed("a").
					WithContentsStrings(html.EscapeString(block.Display())).
					WithAttrs(map[string]string{
						"class": "rocketlink " + block.Classes(),
						"href":  block.Href(),
					}))

	case blocks.LaunchPad:
		block.ColorRockets() // TODO: fumu fumu
		var rockets []tag.Tag
		for _, rocket := range block.Rockets {
			rockets = append(rockets, BlockToTag(ctx, rocket, counter))
		}
		attrs["class"] = "launchpad"
		return tag.NewClosed("ul").WithAttrs(attrs).WithChildren(rockets...)

	case blocks.CodeBlock:
		var contentsLines []lines.Line
		for _, line := range strings.Split(block.Contents(), "\n") {
			contentsLines = append(contentsLines, lines.UnindentableFrom(line))
		}
		attrs["class"] = "codeblock"
		return tag.NewClosed("pre").
			WithAttrs(attrs).
			WithChildren(
				tag.NewClosed("code").
					WithContentsLines(contentsLines...).
					WithAttrs(map[string]string{"class": "language-" + block.Language()}))

	case blocks.Img:
		block = block.WithExistingTargetsMarked()
		var children []tag.Tag
		for _, entry := range block.Entries {
			children = append(children, BlockToTag(ctx, entry, counter))
		}
		attrs["class"] = "img-gallery " + util.TernaryConditionString(block.HasOneImage(), "img-gallery_one-image", "img-gallery_many-images")
		return tag.NewClosed("section").WithAttrs(attrs).WithChildren(children...)

	case blocks.ImgEntry:
		var children []tag.Tag

		if block.Target.IsBlueLink() {
			imgAttrs := map[string]string{
				"src": block.Target.ImgSrc(),
			}
			if block.GetWidth() != "" {
				imgAttrs["width"] = block.GetWidth()
			}
			if block.GetHeight() != "" {
				imgAttrs["height"] = block.GetHeight()
			}
			children = append(
				children,
				tag.NewClosed("a").
					WithAttrs(map[string]string{"href": block.Target.Href()}).
					WithChildren(
						tag.NewUnclosed("img").WithAttrs(imgAttrs)),
			)
		} else {
			children = append(
				children,
				tag.NewClosed("a").
					WithContentsStrings(fmt.Sprintf("Hypha <i>%s</i> does not exist", block.Target.TargetHypha())).
					WithAttrs(map[string]string{
						"class": block.Target.Classes(),
						"href":  block.Target.Href(),
					}),
			)
		}
		if block.Description != "" {
			figcaption := tag.NewClosed("figcaption").
				WithChildren(BlockToTag(
					ctx,
					parser.MakeFormatted(block.Description, ctx.HyphaName()),
					counter,
				))
			children = append(children, figcaption)
		}
		return tag.NewClosed("figure").
			WithAttrs(map[string]string{"class": "img-gallery__entry"}).
			WithChildren(children...)

	case blocks.HorizontalLine:
		return tag.NewUnclosed("hr").WithAttrs(attrs)

	default:
		return tag.NewUnclosed("error").WithAttrs(attrs)
	}
}
