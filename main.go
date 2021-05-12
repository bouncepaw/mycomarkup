package main

import (
	"fmt"
	markup "github.com/bouncepaw/mycomarkup/legacy"
)

func text() string {
	return `# I am an internet
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
}

func main() {
	markup.HyphaExists = func(s string) bool {
		return true
	}
	markup.HyphaAccess = func(s string) (rawText, binaryHtml string, err error) {
		return "aaaaaaaa,", "aaaaaaaaaaaaa", nil
	}
	markup.HyphaIterate = func(f func(string)) {
		fmt.Println("hello")
	}

	doc := markup.Doc("Example", text())
	fmt.Println(doc.AsHTML())
}
