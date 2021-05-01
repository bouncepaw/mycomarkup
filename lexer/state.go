package lexer

import (
	"bytes"
)

type SourceText struct {
	// General:
	b *bytes.Buffer

	// Configuration:
	allowMultilineParagraph bool
	terminateOnCloseBrace   bool
}
