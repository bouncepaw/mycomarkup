// Package mycomarkup provides an API for processing Mycomarkup-formatted documents.
package mycomarkup

import (
	"sync"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"github.com/bouncepaw/mycomarkup/parser"
)

// BlockTree returns a slice of blocks parsed from the Mycomarkup document contained in ctx.
//
// Pass visitors. Visitors are functions (usually closures) that are called on every found block.
//
// Visitors included with mycomarkup can be gotten from OpenGraphVisitors. More visitors coming soon.
func BlockTree(ctx mycocontext.Context, visitors ...func(block blocks.Block)) []blocks.Block {
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

		for _, visitor := range visitors {
			visitor := visitor
			visitor(token)
		}
	}

	wg.Wait()
	return ast
}

// BlocksToHTML turns the blocks into their HTML representation.
func BlocksToHTML(_ mycocontext.Context, ast []blocks.Block) string {
	counter := &blocks.IDCounter{
		ShouldUseResults: true,
	}
	return generateHTML(ast, 0, counter)
}
