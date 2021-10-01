package main

import (
	"flag"
	"fmt"
	"github.com/bouncepaw/mycomarkup/v3"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"io/ioutil"
)

func main() {
	hyphaName, filename := parseFlags()
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		_ = fmt.Errorf("%s\n", err)
	}

	// TODO: provide a similar function but for []byte and use it here.
	ctx, _ := mycocontext.ContextFromStringInput(hyphaName, string(contents))
	ast := mycomarkup.BlockTree(ctx)
	fmt.Println(mycomarkup.BlocksToHTML(ctx, ast))
}

func parseFlags() (hyphaName, filename string) {
	flag.StringVar(&hyphaName, "hypha-name", "", "Set hypha name. Relative links depend on it.")
	flag.StringVar(&filename, "filename", "/dev/stdin", "File with mycomarkup.")
	flag.Parse()

	return
}
