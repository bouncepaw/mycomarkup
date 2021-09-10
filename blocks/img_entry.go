package blocks

import (
	"strings"

	"github.com/bouncepaw/mycomarkup/links"
)

// ImgEntry is an entry of an image gallery. It can only be nested into Img.
type ImgEntry struct {
	Srclink   *links.Link
	HyphaName string
	Path      strings.Builder
	SizeW     strings.Builder
	SizeH     strings.Builder
	Desc      strings.Builder
}

// ID returns an empty string because images do not have ids. Image galleries do have them, by the way, see Img.
func (entry ImgEntry) ID(_ *IDCounter) string {
	return ""
}

func (entry ImgEntry) isBlock() {}

// SizeWAsAttr returns either an empty string or the width attribute for the image, depending on what has been written in the markup.
func (entry *ImgEntry) SizeWAsAttr() string {
	if entry.SizeW.Len() == 0 {
		return ""
	}
	return ` width="` + entry.SizeW.String() + `"`
}

// SizeHAsAttr returns either an empty string or the height attribute for the image, depending on what has been written in the markup.
func (entry *ImgEntry) SizeHAsAttr() string {
	if entry.SizeH.Len() == 0 {
		return ""
	}
	return ` height="` + entry.SizeH.String() + `"`
}
