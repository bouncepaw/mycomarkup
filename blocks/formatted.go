package blocks

import (
	"bytes"
)

// Formatted is a piece of formatted text. It is always part of a bigger block, such as Paragraph.
type Formatted struct {
	// HyphaName is the name of the hypha that contains the formatted text.
	HyphaName string
	Html      string
	*bytes.Buffer
	Spans []interface{} // Forgive me, for I have sinned
}

func (p Formatted) isBlock() {}

// ID returns an empty string because Formatted is always part of a bigger block.
func (p Formatted) ID(_ *IDCounter) string {
	return ""
}
