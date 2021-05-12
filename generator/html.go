package generator

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/blocks"
)

func BlockToHTML(block interface{}) string {
	switch b := block.(type) {
	case blocks.HorizontalLine:
		return fmt.Sprintf(`<hr id="%s"/>`, b.ID())
	}
	return ""
}
