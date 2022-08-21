package blocks

import (
	"fmt"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/links"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
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
	line = strings.TrimSpace(line[2:])
	switch {
	case !ctx.TransclusionSupported():
		return Transclusion{
			"",
			false,
			SelectorOverview,
			TransclusionError{TransclusionNotSupported},
		}
	case line == "" || strings.HasPrefix(line, "|"):
		// The second operand matches transclusion like <= | tututu or <= |
		return Transclusion{
			"",
			false,
			SelectorOverview,
			TransclusionError{TransclusionErrorNoTarget},
		}
	}

	var target string
	if pipepos := strings.IndexRune(line, '|'); pipepos >= 0 {
		target = line[0:pipepos]
	} else {
		target = line
	}

	link := links.LinkFrom(ctx, target, "") // Don't care about display, never happens
	switch link := link.(type) {
	case *links.InterwikiLink:
		return Transclusion{
			Target:            target,
			Blend:             false,
			Selector:          SelectorOverview,
			TransclusionError: TransclusionError{TransclusionCannotTranscludeInterwiki},
		}
	case *links.URLLink, *links.LocalRootedLink:
		return Transclusion{
			Target:            target,
			Blend:             false,
			Selector:          SelectorOverview,
			TransclusionError: TransclusionError{TransclusionCannotTranscludeURL},
		}
	case *links.LocalLink:
		target = link.Target(ctx)
		if !mycocontext.HyphaExists(ctx, target) {
			return Transclusion{
				Target:            target,
				Blend:             false,
				Selector:          SelectorOverview,
				TransclusionError: TransclusionError{TransclusionErrorNotExists},
			}
		}

		var (
			selector TransclusionSelector
			blend    bool
		)
		if pipepos := strings.IndexRune(line, '|'); pipepos >= 0 {
			selector = selectorFrom(line[pipepos+1:])
			blend = strings.Contains(line[pipepos+1:], "blend")
		} else {
			selector = SelectorOverview
		}

		return Transclusion{
			Target:   target,
			Blend:    blend,
			Selector: selector,
		}
	}
	panic("unreachable")
}

// TransclusionErrorReason is the reason why the transclusion failed during parsing.
type TransclusionErrorReason int

const (
	// TransclusionNoError means there is no error.
	TransclusionNoError TransclusionErrorReason = iota
	// TransclusionNotSupported means that transclusion is not supported.
	TransclusionNotSupported
	// TransclusionErrorNoTarget means that no target hypha was specified.
	TransclusionErrorNoTarget
	// TransclusionCannotTranscludeURL means : was found in the target.
	TransclusionCannotTranscludeURL
	// TransclusionCannotTranscludeInterwiki means an interwiki transclusion was attempted.
	TransclusionCannotTranscludeInterwiki
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
