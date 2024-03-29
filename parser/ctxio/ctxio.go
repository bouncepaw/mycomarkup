// Package ctxio provides input and output operations on the input buffer in mycocontext.Context.
package ctxio

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"strings"
)

// EatUntilSpace reads characters until it encounters a non-space character. The read characters are returned. No errors are reported even if there are any, be bold.
func EatUntilSpace(ctx mycocontext.Context) (line string) {
	// We do not care what is read, therefore we drop the read line.
	// We know that there //is// a space beforehand, therefore we drop the error.
	line, _ = ctx.Input().ReadString(' ')
	return line
}

// TODO: Gotta make those Next* functions return just one thing and no bools. It is awkward.

// NextByte returns the next byte in the inputFrom. The CR byte (\r) is never returned, if there is a CR in the inputFrom, the byte after it is returned. If there is no next byte, the NL byte (\n) is returned and eof is true.
func NextByte(ctx mycocontext.Context) (b byte, eof bool) {
	b, err := ctx.Input().ReadByte()
	if err != nil {
		return '\n', true
	}
	if b == '\r' {
		return NextByte(ctx)
	}
	return b, false
}

// UnreadRune unreads the previous rune. Pray so it doesn't throw any errors, because they are ignored.
func UnreadRune(ctx mycocontext.Context) {
	_ = ctx.Input().UnreadRune()
}

// NextRune is like NextByte, but for runes.
func NextRune(ctx mycocontext.Context) (r rune, eof bool) {
	r, _, err := ctx.Input().ReadRune()
	if err != nil {
		return '\n', true
	}
	if r == '\r' {
		return NextRune(ctx)
	}
	return r, false
}

// NextLine returns the text in the inputFrom up to the next newline. The characters are gotten using nextByte.
func NextLine(ctx mycocontext.Context) (line string, done bool) {
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

// IsEof is true if there is nothing left to read in the input. It does not handle the case when all next characters are \r, which are never returned by NextRune, thus making this function lie.
//
// Be not afraid because everyone lies. Not a good idea to trust a //function// anyway.
func IsEof(ctx mycocontext.Context) bool {
	return ctx.Input().Len() == 0
}
