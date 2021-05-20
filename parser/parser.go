// Package parser turns the source text into a sequence of blocks.
package parser

import (
	"context"
	"strings"
)

// Parse parses the Mycomarkup document in the given context. All parsed blocks are written to out.
func Parse(ctx context.Context, out chan interface{}) {
	var (
		state = parserState{}
		token interface{}
		done  bool
	)
	defer close(out)
	for !done {
		select {
		case <-ctx.Done():
			return
		default:
			token, done = nextToken(ctx, &state)
			if token != nil {
				out <- token
			}
		}
	}
}

// nextByte returns the next byte in the inputFrom. The CR byte (\r) is never returned, if there is a CR in the inputFrom, the byte after it is returned. If there is no next byte, the NL byte (\n) is returned and done is true.
func nextByte(ctx context.Context) (b byte, done bool) {
	b, err := inputFrom(ctx).ReadByte()
	if err != nil {
		return '\n', true
	}
	if b == '\r' {
		return nextByte(ctx)
	}
	return b, false
}

// nextLine returns the text in the inputFrom up to the next newline. The characters are gotten using nextByte.
func nextLine(ctx context.Context) (line string, done bool) {
	var (
		lineBuffer strings.Builder
		b          byte
	)
	b, done = nextByte(ctx)
	for b != '\n' {
		lineBuffer.WriteByte(b)
		b, done = nextByte(ctx)
	}
	return lineBuffer.String(), done
}
