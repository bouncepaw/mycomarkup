package util

import (
	"bytes"
	"context"
	"strings"
)

// Key is used for setting and getting from context.Context in this project.
type Key int

// These are keys for the context that floats around.
const (
	// KeyHyphaName is for storing current hypha name as a string here.
	KeyHyphaName Key = iota
	// KeyInputBuffer is for storing *bytes.Buffer with unread bytes of the source document.
	KeyInputBuffer
	// KeyRecursionLevel stores current level of transclusion recursion.
	KeyRecursionLevel
)

// HyphaNameFrom retrieves current hypha name from the given context.
func HyphaNameFrom(ctx context.Context) string {
	return ctx.Value(KeyHyphaName).(string)
}

// InputFrom retrieves the current bytes buffer from the given context.
func InputFrom(ctx context.Context) *bytes.Buffer {
	return ctx.Value(KeyInputBuffer).(*bytes.Buffer)
}

// EatUntilSpace reads characters until it encounters a non-space character.
func EatUntilSpace(ctx context.Context) {
	// We do not care what is read, therefore we drop the read line.
	// We know that there //is// a space beforehand, therefore we drop the error.
	_, _ = InputFrom(ctx).ReadString(' ')
}

// NextByte returns the next byte in the inputFrom. The CR byte (\r) is never returned, if there is a CR in the inputFrom, the byte after it is returned. If there is no next byte, the NL byte (\n) is returned and done is true.
func NextByte(ctx context.Context) (b byte, done bool) {
	b, err := InputFrom(ctx).ReadByte()
	if err != nil {
		return '\n', true
	}
	if b == '\r' {
		return NextByte(ctx)
	}
	return b, false
}

// NextLine returns the text in the inputFrom up to the next newline. The characters are gotten using nextByte.
func NextLine(ctx context.Context) (line string, done bool) {
	var (
		lineBuffer strings.Builder
		b          byte
	)
	b, done = NextByte(ctx)
	for b != '\n' {
		lineBuffer.WriteByte(b)
		b, done = NextByte(ctx)
	}
	return lineBuffer.String(), done
}
