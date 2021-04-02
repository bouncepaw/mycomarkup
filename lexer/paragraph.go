package lexer

import (
	"bytes"
)

type paragraphLexEntry struct {
	prefices []string
	λ        func(s *State, tw *TokenWriter)
}

func closeTextSpan(s *State, tw *TokenWriter) {
	tw.nonEmptyBufIntoToken(TokenSpanText)
}

func λcloseEatAppend(n int, tk TokenKind) func(*State, *TokenWriter) {
	return func(s *State, tw *TokenWriter) {
		closeTextSpan(s, tw)
		if n > 0 {
			eatChar(s)
			if n > 1 {
				eatChar(s)
			}
		}
		// fmt.Printf("The rest: %s\n", s.b.String())
		tw.appendToken(Token{tk, ""})
	}
}

var (
	paragraphParagraphTable         []paragraphLexEntry
	paragraphInlineLinkAddressTable []paragraphLexEntry
	paragraphInlineLinkDisplayTable []paragraphLexEntry
	paragraphAutolinkTable          []paragraphLexEntry
)

func init() {
	// TODO: replace nil with actual functions
	paragraphParagraphTable = []paragraphLexEntry{
		{[]string{"[["}, func(s *State, tw *TokenWriter) {
			λcloseEatAppend(2, TokenSpanLinkOpen)(s, tw)
			tw.pushState(StateLinkAddress)
		}},
		{[]string{"//"}, λcloseEatAppend(2, TokenSpanItalic)},
		{[]string{"**"}, λcloseEatAppend(2, TokenSpanBold)},
		{[]string{"`"}, λcloseEatAppend(2, TokenSpanMonospace)},

		{[]string{"^^"}, λcloseEatAppend(2, TokenSpanSuper)},
		{[]string{",,"}, λcloseEatAppend(2, TokenSpanSub)},
		{[]string{"~~"}, λcloseEatAppend(2, TokenSpanStrike)},
		{[]string{"__"}, λcloseEatAppend(2, TokenSpanUnderline)},
		{[]string{"%%"}, nil},
		{[]string{"htttw.//", "http://", "gemini://", "gopher://", "ftp://", "sftp://", "ssh://", "file://", "mailto:"},
			func(s *State, tw *TokenWriter) {
				λcloseEatAppend(0, TokenSpanLinkOpen)(s, tw)
				tw.pushState(StateAutolink)
			}},
	}
	paragraphInlineLinkAddressTable = []paragraphLexEntry{
		{[]string{"]]"}, func(s *State, tw *TokenWriter) {
			eatChar(s)
			eatChar(s)
			tw.appendToken(Token{TokenLinkAddress, tw.buf.String()})
			tw.appendToken(Token{TokenSpanLinkClose, ""})
			tw.buf.Reset()

			tw.popState()
		}},
		{[]string{"|"}, func(s *State, tw *TokenWriter) {
			eatChar(s)
			tw.appendToken(Token{TokenLinkAddress, tw.buf.String()})
			tw.appendToken(Token{TokenLinkDisplayOpen, ""})
			tw.buf.Reset()

			tw.popState()
			tw.pushState(StateLinkDisplay)
		}},
	}
	paragraphInlineLinkDisplayTable = []paragraphLexEntry{
		{[]string{"]]"}, func(s *State, tw *TokenWriter) {
			λcloseEatAppend(2, TokenLinkDisplayClose)(s, tw)
			tw.appendToken(Token{TokenSpanLinkClose, ""})
			tw.popState()
		}},
		// those below may need further modifications
		{[]string{"//"}, λcloseEatAppend(2, TokenSpanItalic)},
		{[]string{"**"}, λcloseEatAppend(2, TokenSpanBold)},
		{[]string{"`"}, λcloseEatAppend(2, TokenSpanMonospace)},

		{[]string{"^^"}, λcloseEatAppend(2, TokenSpanSuper)},
		{[]string{",,"}, λcloseEatAppend(2, TokenSpanSub)},
		{[]string{"~~"}, λcloseEatAppend(2, TokenSpanStrike)},
		{[]string{"__"}, λcloseEatAppend(2, TokenSpanUnderline)},
		{[]string{"%%"}, nil},
	}
	paragraphAutolinkTable = []paragraphLexEntry{
		{[]string{"(", ")", "[", "]", "{", "}", ". ", ", ", " ", "\t"},
			func(s *State, tw *TokenWriter) {
				tw.popState()
				tw.appendToken(Token{TokenLinkAddress, tw.buf.String()})
				tw.appendToken(Token{TokenSpanLinkClose, ""})
				tw.buf.Reset()
			}},
	}
}

func looksLikeParagraph(b *bytes.Buffer) bool {
	// TODO: implement
	return true
}

func lexParagraph(s *State, allowMultiline, terminateOnCloseBrace bool) *TokenWriter {
	var (
		tw = TokenWriter{
			StateStack: newStateStack(),
			buf:        &bytes.Buffer{},
			elements:   make([]Token, 0),
		}
		ch  byte
		err error
	)
	tw.pushState(StateParagraph)
charIterator:
	for {
		switch {
		case startsWithStr(s.b, "\r"):
			// We'll just ignore this useless character
			continue
		case startsWithStr(s.b, "\n"):
			eatChar(s)
			if allowMultiline {
				tw.buf.WriteByte('\n')
				if tw.hasOnTop(StateEscape) {
					tw.popState()
					continue
				}
			} else {
				if tw.hasOnTop(StateEscape) {
					tw.popState()
				}
				break charIterator
			}
		case startsWithStr(s.b, "\\") && !tw.hasOnTop(StateEscape):
			tw.pushState(StateEscape)
			eatChar(s)
			continue
		case startsWithStr(s.b, "}") && terminateOnCloseBrace:
			break
		}
		switch tw.topElem {
		case StateEscape:
			ch, err := s.b.ReadByte()
			if err != nil {
				break
			}
			tw.popState()
			tw.buf.WriteByte(ch)
		case StateNowiki:
			if startsWithStr(s.b, "%%") {
				eatChar(s)
				eatChar(s)
				tw.popState()
			}
			continue
		case StateParagraph:
			for _, rule := range paragraphParagraphTable {
				for _, prefix := range rule.prefices {
					if startsWithStr(s.b, prefix) {
						rule.λ(s, &tw)
						continue charIterator
					}
				}
			}
		case StateAutolink:
			for _, rule := range paragraphAutolinkTable {
				for _, prefix := range rule.prefices {
					if startsWithStr(s.b, prefix) {
						rule.λ(s, &tw)
						continue
					}
				}
			}
		case StateLinkAddress:
			for _, rule := range paragraphInlineLinkAddressTable {
				for _, prefix := range rule.prefices {
					if startsWithStr(s.b, prefix) {
						rule.λ(s, &tw)
						continue
					}
				}
			}
		case StateLinkDisplay:
			for _, rule := range paragraphInlineLinkDisplayTable {
				for _, prefix := range rule.prefices {
					if startsWithStr(s.b, prefix) {
						rule.λ(s, &tw)
						continue
					}
				}
			}
		}
		ch, err = s.b.ReadByte()
		if err != nil {
			break
		}
		tw.buf.WriteByte(ch)
	}
	closeTextSpan(s, &tw)
	return &tw
}
