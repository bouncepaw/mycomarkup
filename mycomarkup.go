// Package mycomarkup provides an API for processing Mycomarkup-formatted documents.
package mycomarkup

import (
	"errors"
	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/genhtml"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"github.com/bouncepaw/mycomarkup/v3/parser"
	"github.com/bouncepaw/mycomarkup/v3/temporary_workaround"
)

// BlockTree returns a slice of blocks parsed from the Mycomarkup document contained in ctx.
//
// Pass visitors. Visitors are functions (usually closures) that are called on every top-level found block.
//
// Some pre-implemented visitors are in the tools package.
func BlockTree(ctx mycocontext.Context, visitors ...func(block blocks.Block)) []blocks.Block {
	var (
		tokens = make([]blocks.Block, 0)
		token  blocks.Block
		done   bool
	)

	for !done {
		select {
		case <-ctx.Done():
			return tokens
		default:
			token, done = parser.NextToken(ctx)
			if token != nil {
				tokens = append(tokens, token)

				for _, visitor := range visitors {
					visitor := visitor
					visitor(token)
				}
			}
		}
	}

	return tokens
}

// BlocksToHTML turns the blocks into their HTML representation.
func BlocksToHTML(ctx mycocontext.Context, ast []blocks.Block) string {
	counter := blocks.NewIDCounter()
	var res string
	for _, block := range ast {
		res += genhtml.BlockToTag(ctx, block, counter).String()
	}
	return res
}

// transclusionVisitor returns a visitor to pass to BlockTree and a function to get the results.
func transclusionVisitor(xcl blocks.Transclusion) (
	visitor func(block blocks.Block),
	result func() ([]blocks.Block, error),
) {
	var (
		collected             []blocks.Block
		metDescriptionAlready = false
	)
	visitor = func(block blocks.Block) {
		switch xcl.Selector {
		case blocks.SelectorAttachment:
			// We don't need any of that when we only transclude attachments.
		case blocks.SelectorText, blocks.SelectorFull:
			collected = append(collected, block)
		case blocks.SelectorOverview, blocks.SelectorDescription:
			switch block.(type) {
			case blocks.Paragraph:
				if metDescriptionAlready {
					break
				}
				metDescriptionAlready = true
				collected = append(collected, block)
			}
		}
	}
	result = func() ([]blocks.Block, error) {
		if len(collected) == 0 {
			switch xcl.Selector {
			case blocks.SelectorDescription:
				// Asked for a description, got no description.
				return nil, errors.New("no description")
			case blocks.SelectorText:
				// Asked for a text, found emptiness...
				return nil, errors.New("no text")
			}
		}

		return collected, nil
	}
	return
}

func init() {
	temporary_workaround.BlockTree = BlockTree
	temporary_workaround.TransclusionVisitor = transclusionVisitor
}
