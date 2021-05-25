package mycomarkup

import (
	"strings"
	"sync"
	"testing"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"github.com/bouncepaw/mycomarkup/parser"
)

const input = `What will you give me for this simple dimple?
img {
  ./squish
}

I'll give you //this// squish.

I agree ✅

I agree ✅`

func TestOpenGraphHTML(t *testing.T) {
	var (
		ast      = []blocks.Block{}
		ctx, _   = mycocontext.ContextFromStringInput("test", input)
		wg       sync.WaitGroup
		blocksCh = make(chan blocks.Block)
	)
	wg.Add(1)
	go func() {
		parser.Parse(ctx, blocksCh)
		wg.Done()
	}()
	for block := range blocksCh {
		ast = append(ast, block)
	}
	wg.Wait()

	html := OpenGraphHTML(ctx, ast)
	if !strings.Contains(html, `<meta property="og:image" content="/binary/test/squish"/>`) || !strings.Contains(html, `<meta property="og:description" content="What will you give me for this simple dimple?"/>`) {
		t.Errorf("Wrong output:\n%s\n", html)
	}
}
