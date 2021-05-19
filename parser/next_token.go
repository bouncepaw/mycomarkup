package parser

import (
	"context"
	"fmt"
	"html"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
)

// ParserState is used by markup parser to remember what is going on.
type ParserState struct {
	// Target of hypha being lexed
	Name  string
	where string // "", "list", "pre", etc.
	// Temporaries
	code      *blocks.CodeBlock
	table     *blocks.Table
	list      *blocks.List
	launchpad *blocks.LaunchPad
	paragraph *blocks.Paragraph
}

func startsWith(ctx context.Context, s string) bool {
	return strings.HasPrefix(inputFrom(ctx).String(), s)
}

func nextImg(ctx context.Context, state *ParserState, line string, doneBefore bool) (img blocks.Img, doneAfter bool) {
	var b byte
	img, imgDone := blocks.MakeImg(line, state.Name)
	if imgDone {
		return img, doneBefore
	}

	for !imgDone {
		b, doneAfter = nextByte(ctx)
		imgDone = img.ProcessRune(rune(b))
	}

	defer nextLine(ctx) // Characters after the final } of img are ignored.
	return img, doneAfter
}

func nextCodeBlock(ctx context.Context, firstLine string) (code blocks.CodeBlock, done bool) {
	code = blocks.MakeCodeBlock(strings.TrimPrefix(firstLine, "```"), "")

	var line string
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

// Lex `line` in markup and maybe return a token.
func nextToken(ctx context.Context, state *ParserState) (interface{}, bool) {
	var (
		line, done = nextLine(ctx)
		ret        interface{}
	)
	addParagraphIfNeeded := func() {
		if state.where == "p" {
			state.where = ""
			ret = *state.paragraph
		}
	}
	startsWith := func(token string) bool {
		return strings.HasPrefix(line, token)
	}
	addHeading := func(i int) {
		ret = blocks.MakeHeading(line, state.Name, uint(i))
	}
	// Beware! Usage of goto. Some may say it is considered evil but in this case it helped to make a better-structured code.
	switch state.where {
	case "table":
		goto tableState
	case "list":
		goto listState
	case "launchpad":
		goto launchpadState
	default: // "p" or ""
		goto normalState
	}

tableState:
	if done := state.table.ProcessLine(line); done {
		state.where = ""
		ret = *state.table
	}
	goto end

listState:
	if done := state.list.Parse(line); done {
		state.list.Finalize()
		state.where = ""
		goto normalState
	}
	goto end

launchpadState:
	switch {
	case "" == strings.TrimSpace(line):
		state.where = ""
		ret = *state.launchpad
		state.launchpad = nil
	case startsWith("=>"):
		state.launchpad.AddRocket(blocks.MakeRocketLink(line, state.Name))
	case startsWith("```"):
		return nextCodeBlock(ctx, line)
	default:
		fmt.Println("night call")
		ret = *state.launchpad
		state.where = ""
		goto normalState
	}
	goto end

normalState:
	switch {
	case "" == strings.TrimSpace(line):
		addParagraphIfNeeded()
	case startsWith("```"):
		addParagraphIfNeeded()
		return nextCodeBlock(ctx, line)

	case startsWith("###### "):
		addParagraphIfNeeded()
		addHeading(6)
	case startsWith("##### "):
		addParagraphIfNeeded()
		addHeading(5)
	case startsWith("#### "):
		addParagraphIfNeeded()
		addHeading(4)
	case startsWith("### "):
		addParagraphIfNeeded()
		addHeading(3)
	case startsWith("## "):
		addParagraphIfNeeded()
		addHeading(2)
	case startsWith("# "):
		addParagraphIfNeeded()
		addHeading(1)

	case startsWith(">"):
		addParagraphIfNeeded()
		ret = blocks.MakeQuote(line, state.Name)
	case startsWith("=>"):
		addParagraphIfNeeded()
		state.where = "launchpad"
		lp := blocks.MakeLaunchPad()
		state.launchpad = &lp
		goto launchpadState

	case startsWith("<="):
		addParagraphIfNeeded()
		ret = blocks.MakeTransclusion(line, state.Name)
	case startsWith("----"):
		addParagraphIfNeeded()
		ret = blocks.MakeHorizontalLine(line)
	case blocks.MatchesList(line):
		addParagraphIfNeeded()
		list, _ := blocks.NewList(line, state.Name)
		state.where = "list"
		state.list = list
		ret = state.list
	case blocks.MatchesImg(line):
		addParagraphIfNeeded()
		return nextImg(ctx, state, line, done)

	case blocks.MatchesTable(line):
		addParagraphIfNeeded()
		state.where = "table"
		state.table = blocks.TableFromFirstLine(line, state.Name)

	case state.where == "p":
		state.paragraph.AddLine(line)
	default:
		state.where = "p"
		p := blocks.MakeParagraph(line, state.Name)
		state.paragraph = &blocks.Paragraph{Formatted: p}
	}

end:
	return ret, done
}
