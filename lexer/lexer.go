package lexer

import (
	"bytes"
	"strings"
)

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
