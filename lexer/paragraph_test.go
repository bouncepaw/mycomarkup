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
		t.Errorf("Failure! See the lexeme printouts below!")
		t.Logf("Wanted this:\n")
		for i, e := range expected {
			t.Logf("%d	%s\n", i, e.String())
		}
		t.Logf("Got this instead:\n")
		for i, e := range s.elements {
			t.Logf("%d	%s\n", i, e.String())
		}
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

func TestParagraphWithAutoLink1(t *testing.T) {
	paraTestHelper(
		t,
		"ftp://example.org gemini://lesarbr.es two nice links\nfor two nice people here",
		true, true,
		[]Token{
			Token{TokenSpanLinkOpen, ""},
			Token{TokenLinkAddress, "ftp://example.org"},
			Token{TokenSpanLinkClose, ""},
			Token{TokenSpanText, " "},
			Token{TokenSpanLinkOpen, ""},
			Token{TokenLinkAddress, "gemini://lesarbr.es"},
			Token{TokenSpanLinkClose, ""},
			Token{TokenSpanText, " two nice links\nfor two nice people here"},
		})
}

func TestParagraphWithAutoLink2(t *testing.T) {
	paraTestHelper(
		t,
		"Do not hesitate to contact me (mailto:nikołaj.przewalski@example.org). I will not bite you!",
		false, true,
		[]Token{
			Token{TokenSpanText, "Do not hesitate to contact me ("},
			Token{TokenSpanLinkOpen, ""},
			Token{TokenLinkAddress, "mailto:nikołaj.przewalski@example.org"},
			Token{TokenSpanLinkClose, ""},
			Token{TokenSpanText, "). I will not bite you!"},
		})
}

func TestParagraphNewLine1(t *testing.T) {
	paraTestHelper(
		t,
		"line that is consumed\nline that is not consumed yet",
		false, false,
		[]Token{
			Token{TokenSpanText, "line that is consumed"},
		})
}

func TestParagraphNewLine2(t *testing.T) {
	paraTestHelper(
		t,
		"line that is consumed\nline that is consumed too",
		true, false,
		[]Token{
			Token{TokenSpanText, "line that is consumed\nline that is consumed too"},
		})
}

func TestParagraphWithEscaping1(t *testing.T) {
	paraTestHelper(
		t,
		"\\Escape first char",
		false, false,
		[]Token{
			Token{TokenSpanText, "Escape first char"},
		})
}

func TestParagraphWithEscaping2(t *testing.T) {
	// multiline on
	paraTestHelper(
		t,
		"Escape last char on line\\\nHmm",
		true, false,
		[]Token{
			Token{TokenSpanText, "Escape last char on line\nHmm"},
		})
}

func TestParagraphWithEscaping3(t *testing.T) {
	// multiline off
	paraTestHelper(
		t,
		"Escape last char on line\\\nHmm",
		false, false,
		[]Token{
			Token{TokenSpanText, "Escape last char on line"},
		})
}

func TestParagraphWithEscaping4(t *testing.T) {
	paraTestHelper(
		t,
		"\\",
		false, false,
		[]Token{})
}
