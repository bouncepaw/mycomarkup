package lexer

type tableEntry struct {
	prefix    string
	callback  func(*State)
	condition Condition
}

var (
	table []tableEntry
	// For some limited times we use smaller subtables
	rocketLinkTable []tableEntry
)

func init() {
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
	}
}
