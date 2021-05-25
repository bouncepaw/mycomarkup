// Package parser turns the source text into a sequence of blocks.
package parser

import (
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/mycocontext"
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
