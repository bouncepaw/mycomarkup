package parser

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/links"
	"strings"
)

// PushImgEntry pushes the most recent entry of Img to Img.Entries and creates a new entry. What an ugly function!
func PushImgEntry(img *blocks.Img) {
	if strings.TrimSpace(img.CurrEntry.Path.String()) != "" {
		img.CurrEntry.Srclink = links.From(img.CurrEntry.Path.String(), "", img.HyphaName)
		// img.currEntry.Srclink.DoubtExistence()
		img.Entries = append(img.Entries, img.CurrEntry)
		img.CurrEntry = blocks.ImgEntry{HyphaName: img.HyphaName}
		img.CurrEntry.Path.Reset()
	}
}

// ProcessImgLine parses the line and tells if the gallery is finished.
func ProcessImgLine(img *blocks.Img, line string) (imgFinished bool) {
	for _, r := range line {
		if shouldReturnTrue := ProcessImgRune(img, r); shouldReturnTrue {
			return true
		}
	}
	// We do that because \n are not part of line we receive as the argument:
	ProcessImgRune(img, '\n')

	return false
}

// ProcessImgRune parses the rune.
func ProcessImgRune(img *blocks.Img, r rune) (done bool) {
	// TODO: move to the parser module.
	if r == '\r' {
		return false
	}
	switch img.State {
	case blocks.InRoot:
		return ProcessInRoot(img, r)
	case blocks.InName:
		return ProcessInName(img, r)
	case blocks.InDimensionsW:
		return ProcessInDimensionsW(img, r)
	case blocks.InDimensionsH:
		return ProcessInDimensionsH(img, r)
	case blocks.InDescription:
		return ProcessInDescription(img, r)
	}
	fmt.Println("ProcessImgRune: unreachable state", r)
	return true
}

func ProcessInDescription(img *blocks.Img, r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		img.State = blocks.InName
	default:
		img.CurrEntry.Desc.WriteRune(r)
	}
	return false
}

func ProcessInRoot(img *blocks.Img, r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		PushImgEntry(img)
		return true
	case '\n':
		PushImgEntry(img)
	case ' ', '\t':
	default:
		img.State = blocks.InName
		img.CurrEntry = blocks.ImgEntry{}
		img.CurrEntry.Path.Reset()
		img.CurrEntry.Path.WriteRune(r)
	}
	return false
}

func ProcessInName(img *blocks.Img, r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		PushImgEntry(img)
		return true
	case '|':
		img.State = blocks.InDimensionsW
	case '{':
		img.State = blocks.InDescription
	case '\n':
		PushImgEntry(img)
		img.State = blocks.InRoot
	default:
		img.CurrEntry.Path.WriteRune(r)
	}
	return false
}

func ProcessInDimensionsW(img *blocks.Img, r rune) (shouldReturnTrue bool) {
	switch r {
	case '}':
		PushImgEntry(img)
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

func ProcessInDimensionsH(img *blocks.Img, r rune) (imgFinished bool) {
	switch r {
	case '}':
		PushImgEntry(img)
		return true
	case ' ', '\t', '\n':
	case '{':
		img.State = blocks.InDescription
	default:
		img.CurrEntry.SizeH.WriteRune(r)
	}
	return false
}

// ParseImgFirstLine parses the image gallery on the line and returns it. It also tells if the gallery is finished or not.
func ParseImgFirstLine(line, hyphaName string) (img blocks.Img, imgFinished bool) {
	img = blocks.Img{
		HyphaName: hyphaName,
		Entries:   make([]blocks.ImgEntry, 0),
	}
	line = line[strings.IndexRune(line, '{')+1:]
	return img, ProcessImgLine(&img, line)
}
