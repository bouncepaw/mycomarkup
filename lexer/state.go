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
	inHeading              Ternary
	inSpan                 bool
	inInlineLink           bool

	stateImg *imgState
}

func (s *State) onNewLine() Ternary {
	if s.lastElement == nil {
		return True
	}
	switch s.lastElement.kind {
	case TokenNewLine, TokenHeadingClose, TokenRocketLinkClose:
		return True
	}
	return False
}

func (s *State) onImg() Ternary {
	return s.stateImg != nil
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
