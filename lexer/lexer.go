package lexer

import (
	"bytes"

	"github.com/bouncepaw/mycomarkup/blocks"
)

type Lexeme struct {
	pos   Position
	block blocks.Block
}

func startsWithStrFromI(b []byte, i int, s string) bool {
	return bytes.HasPrefix(b[i:], []byte(s))
}
