// Package parser turns the source text into a sequence of blocks.
package parser

import (
	"bytes"
	"lesarbr.es/mycomarkup/v5/blocks"
	"lesarbr.es/mycomarkup/v5/mycocontext"
)

// parseSubdocumentForEachBlock replaces the buffer in the given context and parses the document contained in the buffer. The function is called on every block.
func parseSubdocumentForEachBlock(ctx mycocontext.Context, buf *bytes.Buffer, f func(block blocks.Block)) {
	ctx = mycocontext.WithBuffer(ctx, buf)
	var (
		done  bool
		token blocks.Block
	)
	for !done {
		token, done = NextToken(ctx)
		if token != nil {
			f(token)
		}
	}
}

// parseHeading parses the =heading on the given line and returns it. Find its level by yourself though.
func parseHeading(ctx mycocontext.Context, line string, level uint) blocks.Heading {
	// level is the number of =, then there is a space
	return blocks.NewHeading(level, MakeFormatted(ctx, line[level+1:]), line)
	// TODO: figure out the level here. Maybe?
}
