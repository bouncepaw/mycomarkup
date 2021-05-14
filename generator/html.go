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
	case blocks.LaunchPad:
		return launchpadToHTML(b)
	case blocks.RocketLink:
		return fmt.Sprintf(`
	<li class="launchpad__entry"><a href="%s" class="rocketlink %s">%s</a></li>`, b.Href(), b.Classes(), b.Display())
	case blocks.Heading:
		return fmt.Sprintf(`<h%[1]d id='%[2]d'>%[3]s<a href="#%[4]s" id="%[4]s" class="heading__link"></a></h%[1]d>
`, b.Level, b.LegacyID, b.ContentsHTML, b.ID())
	}
	return ""
}

func launchpadToHTML(lp blocks.LaunchPad) string {
	var ret string
	for _, rocket := range lp.Rockets {
		ret += BlockToHTML(rocket)
	}
	return fmt.Sprintf(`<ul class="launchpad">%s
</ul>`, ret)
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
