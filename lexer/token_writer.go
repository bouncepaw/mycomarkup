package lexer

import (
	"bytes"
)

// TokenWriter stores tokens found during lexing and stores the state needed for that. Short: tw.
type TokenWriter struct {
	StateStack
	buf *bytes.Buffer

	elements    []Token
	lastElement *Token
}

// appendToken appends the given token and saves a pointer to that token for quick access.
func (tw *TokenWriter) appendToken(token Token) {
	tw.lastElement = &token
	tw.elements = append(tw.elements, token)
}

// bufIntoToken takes the buffer contents and creates a token of the given kind and with the contents of the buffer. Then it resets the buffer.
func (tw *TokenWriter) bufIntoToken(kind TokenKind) {
	tw.appendToken(Token{kind, tw.buf.String()})
	tw.buf.Reset()
}

// nonEmptyBufIntoToken calls bufIntoToken if the buffer is not empty.
func (tw *TokenWriter) nonEmptyBufIntoToken(kind TokenKind) {
	if tw.buf.Len() != 0 {
		tw.bufIntoToken(kind)
	}
}
