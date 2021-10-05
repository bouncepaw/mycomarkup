package tools

import (
	"testing"

	"github.com/bouncepaw/mycomarkup/v3"
	"github.com/bouncepaw/mycomarkup/v3/links"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
)

const inputLinks = `[[ TODO ]]

=> links
=> links/Games Games

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
	)
	ctx, _ := mycocontext.ContextFromStringInput(hyphaName, inputLinks)
	ctx = mycocontext.WithCalledInShell(ctx)
	linkVisitor, getLinks := LinkVisitor(ctx)
	mycomarkup.BlockTree(ctx, linkVisitor)
	foundLinks := getLinks()

	expectedLinks := []links.Link{
		links.From("TODO", "", hyphaName),
		links.From("links", "", hyphaName),
		links.From("links/Games", "Games", hyphaName),
		links.From("ideas", "", hyphaName),
		links.From("links/Anime", "", hyphaName),
		links.From("./kittens", "", hyphaName),
		links.From("../puppies", "", hyphaName),
		links.From("https://example.com/favicon.ico", "", hyphaName),
		links.From("home", "", hyphaName),
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
