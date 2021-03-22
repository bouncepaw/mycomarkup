package lexer

import (
	"bytes"
)

type State struct {
	// General:

	b           *bytes.Buffer
	line        uint
	column      uint
	elements    []Token
	lastElement *Token

	// Stuff:

	gottaGoFurtherNextTime bool
	inHeading              bool
	inRocket               bool
}

func (s *State) onNewLine() Ternary {
	if s.lastElement.kind == TokenNewLine || s.lastElement == nil {
		return True
	}
	return False
}

// TODO:
// Some elements are available in some places only
func (s *State) inGeneralText() Ternary {
	return True
}

// When the line is - only
func (s *State) okForHorizontalLine() Ternary {
	for _, r := range s.b.Bytes() {
		switch r {
		case '-', '\r':
			continue
		case '\n':
			return True
		default:
			break
		}
	}
	return False
}

func (s *State) appendToken(token Token) {
	s.lastElement = &token
	s.elements = append(s.elements, token)
}
