package blocks

import (
	"lesarbr.es/mycomarkup/v5/util"
)

// Heading is a formatted heading in the document.
type Heading struct {
	// level is a number between 1 and 4.
	level    uint
	contents Formatted
	srcLine  string
}

// NewHeading returns a Heading with the given data.
func NewHeading(level uint, contents Formatted, srcLine string) Heading {
	return Heading{
		level:    level,
		contents: contents,
		srcLine:  srcLine,
	}
}

// Level returns the Heading's level, 1 from 6.
//
//     Prefix  | Level
//     =      | 1
//     ==     | 2
//     ===    | 3
//     ====   | 4
func (h Heading) Level() uint {
	return h.level
}

// Contents returns the Heading's contents.
func (h Heading) Contents() Formatted {
	return h.contents
}

// ID returns the Heading's id which is basically a stripped version of its contents. See util.StringID.
func (h Heading) ID(_ *IDCounter) string {
	return util.StringID(h.srcLine[h.level:])
}
