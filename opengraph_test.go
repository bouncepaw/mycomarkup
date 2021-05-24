package mycomarkup

import (
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"strings"
	"sync"
	"testing"

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
		blocks   = []interface{}{}
		ctx, _   = mycocontext.ContextFromStringInput("test", input)
		wg       sync.WaitGroup
		blocksCh = make(chan interface{})
	)
	wg.Add(1)
	go func() {
		parser.Parse(ctx, blocksCh)
		wg.Done()
	}()
	for block := range blocksCh {
		blocks = append(blocks, block)
	}
	wg.Wait()

	html := OpenGraphHTML(ctx, blocks)
	if !strings.Contains(html, `<meta property="og:image" content="/binary/test/squish"/>`) || !strings.Contains(html, `<meta property="og:description" content="What will you give me for this simple dimple?"/>`) {
		t.Errorf("Wrong output:\n%s\n", html)
	}
}
