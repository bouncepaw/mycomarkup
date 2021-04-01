package lexer

import (
	"bytes"
)

type LexerState int

const (
	StateErr LexerState = iota
	StateNil

	StateParagraph
	StateEscape
	StateNowiki
	StateAutolink
	StateLinkAddress
	StateLinkDisplay
)

type ParagraphState struct {
	stackState *stateStack
	buf        *bytes.Buffer
}

func (ps *ParagraphState) hasOnTop(ls LexerState) bool {
	return ps.stackState.topElem == &ls
}

type paragraphLexEntry struct {
	prefices []string
	λ        func(s *State, ps *ParagraphState)
}

func closeTextSpan(s *State, ps *ParagraphState) {
	if ps.buf.Len() != 0 {
		s.appendToken(Token{TokenSpanText, ps.buf.String()})
		ps.buf.Reset()
	}
}

func λcloseEatAppend(n int, tk TokenKind) func(*State, *ParagraphState) {
	return func(s *State, ps *ParagraphState) {
		closeTextSpan(s, ps)
		if n > 0 {
			eatChar(s)
			if n > 1 {
				eatChar(s)
			}
		}
		// fmt.Printf("The rest: %s\n", s.b.String())
		s.appendToken(Token{tk, ""})
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
		{[]string{"\\"}, nil},
		{[]string{"[["}, func(s *State, ps *ParagraphState) {
			λcloseEatAppend(2, TokenSpanLinkOpen)(s, ps)
			ps.stackState.push(StateLinkAddress)
		}},
		{[]string{"//"}, λcloseEatAppend(2, TokenSpanItalic)},
		{[]string{"**"}, λcloseEatAppend(2, TokenSpanBold)},
		{[]string{"`"}, λcloseEatAppend(2, TokenSpanMonospace)},

		{[]string{"^^"}, λcloseEatAppend(2, TokenSpanSuper)},
		{[]string{",,"}, λcloseEatAppend(2, TokenSpanSub)},
		{[]string{"~~"}, λcloseEatAppend(2, TokenSpanStrike)},
		{[]string{"__"}, λcloseEatAppend(2, TokenSpanUnderline)},
		{[]string{"%%"}, nil},
		{[]string{"https://", "http://", "gemini://", "gopher://", "ftp://", "sftp://", "ssh://", "file://", "mailto:"},
			func(s *State, ps *ParagraphState) {
				λcloseEatAppend(0, TokenSpanLinkOpen)(s, ps)
				ps.stackState.push(StateAutolink)
			}},
	}
	paragraphInlineLinkAddressTable = []paragraphLexEntry{
		{[]string{"\\"}, nil},
		{[]string{"]]"}, func(s *State, ps *ParagraphState) {
			eatChar(s)
			eatChar(s)
			s.appendToken(Token{TokenLinkAddress, ps.buf.String()})
			s.appendToken(Token{TokenSpanLinkClose, ""})
			ps.buf.Reset()

			ps.stackState.pop()
		}},
		{[]string{"|"}, func(s *State, ps *ParagraphState) {
			eatChar(s)
			s.appendToken(Token{TokenLinkAddress, ps.buf.String()})
			s.appendToken(Token{TokenLinkDisplayOpen, ""})
			ps.buf.Reset()

			ps.stackState.pop()
			ps.stackState.push(StateLinkDisplay)
		}},
	}
	paragraphInlineLinkDisplayTable = []paragraphLexEntry{
		{[]string{"\\"}, nil},
		{[]string{"]]"}, func(s *State, ps *ParagraphState) {
			λcloseEatAppend(2, TokenLinkDisplayClose)(s, ps)
			s.appendToken(Token{TokenSpanLinkClose, ""})
			ps.stackState.pop()
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
			func(s *State, ps *ParagraphState) {
				ps.stackState.pop()
				s.appendToken(Token{TokenLinkAddress, s.b.String()})
				s.b.Reset()
			}},
	}
}

func looksLikeParagraph(b *bytes.Buffer) bool {
	// TODO: implement
	return true
}

func lexParagraph(s *State, allowMultiline, terminateOnCloseBrace bool) []Token {
	var (
		paragraphState = ParagraphState{
			stackState: newStateStack(),
			buf:        &bytes.Buffer{},
		}
		ch  byte
		err error
	)
	paragraphState.stackState.push(StateParagraph)
	for {
		switch {
		case startsWithStr(s.b, "\n"):
			eatChar(s)
			if looksLikeParagraph(s.b) && allowMultiline {
				paragraphState.buf.WriteByte('\n')
			} else {
				break
			}
		case startsWithStr(s.b, "\\") && !paragraphState.hasOnTop(StateEscape):
			paragraphState.stackState.push(StateEscape)
			eatChar(s)
			continue
		case startsWithStr(s.b, "}") && terminateOnCloseBrace:
			break
		}
		switch *(paragraphState.stackState.topElem) {
		case StateEscape:
			ch, err := s.b.ReadByte()
			if err != nil {
				break
			}
			paragraphState.stackState.pop()
			paragraphState.buf.WriteByte(ch)
		case StateNowiki:
			if startsWithStr(s.b, "%%") {
				eatChar(s)
				eatChar(s)
				paragraphState.stackState.pop()
			}
			continue
		case StateParagraph:
			for _, rule := range paragraphParagraphTable {
				for _, prefix := range rule.prefices {
					if startsWithStr(s.b, prefix) {
						rule.λ(s, &paragraphState)
						goto next // this is required because of the nested loop
					}
				}
			}
		case StateAutolink:
			for _, rule := range paragraphAutolinkTable {
				for _, prefix := range rule.prefices {
					if startsWithStr(s.b, prefix) {
						rule.λ(s, &paragraphState)
						continue
					}
				}
			}
		case StateLinkAddress:
			for _, rule := range paragraphInlineLinkAddressTable {
				for _, prefix := range rule.prefices {
					if startsWithStr(s.b, prefix) {
						rule.λ(s, &paragraphState)
						continue
					}
				}
			}
		case StateLinkDisplay:
			for _, rule := range paragraphInlineLinkDisplayTable {
				for _, prefix := range rule.prefices {
					if startsWithStr(s.b, prefix) {
						rule.λ(s, &paragraphState)
						continue
					}
				}
			}
		}
		ch, err = s.b.ReadByte()
		if err != nil {
			break
		}
		paragraphState.buf.WriteByte(ch)
	next:
	}
	closeTextSpan(s, &paragraphState)
	return nil
}
