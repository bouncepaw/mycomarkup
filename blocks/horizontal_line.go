package blocks

import (
	"fmt"
)

type HorizontalLine struct {
	length uint
}

func MakeHorizontalLine(length uint) *HorizontalLine {
	return &HorizontalLine{length}
}

func (hr *HorizontalLine) String() string {
	return fmt.Sprintf(`HorizontalLine(%d);`, hr.length)
}

func (hr *HorizontalLine) IsNesting() bool {
	return false
}

func (hr *HorizontalLine) Kind() BlockKind {
	return KindHorizontalLine
}

func (hr *HorizontalLine) Length() uint {
	return hr.length
}
