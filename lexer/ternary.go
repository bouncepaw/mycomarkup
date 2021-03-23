package lexer

type Ternary int

const (
	Unknown Ternary = iota
	True
	False
)

func (t Ternary) isTrue() bool {
	return t == True
}

func (t Ternary) notUnknown() bool {
	return t != Unknown
}

type Condition struct {
	onNewLine           Ternary
	okForHorizontalLine Ternary
	inGeneralText       Ternary
	inHeading           Ternary
}

func (c *Condition) fullfilledBy(s *State) Ternary {
	switch {
	case c.onNewLine.notUnknown() && c.onNewLine != s.onNewLine():
		return False
	case c.okForHorizontalLine.notUnknown() && c.okForHorizontalLine != s.okForHorizontalLine():
		return False
	case c.inGeneralText.notUnknown() && c.inGeneralText != s.inGeneralText():
		return False
	case c.inHeading.notUnknown() && c.inHeading != s.inHeading:
		return False
	}
	return True
}
