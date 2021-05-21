package parser

import (
	"bytes"
	"context"
	"github.com/bouncepaw/mycomarkup/util"
	"html"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
)

func isPrefixedBy(ctx context.Context, s string) bool { // This function has bugs in it!
	return bytes.HasPrefix(util.InputFrom(ctx).Bytes(), []byte(s))
}

func nextLaunchPad(ctx context.Context) (blocks.LaunchPad, bool) {
	var (
		hyphaName = util.HyphaNameFrom(ctx)
		launchPad = blocks.MakeLaunchPad()
		line      string
		done      bool
	)
	for isPrefixedBy(ctx, "=>") {
		line, done = util.NextLine(ctx)
		launchPad.AddRocket(blocks.MakeRocketLink(line, hyphaName))
	}
	return launchPad, done
}

func nextImg(ctx context.Context) (img blocks.Img, done bool) {
	var b byte
	line, done := util.NextLine(ctx)
	img, imgDone := blocks.MakeImg(line, util.HyphaNameFrom(ctx))
	if imgDone {
		return img, done
	}

	for !imgDone {
		b, done = util.NextByte(ctx)
		imgDone = img.ProcessRune(rune(b))
	}

	defer util.NextLine(ctx) // Characters after the final } of img are ignored.
	return img, done
}

func nextCodeBlock(ctx context.Context) (code blocks.CodeBlock, done bool) {
	line, done := util.NextLine(ctx)
	code = blocks.MakeCodeBlock(strings.TrimPrefix(line, "```"), "")

	for {
		line, done = util.NextLine(ctx)
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

func nextTable(ctx context.Context) (t blocks.Table, done bool) {
	line, done := util.NextLine(ctx)
	t = blocks.TableFromFirstLine(line, util.HyphaNameFrom(ctx))
	for {
		line, done = util.NextLine(ctx)
		if t.ProcessLine(line) {
			break
		}
	}
	return t, done
}

func nextParagraph(ctx context.Context) (p blocks.Paragraph, done bool) {
	line, done := util.NextLine(ctx)
	p = blocks.Paragraph{blocks.MakeFormatted(line, util.HyphaNameFrom(ctx))}
	if nextLineIsSomething(ctx) {
		return
	}
	for {
		line, done = util.NextLine(ctx)
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

func nextLineIsSomething(ctx context.Context) bool {
	prefices := []string{"=>", "<=", "```", "* ", "*. ", "*v ", "*x ", "# ", "## ", "### ", "#### ", "##### ", "###### ", ">", "----"}
	for _, prefix := range prefices {
		if isPrefixedBy(ctx, prefix) {
			return true
		}
	}
	return emptyLine(ctx) || blocks.MatchesImg(util.InputFrom(ctx).String()) || blocks.MatchesTable(util.InputFrom(ctx).String())
}

func emptyLine(ctx context.Context) bool {
	for _, b := range util.InputFrom(ctx).Bytes() {
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
func nextToken(ctx context.Context) (interface{}, bool) {
	switch {
	case emptyLine(ctx):
		_, done := util.NextLine(ctx)
		return nil, done
	case looksLikeList(ctx):
		//case isPrefixedBy(ctx, "* "), isPrefixedBy(ctx, "*. "), isPrefixedBy(ctx, "*v "), isPrefixedBy(ctx, "*x "): — alternative way?
		return nextList(ctx)
	case isPrefixedBy(ctx, "```"):
		return nextCodeBlock(ctx)
	case isPrefixedBy(ctx, "=>"):
		return nextLaunchPad(ctx)
	case isPrefixedBy(ctx, "<="):
		line, done := util.NextLine(ctx)
		return blocks.MakeTransclusion(line, util.HyphaNameFrom(ctx)), done
	case isPrefixedBy(ctx, "----"):
		line, done := util.NextLine(ctx)
		return blocks.MakeHorizontalLine(line), done

	case isPrefixedBy(ctx, "###### "):
		line, done := util.NextLine(ctx)
		return blocks.MakeHeading(line, util.HyphaNameFrom(ctx), 6), done
	case isPrefixedBy(ctx, "##### "):
		line, done := util.NextLine(ctx)
		return blocks.MakeHeading(line, util.HyphaNameFrom(ctx), 5), done
	case isPrefixedBy(ctx, "#### "):
		line, done := util.NextLine(ctx)
		return blocks.MakeHeading(line, util.HyphaNameFrom(ctx), 4), done
	case isPrefixedBy(ctx, "### "):
		line, done := util.NextLine(ctx)
		return blocks.MakeHeading(line, util.HyphaNameFrom(ctx), 3), done
	case isPrefixedBy(ctx, "## "):
		line, done := util.NextLine(ctx)
		return blocks.MakeHeading(line, util.HyphaNameFrom(ctx), 2), done
	case isPrefixedBy(ctx, "# "):
		line, done := util.NextLine(ctx)
		return blocks.MakeHeading(line, util.HyphaNameFrom(ctx), 1), done

	case isPrefixedBy(ctx, ">"): // TODO: implement proper fractal quotes
		line, done := util.NextLine(ctx)
		return blocks.MakeQuote(line, util.HyphaNameFrom(ctx)), done
	case blocks.MatchesImg(util.InputFrom(ctx).String()):
		return nextImg(ctx)
	case blocks.MatchesTable(util.InputFrom(ctx).String()):
		return nextTable(ctx)
	}
	return nextParagraph(ctx)
}
