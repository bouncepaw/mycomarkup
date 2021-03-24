package main

import (
	"bytes"
	"fmt"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/lexer"
	"github.com/bouncepaw/mycomarkup/links"
)

func demoBlocks() {
	doc := []blocks.Block{
		blocks.MakeHeading(1, blocks.MakeFormatted("This is a test tree of mycomarkup blocks")),
		blocks.MakeHorizontalLine(4),
		blocks.MakeRocketLink(links.From("apple", "яблоко", "home")),
	}

	for _, e := range doc {
		fmt.Println(e.String())
	}

}

func main() {
	doc := `# I am an internet
Why the life is so rough with me?, I wonder.
=> link
=> link display
=> link_link display
=> link\ link display
=> https://superlink
`
	tokens := lexer.Lex(bytes.NewBufferString(doc))
	fmt.Printf("↓ Source text:\n%s\n↓ Tokens got:\n", doc)
	for _, token := range tokens {
		fmt.Println(token.String())
	}
}
