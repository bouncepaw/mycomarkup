// Package genhtml provides utilities for converting Mycomarkup blocks into HTML documents. As of now, some parts of HTML generation are in other parts of the library, WIP.
package genhtml

import (
	"fmt"
	"html"

	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/genhtml/tag"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
)

// This package shall not depend on anything other than blocks, links, globals, mycocontext, util, tag.

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
					contents += tag.NewClosed("a", map[string]string{
						"href":  s.Href(),
						"class": s.Classes(),
					}, s.Display()).String()

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
		return tag.NewWrapper(contents)
	case blocks.Paragraph:
		return tag.NewClosed("p", attrs, "", BlockToTag(ctx, block.Formatted, counter))
	case blocks.Heading:
		return tag.NewClosed(
			fmt.Sprintf("h%d", block.Level()),
			attrs,
			"",
			BlockToTag(ctx, block.Contents(), counter),
			tag.NewClosed("a", map[string]string{"href": "#" + attrs["id"], "class": "heading__link"}, ""),
		)
	case blocks.RocketLink:
		return tag.NewClosed(
			"li",
			map[string]string{
				"class": "launchpad__entry",
			}, "",
			tag.NewClosed(
				"a",
				map[string]string{
					"class": "rocketlink " + block.Classes(),
					"href":  block.Href(),
				},
				html.EscapeString(block.Display()),
			),
		)
	case blocks.LaunchPad:
		block.ColorRockets()
		var rockets []tag.Tag
		for _, rocket := range block.Rockets {
			rockets = append(rockets, BlockToTag(ctx, rocket, counter))
		}
		attrs["class"] = "launchpad"
		return tag.NewClosed("ul", attrs, "", rockets...)
	case blocks.CodeBlock:
		attrs["class"] = "codeblock"
		return tag.NewClosed("pre", attrs, "",
			tag.NewClosed(
				"code",
				map[string]string{"class": "language-" + block.Language()},
				block.Contents(),
			))
	case blocks.HorizontalLine:
		return tag.NewUnclosed("hr", attrs)
	default:
		return tag.NewUnclosed("error", attrs)
	}
}
