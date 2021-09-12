package parser

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/v2/blocks"
	"github.com/bouncepaw/mycomarkup/v2/links"
	"github.com/bouncepaw/mycomarkup/v2/mycocontext"
	"regexp"
	"strings"
)

var imgRe = regexp.MustCompile(`^img\s*[a-z\s]*{`)

func matchesImg(ctx mycocontext.Context) bool {
	return imgRe.Match(ctx.Input().Bytes())
}

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

func nextImg(ctx mycocontext.Context) (img blocks.Img, eof bool) {
	img = parseImgUntilCurlyBrace(ctx)
	var (
		r       rune
		imgDone bool
	)
	for !imgDone && !eof {
		r, eof = mycocontext.NextRune(ctx)
		imgDone = processImgRune(&img, r)
	}

	defer mycocontext.NextLine(ctx) // Characters after the final } of img are ignored.
	return img, eof
}

// Call this function if and only if matchesImg(ctx) == true.
func parseImgUntilCurlyBrace(ctx mycocontext.Context) (img blocks.Img) {
	// Input:
	// img<stuff>{<rest...>

	// Read img first. Sorry for party rocking 😎
	_, _ = mycocontext.NextRune(ctx)
	_, _ = mycocontext.NextRune(ctx)
	_, _ = mycocontext.NextRune(ctx)

	var stuff strings.Builder
	for {
		// It must be safe to ignore the error as long as parseImgUntilCurlyBrace is called correctly.
		r, _ := mycocontext.NextRune(ctx)
		if r == '{' {
			break
		}
		_, _ = stuff.WriteRune(r)
	}

	// Ignore stuff for now. TODO: https://github.com/bouncepaw/mycomarkup/issues/6
	_ = stuff

	return blocks.Img{
		Entries:   make([]blocks.ImgEntry, 0),
		CurrEntry: blocks.ImgEntry{},
		HyphaName: ctx.HyphaName(),
		State:     blocks.InRoot,
	}
}
