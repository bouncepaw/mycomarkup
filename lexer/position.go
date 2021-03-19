package lexer

import (
	"strings"
)

type Position struct {
	lineFrom   uint
	columnFrom uint
	lineTo     uint
	columnTo   uint
	text       strings.Builder
}
