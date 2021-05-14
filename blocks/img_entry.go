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

func (entry *ImgEntry) DescriptionAsHtml() (html string) {
	if entry.desc.Len() == 0 {
		return ""
	}
	lines := strings.Split(entry.desc.String(), "\n")
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			if html != "" {
				html += `<br>`
			}
			html += ParagraphToHtml(entry.hyphaName, line)
		}
	}
	return `<figcaption>` + html + `</figcaption>`
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
