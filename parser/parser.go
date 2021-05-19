// Package parser turns the source text into a sequence of blocks.
package parser

import (
	"context"
	"strings"
)

// Parse parses the Mycomarkup document in the given context. All parsed blocks are written to out.
func Parse(ctx context.Context, out chan interface{}) {
	for {
		block, ok := getNextBlock(ctx)
		select {
		case <-ctx.Done():
			close(out)
			return
		default:
			if !ok {
				close(out)
				return
			}
			out <- block
		}
	}
}

// getNextBlock is ok when there is a block to return. The block is one of the types in the blocks package.
func getNextBlock(ctx context.Context) (block interface{}, ok bool) {
	return nil, false
}

// nextByte returns the next byte in the input. The CR byte (\r) is never returned, if there is a CR in the input, the byte after it is returned. If there is no next byte, the NL byte (\n) is returned and done is true.
func nextByte(ctx context.Context) (b byte, done bool) {
	b, err := input(ctx).ReadByte()
	if err != nil {
		return '\n', true
	}
	if b == '\r' {
		return nextByte(ctx)
	}
	return b, false
}

// nextLine returns the text in the input up to the next newline. The characters are gotten using nextByte.
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
