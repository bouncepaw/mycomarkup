package blocks

import (
	"fmt"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/links"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/util"
)

// Formatted is a piece of formatted text. It is always part of a bigger block, such as Paragraph.
type Formatted struct {
	// HyphaName is the name of the hypha that contains the Formatted text.
	HyphaName string
	// Lines are parsed lines of the Formatted text.
	Lines [][]Span
}

// ID returns an empty string because Formatted is always part of a bigger block that has a meaningful ID.
func (p Formatted) ID(_ *IDCounter) string {
	return ""
}

// AddLine stores an additional line of the formatted text. V3
func (p *Formatted) AddLine(line []Span) {
	p.Lines = append(p.Lines, line)
}

// Span is a piece of Formatted text. There are three implementors of this interface: SpanTableEntry (styles), InlineLink ([[links]]), InlineText (usual text).
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
)

// SpanTableEntry is an entry of SpanTable and simultaneously a Span.
type SpanTableEntry struct {
	kind    SpanKind
	Token   string
	htmlTag string
}

// Kind returns one of SpanKind. See the first column of SpanTable to learn what values are possible.
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

// InlineLink is a link that is part of a Formatted text.
type InlineLink struct {
	links.Link
}

// Kind returns SpanLink.
func (il InlineLink) Kind() SpanKind {
	return SpanLink
}

// InlineText is the most wholesome thing in Mycomarkup, just a bunch of characters with no formatting or any other special meaning.
type InlineText struct {
	Contents string
}

// Kind returns SpanText.
func (it InlineText) Kind() SpanKind {
	return SpanText
}

// TagFromState returns an appropriate tag half (<left> or </right>) depending on tagState and also mutates it. V3
//
// TODO: get rid of.
func TagFromState(stt SpanKind, tagState map[SpanKind]bool) string {
	tagName := util.TernaryConditionString(stt == SpanLink, "a", entryForSpan(stt).htmlTag)
	if tagState[stt] {
		tagState[stt] = false
		return fmt.Sprintf("</%s>", tagName)
	} else {
		tagState[stt] = true
		return fmt.Sprintf("<%s>", tagName)
	}
}

// CleanStyleState returns a map where keys are SpanKind representing inline style and values are booleans. Mutate this map to keep track of active and inactive styles.
//
// Values:
// `false`: the style is not active
// `true`: the style is active
//
// For example, for a Formatted line like that:
//     **Welcome** to //California
// `CleanStyleState()[SpanItalic] == true` at the end of string.
func CleanStyleState() map[SpanKind]bool {
	return map[SpanKind]bool{
		SpanItalic:    false,
		SpanBold:      false,
		SpanMono:      false,
		SpanSuper:     false,
		SpanSub:       false,
		SpanMark:      false,
		SpanStrike:    false,
		SpanUnderline: false,
		SpanLink:      false,
	}
}
