package tools

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5/blocks"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
)

// HeadingVisitor creates a visitor that visits all headings.
func HeadingVisitor(ctx mycocontext.Context) (
	visitor func(block blocks.Block),
	result func() []blocks.Heading,
) {
	var (
		collected []blocks.Heading
	)
	visitor = func(block blocks.Block) {
		switch block := block.(type) {
		case blocks.Heading:
			collected = append(collected, block)
		case blocks.Quote:
			for _, subblock := range block.Contents() {
				visitor(subblock)
			}
		case blocks.Table:
			for _, row := range block.Rows() {
				for _, cell := range row.Cells() {
					visitor(cell)
				}
			}
		case blocks.TableCell:
			for _, subblock := range block.Contents() {
				visitor(subblock)
			}
		}
	}
	result = func() []blocks.Heading {
		return collected
	}
	return
}
