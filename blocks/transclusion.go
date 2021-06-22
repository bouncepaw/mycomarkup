package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/links"
	"strings"
)

// Transclusion is the block representing an extract from a different document.
// TODO: visitors for transclusion.
type Transclusion struct {
	// Target is the name of the hypha to be transcluded.
	Target   string
	Blend    bool
	Selector TransclusionSelector
}

// ID returns the transclusion's id which is transclusion- and a number.
func (t Transclusion) ID(counter *IDCounter) string {
	counter.transclusions++
	return fmt.Sprintf("transclusion-%d", counter.transclusions)
}

func (t Transclusion) isBlock() {}

// MakeTransclusion parses the line and returns a transclusion block.
func MakeTransclusion(line, hyphaName string) Transclusion {
	// TODO: move to the parser module.
	line = strings.TrimSpace(line[2:])
	if line == "" {
		return Transclusion{"", false, SelectorOverview}
	}

	if strings.ContainsRune(line, '|') {
		parts := strings.SplitN(line, "|", 2)
		return Transclusion{
			Target:   links.From(strings.TrimSpace(parts[0]), "", hyphaName).TargetHypha(),
			Blend:    strings.Contains(parts[1], "blend"),
			Selector: selectorFrom(parts[1]),
		}
	}

	return Transclusion{
		Target:   links.From(strings.TrimSpace(line), "", hyphaName).TargetHypha(),
		Blend:    false,
		Selector: SelectorOverview,
	}
}

// TransclusionSelector is the thing that specifies what parts of the document shall be transcluded.
type TransclusionSelector int

const (
	// SelectorOverview is SelectorAttachment and SelectorDescription combined.
	SelectorOverview TransclusionSelector = iota
	// SelectorAttachment selects the attachment of the target hypha.
	SelectorAttachment
	// SelectorDescription selects the description of the target hypha. The first paragraph of the text part of the hypha is considered its description.
	SelectorDescription
	// SelectorText selects all of the text in the hypha, but not the attachment.
	SelectorText
	// SelectorFull selects everything in the hypha, including the attachment.
	SelectorFull
)

func selectorFrom(s string) TransclusionSelector {
	switch {
	case strings.Contains(s, "full"):
		return SelectorFull
	case strings.Contains(s, "text"):
		return SelectorText
	case strings.Contains(s, "overview"):
		return SelectorOverview
	case strings.Contains(s, "description"):
		if strings.Contains(s, "attachment") {
			return SelectorOverview
		}
		return SelectorDescription
	case strings.Contains(s, "attachment"):
		return SelectorAttachment
	default:
		return SelectorOverview
	}
}
