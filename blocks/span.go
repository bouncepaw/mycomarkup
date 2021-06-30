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
		} else if strings.IndexByte("/*`^,+[~", ch) >= 0 {
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
