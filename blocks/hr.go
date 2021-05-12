package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/util"
)

type HorizontalLine struct {
	TerminalBlock
	src string
}

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
