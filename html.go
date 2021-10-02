package mycomarkup

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/v3/blocks"
)

func idAttribute(b blocks.Block, counter *blocks.IDCounter) string {
	switch id := b.ID(counter); {
	case !counter.ShouldUseResults(), id == "":
		return ""
	default:
		return fmt.Sprintf(` id="%s"`, id)
	}
}
