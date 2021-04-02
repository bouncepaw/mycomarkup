package lexer

// LexerState is a type representing //state// (as in terms of finite-state automata). Depending on the lexer state, different lexing tables may be used.
type LexerState int

const (
	StateErr LexerState = iota
	StateNil

	StateParagraph
	StateEscape
	StateNowiki
	StateAutolink
	StateLinkAddress
	StateLinkDisplay
)

// StateStack is a stack of LexerState.
type StateStack struct {
	stack   []LexerState
	topElem LexerState
}

// newStateStack returns a state stack that already has one element: StateNil.
func newStateStack() StateStack {
	ss := StateStack{
		stack: []LexerState{StateNil},
	}
	ss.topElem = StateNil
	return ss
}

// hasOnTop compares the top element with the given lexer state.
func (ss *StateStack) hasOnTop(ls LexerState) bool {
	return ss.topElem == ls
}

// pushState pushes the given new lexer state to the stack.
func (ss *StateStack) pushState(ls LexerState) {
	ss.topElem = ls
	ss.stack = append(ss.stack, ls)
}

// popState pops the top lexer state from the stack. It is O(1), I think. If there are no lexer states in the stack, it may panic, I guess; don't do that.
func (ss *StateStack) popState() LexerState {
	lastElem := ss.stack[ss.lastElemPos()]
	ss.stack = ss.stack[:ss.lastElemPos()]
	if len(ss.stack) > 0 {
		ss.topElem = ss.stack[ss.lastElemPos()]
	}
	return lastElem
}

func (ss *StateStack) lastElemPos() int {
	return len(ss.stack) - 1
}
