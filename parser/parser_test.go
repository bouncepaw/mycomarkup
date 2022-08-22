package parser

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/parser/ctxio"
	"testing"
)

var opts = options.Options{
	HyphaName:             "test",
	WebSiteURL:            "",
	TransclusionSupported: false,
}.FillTheRest()

func TestIsPrefixedBy(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("input input", opts)
	res := isPrefixedBy(ctx, "input")
	if !res {
		t.Errorf("wrong")
	}
}

func TestLooksLikeList1(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("* i got drip", opts)
	res := looksLikeList(ctx)
	if !res {
		t.Errorf("wrong")
	}
}

func TestLooksLikeList2(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("* { what you gonna do\n when they come for you }", opts)
	res := looksLikeList(ctx)
	if !res {
		t.Errorf("wrong")
	}
}

func TestLooksLikeList3(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("*{ what you gonna do\n when they come for you }", opts)
	res := looksLikeList(ctx)
	if res {
		t.Errorf("wrong")
	}
}

func TestList1(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("* li", opts)
	ctxio.EatUntilSpace(ctx)
	text, _ := readNextListItemsContents(ctx)
	if text.String() != "li" {
		t.Errorf("wrong")
	}
}

func TestList2(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("* {dreamy\n   sky} ", opts)
	ctxio.EatUntilSpace(ctx)
	text, _ := readNextListItemsContents(ctx)
	if text.String() != "dreamy\nsky " {
		t.Errorf("wrong %q", text.String())
	}
}

func TestNextLineIsSomething1(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("=> space", opts)
	res := nextLineIsSomething(ctx)
	if !res {
		t.Errorf("wrong")
	}
}

func TestNextLineIsSomething2(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("* line", opts)
	res := nextLineIsSomething(ctx)
	if !res {
		t.Errorf("wrong")
	}
}

func TestNextLineIsSomething3(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput(" \r\n", opts)
	res := nextLineIsSomething(ctx)
	if !res {
		t.Errorf("wrong")
	}
}

func TestEmptyLine(t *testing.T) {
	ctx1, _ := mycocontext.ContextFromStringInput("", opts)
	if !matchesEmptyLine(ctx1) {
		t.Errorf("Wrong 1")
	}

	ctx2, _ := mycocontext.ContextFromStringInput("\r\n", opts)
	if !matchesEmptyLine(ctx2) {
		t.Errorf("Wrong 2")
	}

	ctx3, _ := mycocontext.ContextFromStringInput("aboba\r\n", opts)
	if matchesEmptyLine(ctx3) {
		t.Errorf("Wrong 3")
	}

	ctx4, _ := mycocontext.ContextFromStringInput(" \r\n", opts)
	if !matchesEmptyLine(ctx4) {
		t.Errorf("Wrong 4")
	}

	ctx5, _ := mycocontext.ContextFromStringInput("aboba\r\n", opts)
	if matchesEmptyLine(ctx5) {
		t.Errorf("Wrong 5")
	}
}
