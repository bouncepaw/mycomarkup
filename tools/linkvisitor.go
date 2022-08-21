// Package tools provides visitors for external usage and will provide other tools one day.
package tools

import (
	"lesarbr.es/mycomarkup/v5/blocks"
	"lesarbr.es/mycomarkup/v5/links"
	"lesarbr.es/mycomarkup/v5/mycocontext"
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
			link := links.LinkFrom(ctx, b.Target, "")
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
		case blocks.Table:
			for _, row := range b.Rows() {
				for _, cell := range row.Cells() {
					extractLinks(cell)
				}
			}
		case blocks.TableCell:
			for _, block := range b.Contents() {
				extractLinks(block)
			}
		case blocks.Quote:
			for _, block := range b.Contents() {
				extractLinks(block)
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
