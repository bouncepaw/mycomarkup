package lexer

// byteworks is for byte-related operations.

import (
	"bytes"
	"strings"
)

func startsWithStr(b *bytes.Buffer, s string) bool {
	return strings.HasPrefix(b.String(), s)
}

func eatChar(s *SourceText) {
	// how confident i am
	_, _ = s.b.ReadByte()
}
