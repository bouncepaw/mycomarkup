package lexer

import ()

// See table.go for imgTable, see lexer.go for the table usage

type imgStatePosition int

const (
	imgStart imgStatePosition = iota
	imgLineBegin
	imgAddress
	imgDescription
	imgDimensionBegin
	imgDimensionHorizontal
	imgDimensionVertical
	imgPreEnd
	imgEnd
)

type imgState struct {
	position imgStatePosition
}

func (is *imgState) transition(isp imgStatePosition) {
	is.position = isp
}

func imgStartToLineBegin(s *State) {
	eatChar(s)
}
