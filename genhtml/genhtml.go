// Package genhtml provides utilities for converting Mycomarkup blocks into HTML documents. As of now, some parts of HTML generation are in other parts of the library, WIP.
package genhtml

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v3/globals"
	"github.com/bouncepaw/mycomarkup/v3/options"
	"github.com/bouncepaw/mycomarkup/v3/temporary_workaround"
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
					contents += strings.TrimSuffix(tag.NewClosed("a").
						WithAttrs(map[string]string{
							"href":  s.Href(),
							"class": s.Classes(),
						}).
						WithContentsStrings(html.EscapeString(s.Display())).
						String(), "\n")

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
		attrs["class"] = fmt.Sprintf(
			`img-gallery img-gallery_%s img-gallery_layout-%s`,
			util.TernaryConditionString(block.HasOneImage(), "one-image", "many-images"),
			block.Layout().String(),
		)
		return tag.NewClosed("section").WithAttrs(attrs).WithChildren(children...)

	case blocks.ImgEntry:
		var children []tag.Tag

		if block.Target.IsBlueLink() {
			imgAttrs := map[string]string{
				"src": block.Target.ImgSrc(),
			}
			if block.Width() != "" {
				imgAttrs["width"] = block.Width()
			}
			if block.Height() != "" {
				imgAttrs["height"] = block.Height()
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
		if block.Description() != "" {
			figcaption := tag.NewClosed("figcaption").
				WithChildren(BlockToTag(
					ctx,
					parser.MakeFormatted(ctx, block.Description()),
					counter,
				))
			children = append(children, figcaption)
		}
		return tag.NewClosed("figure").
			WithAttrs(map[string]string{"class": "img-gallery__entry"}).
			WithChildren(children...)

	case blocks.Quote:
		var children []tag.Tag
		for _, child := range block.Contents() {
			children = append(children, BlockToTag(ctx, child, counter.UnusableCopy()))
		}
		return tag.NewClosed("blockquote").
			WithAttrs(map[string]string{"id": block.ID(counter)}).
			WithChildren(children...)

	case blocks.List:
		var items []tag.Tag
		for _, item := range block.Items {
			var children []tag.Tag
			for _, child := range item.Contents {
				children = append(children, BlockToTag(ctx, child, counter.UnusableCopy()))
			}
			items = append(items, listMarkerToItemTag(item.Marker).WithChildren(children...))
		}
		return listMarkerToParentTag(block.Marker).
			WithAttrs(map[string]string{
				"id": block.ID(counter),
			}).WithChildren(items...)

	case blocks.HorizontalLine:
		return tag.NewUnclosed("hr").WithAttrs(attrs)

	case blocks.TableCell:
		cell := tag.NewClosed(util.TernaryConditionString(block.IsHeaderCell(), "th", "td"))

		if block.Colspan() > 1 {
			cell = cell.WithAttrs(map[string]string{
				"colspan": fmt.Sprintf("%d", block.Colspan()),
			})
		}

		var children []tag.Tag
		for _, child := range block.Contents() {
			children = append(children, BlockToTag(ctx, child, counter))
		}
		return cell.WithChildren(children...)

	case blocks.TableRow:
		var cells []tag.Tag
		for _, cell := range block.Cells() {
			cells = append(cells, BlockToTag(ctx, cell, counter))
		}
		return tag.NewClosed("tr").WithChildren(cells...)

	case blocks.Table:
		var children []tag.Tag
		if block.Caption() != "" {
			children = append(children, tag.NewClosed("caption").WithContentsStrings(html.EscapeString(block.Caption())))
		}

		var rows []tag.Tag
		for _, row := range block.Rows() {
			rows = append(rows, BlockToTag(ctx, row, counter))
		}
		if len(rows) > 0 {
			children = append(children, tag.NewClosed("tbody").WithChildren(rows...))
		}

		return tag.NewClosed("table").WithAttrs(map[string]string{
			"id": block.ID(counter.UnusableCopy()),
		}).WithChildren(children...)

	case blocks.Transclusion:
		xcl := block
		if ctx.RecursionLevel() > 3 {
			return tag.NewClosed("section").
				WithAttrs(map[string]string{
					"class": "transclusion transclusion_failed transclusion_not-exists",
				}).
				WithChildren(tag.NewClosed("p").
					WithContentsStrings("Transclusion depth limit"))
		}

		if xcl.HasError() {
			return MapTransclusionErrorToTag(xcl)
		}

		// V4 This part is awful, bloody hell. Move to the parser module
		// Now, to real transclusion:
		rawText, binaryHtml, err := globals.HyphaAccess(xcl.Target)
		if err != nil {
			xcl.TransclusionError.Reason = blocks.TransclusionErrorNotExists
			return MapTransclusionErrorToTag(xcl)
		}
		xclVisitor, result := temporary_workaround.TransclusionVisitor(xcl)
		xclctx, _ := mycocontext.ContextFromStringInput(rawText, options.Options{
			HyphaName:             xcl.Target,
			WebSiteURL:            ctx.WebSiteURL(),
			TransclusionSupported: true,
		}.FillTheRest()) // FIXME: it will bite us one day UPDATE: is it the day? I don't feel the bite.
		_ = temporary_workaround.BlockTree(xclctx, xclVisitor) // Call for side-effects

		var children []tag.Tag
		collected, err := result()
		if err != nil {
			switch err.Error() {
			case "no description":
				return tag.NewClosed("section").
					WithAttrs(map[string]string{
						"class": "transclusion transclusion_failed transclusion_no-description",
					}).
					WithChildren(tag.NewClosed("p").
						WithContentsStrings(
							fmt.Sprintf(
								`Hypha <a href="/hypha/%s" class="wikilink">%s</a> has no description`,
								xcl.Target,
								util.BeautifulName(xcl.Target),
							)))
			case "no text":
				return tag.NewClosed("section").
					WithAttrs(map[string]string{
						"class": "transclusion transclusion_failed transclusion_no-text",
					}).
					WithChildren(tag.NewClosed("p").
						WithContentsStrings(
							fmt.Sprintf(
								`Hypha <a href="/hypha/%s" class="wikilink">%s</a> has no text`,
								xcl.Target,
								util.BeautifulName(xcl.Target),
							)))
			}
		}
		for _, child := range collected {
			children = append(children, BlockToTag(ctx.WithIncrementedRecursionLevel(), child, counter.UnusableCopy()))
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
					WithContentsStrings(util.TernaryConditionString(
						xcl.Selector == blocks.SelectorAttachment || xcl.Selector == blocks.SelectorFull || xcl.Selector == blocks.SelectorOverview,
						binaryHtml,
						"",
					)).
					WithAttrs(map[string]string{
						"class": "transclusion__content",
					}).
					WithChildren(children...),
			)

	default:
		return tag.NewUnclosed("error").WithAttrs(attrs)
	}
}

