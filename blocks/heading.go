package blocks

import (
	"github.com/bouncepaw/mycomarkup/util"
)

type Heading struct {
	Level        uint
	ContentsHTML string
	src          string
	LegacyID     int
}

func MakeHeading(line, hyphaName string, level uint, legacyID int) Heading {
	h := Heading{
		Level:        level,
		ContentsHTML: ParagraphToHtml(hyphaName, line[level+1:]),
		src:          line,
		LegacyID:     legacyID,
	}
	return h
}

func (h *Heading) ID() string {
	return util.StringID(h.src[h.Level+1:])
}
