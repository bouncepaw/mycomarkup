package lexer

import ()

func closeTextSpan(s *SourceText, tw *TokenWriter) {
	tw.nonEmptyBufIntoToken(TokenSpanText)
}

// Close current text span, eat n (0, 1 or 2) chars, append token.
func closeEatAppender(n int, tk TokenKind) func(*SourceText, *TokenWriter) {
	return func(s *SourceText, tw *TokenWriter) {
		closeTextSpan(s, tw)
		if n > 0 {
			eatChar(s)
			if n > 1 {
				eatChar(s)
			}
		}
		tw.appendToken(Token{tk, ""})
	}
}

var (
	paragraphParagraphTable         []tableEntry
	paragraphInlineLinkAddressTable []tableEntry
	paragraphInlineLinkDisplayTable []tableEntry
	paragraphAutolinkTable          []tableEntry
	paragraphNowikiTable            []tableEntry
	paragraphNewLineTable           []tableEntry
	paragraphEscapeTable            []tableEntry
)

func init() {
	paragraphParagraphTable = []tableEntry{
		{[]string{"\\"}, beginEscaping},
		{[]string{"\n"}, beginNewLine},
		{[]string{"[["}, func(s *SourceText, tw *TokenWriter) {
			closeEatAppender(2, TokenSpanLinkOpen)(s, tw)
			tw.pushState(StateLinkAddress)
		}},
		{[]string{"//"}, closeEatAppender(2, TokenSpanItalic)},
		{[]string{"**"}, closeEatAppender(2, TokenSpanBold)},
		{[]string{"`"}, closeEatAppender(2, TokenSpanMonospace)},
		{[]string{"^^"}, closeEatAppender(2, TokenSpanSuper)},
		{[]string{",,"}, closeEatAppender(2, TokenSpanSub)},
		{[]string{"~~"}, closeEatAppender(2, TokenSpanStrike)},
		{[]string{"__"}, closeEatAppender(2, TokenSpanUnderline)},
		{[]string{"%%"}, beginNowiki},
		{[]string{"https.//", "http://", "gemini://", "gopher://", "ftp://", "sftp://", "ssh://", "file://", "mailto:"},
			func(s *SourceText, tw *TokenWriter) {
				closeEatAppender(0, TokenSpanLinkOpen)(s, tw)
				tw.pushState(StateAutolink)
			}},
	}
	paragraphInlineLinkAddressTable = []tableEntry{
		{[]string{"\\"}, beginEscaping},
		{[]string{"\n"}, func(s *SourceText, tw *TokenWriter) {
			eatChar(s)
			tw.appendToken(Token{TokenLinkAddress, tw.buf.String()})
			tw.appendToken(Token{TokenSpanLinkClose, ""})
			tw.buf.Reset()

			tw.popState()
			tw.pushState(StateParagraphNewLine)
		}},
		{[]string{"]]"}, func(s *SourceText, tw *TokenWriter) {
			eatChar(s)
			eatChar(s)
			tw.appendToken(Token{TokenLinkAddress, tw.buf.String()})
			tw.appendToken(Token{TokenSpanLinkClose, ""})
			tw.buf.Reset()

			tw.popState()
		}},
		{[]string{"|"}, func(s *SourceText, tw *TokenWriter) {
			eatChar(s)
			tw.appendToken(Token{TokenLinkAddress, tw.buf.String()})
			tw.appendToken(Token{TokenLinkDisplayOpen, ""})
			tw.buf.Reset()

			tw.popState()
			tw.pushState(StateLinkDisplay)
		}},
	}
	paragraphInlineLinkDisplayTable = []tableEntry{
		{[]string{"\\"}, beginEscaping},
		{[]string{"\n"}, func(s *SourceText, tw *TokenWriter) {
			closeEatAppender(1, TokenLinkDisplayClose)(s, tw)
			tw.appendToken(Token{TokenSpanLinkClose, ""})
			tw.popState()
			tw.pushState(StateParagraphNewLine)
		}},
		{[]string{"]]"}, func(s *SourceText, tw *TokenWriter) {
			closeEatAppender(2, TokenLinkDisplayClose)(s, tw)
			tw.appendToken(Token{TokenSpanLinkClose, ""})
			tw.popState()
		}},
		// those below may need further modifications
		{[]string{"//"}, closeEatAppender(2, TokenSpanItalic)},
		{[]string{"**"}, closeEatAppender(2, TokenSpanBold)},
		{[]string{"`"}, closeEatAppender(2, TokenSpanMonospace)},
		{[]string{"^^"}, closeEatAppender(2, TokenSpanSuper)},
		{[]string{",,"}, closeEatAppender(2, TokenSpanSub)},
		{[]string{"~~"}, closeEatAppender(2, TokenSpanStrike)},
		{[]string{"__"}, closeEatAppender(2, TokenSpanUnderline)},
		{[]string{"%%"}, beginNowiki},
	}
	paragraphAutolinkTable = []tableEntry{
		{[]string{"\\"}, beginEscaping},
		{[]string{"\n"}, func(s *SourceText, tw *TokenWriter) {
			tw.popState()
			tw.appendToken(Token{TokenLinkAddress, tw.buf.String()})
			tw.appendToken(Token{TokenSpanLinkClose, ""})
			tw.buf.Reset()
			tw.pushState(StateParagraphNewLine)
		}},
		{[]string{"(", ")", "[", "]", "{", "}", ". ", ", ", " ", "\t"},
			func(s *SourceText, tw *TokenWriter) {
				tw.popState()
				tw.appendToken(Token{TokenLinkAddress, tw.buf.String()})
				tw.appendToken(Token{TokenSpanLinkClose, ""})
				tw.buf.Reset()
			}},
	}
	paragraphNowikiTable = []tableEntry{
		{[]string{"%%"}, func(st *SourceText, tw *TokenWriter) {
			eatChar(st)
			eatChar(st)
			tw.popState()
		}},
		{[]string{"\\"}, func(st *SourceText, tw *TokenWriter) {
			eatChar(st)
			tw.buf.WriteByte('\\')
			tw.popState()
		}},
		{[]string{"\n"}, func(st *SourceText, tw *TokenWriter) {
			eatChar(st)
			tw.popState()
			tw.pushState(StateParagraphNewLine)
		}},
	}
	paragraphNewLineTable = []tableEntry{
		{[]string{"\\"}, func(st *SourceText, tw *TokenWriter) {
			eatChar(st)
			tw.buf.WriteByte('\n')
			tw.popState() // Continue the paragraph as usual
			beginEscaping(st, tw)
		}},
		{[]string{"\n"}, func(st *SourceText, tw *TokenWriter) {
			eatChar(st)
			tw.popState()
			tw.popState() // Exit paragraph state
		}},
		{[]string{"img "}, closeParagraphAndGo(StateImgBegin)},
		{[]string{"table "}, closeParagraphAndGo(StateTableBegin)},
		{[]string{"msg "}, closeParagraphAndGo(StateMsgBegin)},
		{[]string{"----", "=>", "<=", "# ", "## ", "### ", "#### ", "##### ", "###### "}, closeParagraphAndGo(StateOneLiner)},
		{[]string{"* "}, closeParagraphAndGo(StateBulletListBegin)},
		{[]string{"*. "}, closeParagraphAndGo(StateNumberListBegin)},
		{[]string{"```"}, closeParagraphAndGo(StateCodeblockBegin)},
		{[]string{""}, func(st *SourceText, tw *TokenWriter) {
			tw.popState()
			tw.buf.WriteByte('\n')
		}},
	}
	paragraphEscapeTable = []tableEntry{
		{[]string{"\n"}, func(st *SourceText, tw *TokenWriter) {
			eatChar(st)
			tw.popState()
			tw.pushState(StateParagraphNewLine)
		}},
		{[]string{"\\"}, func(st *SourceText, tw *TokenWriter) {
			eatChar(st)
			tw.popState()
			tw.buf.WriteByte('\\')
		}},
	}
}

