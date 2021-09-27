// Package tools provides visitors for external usage and will provide other tools one day.
package tools

import (
	"github.com/bouncepaw/mycomarkup/v2/blocks"
	"github.com/bouncepaw/mycomarkup/v2/links"
	"github.com/bouncepaw/mycomarkup/v2/mycocontext"
)

// LinkVisitor creates a visitor which extracts all the links from the document in context.
//
// We consider inline link, rocket link, image gallery and transclusion targets to be links.
func LinkVisitor(ctx mycocontext.Context) (
	visitor func(block blocks.Block),
	result func() []links.Link,
) {
	var (
		collected    []links.Link
		extractLinks func(block blocks.Block)
	)
	extractLinks = func(block blocks.Block) {
		switch b := block.(type) {
		case blocks.Paragraph:
			extractLinks(b.Formatted)
		case blocks.Heading:
			extractLinks(b.Contents())
		case blocks.List:
			for _, item := range b.Items {
				for _, sub := range item.Contents {
					extractLinks(sub)
				}
			}
		case blocks.Img:
			for _, entry := range b.Entries {
				extractLinks(entry)
			}
		case blocks.ImgEntry:
			collected = append(collected, b.Target)
		case blocks.Transclusion:
			link := links.From(b.Target, "", ctx.HyphaName())
			collected = append(collected, link)
		case blocks.LaunchPad:
			for _, rocket := range b.Rockets {
				extractLinks(rocket)
			}
		case blocks.Formatted:
			for _, line := range b.Lines {
				for _, span := range line {
					switch s := span.(type) {
					case blocks.InlineLink:
						collected = append(collected, s.Link)
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
		extractLinks(block)
	}
	result = func() []links.Link {
		return collected
	}
	return
}
