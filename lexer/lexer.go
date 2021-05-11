// Package lexer contains a lexer for mycomarkup and the stuff to make it work.
//
// The lexing is based on finite-state automata.
package lexer

import "bytes"

// Lex generates tokens from the source text.
func Lex(st SourceText, hyphaName string) []Token {
	panic("not implemented")
	return nil
}

// LexParagraph is a stripped down version of LexString that only supports a subset of mycomarkup: paragraphs. It returns the bytes it didn't consume during the lexing.
func LexParagraph(st SourceText) (tokens []Token, rest []byte) {
	st.lexingParagraphOnly = true
	tw := &TokenWriter{
		StateStack:  newStateStack(),
		buf:         &bytes.Buffer{},
		savedTokens: make([]Token, 0),
	}
	lexParagraph(&st, tw)
	return tw.savedTokens, st.b.Bytes()
}
