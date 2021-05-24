// This is not done yet
package doc

import (
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"github.com/bouncepaw/mycomarkup/parser"
	"sync"
)

// TODO: remove this abomination
var cfg = struct {
	URL string
}{
	URL: "PLACEHOLDER",
}

// MycoDoc is a mycomarkup-formatted document.
//
// TODO: remove or rethink this type because I hate it.
type MycoDoc struct {
	// data
	HyphaName string
	Contents  string
	// indicators
	parsedAlready bool
}

// Doc returns a mycomarkup document with the given name and mycomarkup-formatted contents.
//
// The returned document is not lexed not parsed yet. You have to do that separately.
func Doc(hyphaName, contents string) *MycoDoc {
	md := &MycoDoc{
		HyphaName: hyphaName,
		Contents:  contents,
	}
	return md
}

func (md *MycoDoc) Lex() []interface{} {
	var (
		ctx, _ = mycocontext.ContextFromStringInput(md.HyphaName, md.Contents)
		tokens = make(chan interface{})
		ast    = []interface{}{}
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

// AsHTML returns an html representation of the document
func (md *MycoDoc) AsHTML() string {
	return GenerateHTML(md.Lex(), 0)
}

// AsGemtext returns a gemtext representation of the document. Currently really limited, just returns source text
func (md *MycoDoc) AsGemtext() string {
	return md.Contents
}
