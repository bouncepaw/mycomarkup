package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/util"
)

// HorizontalLine represents the horizontal line block.
//
// In Mycomarkup it is written like that:
//     ----
type HorizontalLine struct {
	src string
}

func (h HorizontalLine) IsBlock() {}

func MakeHorizontalLine(src string) HorizontalLine {
	return HorizontalLine{
		src: src,
	}
}

func (h HorizontalLine) ID(counter *IDCounter) string {
	counter.hrs++
	if len(h.src) > 4 {
		return util.StringID(h.src[4:])
	}
	return fmt.Sprintf("hr-%d", counter.hrs)
}