func listMarkerToParentTag(m blocks.ListMarker) tag.Tag {
	switch m {
	case blocks.MarkerOrdered:
		return tag.NewClosed("ol")
	default:
		return tag.NewClosed("ul")
	}
}

func listMarkerToItemTag(m blocks.ListMarker) tag.Tag {
	t := tag.NewClosed("li")
	switch m {
	case blocks.MarkerUnordered:
		return t.WithAttrs(map[string]string{
			"class": "item_unordered",
		})
	case blocks.MarkerOrdered:
		return t.WithAttrs(map[string]string{
			"class": "item_ordered",
		})
	case blocks.MarkerTodo:
		return t.WithAttrs(map[string]string{
			"class": "item_todo",
		}).WithContentsLines(
			tag.NewUnclosed("input").
				WithAttrs(map[string]string{
					"type":     "checkbox",
					"disabled": "disabled",
				}).
				Lines()...,
		)
	case blocks.MarkerTodoDone:
		return t.WithAttrs(map[string]string{
			"class": "item_todo item_todo-done",
		}).WithContentsLines(
			tag.NewUnclosed("input").
				WithAttrs(map[string]string{
					"type":     "checkbox",
					"disabled": "disabled",
					"checked":  "checked",
				}).
				Lines()...,
		)
	}
	panic("unreachable")
}
