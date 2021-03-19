package main

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/links"
)

func main() {
	doc := []blocks.Block{
		blocks.MakeHeading(1, "", blocks.MakeFormatted("This is a test tree of mycomarkup blocks")),
		blocks.MakeHorizontalLine(4),
		blocks.MakeRocketLink(links.From("apple", "яблоко", "home")),
	}

	for _, e := range doc {
		fmt.Println(e.String())
	}
}
