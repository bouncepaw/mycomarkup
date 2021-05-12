package main

import (
	"flag"
	"fmt"
	"github.com/bouncepaw/mycomarkup/doc"
	"github.com/bouncepaw/mycomarkup/opts"
	"io/ioutil"
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

func init() {
	doc.HyphaExists = func(s string) bool {
		return true
	}
	doc.HyphaAccess = func(s string) (rawText, binaryHtml string, err error) {
		return "aaaaaaaa,", "aaaaaaaaaaaaa", nil
	}
	doc.HyphaIterate = func(f func(string)) {
		fmt.Println("hello")
	}
}

func main() {
	hyphaName, filename := parseFlags()
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Errorf("%s\n", err)
	}

	dok := doc.Doc(hyphaName, string(contents))
	fmt.Println(dok.AsHTML())
}

func parseFlags() (hyphaName, filename string) {
	opts.UseBatch = true

	flag.StringVar(&hyphaName, "hypha-name", "", "Set hypha name. Relative links depend on it.")
	flag.StringVar(&filename, "filename", "/dev/stdin", "File with mycomarkup.")
	flag.Parse()

	return
}
