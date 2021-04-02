package lexer

import (
	"bytes"
)

type State struct {
	// General:

	b *bytes.Buffer
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
