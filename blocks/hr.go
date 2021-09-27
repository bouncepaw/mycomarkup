package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v2/util"
)

// HorizontalLine represents the horizontal line block.
//
// In Mycomarkup it is written like that:
//     ----
type HorizontalLine struct {
	src string
}

// NewHorizontalLine parses the horizontal line block on the given text line and returns it.
func NewHorizontalLine(line string) HorizontalLine {
	return HorizontalLine{
		src: line,
	}
}

// ID returns the line's id. By default, it is hr- and a number. If the line was written like that:
//    ----id
// , the specified id is returned instead.
func (h HorizontalLine) ID(counter *IDCounter) string {
	counter.hrs++
	if len(h.src) > 4 {
		return util.StringID(h.src[4:])
	}
	return fmt.Sprintf("hr-%d", counter.hrs)
}
