package parser

import (
	"bytes"
	"context"
	"html"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
)

func isPrefixedBy(ctx context.Context, s string) bool { // This function has bugs in it!
	return bytes.HasPrefix(inputFrom(ctx).Bytes(), []byte(s))
}

func nextLaunchPad(ctx context.Context) (blocks.LaunchPad, bool) {
	var (
		hyphaName = hyphaNameFrom(ctx)
		launchPad = blocks.MakeLaunchPad()
		line      string
		done      bool
	)
	for isPrefixedBy(ctx, "=>") {
		line, done = nextLine(ctx)
		launchPad.AddRocket(blocks.MakeRocketLink(line, hyphaName))
	}
	return launchPad, done
}

func nextImg(ctx context.Context) (img blocks.Img, done bool) {
	var b byte
	line, done := nextLine(ctx)
	img, imgDone := blocks.MakeImg(line, hyphaNameFrom(ctx))
	if imgDone {
		return img, done
	}

	for !imgDone {
		b, done = nextByte(ctx)
		imgDone = img.ProcessRune(rune(b))
	}

	defer nextLine(ctx) // Characters after the final } of img are ignored.
	return img, done
}

func nextCodeBlock(ctx context.Context) (code blocks.CodeBlock, done bool) {
	line, done := nextLine(ctx)
	code = blocks.MakeCodeBlock(strings.TrimPrefix(line, "```"), "")

	for {
		line, done = nextLine(ctx)
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
	line, done := nextLine(ctx)
	t = blocks.TableFromFirstLine(line, hyphaNameFrom(ctx))
	for {
		line, done = nextLine(ctx)
		if t.ProcessLine(line) {
			break
		}
	}
	return t, done
}

func nextParagraph(ctx context.Context) (p blocks.Paragraph, done bool) {
	line, done := nextLine(ctx)
	p = blocks.Paragraph{blocks.MakeFormatted(line, hyphaNameFrom(ctx))}
	if nextLineIsSomething(ctx) {
		return
	}
	for {
		line, done = nextLine(ctx)
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
	return emptyLine(ctx) || blocks.MatchesImg(inputFrom(ctx).String()) || blocks.MatchesTable(inputFrom(ctx).String())
}

func emptyLine(ctx context.Context) bool {
	for _, b := range inputFrom(ctx).Bytes() {
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
		_, done := nextLine(ctx)
		return nil, done
	case looksLikeList(ctx):
		//case isPrefixedBy(ctx, "* "), isPrefixedBy(ctx, "*. "), isPrefixedBy(ctx, "*v "), isPrefixedBy(ctx, "*x "): — alternative way?
		return nextList(ctx)
	case isPrefixedBy(ctx, "```"):
		return nextCodeBlock(ctx)
	case isPrefixedBy(ctx, "=>"):
		return nextLaunchPad(ctx)
	case isPrefixedBy(ctx, "<="):
		line, done := nextLine(ctx)
		return blocks.MakeTransclusion(line, hyphaNameFrom(ctx)), done
	case isPrefixedBy(ctx, "----"):
		line, done := nextLine(ctx)
		return blocks.MakeHorizontalLine(line), done

	case isPrefixedBy(ctx, "###### "):
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 6), done
	case isPrefixedBy(ctx, "##### "):
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 5), done
	case isPrefixedBy(ctx, "#### "):
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 4), done
	case isPrefixedBy(ctx, "### "):
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 3), done
	case isPrefixedBy(ctx, "## "):
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 2), done
	case isPrefixedBy(ctx, "# "):
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 1), done

	case isPrefixedBy(ctx, ">"): // TODO: implement proper fractal quotes
		line, done := nextLine(ctx)
		return blocks.MakeQuote(line, hyphaNameFrom(ctx)), done
	case blocks.MatchesImg(inputFrom(ctx).String()):
		return nextImg(ctx)
	case blocks.MatchesTable(inputFrom(ctx).String()):
		return nextTable(ctx)
	}
	return nextParagraph(ctx)
}
