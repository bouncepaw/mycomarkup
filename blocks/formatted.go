package blocks

import (
	"bytes"
)

// Formatted is a piece of formatted text. It is always part of a bigger block, such as Paragraph.
type Formatted struct {
	// HyphaName is the name of the hypha that contains the formatted text.
	HyphaName string
	Html      string
	*bytes.Buffer
	Spans []interface{} // Forgive me, for I have sinned
}

func (p Formatted) isBlock() {}

// ID returns an empty string because Formatted is always part of a bigger block.
func (p Formatted) ID(_ *IDCounter) string {
	return ""
}

// SpanKind is a kind of a span, such as italic, bold, etc.
type SpanKind int

const (
	SpanItalic SpanKind = iota
	SpanBold
	SpanMono
	SpanSuper
	SpanSub
	SpanMark
	SpanStrike
	SpanUnderline
	SpanLink
	// SpanNewLine represents a linebreak (\n) in the formatted text.
	SpanNewLine
)

type SpanTableEntry struct {
	Kind        SpanKind
	Token       string
	TokenLength int
	HtmlTagName string
}

// SpanTable is a table for easier Span lexing.
var SpanTable = []SpanTableEntry{
	{SpanItalic, "//", 2, "em"},
	{SpanBold, "**", 2, "strong"},
	{SpanMono, "`", 1, "code"},
	{SpanSuper, "^^", 2, "sup"},
	{SpanSub, ",,", 2, "sub"},
	{SpanMark, "++", 2, "mark"},
	{SpanStrike, "~~", 2, "s"},
	{SpanUnderline, "__", 2, "u"},
}

func entryForSpan(kind SpanKind) SpanTableEntry {
	for _, entry := range SpanTable {
		if entry.Kind == kind {
			return entry
		}
	}
	// unreachable state, supposedly
	panic("unknown kind of Span")
}

func TagNameForSpan(kind SpanKind) string {
	return entryForSpan(kind).HtmlTagName
}
