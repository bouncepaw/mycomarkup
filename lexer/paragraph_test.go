package lexer

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSimpleParagraph(t *testing.T) {
	s := &State{
		b:        bytes.NewBufferString("use it!"),
		elements: make([]Token, 0),
	}
	lexParagraph(s, false, false)
	expected := []Token{
		Token{TokenSpanText, "use it!"},
	}
	if !reflect.DeepEqual(expected, s.elements) {
		t.Errorf("Failure! Wanted %v, got %v", expected, s.elements)
	}
}
