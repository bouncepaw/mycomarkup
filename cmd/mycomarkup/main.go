package main

import (
	"flag"
	"fmt"
	"git.sr.ht/~bouncepaw/mycomarkup/v5"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
	"io/ioutil"
)

func main() {
	hyphaName, fileName := parseFlags()
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		_ = fmt.Errorf("%s\n", err)
	}

	// TODO: provide a similar function but for []byte and use it here.
	ctx, _ := mycocontext.ContextFromBytes(contents, options.Options{
		HyphaName:             hyphaName,
		WebSiteURL:            "",
		TransclusionSupported: false,
	}.FillTheRest())
	ast := mycomarkup.BlockTree(ctx)
	fmt.Println(mycomarkup.BlocksToHTML(ctx, ast))
}

func parseFlags() (hyphaName, fileName string) {
	flag.StringVar(&hyphaName, "hypha-name", "", "Set hypha name. Relative links depend on it.")
	flag.StringVar(&fileName, "file-name", "/dev/stdin", "File with a Mycomarkup document.")
	flag.Parse()

	return
}
