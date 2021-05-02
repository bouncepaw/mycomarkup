package lexer

// byteworks is for byte-related operations.

import (
	"bytes"
	"strings"
)

func startsWithStr(b *bytes.Buffer, s string) bool {
	return strings.HasPrefix(b.String(), s)
}

func eatChar(st *SourceText) {
	// how confident i am
	_, _ = st.b.ReadByte()
}

func eatN(st *SourceText, n uint) {
	for n != 0 {
		eatChar(st)
		n--
	}
}
