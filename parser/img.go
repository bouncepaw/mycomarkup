package parser

import (
	"git.sr.ht/~bouncepaw/mycomarkup/v5/blocks"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/links"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"regexp"
	"strings"
)

var imgRe = regexp.MustCompile(`^img\s*[a-z\s]*{`)

func matchesImg(ctx mycocontext.Context) bool {
	return imgRe.Match(ctx.Input().Bytes())
}

func nextImg(ctx mycocontext.Context) (img blocks.Img, eof bool) {
	img = parseImgUntilCurlyBrace(ctx)
	for {
		imgEntry, found, imgFinished := nextImgEntry(ctx)
		if found {
			img.Entries = append(img.Entries, imgEntry)
		}
		if imgFinished {
			break
		}
	}

	return img, mycocontext.IsEof(ctx)
}

type imgEntryParsingState int

const (
	imgEntryOnStart imgEntryParsingState = iota
	imgEntryCollectingTarget
	imgEntryCollectingDimensionWidth
	imgEntryCollectingDimensionHeight
)

func nextImgEntryDescription(ctx mycocontext.Context) string {
	var (
		r               rune
		eof             bool
		curlyBracesOpen = 0
		res             strings.Builder
	)
	for {
		r, eof = mycocontext.NextRune(ctx)
		if eof {
			return res.String()
		}

		switch r {
		case '{':
			curlyBracesOpen++
			res.WriteRune('{')
		case '}':
			if curlyBracesOpen == 0 {
				return res.String()
			}
			if curlyBracesOpen > 0 {
				curlyBracesOpen--
			}
			res.WriteRune('{')
		default: // Including \n
			res.WriteRune(r)
		}
	}
}

func nextImgEntry(ctx mycocontext.Context) (
	imgEntry blocks.ImgEntry,
	entryFound bool, // true if an entry was found
	imgDone bool, // true if final img } found
) {
	var (
		r     rune
		eof   bool
		state = imgEntryOnStart

		target, width, height strings.Builder
		descbuf               string
	)
	entryFound = true

runewalker:
	for {
		r, eof = mycocontext.NextRune(ctx)
		if eof {
			imgDone = true // Just to be sure
			break
		}

	runechecker: // TODO: add escaping \
		switch state {
		case imgEntryOnStart:
			switch r {
			case '}':
				entryFound, imgDone = false, true
				_, _ = mycocontext.NextLine(ctx) // After closing }
				break runewalker
			case '\n':
				entryFound, imgDone = false, false
				break runewalker
			case ' ', '\t': // Ignore the leading whitespace
			case '|': // Empty target, so it seems. This entry becomes invalid.
				entryFound = false
				state = imgEntryCollectingDimensionWidth
			default:
				state = imgEntryCollectingTarget
				goto runechecker // uwu
			}
		case imgEntryCollectingTarget:
			switch r {
			case '}':
				entryFound, imgDone = true, true
				break runewalker
			case '\n':
				entryFound, imgDone = true, false
				break runewalker
			case '|':
				state = imgEntryCollectingDimensionWidth
			case '{':
				descbuf = nextImgEntryDescription(ctx)
				break runewalker
			default:
				// I am confident in myself, thus I ignore errors
				_, _ = target.WriteRune(r)
			}
		case imgEntryCollectingDimensionWidth:
			switch r {
			case '}':
				entryFound, imgDone = true, true
				break runewalker
			case '\n':
				entryFound, imgDone = true, false
				break runewalker
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				_, _ = width.WriteRune(r)
			case '*':
				state = imgEntryCollectingDimensionHeight
			case '{':
				descbuf = nextImgEntryDescription(ctx)
				break runewalker
			default: // Ignore the garbage!
			}
		case imgEntryCollectingDimensionHeight:
			switch r {
			case '}':
				entryFound, imgDone = true, true
				break runewalker
			case '\n':
				entryFound, imgDone = true, false
				break runewalker
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				_, _ = height.WriteRune(r)
			case '{':
				descbuf = nextImgEntryDescription(ctx)
				break runewalker
			default: // Ignore the garbage!
			}
		default:
			panic("warning warning warning!!!!!!!!")
		}
	}

	return blocks.NewImgEntry(
			links.LinkFrom(ctx, target.String(), ""),
			ctx.HyphaName(),
			width.String(),
			height.String(),
			descbuf),
		entryFound,
		imgDone
}

// Call this function if and only if matchesImg(ctx) == true.
func parseImgUntilCurlyBrace(ctx mycocontext.Context) (img blocks.Img) {
	// Input:
	// img<attr>{<rest...>

	// Read img first. Sorry for party rocking 😎
	_, _ = mycocontext.NextRune(ctx)
	_, _ = mycocontext.NextRune(ctx)
	_, _ = mycocontext.NextRune(ctx)

	var attr strings.Builder
	for {
		// It must be safe to ignore the error as long as parseImgUntilCurlyBrace is called correctly.
		r, _ := mycocontext.NextRune(ctx)
		if r == '{' {
			break
		}
		_, _ = attr.WriteRune(r)
	}

	return blocks.NewImg(make([]blocks.ImgEntry, 0), parseImgLayout(attr.String()))
}

func parseImgLayout(attr string) blocks.ImgLayout {
	switch {
	case strings.Contains(attr, "side"):
		return blocks.ImgLayoutSide
	case strings.Contains(attr, "grid"):
		return blocks.ImgLayoutGrid
	default:
		return blocks.ImgLayoutNormal
	}
}
