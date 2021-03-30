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
	inSpan              Ternary
}

func (c *Condition) fullfilledBy(s *State) Ternary {
	return True
}
