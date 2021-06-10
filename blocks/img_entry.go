package blocks

import (
	"strings"

	"github.com/bouncepaw/mycomarkup/links"
)

// ImgEntry is an entry of an image gallery. It can only be nested into Img.
type ImgEntry struct {
	Srclink   *links.Link
	hyphaName string
	path      strings.Builder
	sizeW     strings.Builder
	sizeH     strings.Builder
	desc      strings.Builder
}

// ID returns an empty string because images do not have ids. Image galleries do have them, by the way, see Img.
func (entry ImgEntry) ID(_ *IDCounter) string {
	return ""
}

func (entry ImgEntry) isBlock() {}

// Description returns the description of the image.
func (entry *ImgEntry) Description() Formatted {
	return MakeFormatted(entry.desc.String(), entry.hyphaName)
}

// SizeWAsAttr returns either an empty string or the width attribute for the image, depending on what has been written in the markup.
func (entry *ImgEntry) SizeWAsAttr() string {
	if entry.sizeW.Len() == 0 {
		return ""
	}
	return ` width="` + entry.sizeW.String() + `"`
}

// SizeHAsAttr returns either an empty string or the height attribute for the image, depending on what has been written in the markup.
func (entry *ImgEntry) SizeHAsAttr() string {
	if entry.sizeH.Len() == 0 {
		return ""
	}
	return ` height="` + entry.sizeH.String() + `"`
}
