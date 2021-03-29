package lexer

type stateStack struct {
	stack   []LexerState
	topElem *LexerState
}

func newStateStack() *stateStack {
	ss := stateStack{
		stack:   make([]LexerState, 1),
		topElem: &StateNil,
	}
	return &ss
}

func (ss *stateStack) lastElemPos() int {
	return len(ss.stack) - 1
}

func (ss *stateStack) push(ls LexerState) {
	ss.topElem = &ls
	ss.stack = append(ss.stack, ls)
}

func (ss *stateStack) pop() LexerState {
	lastElem := ss[ss.lastElemPos()]
	ss.stack = ss.stack[:ss.lastElemPos()]
	if len(ss.stack) == 0 {
		ss.topElem = &StateNil
	} else {
		ss.topElem = &ss.stack[ss.lastElemPos()]
	}
	return lastElem
}
