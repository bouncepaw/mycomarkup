package blocks

import (
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
			Target:   util.XclCanonicalName(hyphaName, strings.TrimSpace(parts[0])),
			Selector: strings.TrimSpace(parts[1]),
		}
	}

	return Transclusion{
		Target:   util.XclCanonicalName(hyphaName, strings.TrimSpace(line)),
		Selector: "",
	}
}

func (t *Transclusion) String() string {
	panic("implement me")
}

func (t *Transclusion) ID() string {
	panic("implement me")
}
