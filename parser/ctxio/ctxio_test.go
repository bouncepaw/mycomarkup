package ctxio

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/options"
	"testing"
)

const input = "a\r\nsamantha\n\r\n"

var opts = options.Options{HyphaName: "🐙"}.FillTheRest()

func TestNextByte(t *testing.T) {
	var ctx, _ = mycocontext.ContextFromStringInput(input, opts)

	b1, _ := NextByte(ctx)
	if b1 != 'a' {
		t.Errorf("Expected a, got %c\n", b1)
	}

	b2, _ := NextByte(ctx)
	if b2 != '\n' {
		t.Errorf("Expected \\n, got %c\n", b2)
	}
}

func TestNextLine(t *testing.T) {
	var ctx, _ = mycocontext.ContextFromStringInput(input, opts)

	line1, _ := NextLine(ctx)
	if line1 != "a" {
		t.Errorf("Expected a, got %q\n", line1)
	}

	line2, _ := NextLine(ctx)
	if line2 != "samantha" {
		t.Errorf("Expected samantha, got %q\n", line2)
	}

	line3, _ := NextLine(ctx)
	if line3 != "" {
		t.Errorf("Expected empty line, got %q\n", line3)
	}
}
