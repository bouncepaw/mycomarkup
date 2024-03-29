// Package temporary_workaround is a temporary workaround to break import cycle for some transclusion tricks.
//
// It is planned to get rid of it one day.
package temporary_workaround

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5/blocks"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
)

var TransclusionVisitor func(xcl blocks.Transclusion) (
	visitor func(block blocks.Block),
	result func() ([]blocks.Block, error),
)

var BlockTree func(ctx mycocontext.Context, visitors ...func(block blocks.Block)) []blocks.Block
