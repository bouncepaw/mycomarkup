package lexer

import (
	"bytes"
	"reflect"
	"testing"
)

func paraTestHelper(t *testing.T, instr string, allowMultiline, terminateOnCloseBrace bool, expected []Token) {
	s := &State{
		b:        bytes.NewBufferString(instr),
		elements: make([]Token, 0),
	}
	lexParagraph(s, allowMultiline, terminateOnCloseBrace)
	if !reflect.DeepEqual(expected, s.elements) {
		t.Errorf("Failure! Wanted %v, got %v", expected, s.elements)
	}
}

func TestSimpleParagraph(t *testing.T) {
	paraTestHelper(
		t,
		"use it!",
		false, false,
		[]Token{
			Token{TokenSpanText, "use it!"},
		})
}

func TestParagraphWithItalic(t *testing.T) {
	paraTestHelper(
		t,
		"//italic text//",
		false, false,
		[]Token{
			Token{TokenSpanItalic, ""},
			Token{TokenSpanText, "italic text"},
			Token{TokenSpanItalic, ""},
		})
}

func TestParagraphWithMultipleStyles(t *testing.T) {
	paraTestHelper(
		t,
		"adventure: //visit Italy **and look at the Colosseo**//. Good?",
		true, true,
		[]Token{
			Token{TokenSpanText, "adventure: "},
			Token{TokenSpanItalic, ""},
			Token{TokenSpanText, "visit Italy "},
			Token{TokenSpanBold, ""},
			Token{TokenSpanText, "and look at the Colosseo"},
			Token{TokenSpanBold, ""},
			Token{TokenSpanItalic, ""},
			Token{TokenSpanText, ". Good?"},
		})

}

func TestParagraphWithLink(t *testing.T) {
	paraTestHelper(
		t,
		"see these resources: [[hypha|ὑφή]]",
		true, false,
		[]Token{
			Token{TokenSpanText, "see these resources: "},
			Token{TokenSpanLinkOpen, ""},
			Token{TokenLinkAddress, "hypha"},
			Token{TokenLinkDisplayOpen, ""},
			Token{TokenSpanText, "ὑφή"},
			Token{TokenLinkDisplayClose, ""},
			Token{TokenSpanLinkClose, ""},
		})

}
