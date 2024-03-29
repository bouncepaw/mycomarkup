package parser

import (
	"bytes"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/parser/ctxio"
	"strings"
	"unicode"

	"git.sr.ht/~bouncepaw/mycomarkup/v5/blocks"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/links"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
)

func nextParagraph(ctx mycocontext.Context) (p blocks.Paragraph, done bool) {
	line, done := ctxio.NextLine(ctx)
	p = blocks.Paragraph{MakeFormatted(ctx, line)}
	if nextLineIsSomething(ctx) {
		return
	}
	for {
		line, done = ctxio.NextLine(ctx)
		if done && line == "" {
			break
		}
		spans := spansFromLine(ctx, p.HyphaName, line)
		p.AddLine(spans)
		if nextLineIsSomething(ctx) {
			break
		}
	}
	return
}

// nextInlineLink returns an HTML representation of the next link in the input. Set isBracketedLink if the input starts with [[.
func nextInlineLink(ctx mycocontext.Context, input *bytes.Buffer, hyphaName string, isBracketedLink bool) blocks.InlineLink {
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

	link := links.LinkFrom(ctx, addrBuf.String(), displayBuf.String())
	// Color if needed. One hypha storage iteration.
	if probe := link.HyphaProbe(ctx); probe != nil {
		ctx.Options().IterateHyphaNamesWith(probe)

	}
	return blocks.InlineLink{Link: link}
}

// MakeFormatted parses the formatted text in the input and returns it. Does it?
func MakeFormatted(ctx mycocontext.Context, firstLine string) blocks.Formatted {
	hyphaName := ctx.HyphaName()
	return blocks.Formatted{
		HyphaName: hyphaName,
		Lines:     [][]blocks.Span{spansFromLine(ctx, hyphaName, firstLine)},
	}
}

func spansFromLine(ctx mycocontext.Context, hyphaName, line string) []blocks.Span {
	var (
		input      = bytes.NewBufferString(line)
		spans      = make([]blocks.Span, 0)
		tagState   = blocks.CleanStyleState()
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
			spans = append(spans, nextInlineLink(ctx, input, hyphaName, true))
		case (startsWith("https://") || startsWith("http://") || startsWith("gemini://") || startsWith("gopher://") || startsWith("ftp://")) && noTagsActive():
			spans = append(spans, nextInlineLink(ctx, input, hyphaName, false))
		default:
			spans = append(spans, nextInlineText(input))
		}
	}

	return spans
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
