package lexer

import (
	"bytes"
	"strings"
	//	"github.com/bouncepaw/mycomarkup/blocks"
)

// For debug purposes, ig
func callbackNop(s *State, b *bytes.Buffer) {
}

func λcallbackConsumeHeading(level uint) func(*State) {
	return func(s *State) {
		heading := Token{
			startLine:   s.line,
			startColumn: s.column,
			kind:        TokenHeading,
		}
		// We know that there is no EOF and there is a space, so we drop error.
		value, _ := s.b.ReadString(' ')
		heading.value = value

		s.column += level + 1
		s.inHeading = true
		s.appendToken(heading)
	}
}

func callbackConsumeRocket(s *State) {
	rocket := Token{
		startLine:   s.line,
		startColumn: s.column,
		kind:        TokenRocketLink,
	}
	_, _ = s.b.ReadString('>')
	s.column += 2
	s.inRocket = true
	s.appendToken(rocket)
}

func callbackConsumeHorizontalLine(s *State) {
	hr := Token{
		startLine:   s.line,
		startColumn: s.column,
		kind:        TokenHorizontalLine,
	}
	value, _ := s.b.ReadString('\n')
	s.column += uint(len(value))
	hr.value = value
	s.appendToken(hr)
}

func callbackConsumeBlockquote(s *State) {
	s.appendToken(Token{
		startLine:   s.line,
		startColumn: s.column,
		kind:        TokenBraceOpen,
	})

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
	s.elements = append(s.elements, Lex(buf)...)

	// End quote
	s.appendToken(Token{
		startLine:   s.line,
		startColumn: s.column,
		kind:        TokenBraceClose,
	})
}

var table = []struct {
	prefix    string
	callback  func(*State)
	condition Condition
}{
	{"# ", λcallbackConsumeHeading(1),
		Condition{onNewLine: True, inGeneralText: True}},
	{"## ", λcallbackConsumeHeading(2),
		Condition{onNewLine: True, inGeneralText: True}},
	{"### ", λcallbackConsumeHeading(3),
		Condition{onNewLine: True, inGeneralText: True}},
	{"#### ", λcallbackConsumeHeading(4),
		Condition{onNewLine: True, inGeneralText: True}},
	{"##### ", λcallbackConsumeHeading(5),
		Condition{onNewLine: True, inGeneralText: True}},
	{"###### ", λcallbackConsumeHeading(6),
		Condition{onNewLine: True, inGeneralText: True}},
	{"=>", callbackConsumeRocket,
		Condition{onNewLine: True, inGeneralText: True}},
	{"----", callbackConsumeHorizontalLine,
		Condition{onNewLine: True, inGeneralText: True, okForHorizontalLine: True}},
	{">", callbackConsumeBlockquote,
		Condition{onNewLine: True, inGeneralText: True}},
}

func startsWithStr(b *bytes.Buffer, s string) bool {
	return strings.HasPrefix(b.String(), s)
}

func Lex(b bytes.Buffer) []Token {
	var (
		state = &State{}
		_     = state
	)
	// TODO:
	return state.elements
}
