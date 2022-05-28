package blocks

import (
	"github.com/bouncepaw/mycomarkup/v4/links"
)

// ImgEntry is an entry of an image gallery. It can only be nested into Img. V3: proper readers, encapsulate
type ImgEntry struct {
	Target      links.LegacyLink
	hyphaName   string
	width       string
	height      string
	description string // TODO: change to Formatted type.
}

// NewImgEntry returns a new ImgEntry.
func NewImgEntry(target links.LegacyLink, hyphaName, width, height, description string) ImgEntry {
	return ImgEntry{
		Target:      target,
		hyphaName:   hyphaName,
		width:       width,
		height:      height,
		description: description,
	}
}

// ID returns an empty string because images do not have ids. Image galleries do have them, by the way, see Img.
func (entry ImgEntry) ID(_ *IDCounter) string { return "" }

// Width returns the width property of the entry.
func (entry ImgEntry) Width() string { return entry.width }

// Height returns the height property of the entry.
func (entry ImgEntry) Height() string { return entry.height }

// Description returns the description of the entry. The description is unparsed Mycomarkup string.
func (entry ImgEntry) Description() string { return entry.description }
