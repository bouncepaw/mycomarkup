// Package parser turns the source text into a sequence of blocks.
package parser

import (
	"bytes"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"sync"
)

// Parse parses the Mycomarkup document in the given context. All parsed blocks are written to out.
//
// TODO: decide whether using the channel is really a good idea 🤔
func Parse(ctx mycocontext.Context, out chan blocks.Block) {
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
