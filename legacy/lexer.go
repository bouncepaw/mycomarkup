package markup

import (
	"fmt"
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/generator"
	"html"
	"strings"

	"github.com/bouncepaw/mycomarkup/util"
)

// HyphaExists holds function that checks that a hypha is present.
var HyphaExists func(string) bool

// HyphaAccess holds function that accesses a hypha by its name.
var HyphaAccess func(string) (rawText, binaryHtml string, err error)

// HyphaIterate is a function that iterates all hypha names existing.
var HyphaIterate func(func(string))

// GemLexerState is used by markup parser to remember what is going on.
type GemLexerState struct {
	// Name of hypha being parsed
	name  string
	where string // "", "list", "pre"
	// Line id
	id  int
	buf string
	// Temporaries
	img   *Img
	table *Table
	list  *List
}

type Line struct {
	id int
	// interface{} may be bad. TODO: a proper type
	contents interface{}
}

func (md *MycoDoc) lex() (ast []Line) {
	blocks.HyphaExists = HyphaExists
	var state = GemLexerState{name: md.hyphaName}

	for _, line := range append(strings.Split(md.contents, "\n"), "") {
		lineToAST(line, &state, &ast)
	}
	return ast
}

// Lex `line` in markup and save it to `ast` using `state`.
func lineToAST(line string, state *GemLexerState, ast *[]Line) {
	addLine := func(text interface{}) {
		*ast = append(*ast, Line{id: state.id, contents: text})
	}
	addParagraphIfNeeded := func() {
		if state.where == "p" {
			state.where = ""
			addLine(fmt.Sprintf("<p id='%d'>%s</p>", state.id, strings.ReplaceAll(blocks.ParagraphToHtml(state.name, state.buf), "\n", "<br>")))
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
			addLine(state.buf + "</ul>")
		case "p":
			addParagraphIfNeeded()
		}
		return
	}

	startsWith := func(token string) bool {
		return strings.HasPrefix(line, token)
	}
	addHeading := func(i int) {
		id := util.StringID(line[i+1:])
		addLine(fmt.Sprintf(`<h%d id='%d'>%s<a href="#%s" id="%s" class="heading__link"></a></h%d>`, i, state.id, blocks.ParagraphToHtml(state.name, line[i+1:]), id, id, i))
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
		state.buf = strings.TrimSuffix(state.buf, "\n")
		addLine(state.buf + "</code></pre>")
		state.buf = ""
	default:
		state.buf += html.EscapeString(line) + "\n"
	}
	return

launchpadState:
	switch {
	case startsWith("=>"):
		href, text, class := blocks.Rocketlink(line, state.name)
		state.buf += fmt.Sprintf(`	<li class="launchpad__entry"><a href="%s" class="rocketlink %s">%s</a></li>`, href, class, text)
	case startsWith("```"):
		state.where = "pre"
		addLine(state.buf + "</ul>")
		state.id++
		state.buf = fmt.Sprintf("<pre id='%d' alt='%s' class='codeblock'><code>", state.id, strings.TrimPrefix(line, "```"))
	default:
		state.where = ""
		addLine(state.buf + "</ul>")
		goto normalState
	}
	return

normalState:
	state.id++
	switch {

	case startsWith("```"):
		addParagraphIfNeeded()
		state.where = "pre"
		state.buf = fmt.Sprintf("<pre id='%d' alt='%s' class='codeblock'><code>", state.id, strings.TrimPrefix(line, "```"))

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
				"<blockquote id='%d'>%s</blockquote>",
				state.id,
				blocks.ParagraphToHtml(state.name, remover(">")(line)),
			),
		)
	case startsWith("=>"):
		addParagraphIfNeeded()
		state.where = "launchpad"
		state.buf = fmt.Sprintf("<ul class='launchpad' id='%d'>\n", state.id)
		goto launchpadState

	case startsWith("<="):
		addParagraphIfNeeded()
		addLine(parseTransclusion(line, state.name))
	case startsWith("----"):
		addParagraphIfNeeded()
		hr := blocks.MakeHorizontalLine(line)
		*ast = append(*ast, Line{id: -1, contents: generator.BlockToHTML(hr)})
	case MatchesList(line):
		addParagraphIfNeeded()
		list, _ := NewList(line, state.name)
		state.where = "list"
		state.list = list
		addLine(state.list)
	case MatchesImg(line):
		addParagraphIfNeeded()
		img, shouldGoBackToNormal := ImgFromFirstLine(line, state.name)
		if shouldGoBackToNormal {
			addLine(*img)
		} else {
			state.where = "img"
			state.img = img
		}
	case MatchesTable(line):
		addParagraphIfNeeded()
		state.where = "table"
		state.table = TableFromFirstLine(line, state.name)

	case state.where == "p":
		state.buf += "\n" + line
	default:
		state.where = "p"
		state.buf = line
	}
}
