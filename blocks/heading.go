package blocks

import (
	"github.com/bouncepaw/mycomarkup/v2/util"
)

// Heading is a formatted heading in the document.
type Heading struct {
	// Level is a number between 1 and 6.
	Level    uint
	Contents Formatted
	Src      string
}

func (h Heading) isBlock() {}

// GetContents returns the heading's contents.
func (h *Heading) GetContents() Formatted {
	return h.Contents
}

// ID returns the heading's id which is basically a stripped version of its contents. See util.StringID.
func (h Heading) ID(_ *IDCounter) string {
	return util.StringID(h.Src[h.Level+1:])
}
