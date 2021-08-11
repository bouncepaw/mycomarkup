package blocks

import (
	"bytes"
)

// Formatted is a piece of formatted text. It is always part of a bigger block, such as Paragraph.
type Formatted struct {
	// HyphaName is the name of the hypha that contains the formatted text.
	HyphaName string
	Lines     []string
	*bytes.Buffer
	Spans []interface{} // Forgive me, for I have sinned
}

func (p Formatted) isBlock() {}

// ID returns an empty string because Formatted is always part of a bigger block.
func (p Formatted) ID(_ *IDCounter) string {
	return ""
}

// AddLine stores an additional line of the formatted text.
func (p *Formatted) AddLine(line string) {
	p.Lines = append(p.Lines, line)
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

// SpanTableEntry is an entry of SpanTable.
type SpanTableEntry struct {
	Kind    SpanKind
	Token   string
	HTMLTag string
}

// SpanTable is a table for easier Span lexing, its entries are also nice to fit into Formatted.Spans.
var SpanTable = []SpanTableEntry{
	{SpanItalic, "//", "em"},
	{SpanBold, "**", "strong"},
	{SpanMono, "`", "code"},
	{SpanSuper, "^^", "sup"},
	{SpanSub, ",,", "sub"},
	{SpanMark, "++", "mark"},
	{SpanStrike, "~~", "s"},
	{SpanUnderline, "__", "u"},
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

// TagNameForStyleSpan returns an appropriate HTML tag for the span. Note that the <a> tag is not in the table.
func TagNameForStyleSpan(kind SpanKind) string {
	return entryForSpan(kind).HTMLTag
}
