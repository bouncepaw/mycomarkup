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
	// Test if any of the ternary variables are not ok
	cs := []struct {
		requirement Ternary
		stateState  Ternary
	}{
		{c.inHeading, s.inHeading},
	}
	for _, cp := range cs {
		if cp.requirement.notUnknown() && cp.requirement != cp.stateState {
			return False
		}
	}

	// Then test the same for those that need a function invocation
	cf := []struct {
		requirement Ternary
		stateTest   func() Ternary
	}{
		{c.onNewLine, s.onNewLine},
		{c.inGeneralText, s.inGeneralText},
		{c.okForHorizontalLine, s.okForHorizontalLine},
	}
	for _, cp := range cf {
		if cp.requirement.notUnknown() && cp.requirement != cp.stateTest() {
			return False
		}
	}
	return True
}
