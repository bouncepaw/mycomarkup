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

func (entry ImgEntry) IsBlock() {}

func (entry *ImgEntry) Description() Formatted {
	return MakeFormatted(entry.desc.String(), entry.hyphaName)
}

func (entry *ImgEntry) SizeWAsAttr() string {
	if entry.sizeW.Len() == 0 {
		return ""
	}
	return ` width="` + entry.sizeW.String() + `"`
}

func (entry *ImgEntry) SizeHAsAttr() string {
	if entry.sizeH.Len() == 0 {
		return ""
	}
	return ` height="` + entry.sizeH.String() + `"`
}
