package lexer

import (
	"bytes"
)

// TokenWriter stores tokens found during lexing and stores the state needed for that. Short: tw.
type TokenWriter struct {
	StateStack
	// Current buffer of characters. Emptied and filled at will.
	buf *bytes.Buffer

	// List of saved tokens. New tokens could be added, older should not be modified.
	savedTokens []Token
	// Pointer to last saved token.
	recentToken *Token
}

// appendToken appends the given token and saves a pointer to that token for quick access.
func (tw *TokenWriter) appendToken(token Token) {
	tw.recentToken = &token
	tw.savedTokens = append(tw.savedTokens, token)
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