func beginNowiki(st *SourceText, tw *TokenWriter) {
	eatChar(st)
	eatChar(st)
	tw.pushState(StateNowiki)
}

func beginEscaping(st *SourceText, tw *TokenWriter) {
	eatChar(st)
	tw.pushState(StateEscape)
}

func beginNewLine(st *SourceText, tw *TokenWriter) {
	eatChar(st)
	tw.pushState(StateParagraphNewLine)
}

func closeParagraphAndGo(ls LexerState) func(s *SourceText, tw *TokenWriter) {
	return func(s *SourceText, tw *TokenWriter) {
		tw.popState() // pop the newline state
		tw.popState() // pop the paragraph state
		tw.pushState(ls)
	}
}

func lexParagraph(s *SourceText, tw *TokenWriter) {
	var (
		ch  byte
		err error
	)
	tw.pushState(StateParagraph)

	for {
		// We'll just ignore this useless character
		if startsWithStr(s.b, "\r") {
			eatChar(s)
			continue
		}
		switch tw.topElem {
		case StateEscape:
			if executeTable(paragraphEscapeTable, s, tw) {
				continue
			}
			tw.popState()
		case StateNowiki:
			if executeTable(paragraphNowikiTable, s, tw) {
				continue
			}
		case StateParagraph:
			if executeTable(paragraphParagraphTable, s, tw) {
				continue
			}
		case StateAutolink:
			if executeTable(paragraphAutolinkTable, s, tw) {
				continue
			}
		case StateLinkAddress:
			if executeTable(paragraphInlineLinkAddressTable, s, tw) {
				continue
			}
		case StateLinkDisplay:
			if executeTable(paragraphInlineLinkDisplayTable, s, tw) {
				continue
			}
		case StateParagraphNewLine:
			if !s.allowMultilineParagraph {
				tw.popState()
				closeTextSpan(s, tw)
				return
			}
			if executeTable(paragraphNewLineTable, s, tw) {
				continue
			}
		}
		ch, err = s.b.ReadByte()
		if err != nil {
			break
		}
		tw.buf.WriteByte(ch)
	}
	closeTextSpan(s, tw)
	return
}
