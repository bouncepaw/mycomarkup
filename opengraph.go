package mycomarkup

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/generator"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"regexp"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/util"
)

// Used to clear opengraph description from html tags. This method is usually bad because of dangers of malformed HTML, but I'm going to use it only for Mycorrhiza-generated HTML, so it's okay. The question mark is required; without it the whole string is eaten away.
var htmlTagRe = regexp.MustCompile(`<.*?>`)

// OpenGraphHTML returns an html representation of og: meta tags.
func OpenGraphHTML(ctx mycocontext.Context, ast []interface{}) string {
	ogImage, ogDescription := openGraphImageAndDescription(ast)
	return strings.Join([]string{
		ogTag("title", util.BeautifulName(ctx.HyphaName())),
		ogTag("type", "article"),
		ogTag("image", ogImage),
		// TODO: there should be a full URL ⤵︎. Requires a different API for the lib.
		ogTag("url", "/hypha/"+util.BeautifulName(ctx.HyphaName())),
		ogTag("determiner", ""),
		ogTag("description", ogDescription),
	}, "\n")
}

// return image and description of the document for including in open graph.
func openGraphImageAndDescription(ast []interface{}) (ogImage, ogDescription string) {
	// TODO: there should be a full URL ⤵︎
	ogImage = "/favicon.ico"
	foundDesc := false
	foundImg := false
	for _, block := range ast {
		switch v := block.(type) {
		case blocks.Paragraph:
			if !foundDesc {
				ogDescription = strings.TrimSpace(htmlTagRe.ReplaceAllString(generator.BlockToHTML(v), ""))
				foundDesc = true
			}
		case blocks.Img:
			if !foundImg && len(v.Entries) > 0 {
				ogImage = v.Entries[0].Srclink.ImgSrc()
				/*if !v.Entries[0].Srclink.OfKind(links.LinkExternal) {
					// TODO: there should be a full URL ⤵︎
					ogImage = doc.cfg.URL + ogImage
				}*/
				foundImg = true
			}
		}
	}
	return ogImage, ogDescription
}

func ogTag(property, content string) string {
	return fmt.Sprintf(`<meta property="og:%s" content="%s"/>`, property, content)
}
