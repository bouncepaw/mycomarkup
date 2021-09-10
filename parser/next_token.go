package parser

import (
	"bytes"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"html"
	"strings"
)

func isPrefixedBy(ctx mycocontext.Context, s string) bool { // This function has bugs in it!
	return bytes.HasPrefix(ctx.Input().Bytes(), []byte(s))
}

func nextLaunchPad(ctx mycocontext.Context) (blocks.LaunchPad, bool) {
	var (
		hyphaName = ctx.HyphaName()
		launchPad = blocks.MakeLaunchPad()
		line      string
		done      bool
	)
	for isPrefixedBy(ctx, "=>") {
		line, done = mycocontext.NextLine(ctx)
		launchPad.AddRocket(blocks.MakeRocketLink(line, hyphaName))
	}
	return launchPad, done
}

func nextImg(ctx mycocontext.Context) (img blocks.Img, done bool) {
	var r rune
	line, done := mycocontext.NextLine(ctx)
	img, imgDone := ParseImgFirstLine(line, ctx.HyphaName())
	if imgDone {
		return img, done
	}

	for !imgDone && !done {
		r, done = mycocontext.NextRune(ctx)
		imgDone = ProcessImgRune(&img, r)
	}

	defer mycocontext.NextLine(ctx) // Characters after the final } of img are ignored.
	return img, done
}

func nextCodeBlock(ctx mycocontext.Context) (code blocks.CodeBlock, done bool) {
	line, done := mycocontext.NextLine(ctx)
	code = blocks.MakeCodeBlock(strings.TrimPrefix(line, "```"), "")

	for {
		line, done = mycocontext.NextLine(ctx)
		switch {
		case strings.HasPrefix(line, "```"):
			return code, done
		default:
			code.AddLine(html.EscapeString(line))
		}
		if done {
			return code, done
		}
	}
}

func linesForQuote(ctx mycocontext.Context) ([]string, bool) {
	var (
		line  string
		lines []string
		done  bool
	)
	for {
		line, done = mycocontext.NextLine(ctx)
		// Drop >, remove spaces, save this line
		lines = append(lines, strings.TrimSpace(line[1:]))

		// If the next line is not part of the same quote, we break.
		if !isPrefixedBy(ctx, ">") {
			break
		}
	}
	return lines, done
}

func nextQuote(ctx mycocontext.Context) (blocks.Quote, bool) {
	var (
		quote       = blocks.Quote{}
		lines, done = linesForQuote(ctx)
		innerText   bytes.Buffer
	)

	for i, line := range lines {
		if i > 0 {
			innerText.WriteRune('\n')
		}
		innerText.WriteString(line)
	}

	parseSubdocumentForEachBlock(ctx, &innerText, func(block blocks.Block) {
		quote.AddBlock(block)
	})

	return quote, done
}

func nextLineIsSomething(ctx mycocontext.Context) bool {
	prefices := []string{"=>", "<=", "```", "* ", "*. ", "*v ", "*x ", "# ", "## ", "### ", "#### ", "##### ", "###### ", ">", "----"}
	for _, prefix := range prefices {
		if isPrefixedBy(ctx, prefix) {
			return true
		}
	}
	return emptyLine(ctx) || blocks.MatchesImg(ctx.Input().String()) || matchesTable(ctx)
}

func emptyLine(ctx mycocontext.Context) bool {
	for _, b := range ctx.Input().Bytes() {
		switch b {
		case '\n':
			return true
		case '\t', ' ', '\r':
			continue
		default:
			return false
		}
	}
	return true
}

func nextToken(ctx mycocontext.Context) (blocks.Block, bool) {
	switch {
	case emptyLine(ctx):
		_, done := mycocontext.NextLine(ctx)
		return nil, done
	case looksLikeList(ctx):
		//case isPrefixedBy(ctx, "* "), isPrefixedBy(ctx, "*. "), isPrefixedBy(ctx, "*v "), isPrefixedBy(ctx, "*x "): — alternative way?
		return nextList(ctx)
	case isPrefixedBy(ctx, "```"):
		return nextCodeBlock(ctx)
	case isPrefixedBy(ctx, "=>"):
		return nextLaunchPad(ctx)
	case isPrefixedBy(ctx, ">"):
		return nextQuote(ctx)
	case isPrefixedBy(ctx, "<="):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeTransclusion(line, ctx.HyphaName()), done
	case isPrefixedBy(ctx, "----"):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeHorizontalLine(line), done

	case isPrefixedBy(ctx, "###### "):
		line, done := mycocontext.NextLine(ctx)
		return MakeHeading(line, ctx.HyphaName(), 6), done
	case isPrefixedBy(ctx, "##### "):
		line, done := mycocontext.NextLine(ctx)
		return MakeHeading(line, ctx.HyphaName(), 5), done
	case isPrefixedBy(ctx, "#### "):
		line, done := mycocontext.NextLine(ctx)
		return MakeHeading(line, ctx.HyphaName(), 4), done
	case isPrefixedBy(ctx, "### "):
		line, done := mycocontext.NextLine(ctx)
		return MakeHeading(line, ctx.HyphaName(), 3), done
	case isPrefixedBy(ctx, "## "):
		line, done := mycocontext.NextLine(ctx)
		return MakeHeading(line, ctx.HyphaName(), 2), done
	case isPrefixedBy(ctx, "# "):
		line, done := mycocontext.NextLine(ctx)
		return MakeHeading(line, ctx.HyphaName(), 1), done

	case blocks.MatchesImg(ctx.Input().String()):
		return nextImg(ctx)
	case matchesTable(ctx):
		return nextTable(ctx)
	}
	return nextParagraph(ctx)
}
