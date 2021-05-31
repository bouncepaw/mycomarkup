package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/links"
	"strings"

	"github.com/bouncepaw/mycomarkup/util"
)

// Transclusion is the block representing an extract from a different document.
type Transclusion struct {
	// Target is the name of the hypha to be transcluded.
	Target string

	// Selector specifies what parts of the hypha to transclude.
	Selector string
}

func (t Transclusion) ID(counter *IDCounter) string {
	counter.transclusions++
	return fmt.Sprintf("transclusion-%d", counter.transclusions)
}

func (t Transclusion) IsBlock() {}

// TransclusionError is a value that means that the transclusion is wrong.
const TransclusionError = "err"

func MakeTransclusion(line, hyphaName string) Transclusion {
	line = strings.TrimSpace(util.Remover("<=")(line))
	if line == "" {
		return Transclusion{"", TransclusionError}
	}

	if strings.ContainsRune(line, '|') {
		parts := strings.SplitN(line, "|", 2)
		return Transclusion{
			Target:   links.From(strings.TrimSpace(parts[0]), "", hyphaName).Address(),
			Selector: strings.TrimSpace(parts[1]),
		}
	}

	return Transclusion{
		Target:   links.From(strings.TrimSpace(line), "", hyphaName).Address(),
		Selector: "",
	}
}
