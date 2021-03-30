package lexer

type tableEntry struct {
	prefix    string
	callback  func(*State)
	condition Condition
}

type imgTableEntry struct {
	charset   string
	callback  func(*State)
	condition Condition
}

var (
	table []tableEntry
	// For some limited times we use smaller subtables
	imgTable []tableEntry
)

func init() {
	/*
		table = []tableEntry{
			{"# ", λcallbackHeading(1),
				Condition{onNewLine: True, inGeneralText: True}},
			{"## ", λcallbackHeading(2),
				Condition{onNewLine: True, inGeneralText: True}},
			{"### ", λcallbackHeading(3),
				Condition{onNewLine: True, inGeneralText: True}},
			{"#### ", λcallbackHeading(4),
				Condition{onNewLine: True, inGeneralText: True}},
			{"##### ", λcallbackHeading(5),
				Condition{onNewLine: True, inGeneralText: True}},
			{"###### ", λcallbackHeading(6),
				Condition{onNewLine: True, inGeneralText: True}},
			{"\n", callbackHeadingNewLine,
				Condition{inHeading: True}},
			{"=>", callbackRocket,
				Condition{onNewLine: True, inGeneralText: True}},
			{"----", callbackHorizontalLine,
				Condition{onNewLine: True, inGeneralText: True, okForHorizontalLine: True}},
			{">", callbackBlockquote,
				Condition{onNewLine: True, inGeneralText: True}},
			{"img", callbackImg,
				Condition{onNewLine: True, inGeneralText: True}},
		}

		imgTable = []tableEntry{
			{"{", imgStartToLineBegin,
				Condition{onImgStart: True}},
			{"", eatChar,
				Condition{onImgStart: True}},

			{"}", imgToPreEnd,
				Condition{onImgLineBegin: True}},
			{" \t\n", eatChar,
				Condition{onImgLineBegin: True}},
			{"", imgAddAddrCh,
				Condition{onImgLineBegin: True}},

			{"{", imgToPara,
				Condition{onImgAddress: True}},
			{"|", imgToDimension,
				Condition{onImgAddress: True}},
			{"\n", imgNewLineToLineBegin,
				Condition{onImgAddress: True}},
			{"", imgAddAddrCh,
				Condition{onImgAddress: True}},

			//x/142
			{},
		}
	*/
}
