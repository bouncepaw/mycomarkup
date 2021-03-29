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

	TokenSpanLinkOpen
	TokenSpanLinkClose
	TokenLinkAddress
	TokenLinkDisplay
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
)

type Token struct {
	kind  TokenKind
	value string
}

func (t *Token) String() string {
	return fmt.Sprintf(`[%v →%s]`, t.kind, t.value)
}
