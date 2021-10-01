// Package parser turns the source text into a sequence of blocks.
package parser

import (
	"bytes"
	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
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

// parseHeading parses the heading on the given line and returns it. Find its level by yourself though.
func parseHeading(line, hyphaName string, level uint) blocks.Heading {
	return blocks.NewHeading(level, MakeFormatted(line[level+1:], hyphaName), line)
	// TODO: figure out the level here.
}
