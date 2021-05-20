package parser

import (
	"context"
	"html"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
)

// ParserState is used by markup parser to remember what is going on.
type ParserState struct {
	where string // "", "list", "pre", etc.
	// Temporaries
	list      *blocks.List
	paragraph *blocks.Paragraph
}

func isPrefixedBy(ctx context.Context, s string) bool {
	// TODO: make sure that String() does not make an allocation, and if it does, implement the function in a different way.
	return strings.HasPrefix(inputFrom(ctx).String(), s)
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

func nextImg(ctx context.Context, state *ParserState) (img blocks.Img, done bool) {
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

// Lex `line` in markup and maybe return a token.
func nextToken(ctx context.Context, state *ParserState) (interface{}, bool) {
	var ret interface{}
	addParagraphIfNeeded := func() { // This is a bug source, I know it.
		if state.where == "p" {
			state.where = ""
			ret = *state.paragraph
		}
	}
	switch {
	case blocks.MatchesImg(inputFrom(ctx).String()):
		addParagraphIfNeeded()
		return nextImg(ctx, state)
	case blocks.MatchesTable(inputFrom(ctx).String()):
		addParagraphIfNeeded()
		return nextTable(ctx)
	case isPrefixedBy(ctx, "```"):
		addParagraphIfNeeded()
		return nextCodeBlock(ctx)
	case isPrefixedBy(ctx, "=>"):
		addParagraphIfNeeded()
		return nextLaunchPad(ctx)
	case isPrefixedBy(ctx, "<="):
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeTransclusion(line, hyphaNameFrom(ctx)), done
	case isPrefixedBy(ctx, "----"):
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeHorizontalLine(line), done

	case isPrefixedBy(ctx, "###### "):
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 6), done
	case isPrefixedBy(ctx, "##### "):
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 5), done
	case isPrefixedBy(ctx, "#### "):
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 4), done
	case isPrefixedBy(ctx, "### "):
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 3), done
	case isPrefixedBy(ctx, "## "):
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 2), done
	case isPrefixedBy(ctx, "# "):
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeHeading(line, hyphaNameFrom(ctx), 1), done

	case isPrefixedBy(ctx, ">"): // TODO: implement proper fractal quotes
		addParagraphIfNeeded()
		line, done := nextLine(ctx)
		return blocks.MakeQuote(line, hyphaNameFrom(ctx)), done
	}

	var line, done = nextLine(ctx)

	// Beware! Usage of goto. Some may say it is considered evil but in this case it helped to make a better-structured code.
	switch state.where {
	case "list":
		goto listState
	default: // "p" or ""
		goto normalState
	}

listState:
	if done := state.list.ProcessLine(line); done {
		state.list.Finalize()
		state.where = ""
		goto normalState
	}
	goto end

normalState:
	switch {
	case "" == strings.TrimSpace(line):
		addParagraphIfNeeded()

	case blocks.MatchesList(line):
		addParagraphIfNeeded()
		list, _ := blocks.MakeList(line, hyphaNameFrom(ctx))
		state.where = "list"
		state.list = list
		ret = state.list

	case state.where == "p":
		state.paragraph.AddLine(line)
	default:
		state.where = "p"
		p := blocks.MakeParagraph(line, hyphaNameFrom(ctx))
		state.paragraph = &blocks.Paragraph{Formatted: p}
	}

end:
	return ret, done
}
