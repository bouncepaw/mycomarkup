// Package mycomarkup provides an API for processing Mycomarkup-formatted documents.
package mycomarkup

import (
	"sync"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"github.com/bouncepaw/mycomarkup/parser"
)

// BlockTree returns a slice of blocks parsed from the Mycomarkup document contained in ctx.
func BlockTree(ctx mycocontext.Context) []blocks.Block {
	var (
		tokens = make(chan blocks.Block)
		ast    = []blocks.Block{}
		wg     sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		parser.Parse(ctx, tokens)
		wg.Done()
	}()

	for token := range tokens {
		ast = append(ast, token)
	}

	wg.Wait()
	return ast
}

// BlocksToHTML turns the blocks into their HTML representation.
func BlocksToHTML(_ mycocontext.Context, blocks []blocks.Block) string {
	return generateHTML(blocks, 0)
}
