package tools

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bouncepaw/mycomarkup/v2"
	"github.com/bouncepaw/mycomarkup/v2/blocks"
	"github.com/bouncepaw/mycomarkup/v2/mycocontext"
	"github.com/bouncepaw/mycomarkup/v2/util"
)

// OpenGraphVisitors returns visitors you should pass to BlockTree. They will figure out what should go to the final opengraph. Call resultHTML to get that result.
//
// Description is the first root paragraph of the document. If there is no such paragraph, the description is empty string.
func OpenGraphVisitors(ctx mycocontext.Context) (
	resultHTML func() string,
	descVisitor func(blocks.Block),
	imgVisitor func(blocks.Block),
) {
	var (
		imageUrl    = "/favicon.ico"
		description = ""

		foundImg              = false
		foundSomethingTextual = false // Let's have at least something if there is no paragraph.
		foundProperParagraph  = false
	)

	return func() string {
			return strings.Join([]string{
				ogTag("title", util.BeautifulName(ctx.HyphaName())),
				ogTag("type", "article"), // TODO: change depending on content?
				ogTag("image", imageUrl),
				ogTag("url", ctx.WebSiteURL()+"/hypha/"+util.CanonicalName(ctx.HyphaName())),
				ogTag("determiner", ""),
				ogTag("description", htmlTagRe.ReplaceAllString(description, "")),
			}, "\n")
		}, func(block blocks.Block) {
			if foundProperParagraph { // Won't find anything better.
				return
			}
			switch block := block.(type) {
			case blocks.Paragraph:
				foundSomethingTextual, foundProperParagraph = true, true
				description = mycomarkup.BlockToHTML(block, &blocks.IDCounter{ShouldUseResults: false})
			case blocks.Heading, blocks.CodeBlock: // These two seem alright. Primitive enough.
				if !foundSomethingTextual {
					foundSomethingTextual = true
					description = mycomarkup.BlockToHTML(block, &blocks.IDCounter{ShouldUseResults: false})
				}
			}
		}, func(block blocks.Block) {
			if foundImg { // No need for a second image
				return
			}
			switch block := block.(type) {
			case blocks.Img:
				if len(block.Entries) > 0 {
					imageUrl = ctx.WebSiteURL() + block.Entries[0].Srclink.ImgSrc()
				}
			}
		}
}

// Used to clear opengraph description from html tags. This method is usually bad because of dangers of malformed HTML, but I'm going to use it only for Mycorrhiza-generated HTML, so it's okay. The question mark is required; without it the whole string is eaten away.
var htmlTagRe = regexp.MustCompile(`<.*?>`)

func ogTag(property, content string) string {
	return fmt.Sprintf(`<meta property="og:%s" content="%s"/>`, property, content)
}
