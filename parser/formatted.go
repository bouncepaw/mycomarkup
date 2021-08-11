package parser

import (
	"bytes"
	"fmt"
	"html"
	"strings"
	"unicode"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/globals"
	"github.com/bouncepaw/mycomarkup/links"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"github.com/bouncepaw/mycomarkup/util"
)

func nextParagraph(ctx mycocontext.Context) (p blocks.Paragraph, done bool) {
	line, done := mycocontext.NextLine(ctx)
	p = blocks.Paragraph{MakeFormatted(line, ctx.HyphaName())}
	if nextLineIsSomething(ctx) {
		return
	}
	for {
		line, done = mycocontext.NextLine(ctx)
		if done && line == "" {
			break
		}
		p.AddLine(line)
		if nextLineIsSomething(ctx) {
			break
		}
	}
	return
}

// TagFromState returns an appropriate tag half (<left> or </right>) depending on tagState and also mutates it.
func TagFromState(stt blocks.SpanKind, tagState map[blocks.SpanKind]bool) string {
	var tagName string
	if stt == blocks.SpanLink {
		tagName = "a"
	} else {
		tagName = blocks.TagNameForStyleSpan(stt)
	}
	if tagState[stt] {
		tagState[stt] = false
		return fmt.Sprintf("</%s>", tagName)
	} else {
		tagState[stt] = true
		return fmt.Sprintf("<%s>", tagName)
	}
}

// GetLinkNode returns an HTML representation of the next link in the input. Set isBracketedLink if the input starts with [[.
//
// TODO: return something bearable, not HTML.
func GetLinkNode(input *bytes.Buffer, hyphaName string, isBracketedLink bool) string {
	if isBracketedLink {
		input.Next(2) // drop those [[
	}
	var (
		escaping   = false
		addrBuf    = bytes.Buffer{}
		displayBuf = bytes.Buffer{}
		currBuf    = &addrBuf
	)
	for input.Len() != 0 {
		b, _ := input.ReadByte()
		if escaping {
			currBuf.WriteByte(b)
			escaping = false
		} else if isBracketedLink && b == '|' && currBuf == &addrBuf {
			currBuf = &displayBuf
		} else if isBracketedLink && b == ']' && bytes.HasPrefix(input.Bytes(), []byte{']'}) {
			input.Next(1)
			break
		} else if !isBracketedLink && (unicode.IsSpace(rune(b)) || strings.ContainsRune("<>{}|\\^[]`,()", rune(b))) {
			_ = input.UnreadByte()
			break
		} else {
			currBuf.WriteByte(b)
		}
	}

	link := links.From(addrBuf.String(), displayBuf.String(), hyphaName)
	if globals.HyphaExists(util.CanonicalName(link.TargetHypha())) {
		link.MarkAsExisting()
	}
	href, text, class := link.Href(), html.EscapeString(link.Display()), html.EscapeString(link.Classes())
	return fmt.Sprintf(`<a href="%s" class="%s">%s</a>`, href, class, text)
}

// MakeFormatted parses the formatted text in the input and returns it.
func MakeFormatted(firstLine, hyphaName string) blocks.Formatted {
	return blocks.Formatted{
		HyphaName: hyphaName,
		Lines:     []string{firstLine},
	}
}

func stateAtNewLine() map[blocks.SpanKind]bool {
	// false: tag not open
	// true: tag open
	return map[blocks.SpanKind]bool{
		blocks.SpanItalic:    false,
		blocks.SpanBold:      false,
		blocks.SpanMono:      false,
		blocks.SpanSuper:     false,
		blocks.SpanSub:       false,
		blocks.SpanMark:      false,
		blocks.SpanStrike:    false,
		blocks.SpanUnderline: false,
		blocks.SpanLink:      false,
	}
}

// ParagraphToHtml turns the given line of formatted text into HTML by lexing and parsing it in place.
//
// TODO: separate all these steps.
func ParagraphToHtml(hyphaName, input string) string {
	var (
		p = &blocks.Formatted{
			HyphaName: hyphaName,
			Lines:     []string{},
			Buffer:    bytes.NewBufferString(input),
			Spans:     make([]interface{}, 0),
		}
		ret strings.Builder
		// true = tag is opened, false = tag is not opened
		tagState   = stateAtNewLine()
		startsWith = func(t string) bool {
			return bytes.HasPrefix(p.Bytes(), []byte(t))
		}
		noTagsActive = func() bool {
			// This function used to be one boolean expression. I changed it to a loop so it is harder to forger 💀 any span kinds.
			for _, entry := range blocks.SpanTable {
				if tagState[entry.Kind] { // If span is open
					return false
				}
			}
			// All other spans are closed, let's check for link finally.
			return !tagState[blocks.SpanLink]
		}
	)

runeWalker:
	for p.Len() != 0 {
		for _, entry := range blocks.SpanTable {
			if startsWith(entry.Token) {
				p.Spans = append(p.Spans, entry)
				p.Next(len(entry.Token))
				continue runeWalker
			}
		}
		switch {
		case startsWith("[["):
			p.Spans = append(p.Spans, GetLinkNode(p.Buffer, hyphaName, true))
		case (startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")) && noTagsActive():
			p.Spans = append(p.Spans, GetLinkNode(p.Buffer, hyphaName, false))
		default:
			p.Spans = append(p.Spans, GetSpanText(p.Buffer).HTMLWithState(tagState))
		}
	}

	for _, span := range p.Spans {
		switch s := span.(type) {
		case blocks.SpanTableEntry:
			ret.WriteString(TagFromState(s.Kind, tagState))
		case string:
			ret.WriteString(s)
		default:
			panic("unknown Span kind... What do you expect from me?")
		}
	}

	for stt, open := range tagState {
		if open {
			ret.WriteString(TagFromState(stt, tagState))
		}
	}

	return ret.String()
}

type SpanText struct {
	bytes.Buffer
}

func (s SpanText) HTMLWithState(_ map[blocks.SpanKind]bool) string {
	return s.String()
}

// GetSpanText returns the next SpanText there is in input.
func GetSpanText(input *bytes.Buffer) SpanText {
	var (
		st         = SpanText{}
		escaping   = false
		startsWith = func(t string) bool {
			return bytes.HasPrefix(input.Bytes(), []byte(t))
		}
		couldBeLinkStart = func() bool {
			return startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")
		}
	)

	// Always read the first byte in advance to avoid endless loops that kill computers (sad experience)
	if input.Len() != 0 {
		b, _ := input.ReadByte()
		_ = st.WriteByte(b)
	}
	for input.Len() != 0 {
		// We check for length, this should never fail:
		ch, _ := input.ReadByte()
		if escaping {
			st.WriteByte(ch)
			escaping = false
		} else if ch == '\\' {
			escaping = true
		} else if strings.IndexByte("/*`^,+[~_", ch) >= 0 { // TODO: generate that string there dynamically
			input.UnreadByte()
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
