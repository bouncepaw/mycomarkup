package blocks

import (
	"github.com/bouncepaw/mycomarkup/v3/links"
)

// ImgEntry is an entry of an image gallery. It can only be nested into Img. V3: proper readers, encapsulate
type ImgEntry struct {
	Target      links.Link
	HyphaName   string
	Width       string
	Height      string
	Description string // TODO: change to Formatted type.
}

// ID returns an empty string because images do not have ids. Image galleries do have them, by the way, see Img.
func (entry ImgEntry) ID(_ *IDCounter) string {
	return ""
}

// GetWidth returns the width property of the entry. TODO: rename to Width.
func (entry ImgEntry) GetWidth() string { return entry.Width }

// GetHeight returns the height property of the entry. TODO: rename to Height.
func (entry ImgEntry) GetHeight() string { return entry.Height }
