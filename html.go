package mycomarkup

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/parser"
	"html"
)

// BlockToHTML turns the given block into HTML. It supports only a subset of Mycomarkup.
func BlockToHTML(block blocks.Block, counter *blocks.IDCounter) string {
	switch b := block.(type) {
	case blocks.Formatted:
		return b.Html
	case blocks.Paragraph:
		return fmt.Sprintf("\n<p%s>%s</p>", idAttribute(b, counter), b.Html)
	case blocks.HorizontalLine:
		return fmt.Sprintf(`<hr id="%s"/>`, idAttribute(b, counter))
	case blocks.Img:
		return imgToHTML(b, counter)
	case blocks.ImgEntry:
		return imgEntryToHTML(b, counter)
	case blocks.LaunchPad:
		return launchpadToHTML(b, counter)
	case blocks.RocketLink:
		return fmt.Sprintf(`
	<li class="launchpad__entry"><a href="%s" class="rocketlink %s">%s</a></li>`, b.Href(), b.Classes(), html.EscapeString(b.Display()))
	case blocks.Heading:
		return fmt.Sprintf(`
<h%[1]d%[4]s>%[2]s<a href="#%[3]s" id="%[3]s" class="heading__link"></a></h%[1]d>
`, b.Level, BlockToHTML(b.GetContents(), counter), b.ID(counter), idAttribute(b, counter))
	case blocks.CodeBlock:
		return fmt.Sprintf("\n<pre class='codeblock'%s><code class='language-%s'>%s</code></pre>", idAttribute(b, counter), b.Language(), b.Contents())
	}
	fmt.Printf("%q\n", block)
	return "<b>UNKNOWN ELEMENT</b>"
}

func idAttribute(b blocks.Block, counter *blocks.IDCounter) string {
	switch id := b.ID(counter); {
	case !counter.ShouldUseResults, id == "":
		return ""
	default:
		return fmt.Sprintf(` id="%s"`, id)
	}
}

func launchpadToHTML(lp blocks.LaunchPad, counter *blocks.IDCounter) string {
	lp.ColorRockets()
	var ret string
	for _, rocket := range lp.Rockets {
		ret += BlockToHTML(rocket, counter)
	}
	return fmt.Sprintf(`<ul class="launchpad"%s>%s
</ul>`, idAttribute(lp, counter), ret)
}

func imgEntryToHTML(entry blocks.ImgEntry, counter *blocks.IDCounter) string {
	var ret string
	if entry.Srclink.IsBlueLink() {
		ret += fmt.Sprintf(
			`<a href="%s"><img src="%s" %s %s></a>`,
			entry.Srclink.Href(),
			entry.Srclink.ImgSrc(),
			entry.SizeWAsAttr(),
			entry.SizeHAsAttr())
	} else {
		ret += fmt.Sprintf(
			`<a class="%s" href="%s">Hypha <i>%s</i> does not exist</a>`,
			entry.Srclink.Classes(),
			entry.Srclink.Href(),
			entry.Srclink.TargetHypha())
	}
	return fmt.Sprintf(`<figure class="img-gallery__entry">
	%s
	<figcaption>%s</figcaption>
</figure>
`, ret, BlockToHTML(parser.MakeFormatted(entry.Desc.String(), entry.HyphaName), counter))
}

func imgToHTML(img blocks.Img, counter *blocks.IDCounter) string {
	img.MarkExistenceOfSrcLinks()
	var ret string
	for _, entry := range img.Entries {
		ret += BlockToHTML(entry, counter)
	}
	return fmt.Sprintf(`<section class="img-gallery %s"%s>
%s</section>`,
		func() string {
			if img.HasOneImage() {
				return "img-gallery_one-image"
			}
			return "img-gallery_many-images"
		}(),
		idAttribute(img, counter),
		ret)
}
