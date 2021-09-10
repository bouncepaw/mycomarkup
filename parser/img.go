package parser

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/links"
	"strings"
)

// pushImgEntry pushes the most recent entry of Img to Img.Entries and creates a new entry. What an ugly function!
func pushImgEntry(img *blocks.Img) {
	if strings.TrimSpace(img.CurrEntry.Path.String()) != "" {
		img.CurrEntry.Srclink = links.From(img.CurrEntry.Path.String(), "", img.HyphaName)
		// img.currEntry.Srclink.DoubtExistence()
		img.Entries = append(img.Entries, img.CurrEntry)
		img.CurrEntry = blocks.ImgEntry{HyphaName: img.HyphaName}
		img.CurrEntry.Path.Reset()
	}
}

// processImgLine parses the line and tells if the gallery is finished.
func processImgLine(img *blocks.Img, line string) (imgFinished bool) {
	for _, r := range line {
		if shouldReturnTrue := processImgRune(img, r); shouldReturnTrue {
			return true
		}
	}
	// We do that because \n are not part of line we receive as the argument:
	processImgRune(img, '\n')

	return false
}

// processImgRune parses the rune.
func processImgRune(img *blocks.Img, r rune) (done bool) {
	// TODO: move to the parser module.
	if r == '\r' {
		return false
	}
	switch img.State {
	case blocks.InRoot:
		return processInRoot(img, r)
	case blocks.InName:
		return processInName(img, r)
	case blocks.InDimensionsW:
		return processInDimensionsW(img, r)
	case blocks.InDimensionsH:
		return processInDimensionsH(img, r)
	case blocks.InDescription:
		return processInDescription(img, r)
	}
	fmt.Println("processImgRune: unreachable state", r)
	return true
}

func processInDescription(img *blocks.Img, r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.State = blocks.InName
	default:
		img.CurrEntry.Desc.WriteRune(r)
	}
	return false
}

func processInRoot(img *blocks.Img, r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		pushImgEntry(img)
		return true
	case '\n':
		pushImgEntry(img)
	case ' ', '\t':
	default:
		img.State = blocks.InName
		img.CurrEntry = blocks.ImgEntry{}
		img.CurrEntry.Path.Reset()
		img.CurrEntry.Path.WriteRune(r)
	}
	return false
}

func processInName(img *blocks.Img, r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		pushImgEntry(img)
		return true
	case '|':
		img.State = blocks.InDimensionsW
	case '{':
		img.State = blocks.InDescription
	case '\n':
		pushImgEntry(img)
		img.State = blocks.InRoot
	default:
		img.CurrEntry.Path.WriteRune(r)
	}
	return false
}

func processInDimensionsW(img *blocks.Img, r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		pushImgEntry(img)
		return true
	case '*':
		img.State = blocks.InDimensionsH
	case ' ', '\t', '\n':
	case '{':
		img.State = blocks.InDescription
	default:
		img.CurrEntry.SizeW.WriteRune(r)
	}
	return false
}

func processInDimensionsH(img *blocks.Img, r rune) (imgFinished bool) {
	switch r {
	case '}':
		pushImgEntry(img)
		return true
	case ' ', '\t', '\n':
	case '{':
		img.State = blocks.InDescription
	default:
		img.CurrEntry.SizeH.WriteRune(r)
	}
	return false
}

// parseImgFirstLine parses the image gallery on the line and returns it. It also tells if the gallery is finished or not.
func parseImgFirstLine(line, hyphaName string) (img blocks.Img, imgFinished bool) {
	img = blocks.Img{
		HyphaName: hyphaName,
		Entries:   make([]blocks.ImgEntry, 0),
	}
	line = line[strings.IndexRune(line, '{')+1:]
	return img, processImgLine(&img, line)
}
