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
	TokenHeading

	TokenSpanText
	TokenSpanItalic
	TokenSpanBold
	TokenSpanMonospace
	TokenSpanMarker
	TokenSpanSuper
	TokenSpanSub
	TokenSpanStrike
	TokenSpanLinkOpen
	TokenSpanLinkSeparator
	TokenSpanLinkClose
	TokenAutoLink

	TokenRocketLink
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
	startLine   uint
	startColumn uint
	kind        TokenKind
	value       string
}

func (t *Token) String() string {
	return fmt.Sprintf(`%d:%d→%v:"%s"`, t.startLine, t.startColumn, t.kind, t.value)
}
