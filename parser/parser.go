// Package parser turns the source text into a sequence of blocks.
package parser

import (
	"bytes"
	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"sync"
)

// Parse parses the Mycomarkup document in the given context. All parsed blocks are written to out.
func Parse(ctx mycocontext.Context, out chan blocks.Block) {
	// Using a channel seems like a good idea. The downside is that using this function is harder. But does it matter in this case? Not really. Channel supremacy all the way down.
	// Added later: maybe it is time to reconsider!
	var (
		token blocks.Block
		done  bool
	)
	defer close(out)
	for !done {
		select {
		case <-ctx.Done():
			return
		default:
			token, done = nextToken(ctx)
			if token != nil {
				out <- token
			}
		}
	}
}

// parseSubdocumentForEachBlock replaces the buffer in the given context and parses the document contained in the buffer. The function is called on every block.
func parseSubdocumentForEachBlock(ctx mycocontext.Context, buf *bytes.Buffer, f func(block blocks.Block)) {
	var (
		wg       sync.WaitGroup
		blocksCh = make(chan blocks.Block)
	)
	wg.Add(1)
	go func() {
		Parse(mycocontext.WithBuffer(ctx, buf), blocksCh)
		wg.Done()
	}()
	for block := range blocksCh {
		f(block)
	}
	wg.Wait()
}

// parseHeading parses the heading on the given line and returns it. Find its level by yourself though.
func parseHeading(line, hyphaName string, level uint) blocks.Heading {
	return blocks.NewHeading(level, MakeFormatted(line[level+1:], hyphaName), line)
	// TODO: figure out the level here.
}
