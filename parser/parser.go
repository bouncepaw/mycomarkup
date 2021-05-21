// Package parser turns the source text into a sequence of blocks.
package parser

import (
	"bytes"
	"context"
	"github.com/bouncepaw/mycomarkup/util"
)

// Parse parses the Mycomarkup document in the given context. All parsed blocks are written to out.
func Parse(ctx context.Context, out chan interface{}) {
	var (
		token interface{}
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

// ContextFromStringInput returns the context for the given inputFrom.
func ContextFromStringInput(hyphaName, input string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(
		context.WithValue(
			context.WithValue(
				context.WithValue(
					context.Background(),
					util.KeyHyphaName,
					hyphaName),
				util.KeyInputBuffer,
				bytes.NewBufferString(input),
			),
			util.KeyRecursionLevel,
			0,
		),
	)
	return ctx, cancel
}
