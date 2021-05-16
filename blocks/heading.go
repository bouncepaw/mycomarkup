package blocks

import (
	"github.com/bouncepaw/mycomarkup/util"
)

type Heading struct {
	Level        uint
	ContentsHTML string
	src          string
}

func MakeHeading(line, hyphaName string, level uint) Heading {
	h := Heading{
		Level:        level,
		ContentsHTML: ParagraphToHtml(hyphaName, line[level+1:]),
		src:          line,
	}
	return h
}

func (h *Heading) ID() string {
	return util.StringID(h.src[h.Level+1:])
}
