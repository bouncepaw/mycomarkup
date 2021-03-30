package lexer

type stateStack struct {
	stack   []LexerState
	topElem *LexerState
}

func newStateStack() *stateStack {
	ss := stateStack{
		stack: []LexerState{StateNil},
	}
	ss.topElem = &ss.stack[0]
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
	lastElem := ss.stack[ss.lastElemPos()]
	ss.stack = ss.stack[:ss.lastElemPos()]
	if len(ss.stack) == 0 {
		ss.topElem = &ss.stack[0]
	} else {
		ss.topElem = &ss.stack[ss.lastElemPos()]
	}
	return lastElem
}
