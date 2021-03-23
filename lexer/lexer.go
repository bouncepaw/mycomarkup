package lexer

import (
	"bytes"
	"strings"
)

// For debug purposes, ig
func callbackNop(s *State, b *bytes.Buffer) {
}

func λcallbackHeading(level uint) func(*State) {
	return func(s *State) {
		heading := Token{kind: TokenHeadingOpen}
		// We know that there is no EOF and there is a space, so we drop error.
		value, _ := s.b.ReadString(' ')
		heading.value = value

		s.inHeading = True
		s.appendToken(heading)
	}
}

func callbackHeadingNewLine(s *State) {
	_, _ = s.b.ReadByte()
	s.inHeading = False
	s.appendToken(Token{kind: TokenHeadingClose})
}

func callbackRocket(s *State) {
	rocket := Token{kind: TokenRocketLink}
	_, _ = s.b.ReadString('>')
	s.column += 2
	s.inRocket = true
	s.appendToken(rocket)
}

func callbackHorizontalLine(s *State) {
	hr := Token{kind: TokenHorizontalLine}
	value, _ := s.b.ReadString('\n')
	s.column += uint(len(value))
	hr.value = value
	s.appendToken(hr)
}

func callbackBlockquote(s *State) {
	s.appendToken(Token{kind: TokenBraceOpen})

	var (
		buf = bytes.Buffer{}
		nl  = false // prev char was newline
		bq  = false // prev char was > or a space after it
	)
	// Read all next consecutive lines starting with '>', strip them of the '>'s and the leading whitespace.
	for {
		r, err := s.b.ReadByte()
		if err != nil {
			break
		}
		switch r {
		case '\r':
		case '\n':
			nl = true
			bq = false
			buf.WriteByte('\n')
		case ' ', '\t':
			nl = false
			if !bq {
				buf.WriteByte(r)
			}
		default:
			nl = false
			bq = false
			if nl {
				if r == '>' {
					bq = true
					continue
				}
				_ = buf.UnreadByte()
				break
			}
			buf.WriteByte(r)
		}
	}

	// Zoom in what we've collected...
	s.elements = append(s.elements, Lex(&buf)...)

	// End quote
	s.appendToken(Token{kind: TokenBraceClose})
}

type tableEntry struct {
	prefix    string
	callback  func(*State)
	condition Condition
}

var table []tableEntry

func init() {
	table = []tableEntry{
		{"# ", λcallbackHeading(1),
			Condition{onNewLine: True, inGeneralText: True}},
		{"## ", λcallbackHeading(2),
			Condition{onNewLine: True, inGeneralText: True}},
		{"### ", λcallbackHeading(3),
			Condition{onNewLine: True, inGeneralText: True}},
		{"#### ", λcallbackHeading(4),
			Condition{onNewLine: True, inGeneralText: True}},
		{"##### ", λcallbackHeading(5),
			Condition{onNewLine: True, inGeneralText: True}},
		{"###### ", λcallbackHeading(6),
			Condition{onNewLine: True, inGeneralText: True}},
		{"\n", callbackHeadingNewLine,
			Condition{inHeading: True}},
		{"=>", callbackRocket,
			Condition{onNewLine: True, inGeneralText: True}},
		{"----", callbackHorizontalLine,
			Condition{onNewLine: True, inGeneralText: True, okForHorizontalLine: True}},
		{">", callbackBlockquote,
			Condition{onNewLine: True, inGeneralText: True}},
	}
}

func startsWithStr(b *bytes.Buffer, s string) bool {
	return strings.HasPrefix(b.String(), s)
}

func Lex(b *bytes.Buffer) []Token {
	var (
		state = &State{
			b:           b,
			line:        0,
			column:      0,
			elements:    make([]Token, 0),
			lastElement: nil,

			gottaGoFurtherNextTime: false,
			inHeading:              False,
			inRocket:               false,
		}
		textbuf bytes.Buffer
		r       byte
		err     error
	)
	for {
		// Rules are rules
		for _, rule := range table {
			if startsWithStr(state.b, rule.prefix) && rule.condition.fullfilledBy(state).isTrue() {
				// temporary block:
				if textbuf.Len() > 0 {
					state.appendToken(
						Token{kind: TokenSpanText, value: textbuf.String()})
				}
				textbuf.Reset()
				rule.callback(state)
				goto next // I'm sorry
			}
		}
		r, err = state.b.ReadByte()
		if err != nil {
			if textbuf.Len() > 0 {
				state.appendToken(
					Token{kind: TokenSpanText, value: textbuf.String()})
			}
			break
		}
		textbuf.WriteByte(r)
	next:
	}
	return state.elements
}
