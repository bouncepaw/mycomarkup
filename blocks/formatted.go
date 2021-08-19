package blocks

import (
	"bytes"
	"github.com/bouncepaw/mycomarkup/links"
)

// Formatted is a piece of formatted text. It is always part of a bigger block, such as Paragraph.
type Formatted struct {
	// HyphaName is the name of the hypha that contains the formatted text.
	HyphaName string
	// Lines is where lines of the formatted text are stored. They are parsed afterwards. TODO: get rid of maybe
	Lines []string
	*bytes.Buffer
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

type Span interface {
	Kind() SpanKind
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
	SpanText
	// SpanNewLine represents a linebreak (\n) in the formatted text.
	SpanNewLine
)

// SpanTableEntry is an entry of SpanTable and simultaneously
type SpanTableEntry struct {
	kind    SpanKind
	Token   string
	HTMLTag string
}

func (ste SpanTableEntry) Kind() SpanKind {
	return ste.kind
}

// SpanTable is a table for easier span lexing, its entries are also Span too.
var SpanTable = []SpanTableEntry{ // it is so cute so cute
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
		if entry.Kind() == kind {
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

type InlineLink struct {
	*links.Link
}

func (il InlineLink) Kind() SpanKind {
	return SpanLink
}

type InlineText struct {
	contents string
}

func (it InlineText) Kind() SpanKind {
	return SpanText
}
