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

func (h *HorizontalLine) String() string {
	return fmt.Sprintf(`HorizontalLine#%s;`, h.ID())
}

func (h *HorizontalLine) ID() string {
	if len(h.src) > 4 {
		return util.StringID(h.src[4:])
	}
	return ""
}
