package lexer

import (
	"bytes"
	"reflect"
	"testing"
)

func paraTestHelper(t *testing.T, instr string, allowMultilineParagraph, terminateOnCloseBrace bool, expected []Token) {
	st := SourceText{
		b:                       bytes.NewBufferString(instr),
		allowMultilineParagraph: allowMultilineParagraph,
		terminateOnCloseBrace:   terminateOnCloseBrace,
	}
	result, _ := LexParagraph(st)
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Failure! See the lexeme printouts below!")
		t.Logf("Wanted this:\n")
		for i, e := range expected {
			t.Logf("%d	%s\n", i, e.String())
		}
		t.Logf("Got this instead:\n")
		for i, e := range result {
			t.Logf("%d	%s\n", i, e.String())
		}
	}
}

func TestEnsureNewLineAndEscapeAreHandled(t *testing.T) {
	tables := [][]tableEntry{paragraphParagraphTable, paragraphInlineLinkAddressTable, paragraphInlineLinkDisplayTable, paragraphAutolinkTable, paragraphNowikiTable, paragraphNewLineTable, paragraphEscapeTable}

	for i, table := range tables {
		var (
			escapingHandled bool
			newLineHandled  bool
		)
		for _, entry := range table {
			for _, prefix := range entry.prefices {
				if prefix == "\\" {
					escapingHandled = true
				} else if prefix == "\n" {
					newLineHandled = true
				}
			}
		}
		if !(escapingHandled && newLineHandled) {
			t.Errorf("Escaping and newlines not handled in table number %d", i)
		}
	}
}

func TestSimpleParagraph(t *testing.T) {
	paraTestHelper(
		t,
		"use it!",
		false, false,
		[]Token{
			{TokenSpanText, "use it!"},
		})
}

func TestParagraphWithItalic(t *testing.T) {
	paraTestHelper(
		t,
		"//italic text//",
		false, false,
		[]Token{
			{TokenSpanItalic, ""},
			{TokenSpanText, "italic text"},
			{TokenSpanItalic, ""},
		})
}

func TestParagraphWithMultipleStyles(t *testing.T) {
	paraTestHelper(
		t,
		"adventure: //visit Italy **and look at the Colosseo**//. Good?",
		true, true,
		[]Token{
			{TokenSpanText, "adventure: "},
			{TokenSpanItalic, ""},
			{TokenSpanText, "visit Italy "},
			{TokenSpanBold, ""},
			{TokenSpanText, "and look at the Colosseo"},
			{TokenSpanBold, ""},
			{TokenSpanItalic, ""},
			{TokenSpanText, ". Good?"},
		})

}

func TestParagraphWithLink1(t *testing.T) {
	paraTestHelper(
		t,
		"see these resources: [[hypha|ὑφή]]",
		true, false,
		[]Token{
			{TokenSpanText, "see these resources: "},
			{TokenSpanLinkOpen, ""},
			{TokenLinkAddress, "hypha"},
			{TokenLinkDisplayOpen, ""},
			{TokenSpanText, "ὑφή"},
			{TokenLinkDisplayClose, ""},
			{TokenSpanLinkClose, ""},
		})
}

func TestParagraphWithLink2(t *testing.T) {
	paraTestHelper(
		t,
		"[[]]",
		false, false,
		[]Token{
			{TokenSpanLinkOpen, ""},
			{TokenLinkAddress, ""},
			{TokenSpanLinkClose, ""},
		})
}

func TestParagraphWithLink3(t *testing.T) {
	paraTestHelper(
		t,
		"[[|]]",
		false, false,
		[]Token{
			{TokenSpanLinkOpen, ""},
			{TokenLinkAddress, ""},
			{TokenLinkDisplayOpen, ""},
			{TokenLinkDisplayClose, ""},
			{TokenSpanLinkClose, ""},
		})
}

