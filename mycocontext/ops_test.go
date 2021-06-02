package mycocontext

import "testing"

const input = "a\r\nsamantha\n"

func TestNextByte(t *testing.T) {
	var ctx, _ = ContextFromStringInput("🐙", input)

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
	var ctx, _ = ContextFromStringInput("🐙", input)

	line, _ := NextLine(ctx)
	if line != "a" {
		t.Errorf("Expected a, got %q\n", line)
	}
}
