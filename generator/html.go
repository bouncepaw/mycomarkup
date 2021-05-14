// Package generator provides the HTML generator.
package generator

import (
	"fmt"

	"github.com/bouncepaw/mycomarkup/blocks"
)

// BlockToHTML turns the given block into HTML. It supports only a subset of Mycomarkup.
func BlockToHTML(block interface{}) string {
	switch b := block.(type) {
	case blocks.HorizontalLine:
		return fmt.Sprintf(`<hr id="%s"/>`, b.ID())
	case blocks.Img:
		return imgToHTML(b)
	case blocks.ImgEntry:
		return imgEntryToHTML(b)
	}
	return ""
}

func imgEntryToHTML(entry blocks.ImgEntry) string {
	var ret string
	if entry.Srclink.DestinationUnknown {
		ret += fmt.Sprintf(
			`<a class="%s" href="%s">Hypha <i>%s</i> does not exist</a>`,
			entry.Srclink.Classes(),
			entry.Srclink.Href(),
			entry.Srclink.Address())
	} else {
		ret += fmt.Sprintf(
			`<a href="%s"><img src="%s" %s %s></a>`,
			entry.Srclink.Href(),
			entry.Srclink.ImgSrc(),
			entry.SizeWAsAttr(),
			entry.SizeHAsAttr())
	}
	return fmt.Sprintf(`<figure class="img-gallery__entry">
%s
%s
</figure>
`, ret, entry.DescriptionAsHtml())
}

func imgToHTML(img blocks.Img) string {
	img.MarkExistenceOfSrcLinks()
	var ret string
	for _, entry := range img.Entries {
		ret += BlockToHTML(entry)
	}
	return fmt.Sprintf(`<section class="img-gallery %s">
%s</section>`,
		func() string {
			if img.HasOneImage() {
				return "img-gallery_one-image"
			}
			return "img-gallery_many-images"
		}(),
		ret)
}
