package blocks

import (
	"bytes"
	"strings"
)

type span interface {
	tagName() string
	htmlWithState(map[spanTokenType]bool) string
}

type spanText struct {
	bytes.Buffer
}

func (s spanText) tagName() string {
	return "p"
}

func (s spanText) htmlWithState(_ map[spanTokenType]bool) string {
	return s.String()
}

func getSpanText(p *Formatted) spanText {
	var (
		st         = spanText{}
		escaping   = false
		startsWith = func(t string) bool {
			return bytes.HasPrefix(p.Bytes(), []byte(t))
		}
		couldBeLinkStart = func() bool {
			return startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")
		}
	)

	// Always read the first byte in advance to avoid endless loops that kill computers (sad experience)
	if p.Len() != 0 {
		b, _ := p.ReadByte()
		_ = st.WriteByte(b)
	}
	for p.Len() != 0 {
		// We check for length, this should never fail:
		ch, _ := p.ReadByte()
		if escaping {
			st.WriteByte(ch)
			escaping = false
		} else if ch == '\\' {
			escaping = true
		} else if strings.IndexByte("/*`^,+[~_", ch) >= 0 {
			p.UnreadByte()
			break
		} else if couldBeLinkStart() {
			st.WriteByte(ch)
			break
		} else {
			st.WriteByte(ch)
		}
	}

	return st
}

type spanStyle struct {
	kind spanTokenType
}

type spanTableEntry struct {
	kind        spanTokenType
	token       string
	tokenLength int
	htmlTagName string
}

// cool table for cool kids
var spanTable = []spanTableEntry{
	{spanItalic, "//", 2, "em"},
	{spanBold, "**", 2, "strong"},
	{spanMono, "`", 1, "code"},
	{spanSuper, "^^", 2, "sup"},
	{spanSub, ",,", 2, "sub"},
	{spanMark, "++", 2, "mark"},
	{spanStrike, "~~", 2, "s"},
	{spanUnderline, "__", 2, "u"},
}

func entryForSpan(kind spanTokenType) spanTableEntry {
	for _, entry := range spanTable {
		if entry.kind == kind {
			return entry
		}
	}
	// unreachable state, supposedly
	panic("unknown kind of span")
}

func tagNameForSpan(kind spanTokenType) string {
	return entryForSpan(kind).htmlTagName
}
