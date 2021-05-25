package mycomarkup

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/blocks"
)

// BlockToHTML turns the given block into HTML. It supports only a subset of Mycomarkup.
func BlockToHTML(block blocks.Block) string {
	switch b := block.(type) {
	case blocks.Formatted:
		return b.Html
	case blocks.Paragraph:
		return fmt.Sprintf("\n<p>%s</p>", b.Html)
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
	<li class="launchpad__entry"><a href="%s" class="rocketlink %s">%s</a></li>`, b.Href(), b.Classes(), b.Display()) // FIXME: there seems to be a bug, run test1
	case blocks.Heading:
		return fmt.Sprintf(`
<h%[1]d>%[2]s<a href="#%[3]s" id="%[3]s" class="heading__link"></a></h%[1]d>
`, b.Level, BlockToHTML(b.Contents()), b.ID())
	case blocks.Table:
		return tableToHTML(b)
	case blocks.TableRow:
		return tableRowToHTML(b)
	case blocks.TableCell:
		return tableCellToHTML(b)
	case blocks.CodeBlock:
		return fmt.Sprintf("\n<pre class='codeblock'><code class='language-%s'>%s</code></pre>", b.Language(), b.Contents())
	case blocks.Quote:
		return fmt.Sprintf("\n<blockquote>%s</blockquote>", b.Contents())
	}
	fmt.Printf("%q\n", block)
	return "<b>UNKNOWN ELEMENT</b>"
}

func tableCellToHTML(tc blocks.TableCell) string {
	return fmt.Sprintf(
		"\n\t<%[1]s%[2]s>%[3]s</%[1]s>",
		tc.TagName(),
		tc.ColspanAttribute(),
		BlockToHTML(tc.Contents),
	)
}

func tableRowToHTML(tr blocks.TableRow) string {
	var ret string
	for _, tc := range tr.Cells {
		ret += BlockToHTML(*tc)
	}
	return fmt.Sprintf("<tr>%s</tr>", ret)
}

func tableToHTML(t blocks.Table) string {
	var ret string
	if t.Caption != "" {
		ret = fmt.Sprintf("<caption>%s</caption>", t.Caption)
	}
	if len(t.Rows) > 0 && t.Rows[0].LooksLikeThead() {
		ret += fmt.Sprintf("<thead>%s</thead>", BlockToHTML(*t.Rows[0]))
		t.Rows = t.Rows[1:]
	}
	ret += "\n<tbody>\n"
	for _, tr := range t.Rows {
		ret += BlockToHTML(*tr)
	}
	return fmt.Sprintf(`
<table>%s</tbody></table>`, ret)
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
	if entry.Srclink.Exists() {
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
			entry.Srclink.Address())
	}
	return fmt.Sprintf(`<figure class="img-gallery__entry">
	%s
	<figcaption>%s</figcaption>
</figure>
`, ret, BlockToHTML(entry.Description()))
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
