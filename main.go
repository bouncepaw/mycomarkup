package main

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/blocks"
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
	_ = `# I am an internet
Why the life is so rough with me?, I wonder.
=> link
=> link_link display
=> link\ link display
=> [[link]]
=> [[link|display]]
=> [[link|]]
=> [[]]
=> [[|]]
=>
`
	/*	tokens := lexer.Lex(bytes.NewBufferString(doc))
		fmt.Printf("↓ Source text:\n%s\n↓ Tokens got:\n", doc)
		for _, token := range tokens {
			fmt.Println(token.String())
		}*/
}
