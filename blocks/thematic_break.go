package blocks

import (
	"fmt"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/util"
)

// ThematicBreak represents the thematic line block, represented by a horizontal line.
//
// In Mycomarkup it is written like that:
//     ----
type ThematicBreak struct {
	src string
}

// NewThematicBreak parses the horizontal line block on the given text line and returns it.
func NewThematicBreak(line string) ThematicBreak {
	return ThematicBreak{
		src: line,
	}
}

// ID returns the line's id. By default, it is hr- and a number. If the line was written like that:
//    ----id
// , the specified id is returned instead.
func (h ThematicBreak) ID(counter *IDCounter) string {
	counter.hrs++
	if len(h.src) > 4 {
		return util.StringID(h.src[4:])
	}
	return fmt.Sprintf("hr-%d", counter.hrs)
}
