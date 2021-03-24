package lexer

import (
	"bytes"
	"strings"
)

// Read all next spaces and tabs and forget about them
func eatLWS(s *State) {
	for {
		r, err := s.b.ReadByte()
		switch {
		case err != nil: // end of input
			break
		case r == ' ', r == '\t': // lws that we shal eat
			continue
		default: // non-space or a newline
			_ = s.b.UnreadByte()
			break
		}
	}
	return
}

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
	s.appendToken(Token{kind: TokenRocketLinkOpen})
	var (
		srcline, _  = s.b.ReadString('\n')
		line        = strings.TrimSpace(srcline[2:]) // without => and spaces
		hasBrackets = false
		escaping    = false
		inDisplay   = false
		buf         = strings.Builder{}
	)
	// Chop off the [[...]]
	if strings.HasPrefix(line, "[[") {
		hasBrackets = true
		line = line[2:]
		if strings.HasSuffix(line, "]]") {
			line = line[:len(line)-2]
		}
	}

	for _, r := range line {
		switch {
		case r == '\r':
		case inDisplay:
			buf.WriteRune(r)
			// The rest is !inDisplay && ...
		case escaping:
			escaping = false
			buf.WriteRune(r)
		case r == '\\':
			escaping = true
			// If found the separator:
		case hasBrackets && r == '|', !hasBrackets && (r == ' ' || r == '\t'):
			inDisplay = true
			if buf.Len() != 0 {
				s.appendToken(Token{kind: TokenLinkAddress, value: buf.String()})
			}
			buf.Reset()
		default:
			buf.WriteRune(r)
		}
	}

	if buf.Len() != 0 && inDisplay {
		s.appendToken(Token{kind: TokenLinkDisplay, value: buf.String()})
	} else if buf.Len() != 0 {
		s.appendToken(Token{kind: TokenLinkAddress, value: buf.String()})
	}
	s.appendToken(Token{kind: TokenRocketLinkClose})
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