func TestParagraphWithLink4(t *testing.T) {
	paraTestHelper(
		t,
		"[[underground |]]",
		false, false,
		[]Token{
			{TokenSpanLinkOpen, ""},
			{TokenLinkAddress, "underground "},
			{TokenLinkDisplayOpen, ""},
			{TokenLinkDisplayClose, ""},
			{TokenSpanLinkClose, ""},
		})
}

func TestParagraphWithAutoLink1(t *testing.T) {
	paraTestHelper(
		t,
		"ftp://example.org gemini://lesarbr.es two nice links\nfor two nice people here",
		true, true,
		[]Token{
			{TokenSpanLinkOpen, ""},
			{TokenLinkAddress, "ftp://example.org"},
			{TokenSpanLinkClose, ""},
			{TokenSpanText, " "},
			{TokenSpanLinkOpen, ""},
			{TokenLinkAddress, "gemini://lesarbr.es"},
			{TokenSpanLinkClose, ""},
			{TokenSpanText, " two nice links\nfor two nice people here"},
		})
}

func TestParagraphWithAutoLink2(t *testing.T) {
	paraTestHelper(
		t,
		"Do not hesitate to contact me (mailto:nikołaj.przewalski@example.org). I will not bite you!",
		false, true,
		[]Token{
			{TokenSpanText, "Do not hesitate to contact me ("},
			{TokenSpanLinkOpen, ""},
			{TokenLinkAddress, "mailto:nikołaj.przewalski@example.org"},
			{TokenSpanLinkClose, ""},
			{TokenSpanText, "). I will not bite you!"},
		})
}

func TestParagraphNewLine1(t *testing.T) {
	paraTestHelper(
		t,
		"line that is consumed\nline that is not consumed yet",
		false, false,
		[]Token{
			{TokenSpanText, "line that is consumed"},
		})
}

func TestParagraphNewLine2(t *testing.T) {
	paraTestHelper(
		t,
		"line that is consumed\nline that is consumed too",
		true, false,
		[]Token{
			{TokenSpanText, "line that is consumed\nline that is consumed too"},
		})
}

func TestParagraphWithEscaping1(t *testing.T) {
	paraTestHelper(
		t,
		"\\Escape first char",
		false, false,
		[]Token{
			{TokenSpanText, "Escape first char"},
		})
}

func TestParagraphWithEscaping2(t *testing.T) {
	// multiline on
	paraTestHelper(
		t,
		"Escape last char on line\\\nHmm",
		true, false,
		[]Token{
			{TokenSpanText, "Escape last char on line\nHmm"},
		})
}

func TestParagraphWithEscaping3(t *testing.T) {
	// multiline off
	paraTestHelper(
		t,
		"Escape last char on line\\\nHmm",
		false, false,
		[]Token{
			{TokenSpanText, "Escape last char on line"},
		})
}

func TestParagraphWithEscaping4(t *testing.T) {
	paraTestHelper(
		t,
		"\\",
		false, false,
		[]Token{})
}

func TestParagraphWithNowiki1(t *testing.T) {
	paraTestHelper(
		t,
		"Stay %%//safe//%%",
		false, false,
		[]Token{
			{TokenSpanText, "Stay //safe//"},
		})
}

func TestParagraphWithNowiki2(t *testing.T) {
	paraTestHelper(
		t,
		"This is %%more\n//dangerous//",
		true, false,
		[]Token{
			{TokenSpanText, "This is more\n"},
			{TokenSpanItalic, ""},
			{TokenSpanText, "dangerous"},
			{TokenSpanItalic, ""},
		})
}

func TestParagraphWithNowiki3(t *testing.T) {
	paraTestHelper(
		t,
		"the coool [[link| %%!!link!!%%]] has arrived",
		false, false,
		[]Token{
			{TokenSpanText, "the coool "},
			{TokenSpanLinkOpen, ""},
			{TokenLinkAddress, "link"},
			{TokenLinkDisplayOpen, ""},
			{TokenSpanText, " !!link!!"},
			{TokenLinkDisplayClose, ""},
			{TokenSpanLinkClose, ""},
			{TokenSpanText, " has arrived"},
		})
}
