package blocks

import (
	"github.com/bouncepaw/mycomarkup/util"
)

type Heading struct {
	Level    uint
	contents Formatted
	src      string
}

func MakeHeading(line, hyphaName string, level uint) Heading {
	h := Heading{
		Level:    level,
		contents: MakeParagraph(line[level+1:], hyphaName),
		src:      line,
	}
	return h
}

func (h *Heading) Contents() Formatted {
	return h.contents
}

func (h *Heading) ID() string {
	return util.StringID(h.src[h.Level+1:])
}
