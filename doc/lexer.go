package doc

import (
	"fmt"
	"html"
	"strings"

	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/util"
)

// LexerState is used by markup lexer to remember what is going on.
type LexerState struct {
	// Target of hypha being lexed
	name  string
	where string // "", "list", "pre", etc.
	// Temporaries
	buf       string // TODO: get rid of this as soon as we can.
	code      *blocks.CodeBlock
	img       *blocks.Img
	table     *blocks.Table
	list      *blocks.List
	launchpad *blocks.LaunchPad
}

type Token struct {
	// TODO: replace with a proper interface one day, when it's all over.
	Value interface{}
}

// Lex `line` in markup and save it to `ast` using `state`.
func lineToToken(line string, state *LexerState, ast *[]Token) {
	addLine := func(text interface{}) {
		*ast = append(*ast, Token{Value: text})
	}
	addParagraphIfNeeded := func() {
		if state.where == "p" {
			state.where = ""
			addLine(fmt.Sprintf("\n<p>%s</p>", strings.ReplaceAll(blocks.ParagraphToHtml(state.name, state.buf), "\n", "<br>")))
			state.buf = ""
		}
	}

	// Process empty lines depending on the current state
	if "" == strings.TrimSpace(line) {
		switch state.where {
		case "pre":
			state.buf += "\n"
		case "launchpad":
			state.where = ""
			addLine(*state.launchpad)
			state.launchpad = nil
		case "p":
			addParagraphIfNeeded()
		}
		return
	}

	startsWith := func(token string) bool {
		return strings.HasPrefix(line, token)
	}
	addHeading := func(i int) {
		addLine(blocks.MakeHeading(line, state.name, uint(i)))
	}

	// Beware! Usage of goto. Some may say it is considered evil but in this case it helped to make a better-structured code.
	switch state.where {
	case "img":
		goto imgState
	case "table":
		goto tableState
	case "list":
		goto listState
	case "pre":
		goto preformattedState
	case "launchpad":
		goto launchpadState
	default: // "p" or ""
		goto normalState
	}

imgState:
	if shouldGoBackToNormal := state.img.Process(line); shouldGoBackToNormal {
		state.where = ""
		addLine(*state.img)
	}
	return

tableState:
	if shouldGoBackToNormal := state.table.Process(line); shouldGoBackToNormal {
		state.where = ""
		addLine(*state.table)
	}
	return

listState:
	if done := state.list.Parse(line); done {
		state.list.Finalize()
		state.where = ""
		goto normalState
	}
	return

preformattedState:
	switch {
	case startsWith("```"):
		state.where = ""
		addLine(*state.code)
	default:
		fmt.Println("adding line@")
		state.code.AddLine(html.EscapeString(line))
	}
	return

launchpadState:
	switch {
	case startsWith("=>"):
		state.launchpad.AddRocket(blocks.MakeRocketLink(line, state.name))
	case startsWith("```"):
		addLine(*state.launchpad)
		state.launchpad = nil
		state.where = "pre"
		cb := blocks.MakeCodeBlock(strings.TrimPrefix(line, "```"), "")
		state.code = &cb
	default:
		fmt.Println("night call")
		addLine(*state.launchpad)
		state.launchpad = nil
		state.where = ""
		goto normalState
	}
	return

normalState:
	switch {

	case startsWith("```"):
		addParagraphIfNeeded()
		state.where = "pre"
		cb := blocks.MakeCodeBlock(strings.TrimPrefix(line, "```"), "")
		state.code = &cb

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
		addLine(
			fmt.Sprintf(
				"<blockquote>%s</blockquote>",
				blocks.ParagraphToHtml(state.name, util.Remover(">")(line)),
			),
		)
	case startsWith("=>"):
		addParagraphIfNeeded()
		state.where = "launchpad"
		lp := blocks.MakeLaunchPad()
		state.launchpad = &lp
		goto launchpadState

	case startsWith("<="):
		addParagraphIfNeeded()
		addLine(blocks.MakeTransclusion(line, state.name))
	case startsWith("----"):
		addParagraphIfNeeded()
		*ast = append(*ast, Token{Value: blocks.MakeHorizontalLine(line)})
	case blocks.MatchesList(line):
		addParagraphIfNeeded()
		list, _ := blocks.NewList(line, state.name)
		state.where = "list"
		state.list = list
		addLine(state.list)
	case blocks.MatchesImg(line):
		addParagraphIfNeeded()
		img, shouldGoBackToNormal := blocks.ImgFromFirstLine(line, state.name)
		if shouldGoBackToNormal {
			addLine(*img)
		} else {
			state.where = "img"
			state.img = img
		}
	case blocks.MatchesTable(line):
		addParagraphIfNeeded()
		state.where = "table"
		state.table = blocks.TableFromFirstLine(line, state.name)

	case state.where == "p":
		state.buf += "\n" + line
	default:
		state.where = "p"
		state.buf = line
	}
}
