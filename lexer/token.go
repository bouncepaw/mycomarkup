package lexer

import (
	"fmt"
)

type TokenKind int

const (
	TokenErr TokenKind = iota
	TokenBraceOpen
	TokenBraceClose
	TokenNewLine

	TokenHorizontalLine
	TokenPreformattedFence
	TokenPreformattedAltText

	TokenHeadingOpen
	TokenHeadingClose

	TokenSpanNewLine
	TokenSpanText
	TokenSpanItalic
	TokenSpanBold
	TokenSpanMonospace
	TokenSpanMarker
	TokenSpanSuper
	TokenSpanSub
	TokenSpanStrike
	TokenSpanUnderline

	TokenSpanLinkOpen
	TokenSpanLinkClose
	TokenLinkAddress
	TokenLinkDisplayOpen
	TokenLinkDisplayClose
	TokenAutoLink

	TokenRocketLinkOpen
	TokenRocketLinkClose

	TokenBlockQuoteOpen
	TokenBlockQuoteClose

	TokenBulletUnnumbered
	TokenBulletIndent
	TokenBulletNumberedImplicit
	TokenBulletNumberedExplicit

	TokenTagImg
	TokenTagTable

	TokenParagraphOpen
	TokenParagraphClose
)

// Token is a type representing a token. Each token can be of one kind and may have a non-empty value, if it matters.
type Token struct {
	kind  TokenKind
	value string
}

// String is used for debugging output.
func (t *Token) String() string {
	return fmt.Sprintf(`[%v →%s]`, t.kind, t.value)
}
