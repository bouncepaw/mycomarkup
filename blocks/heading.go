package blocks

import (
	"fmt"
)

type Heading struct {
	level   uint
	content *Formatted
}

func MakeHeading(level uint, content *Formatted) *Heading {
	return &Heading{level, content}
}

func (h *Heading) String() string {
	return fmt.Sprintf(`Heading(%d) {
%s
};`, h.level, h.content)
}

func (h *Heading) IsNesting() bool {
	return false
}

func (h *Heading) Kind() BlockKind {
	return KindHeading
}
