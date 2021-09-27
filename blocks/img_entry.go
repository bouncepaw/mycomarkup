package blocks

import (
	"github.com/bouncepaw/mycomarkup/v2/links"
)

// ImgEntry is an entry of an image gallery. It can only be nested into Img.
type ImgEntry struct {
	Srclink     links.Link
	HyphaName   string
	Width       string
	Height      string
	Description string // TODO: change to Formatted type.
}

// ID returns an empty string because images do not have ids. Image galleries do have them, by the way, see Img.
func (entry ImgEntry) ID(_ *IDCounter) string {
	return ""
}

func (entry ImgEntry) isBlock() {}

// WidthAttributeHTML returns either an empty string or the width attribute for the image, depending on what has been written in the markup.
func (entry *ImgEntry) WidthAttributeHTML() string {
	if len(entry.Width) == 0 {
		return ""
	}
	return ` width="` + entry.Width + `"`
}

// HeightAttributeHTML returns either an empty string or the height attribute for the image, depending on what has been written in the markup.
func (entry *ImgEntry) HeightAttributeHTML() string {
	if len(entry.Height) == 0 {
		return ""
	}
	return ` height="` + entry.Height + `"`
}
