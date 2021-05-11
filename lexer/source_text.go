package lexer

import (
	"bytes"
)

type SourceText struct {
	// General:
	hyphaName string
	b         *bytes.Buffer

	// Configuration for paragraph-lexing:
	lexingParagraphOnly     bool
	allowMultilineParagraph bool
	terminateOnCloseBrace   bool
}
