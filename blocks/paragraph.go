package blocks

import (
	"fmt"
)

// Paragraph is a block of formatted text.
type Paragraph struct {
	Formatted
}

// ID returns the paragraph's id which is paragraph- and a number.
func (p Paragraph) ID(counter *IDCounter) string {
	counter.paragraphs++
	return fmt.Sprintf("paragraph-%d", counter.paragraphs)
}
