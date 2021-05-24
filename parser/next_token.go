package parser

import (
	"bytes"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"html"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
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
	var b byte
	line, done := mycocontext.NextLine(ctx)
	img, imgDone := blocks.MakeImg(line, ctx.HyphaName())
	if imgDone {
		return img, done
	}

	for !imgDone {
		b, done = mycocontext.NextByte(ctx)
		imgDone = img.ProcessRune(rune(b))
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

func nextTable(ctx mycocontext.Context) (t blocks.Table, done bool) {
	line, done := mycocontext.NextLine(ctx)
	t = blocks.TableFromFirstLine(line, ctx.HyphaName())
	for {
		line, done = mycocontext.NextLine(ctx)
		if t.ProcessLine(line) {
			break
		}
	}
	return t, done
}

func nextParagraph(ctx mycocontext.Context) (p blocks.Paragraph, done bool) {
	line, done := mycocontext.NextLine(ctx)
	p = blocks.Paragraph{blocks.MakeFormatted(line, ctx.HyphaName())}
	if nextLineIsSomething(ctx) {
		return
	}
	for {
		line, done = mycocontext.NextLine(ctx)
		if done && line == "" {
			break
		}
		p.AddLine(line)
		if nextLineIsSomething(ctx) {
			break
		}
	}
	return
}

func nextLineIsSomething(ctx mycocontext.Context) bool {
	prefices := []string{"=>", "<=", "```", "* ", "*. ", "*v ", "*x ", "# ", "## ", "### ", "#### ", "##### ", "###### ", ">", "----"}
	for _, prefix := range prefices {
		if isPrefixedBy(ctx, prefix) {
			return true
		}
	}
	return emptyLine(ctx) || blocks.MatchesImg(ctx.Input().String()) || blocks.MatchesTable(ctx.Input().String())
}

func emptyLine(ctx mycocontext.Context) bool {
	for _, b := range ctx.Input().Bytes() {
		switch b {
		case '\n':
			return true
		case '\t', ' ':
			continue
		default:
			return false
		}
	}
	return false
}

// Lex `line` in markup and maybe return a token.
func nextToken(ctx mycocontext.Context) (interface{}, bool) {
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
	case isPrefixedBy(ctx, "<="):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeTransclusion(line, ctx.HyphaName()), done
	case isPrefixedBy(ctx, "----"):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeHorizontalLine(line), done

	case isPrefixedBy(ctx, "###### "):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeHeading(line, ctx.HyphaName(), 6), done
	case isPrefixedBy(ctx, "##### "):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeHeading(line, ctx.HyphaName(), 5), done
	case isPrefixedBy(ctx, "#### "):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeHeading(line, ctx.HyphaName(), 4), done
	case isPrefixedBy(ctx, "### "):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeHeading(line, ctx.HyphaName(), 3), done
	case isPrefixedBy(ctx, "## "):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeHeading(line, ctx.HyphaName(), 2), done
	case isPrefixedBy(ctx, "# "):
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeHeading(line, ctx.HyphaName(), 1), done

	case isPrefixedBy(ctx, ">"): // TODO: implement proper fractal quotes
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeQuote(line, ctx.HyphaName()), done
	case blocks.MatchesImg(ctx.Input().String()):
		return nextImg(ctx)
	case blocks.MatchesTable(ctx.Input().String()):
		return nextTable(ctx)
	}
	return nextParagraph(ctx)
}
