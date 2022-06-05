package tools

import (
	"github.com/bouncepaw/mycomarkup/v5/options"
	"testing"

	"github.com/bouncepaw/mycomarkup/v5"
	"github.com/bouncepaw/mycomarkup/v5/links"
	"github.com/bouncepaw/mycomarkup/v5/mycocontext"
)

const inputLinks = `[[ TODO ]]

=> links
=> links/Games | Games

* [[ideas]]
* => links/Anime

img {
	./kittens
	../puppies
	https://example.com/favicon.ico
}

<= home | full`

func TestLinkVisitor(t *testing.T) {
	var (
		hyphaName = "test"
		opts      = options.Options{
			HyphaName: hyphaName,
		}.FillTheRest()
	)
	ctx, _ := mycocontext.ContextFromStringInput(inputLinks, opts)
	linkVisitor, getLinks := LinkVisitor(ctx)
	mycomarkup.BlockTree(ctx, linkVisitor)
	foundLinks := getLinks()

	expectedLinks := []links.LegacyLink{
		links.LegacyFrom("TODO", "", hyphaName),
		links.LegacyFrom("links", "", hyphaName),
		links.LegacyFrom("links/Games", "Games", hyphaName),
		links.LegacyFrom("ideas", "", hyphaName),
		links.LegacyFrom("links/Anime", "", hyphaName),
		links.LegacyFrom("./kittens", "", hyphaName),
		links.LegacyFrom("../puppies", "", hyphaName),
		links.LegacyFrom("https://example.com/favicon.ico", "", hyphaName),
		links.LegacyFrom("home", "", hyphaName),
	}
	// a little dirty hack for destinationKnown
	expectedLinks[0] = expectedLinks[0].CopyMarkedAsExisting()
	expectedLinks[3] = expectedLinks[3].CopyMarkedAsExisting()

	if !(len(expectedLinks) == len(foundLinks)) {
		t.Errorf("Links count mismatch: expected %d, got %d\n", len(expectedLinks), len(foundLinks))
		return
	}
	for i, link := range foundLinks {
		if !(link == expectedLinks[i]) {
			t.Errorf("Link mismatch at %d:\nwanted %#v\ngot    %#v\n", i, expectedLinks[i], link)
		}
	}
}
