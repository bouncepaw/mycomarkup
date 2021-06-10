package blocks

import (
	"github.com/bouncepaw/mycomarkup/util"
)

// Heading is a formatted heading in the document.
type Heading struct {
	// Level is a number between 1 and 6.
	Level    uint
	contents Formatted
	src      string
}

func (h Heading) isBlock() {}

// MakeHeading parses the heading on the given line and returns it. Set its level by yourself though.
func MakeHeading(line, hyphaName string, level uint) Heading {
	// TODO: figure out the level here.
	// TODO: move to the parser module.
	h := Heading{
		Level:    level,
		contents: MakeFormatted(line[level+1:], hyphaName),
		src:      line,
	}
	return h
}

// Contents returns the heading's contents.
func (h *Heading) Contents() Formatted {
	return h.contents
}

// ID returns the heading's id which is basically a stripped version of its contents. See util.StringID.
func (h Heading) ID(_ *IDCounter) string {
	return util.StringID(h.src[h.Level+1:])
}
