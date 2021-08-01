package blocks

import (
	"bytes"
	"fmt"
	"github.com/bouncepaw/mycomarkup/globals"
	"github.com/bouncepaw/mycomarkup/links"
	"github.com/bouncepaw/mycomarkup/util"
	"html"
	"strings"
	"unicode"
)

// Paragraph is a block of formatted text.
type Paragraph struct {
	Formatted
}

// ID returns the paragraph's id which is paragraph- and a number.
func (p Paragraph) ID(counter *IDCounter) string {
	counter.paragraphs++
	return fmt.Sprintf("paragraph-%d", counter.paragraphs)
}

// Formatted is a piece of formatted text.
type Formatted struct {
	HyphaName string
	Html      string
	*bytes.Buffer
	spans []interface{} // Forgive me, for I have sinned
}

func (p Formatted) isBlock() {}

// ID returns an empty string because Formatted is always part of a bigger block.
func (p Formatted) ID(_ *IDCounter) string {
	return ""
}

type spanTokenType int

const (
	_ = iota
	spanItalic
	spanBold
	spanMono
	spanSuper
	spanSub
	spanMark
	spanStrike
	spanUnderline
	spanLink
)

func tagFromState(stt spanTokenType, tagState map[spanTokenType]bool) string {
	var tagName string
	if stt == spanLink {
		tagName = "a"
	} else {
		tagName = tagNameForSpan(stt)
	}
	if tagState[stt] {
		tagState[stt] = false
		return fmt.Sprintf("</%s>", tagName)
	} else {
		tagState[stt] = true
		return fmt.Sprintf("<%s>", tagName)
	}
}

func getLinkNode(input *Formatted, hyphaName string, isBracketedLink bool) string {
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

// AddLine adds a line to the block. The line is prepended with <br>.
func (p *Formatted) AddLine(line string) {
	p.Html += `<br>` + paragraphToHtml(p.HyphaName, line)
}

// MakeFormatted parses the formatted text in the input and returns it.
func MakeFormatted(input, hyphaName string) Formatted {
	return Formatted{
		HyphaName: hyphaName,
		Html:      paragraphToHtml(hyphaName, input),
	}
}

func paragraphToHtml(hyphaName, input string) string {
	var (
		p = &Formatted{
			hyphaName,
			"",
			bytes.NewBufferString(input),
			make([]interface{}, 0),
		}
		ret strings.Builder
		// true = tag is opened, false = tag is not opened
		tagState = map[spanTokenType]bool{
			spanItalic:    false,
			spanBold:      false,
			spanMono:      false,
			spanSuper:     false,
			spanSub:       false,
			spanMark:      false,
			spanStrike:    false,
			spanUnderline: false,
			spanLink:      false,
		}
		startsWith = func(t string) bool {
			return bytes.HasPrefix(p.Bytes(), []byte(t))
		}
		noTagsActive = func() bool {
			return !(tagState[spanItalic] || tagState[spanBold] || tagState[spanMono] || tagState[spanSuper] || tagState[spanSub] || tagState[spanMark] || tagState[spanLink])
		}
	)

	for p.Len() != 0 {
		for _, entry := range spanTable {
			if startsWith(entry.token) {
				p.spans = append(p.spans, entry)
				p.Next(entry.tokenLength)
				continue
			}
		}
		switch {
		case startsWith("[["):
			p.spans = append(p.spans, getLinkNode(p, hyphaName, true))
		case (startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")) && noTagsActive():
			p.spans = append(p.spans, getLinkNode(p, hyphaName, false))
		default:
			p.spans = append(p.spans, getSpanText(p).htmlWithState(tagState))
		}
	}

	for _, span := range p.spans {
		switch s := span.(type) {
		case spanTableEntry:
			ret.WriteString(tagFromState(s.kind, tagState))
		case string:
			ret.WriteString(s)
		default:
			panic("unknown span kind... What do you expect from me?")
		}
	}

	for stt, open := range tagState {
		if open {
			ret.WriteString(tagFromState(stt, tagState))
		}
	}

	return ret.String()
}

// TODO: test for HTML injections thoroughly
