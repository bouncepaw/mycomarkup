// This is not done yet
package doc

import (
	"github.com/bouncepaw/mycomarkup/parser"
)

// TODO: remove this abomination
var cfg = struct {
	URL string
}{
	URL: "PLACEHOLDER",
}

// MycoDoc is a mycomarkup-formatted document.
type MycoDoc struct {
	// data
	HyphaName string
	Contents  string
	// indicators
	parsedAlready bool
}

// Doc returns a mycomarkup document with the given name and mycomarkup-formatted contents.
//
// The returned document is not lexed not parsed yet. You have to do that separately.
func Doc(hyphaName, contents string) *MycoDoc {
	md := &MycoDoc{
		HyphaName: hyphaName,
		Contents:  contents,
	}
	return md
}

func (md *MycoDoc) Lex() []interface{} {
	var (
		ctx, _ = parser.ContextFromStringInput(md.HyphaName, md.Contents)
		state  = parser.ParserState{Name: md.HyphaName}
		ast    = []interface{}{}
	)

	for {
		line, done := parser.NextLine(ctx)
		token := parser.LineToToken(line, &state)
		if token != nil {
			ast = append(ast, token)
		}
		if done {
			break
		}
	}

	return ast
}

// AsHTML returns an html representation of the document
func (md *MycoDoc) AsHTML() string {
	return GenerateHTML(md.Lex(), 0)
}

// AsGemtext returns a gemtext representation of the document. Currently really limited, just returns source text
func (md *MycoDoc) AsGemtext() string {
	return md.Contents
}

/*
/// The rest of the file is OpenGraph-related.

// Used to clear opengraph description from html tags. This method is usually bad because of dangers of malformed HTML, but I'm going to use it only for Mycorrhiza-generated HTML, so it's okay. The question mark is required; without it the whole string is eaten away.
var htmlTagRe = regexp.MustCompile(`<.*?>`)

// OpenGraphHTML returns an html representation of og: meta tags.
func (md *MycoDoc) OpenGraphHTML() string {
	ogImage, ogDescription := md.openGraphImageAndDescription()
	return strings.Join([]string{
		ogTag("title", util.BeautifulName(md.HyphaName)),
		ogTag("type", "article"),
		ogTag("image", ogImage),
		ogTag("url", cfg.URL+"/hypha/"+md.HyphaName),
		ogTag("determiner", ""),
		ogTag("description", ogDescription),
	}, "\n")
}

// return image and description of the document for including in open graph.
func (md *MycoDoc) openGraphImageAndDescription() (ogImage, ogDescription string) {
	ogImage = cfg.URL + "/favicon.ico"
	foundDesc := false
	foundImg := false
	for _, line := range md.ast {
		switch v := line.Value.(type) {
		case string:
			if !foundDesc {
				ogDescription = v
				foundDesc = true
			}
		case blocks.Img:
			if !foundImg && len(v.Entries) > 0 {
				ogImage = v.Entries[0].Srclink.ImgSrc()
				if !v.Entries[0].Srclink.OfKind(links.LinkExternal) {
					ogImage = cfg.URL + ogImage
				}
				foundImg = true
			}
		}
	}
	return ogImage, htmlTagRe.ReplaceAllString(ogDescription, "")
}

func ogTag(property, content string) string {
	return fmt.Sprintf(`<meta property="og:%s" content="%s"/>`, property, content)
}
*/
