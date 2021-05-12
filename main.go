package main

import (
	"fmt"
	doc2 "github.com/bouncepaw/mycomarkup/doc"
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
	doc2.HyphaExists = func(s string) bool {
		return true
	}
	doc2.HyphaAccess = func(s string) (rawText, binaryHtml string, err error) {
		return "aaaaaaaa,", "aaaaaaaaaaaaa", nil
	}
	doc2.HyphaIterate = func(f func(string)) {
		fmt.Println("hello")
	}

	doc := doc2.Doc("Example", text())
	fmt.Println(doc.AsHTML())
}
