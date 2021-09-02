// Package tools provides some helper functions
package tools

import (
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/links"
	"github.com/bouncepaw/mycomarkup/mycocontext"
)

// LinkVisitor creates a visitor which extracts all the links
func LinkVisitor(ctx mycocontext.Context) (
	visitor func(block blocks.Block),
	result func() []links.Link,
) {
	var (
		collected []links.Link
	)
	var extractBlock func(block blocks.Block)
	extractBlock = func(block blocks.Block) {
		switch b := block.(type) {
		case blocks.Paragraph:
			extractBlock(b.Formatted)
		case blocks.Heading:
			extractBlock(b.GetContents())
		case blocks.List:
			for _, item := range b.Items {
				for _, sub := range item.Contents {
					extractBlock(sub)
				}
			}
		case blocks.Img:
			for _, entry := range b.Entries {
				extractBlock(entry)
			}
		case blocks.ImgEntry:
			collected = append(collected, *b.Srclink)
		case blocks.Transclusion:
			link := *links.From(b.Target, "", ctx.HyphaName())
			collected = append(collected, link)
		case blocks.LaunchPad:
			for _, rocket := range b.Rockets {
				extractBlock(rocket)
			}
		case blocks.Formatted:
			for _, line := range b.Lines {
				for _, span := range line {
					switch s := span.(type) {
					case blocks.InlineLink:
						collected = append(collected, *s.Link)
					}
				}
			}
		case blocks.RocketLink:
			if !b.IsEmpty {
				collected = append(collected, b.Link)
			}
		}
	}
	visitor = func(block blocks.Block) {
		extractBlock(block)
	}
	result = func() []links.Link {
		return collected
	}
	return
}
