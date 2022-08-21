package tools

import (
	"lesarbr.es/mycomarkup/v5/options"
	"reflect"
	"testing"

	"lesarbr.es/mycomarkup/v5"
	"lesarbr.es/mycomarkup/v5/links"
	"lesarbr.es/mycomarkup/v5/mycocontext"
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
		opts = options.Options{
			HyphaName:         "test",
			RedLinksSupported: true,
		}.FillTheRest()
	)
	ctx, _ := mycocontext.ContextFromStringInput(inputLinks, opts)
	linkVisitor, getLinks := LinkVisitor(ctx)
	mycomarkup.BlockTree(ctx, linkVisitor)
	foundLinks := getLinks()

	expectedLinks := []links.Link{
		links.LinkFrom(ctx, "TODO", ""),
		links.LinkFrom(ctx, "links", ""),
		links.LinkFrom(ctx, "links/Games", "Games"),
		links.LinkFrom(ctx, "ideas", ""),
		links.LinkFrom(ctx, "links/Anime", ""),
		links.LinkFrom(ctx, "./kittens", ""),
		links.LinkFrom(ctx, "../puppies", ""),
		links.LinkFrom(ctx, "https://example.com/favicon.ico", ""),
		links.LinkFrom(ctx, "home", ""),
	}

	if !(len(expectedLinks) == len(foundLinks)) {
		t.Errorf("Links count mismatch: expected %d, got %d\n", len(expectedLinks), len(foundLinks))
		return
	}
	for i, link := range foundLinks {
		if !(reflect.DeepEqual(link, expectedLinks[i])) {
			t.Errorf("Link mismatch at %d:\nwanted %#v\ngot    %#v\n", i, expectedLinks[i], link)
		}
	}
}
