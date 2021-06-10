package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/links"
	"strings"

	"github.com/bouncepaw/mycomarkup/util"
)

// Transclusion is the block representing an extract from a different document.
// TODO: visitors for transclusion.
type Transclusion struct {
	// Target is the name of the hypha to be transcluded.
	Target string

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
		return Transclusion{"", DefaultSelector()}
	}

	if strings.ContainsRune(line, '|') {
		parts := strings.SplitN(line, "|", 2)
		return Transclusion{
			Target:   links.From(strings.TrimSpace(parts[0]), "", hyphaName).Address(),
			Selector: ParseSelector(parts[1]),
		}
	}

	return Transclusion{
		Target:   links.From(strings.TrimSpace(line), "", hyphaName).Address(),
		Selector: DefaultSelector(),
	}
}

// TransclusionSelector is the thing that specifies what parts of the document shall be transcluded.
type TransclusionSelector struct {
	bound1      string
	dotsPresent bool
	bound2      string
}

// DefaultSelector returns the default selector which is start..description.
func DefaultSelector() TransclusionSelector {
	return TransclusionSelector{"start", true, "description"}
}

// ParseSelector parses the selector according to the following rules.
//
// If the selector is empty, we think of it as of selector start..description and try again.
//
// If there is no .. in the selector, the selector selects just one block with the matching id.
//
// If there is a .. in selector, there are two bounds: left and right. Both bounds are ids of some blocks.
//
// Special bounds:
//
//     attachment: hypha's attachment.
//     start: hypha's text's first block.
//     description: hypha's text's first paragraph.
//     end: hypha's last block.
//
// If the left bound is empty, it is set to start. If the right bound is empty, it is set to end.
func ParseSelector(selector string) TransclusionSelector {
	selector = strings.TrimSpace(selector)
	if selector == "" {
		return ParseSelector("start..description")
	}

	if parts := strings.SplitN(selector, "..", 2); len(parts) == 2 {
		return TransclusionSelector{
			util.DefaultString(strings.TrimRight(parts[0], " "), "start"),
			true,
			util.DefaultString(strings.TrimLeft(parts[1], " "), "end"),
		}
	}

	return TransclusionSelector{
		selector,
		false,
		"",
	}
}
