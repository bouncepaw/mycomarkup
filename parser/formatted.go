package parser

import (
	"bytes"
	"fmt"
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

// nextInlineLink returns an HTML representation of the next link in the input. Set isBracketedLink if the input starts with [[.
func nextInlineLink(input *bytes.Buffer, hyphaName string, isBracketedLink bool) blocks.InlineLink {
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
	return blocks.InlineLink{Link: link}
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

func spansFromLine(hyphaName, line string) ([]blocks.Span, map[blocks.SpanKind]bool) {
	var (
		input      = bytes.NewBufferString(line)
		spans      = make([]blocks.Span, 0)
		tagState   = stateAtNewLine()
		startsWith = func(t string) bool {
			return bytes.HasPrefix(input.Bytes(), []byte(t))
		}
		noTagsActive = func() bool {
			// This function used to be one boolean expression. I changed it to a loop so it is harder to forger 💀 any span kinds.
			for _, entry := range blocks.SpanTable {
				if tagState[entry.Kind()] { // If span is open
					return false
				}
			}
			// All other spans are closed, let's check for link finally.
			return !tagState[blocks.SpanLink]
		}
	)

runeWalker:
	for input.Len() != 0 {
		for _, entry := range blocks.SpanTable {
			if startsWith(entry.Token) {
				spans = append(spans, entry)
				input.Next(len(entry.Token))
				continue runeWalker
			}
		}
		switch {
		case startsWith("[["):
			spans = append(spans, nextInlineLink(input, hyphaName, true))
		case (startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")) && noTagsActive():
			spans = append(spans, nextInlineLink(input, hyphaName, false))
		default:
			spans = append(spans, nextInlineText(input))
		}
	}

	return spans, tagState
}

// FormattedLineToHTML turns the given line of formatted text into HTML by lexing and parsing it in place.
//
// TODO: separate all these steps.
func FormattedLineToHTML(hyphaName, input string) string {
	var (
		spans, tagState = spansFromLine(hyphaName, input)
		ret             = bytes.Buffer{}
	)

	for _, span := range spans {
		switch s := span.(type) {
		case blocks.SpanTableEntry:
			ret.WriteString(TagFromState(s.Kind(), tagState))
		case blocks.InlineLink:
			ret.WriteString(
				fmt.Sprintf(
					`<a href="%s" class="%s">%s</a>`,
					s.Href(),
					s.Classes(),
					s.Display(),
				)) // TODO: test for XSS
		case blocks.InlineText:
			ret.WriteString(s.Contents) // TODO: test for XSS
		default:
			panic("unknown inline block/span... What do you expect from me?") // gotta think about terminology
		}
	}

	for stt, open := range tagState {
		if open {
			ret.WriteString(TagFromState(stt, tagState))
		}
	}

	return ret.String()
}

var protocols [][]byte

func init() {
	protocols = [][]byte{
		[]byte("https://"),
		[]byte("http://"),
		[]byte("gemini://"),
		[]byte("gopher://"),
		[]byte("ftp://")}
	// There was a demand for a way to customize the protocols ^. Do we need that?
}
func bytesStartWithProtocol(b []byte) bool {
	for _, protocol := range protocols {
		if bytes.HasPrefix(b, protocol) {
			return true
		}
	}
	return false
}

// nextInlineText returns the next blocks.InlineText there is in input.
func nextInlineText(input *bytes.Buffer) blocks.InlineText {
	var (
		ret      = bytes.Buffer{}
		escaping = false
	)

	// Always read the first byte in advance to avoid endless loops that kill computers (sad experience)
	if input.Len() != 0 {
		b, _ := input.ReadByte()
		_ = ret.WriteByte(b)
	}
	for input.Len() != 0 {
		// We check for length, this should never fail:
		ch, _ := input.ReadByte()
		if escaping {
			ret.WriteByte(ch)
			escaping = false
		} else if ch == '\\' {
			escaping = true
		} else if strings.IndexByte("/*`^,+[~_", ch) >= 0 { // TODO: generate that string there dynamically
			input.UnreadByte() // sorry, wrong door >_<
			break
		} else if bytesStartWithProtocol(input.Bytes()) {
			ret.WriteByte(ch)
			break
		} else {
			ret.WriteByte(ch)
		}
	}

	return blocks.InlineText{Contents: ret.String()}
}
