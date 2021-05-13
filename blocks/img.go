package blocks

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/globals"
	"github.com/bouncepaw/mycomarkup/util"
	"regexp"
	"strings"

	"github.com/bouncepaw/mycomarkup/links"
)

var imgRe = regexp.MustCompile(`^img\s+{`)

func MatchesImg(line string) bool {
	return imgRe.MatchString(line)
}

type imgState int

const (
	inRoot imgState = iota
	inName
	inDimensionsW
	inDimensionsH
	inDescription
)

type Img struct {
	Entries   []imgEntry
	currEntry imgEntry
	hyphaName string
	state     imgState
}

func (img *Img) pushEntry() {
	if strings.TrimSpace(img.currEntry.path.String()) != "" {
		img.currEntry.Srclink = links.From(img.currEntry.path.String(), "", img.hyphaName)
		// img.currEntry.Srclink.DoubtExistence()
		img.Entries = append(img.Entries, img.currEntry)
		img.currEntry = imgEntry{}
		img.currEntry.path.Reset()
	}
}

func (img *Img) Process(line string) (shouldGoBackToNormal bool) {
	stateToProcessor := map[imgState]func(rune) bool{
		inRoot:        img.processInRoot,
		inName:        img.processInName,
		inDimensionsW: img.processInDimensionsW,
		inDimensionsH: img.processInDimensionsH,
		inDescription: img.processInDescription,
	}
	for _, r := range line {
		if shouldReturnTrue := stateToProcessor[img.state](r); shouldReturnTrue {
			return true
		}
	}
	return false
}

func (img *Img) processInDescription(r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.state = inName
	default:
		img.currEntry.desc.WriteRune(r)
	}
	return false
}

func (img *Img) processInRoot(r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.pushEntry()
		return true
	case '\n', '\r':
		img.pushEntry()
	case ' ', '\t':
	default:
		img.state = inName
		img.currEntry = imgEntry{}
		img.currEntry.path.Reset()
		img.currEntry.path.WriteRune(r)
	}
	return false
}

func (img *Img) processInName(r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.pushEntry()
		return true
	case '|':
		img.state = inDimensionsW
	case '{':
		img.state = inDescription
	case '\n', '\r':
		img.pushEntry()
		img.state = inRoot
	default:
		img.currEntry.path.WriteRune(r)
	}
	return false
}

func (img *Img) processInDimensionsW(r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.pushEntry()
		return true
	case '*':
		img.state = inDimensionsH
	case ' ', '\t', '\n':
	case '{':
		img.state = inDescription
	default:
		img.currEntry.sizeW.WriteRune(r)
	}
	return false
}

func (img *Img) processInDimensionsH(r rune) (shouldGoBackToNormal bool) {
	switch r {
	case '}':
		img.pushEntry()
		return true
	case ' ', '\t', '\n':
	case '{':
		img.state = inDescription
	default:
		img.currEntry.sizeH.WriteRune(r)
	}
	return false
}

func ImgFromFirstLine(line, hyphaName string) (img *Img, shouldGoBackToNormal bool) {
	img = &Img{
		hyphaName: hyphaName,
		Entries:   make([]imgEntry, 0),
	}
	line = line[strings.IndexRune(line, '{')+1:]
	return img, img.Process(line)
}

func (img *Img) pagePathFor(path string) string {
	path = strings.TrimSpace(path)
	if strings.IndexRune(path, ':') != -1 || strings.IndexRune(path, '/') == 0 {
		return path
	} else {
		return "/hypha/" + util.XclCanonicalName(img.hyphaName, path)
	}
}

func parseDimensions(dimensions string) (sizeW, sizeH string) {
	xIndex := strings.IndexRune(dimensions, '*')
	if xIndex == -1 { // If no x in dimensions
		sizeW = strings.TrimSpace(dimensions)
	} else {
		sizeW = strings.TrimSpace(dimensions[:xIndex])
		sizeH = strings.TrimSpace(strings.TrimPrefix(dimensions, dimensions[:xIndex+1]))
	}
	return
}

func (img *Img) markExistenceOfSrcLinks() {
	globals.HyphaIterate(func(hn string) {
		for _, entry := range img.Entries {
			if hn == entry.Srclink.Address() {
				entry.Srclink.DestinationUnknown = false
			}
		}
	})
}

func (img *Img) ToHtml() (html string) {
	img.markExistenceOfSrcLinks()
	isOneImageOnly := len(img.Entries) == 1 && img.Entries[0].desc.Len() == 0
	if isOneImageOnly {
		html += `<section class="img-gallery img-gallery_one-image">`
	} else {
		html += `<section class="img-gallery img-gallery_many-images">`
	}

	for _, entry := range img.Entries {
		html += `<figure>`
		if entry.Srclink.DestinationUnknown {
			html += fmt.Sprintf(
				`<a class="%s" href="%s">Hypha <i>%s</i> does not exist</a>`,
				entry.Srclink.Classes(),
				entry.Srclink.Href(),
				entry.Srclink.Address)
		} else {
			html += fmt.Sprintf(
				`<a href="%s"><img src="%s" %s %s></a>`,
				entry.Srclink.Href(),
				entry.Srclink.ImgSrc(),
				entry.sizeWAsAttr(),
				entry.sizeHAsAttr())
		}
		html += entry.descriptionAsHtml(img.hyphaName)
		html += `</figure>`
	}
	return html + `</section>`
}
