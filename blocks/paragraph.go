package blocks

import (
	"bytes"
	"fmt"
	"github.com/bouncepaw/mycomarkup/links"
	"html"
	"strings"
	"unicode"
)

type Paragraph struct {
	Formatted
}

func (p Paragraph) ID(counter *IDCounter) string {
	counter.paragraphs++
	return fmt.Sprintf("paragraph-%d", counter.paragraphs)
}

type Formatted struct {
	HyphaName string
	Html      string
	*bytes.Buffer
	spans []span
}

func (p Formatted) IsBlock() {}

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
	spanLink
)

func tagFromState(stt spanTokenType, tagState map[spanTokenType]bool, tagName, originalForm string) string {
	if tagState[spanMono] && (stt != spanMono) {
		return originalForm
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
	href, text, class := link.Href(), html.EscapeString(link.Display()), html.EscapeString(link.Classes())
	return fmt.Sprintf(`<a href="%s" class="%s">%s</a>`, href, class, text)
}

func (p *Formatted) AddLine(line string) {
	p.Html += `<br>` + ParagraphToHtml(p.HyphaName, line)
}

func MakeFormatted(input, hyphaName string) Formatted {
	return Formatted{
		HyphaName: hyphaName,
		Html:      ParagraphToHtml(hyphaName, input),
	}
}

func ParagraphToHtml(hyphaName, input string) string {
	var (
		p = &Formatted{
			hyphaName,
			"",
			bytes.NewBufferString(input),
			make([]span, 0),
		}
		ret strings.Builder
		// true = tag is opened, false = tag is not opened
		tagState = map[spanTokenType]bool{
			spanItalic: false,
			spanBold:   false,
			spanMono:   false,
			spanSuper:  false,
			spanSub:    false,
			spanMark:   false,
			spanLink:   false,
		}
		startsWith = func(t string) bool {
			return bytes.HasPrefix(p.Bytes(), []byte(t))
		}
		noTagsActive = func() bool {
			return !(tagState[spanItalic] || tagState[spanBold] || tagState[spanMono] || tagState[spanSuper] || tagState[spanSub] || tagState[spanMark] || tagState[spanLink])
		}
	)

	for p.Len() != 0 {
		switch {
		case startsWith("//"):
			ret.WriteString(tagFromState(spanItalic, tagState, "em", "//"))
			p.Next(2)
		case startsWith("**"):
			ret.WriteString(tagFromState(spanBold, tagState, "strong", "**"))
			p.Next(2)
		case startsWith("`"):
			ret.WriteString(tagFromState(spanMono, tagState, "code", "`"))
			p.Next(1)
		case startsWith("^^"):
			ret.WriteString(tagFromState(spanSuper, tagState, "sup", "^^"))
			p.Next(2)
		case startsWith(",,"):
			ret.WriteString(tagFromState(spanSub, tagState, "sub", ",,"))
			p.Next(2)
		case startsWith("!!"):
			ret.WriteString(tagFromState(spanMark, tagState, "mark", "!!"))
			p.Next(2)
		case startsWith("~~"):
			ret.WriteString(tagFromState(spanMark, tagState, "s", "~~"))
			p.Next(2)
		case startsWith("[["):
			ret.WriteString(getLinkNode(p, hyphaName, true))
		case (startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")) && noTagsActive():
			ret.WriteString(getLinkNode(p, hyphaName, false))
		default:
			ret.WriteString(html.EscapeString(getSpanText(p).htmlWithState(tagState)))
		}
	}

	for stt, open := range tagState {
		if open {
			switch stt {
			case spanItalic:
				ret.WriteString(tagFromState(spanItalic, tagState, "em", "//"))
			case spanBold:
				ret.WriteString(tagFromState(spanBold, tagState, "strong", "**"))
			case spanMono:
				ret.WriteString(tagFromState(spanMono, tagState, "code", "`"))
			case spanSuper:
				ret.WriteString(tagFromState(spanSuper, tagState, "sup", "^^"))
			case spanSub:
				ret.WriteString(tagFromState(spanSub, tagState, "sub", ",,"))
			case spanMark:
				ret.WriteString(tagFromState(spanMark, tagState, "mark", "!!"))
			case spanStrike:
				ret.WriteString(tagFromState(spanMark, tagState, "s", "~~"))
			case spanLink:
				ret.WriteString(tagFromState(spanLink, tagState, "a", "[["))
			}
		}
	}

	return ret.String()
}
