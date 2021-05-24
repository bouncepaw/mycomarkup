package parser

import (
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"testing"
)

func TestIsPrefixedBy(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("test", "input input")
	res := isPrefixedBy(ctx, "input")
	if !res {
		t.Errorf("wrong")
	}
}

func TestLooksLikeList1(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("test", "* i got drip")
	res := looksLikeList(ctx)
	if !res {
		t.Errorf("wrong")
	}
}

func TestLooksLikeList2(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("test", "* { what you gonna do\n when they come for you }")
	res := looksLikeList(ctx)
	if !res {
		t.Errorf("wrong")
	}
}

func TestLooksLikeList3(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("test", "*{ what you gonna do\n when they come for you }")
	res := looksLikeList(ctx)
	if res {
		t.Errorf("wrong")
	}
}

func TestList1(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("test", "* li")
	mycocontext.EatUntilSpace(ctx)
	text, _ := readNextListItemsContents(ctx)
	if text.String() != "li" {
		t.Errorf("wrong")
	}
}

func TestList2(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("test", "* {dreamy\n   sky} ")
	mycocontext.EatUntilSpace(ctx)
	text, _ := readNextListItemsContents(ctx)
	if text.String() != "dreamy\nsky " {
		t.Errorf("wrong %q", text.String())
	}
}

func TestNextLineIsSomething1(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("test", "=> space")
	res := nextLineIsSomething(ctx)
	if !res {
		t.Errorf("wrong")
	}
}

func TestNextLineIsSomething2(t *testing.T) {
	ctx, _ := mycocontext.ContextFromStringInput("test", "* line")
	res := nextLineIsSomething(ctx)
	if !res {
		t.Errorf("wrong")
	}
}
