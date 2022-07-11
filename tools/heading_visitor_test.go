package tools

import (
	"github.com/bouncepaw/mycomarkup/v5"
	"github.com/bouncepaw/mycomarkup/v5/blocks"
	"github.com/bouncepaw/mycomarkup/v5/mycocontext"
	"github.com/bouncepaw/mycomarkup/v5/options"
	"github.com/bouncepaw/mycomarkup/v5/parser"
	"reflect"
	"testing"
)

const headingInput = `
= ok let's go

some text

nothing suspicious going on, just textin'

img {
   notging
moer
wepofjj}

> == heaidng !!!
> quotee
> > = hey hey hey hey
>>> go go go!
table {
| hey you | {
=== out there in the cold
getting lonely getting old
can you hear me ?
}}
`

func TestHeadingVisitor(t *testing.T) {
	var ctx, _ = mycocontext.ContextFromStringInput(headingInput, options.Options{}.FillTheRest())
	var headingExpected = []blocks.Heading{
		blocks.NewHeading(1, parser.MakeFormatted(ctx, "ok let's go"), "= ok let's go"),
		blocks.NewHeading(2, parser.MakeFormatted(ctx, "heaidng !!!"), "== heaidng !!!"),
		blocks.NewHeading(1, parser.MakeFormatted(ctx, "hey hey hey hey"), "= hey hey hey hey"),
		blocks.NewHeading(3, parser.MakeFormatted(ctx, "out there in the cold"), "=== out there in the cold"),
	}
	visitor, result := HeadingVisitor(ctx)
	_ = mycomarkup.BlockTree(ctx, visitor) // Call for side-effect
	foundHeadings := result()
	if !(len(headingExpected) == len(foundHeadings)) {
		t.Errorf("Heading count mismatch: expected %d, got %d\n", len(headingExpected), len(foundHeadings))
		return
	}
	for i, heading := range foundHeadings {
		if !(reflect.DeepEqual(heading, headingExpected[i])) {
			t.Errorf("Heading mismatch at %d:\nwanted %#v\ngot    %#v\n", i, headingExpected[i], heading)
		}
	}
}
