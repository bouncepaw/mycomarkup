package blocks

import (
	"fmt"
)

type Heading struct {
	level   uint
	id      string
	content *Formatted
}

func MakeHeading(level uint, id string, content *Formatted) *Heading {
	return &Heading{level, id, content}
}

func (h *Heading) String() string {
	return fmt.Sprintf(`Heading(%d, "%s") {
%s
};`, h.level, h.id, h.content)
}

func (h *Heading) IsNesting() bool {
	return false
}

func (h *Heading) Kind() BlockKind {
	return KindHeading
}
