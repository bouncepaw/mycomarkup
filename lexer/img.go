package lexer

var (
	imgBeginTable               []tableEntry
	imgAddressTable             []tableEntry
	imgNewLineTable             []tableEntry
	imgHorizontalDimensionTable []tableEntry
	imgVerticalDimensionTable   []tableEntry
	imgDescriptionTable         []tableEntry
)

func closeImgAddress(st *SourceText, tw *TokenWriter) {
	tw.nonEmptyBufIntoToken(TokenLinkAddress)
}

func andCloseImgAddress(f func(*SourceText, *TokenWriter)) func(*SourceText, *TokenWriter) {
	return func(st *SourceText, tw *TokenWriter) {
		closeImgAddress(st, tw)
		f(st, tw)
	}
}

func goToHorizontalDimension(st *SourceText, tw *TokenWriter) {
	eatChar(st) // Eat |
	tw.popState()
	tw.pushState(StateImgHorizontalDimension)
}

func goToDescription(st *SourceText, tw *TokenWriter) {
	eatChar(st) // Eat {
	tw.popState()
	tw.pushState(StateImgDescription)
}

func goToNewLine(st *SourceText, tw *TokenWriter) {
	eatChar(st) // Eat \n
	tw.popState()
	tw.pushState(StateImgNewLine)
}

func goOutOfImg(st *SourceText, tw *TokenWriter) {
	eatChar(st)   // Eat }
	tw.popState() // Pop current state
	tw.popState() // Pop img state
	tw.appendToken(Token{TokenImgClose, ""})
}

func init() {
	imgBeginTable = []tableEntry{
		{[]string{"{"}, func(st *SourceText, tw *TokenWriter) {
			tw.appendToken(Token{TokenImgOpen, ""})
			eatChar(st)
			tw.pushState(StateImgAddress)
			tw.buf.Reset()
		}},
		{[]string{"\\"}, beginEscaping},
		{[]string{"\n"}, func(st *SourceText, tw *TokenWriter) {
			// No { on this line! It was not an img to begin with!
			tw.popState()
			tw.pushState(StateParagraph)
		}},
	}
	imgAddressTable = []tableEntry{
		{[]string{"\\"}, beginEscaping},
		{[]string{"}"}, andCloseImgAddress(goOutOfImg)},
		{[]string{"{"}, andCloseImgAddress(goToDescription)},
		{[]string{"|"}, andCloseImgAddress(goToHorizontalDimension)},
		{[]string{"\n"}, andCloseImgAddress(goToNewLine)},
	}
	imgNewLineTable = []tableEntry{
		{[]string{"|"}, goToHorizontalDimension},
		{[]string{"}"}, goOutOfImg},
		{[]string{"{"}, goToDescription},
		{[]string{"\n"}, goToNewLine},
	}
	imgHorizontalDimensionTable = []tableEntry{}
	imgVerticalDimensionTable = []tableEntry{}
	imgDescriptionTable = []tableEntry{}
}
