package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v3/links"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"github.com/bouncepaw/mycomarkup/v3/util"
	"strings"
)

// Transclusion is the block representing an extract from a different document.
type Transclusion struct {
	// Target is the name of the hypha to be transcluded.
	Target   string
	Blend    bool
	Selector TransclusionSelector
	TransclusionError
}

// ID returns the transclusion's id which is transclusion- and a number.
func (t Transclusion) ID(counter *IDCounter) string {
	counter.transclusions++
	return fmt.Sprintf("transclusion-%d", counter.transclusions)
}

// MakeTransclusion parses the line and returns a transclusion block. V3
func MakeTransclusion(ctx mycocontext.Context, line string) Transclusion {
	if ctx.TransclusionSupported() {
		return Transclusion{
			"",
			false,
			SelectorOverview,
			TransclusionError{TransclusionInTerminal},
		}
	}
	line = strings.TrimSpace(line[2:])
	// The second operand matches transclusion like <= | tututu or <= |
	if line == "" || strings.HasPrefix(line, "|") {
		return Transclusion{
			"",
			false,
			SelectorOverview,
			TransclusionError{TransclusionErrorNoTarget},
		}
	}

	if strings.ContainsRune(line, '|') {
		var (
			parts       = strings.SplitN(line, "|", 2)
			targetHypha = util.CanonicalName(strings.TrimSpace(parts[0]))
		)
		if strings.ContainsRune(targetHypha, ':') {
			return Transclusion{
				Target:            targetHypha,
				Blend:             false,
				Selector:          SelectorOverview,
				TransclusionError: TransclusionError{TransclusionErrorOldSyntax},
			}
		}
		// Sorry for party rocking
		targetHypha = links.From(targetHypha, "", ctx.HyphaName()).TargetHypha()
		if !mycocontext.HyphaExists(ctx, targetHypha) {
			return Transclusion{
				Target:            targetHypha,
				Blend:             false,
				Selector:          SelectorOverview,
				TransclusionError: TransclusionError{TransclusionErrorNotExists},
			}
		}
		return Transclusion{
			Target:   links.From(targetHypha, "", ctx.HyphaName()).TargetHypha(),
			Blend:    strings.Contains(parts[1], "blend"),
			Selector: selectorFrom(parts[1]),
		}
	}

	return Transclusion{
		Target:   links.From(strings.TrimSpace(line), "", ctx.HyphaName()).TargetHypha(),
		Blend:    false,
		Selector: SelectorOverview,
	}
}

// TransclusionErrorReason is the reason why the transclusion failed during parsing.
type TransclusionErrorReason int

const (
	// TransclusionNoError means there is no error.
	TransclusionNoError TransclusionErrorReason = iota
	// TransclusionInTerminal means that Mycomarkup CLI is used. Transclusion is not supported in it.
	TransclusionInTerminal
	// TransclusionErrorNoTarget means that no target hypha was specified.
	TransclusionErrorNoTarget
	// TransclusionErrorOldSyntax means : was found in the target.
	TransclusionErrorOldSyntax
	// TransclusionErrorNotExists means the target hypha does not exist.
	TransclusionErrorNotExists
)

// TransclusionError is the error that occurred during transclusion.
type TransclusionError struct {
	Reason TransclusionErrorReason
}

// HasError is true if there is indeed an error.
func (te *TransclusionError) HasError() bool {
	return te.Reason != TransclusionNoError
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
