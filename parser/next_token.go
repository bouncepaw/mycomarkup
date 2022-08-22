package parser

import (
	"bytes"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/blocks"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/links"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/parser/ctxio"
	"html"
	"strings"
)

func isPrefixedBy(ctx mycocontext.Context, s string) bool { // This function has bugs in it!
	return bytes.HasPrefix(ctx.Input().Bytes(), []byte(s))
}

func nextLaunchPad(ctx mycocontext.Context) (blocks.LaunchPad, bool) {
	var (
		line        string
		done        bool
		rocketLinks = make([]blocks.RocketLink, 0)
	)
	for isPrefixedBy(ctx, "=>") {
		line, done = ctxio.NextLine(ctx)
		rocketLinks = append(rocketLinks, lineToRocketLink(ctx, line))
	}
	return blocks.NewLaunchPad(rocketLinks), done
}

func lineToRocketLink(ctx mycocontext.Context, line string) blocks.RocketLink {
	line = strings.TrimSpace(line[2:]) // Drop =>
	if line == "" {
		return blocks.RocketLink{IsEmpty: true, Link: links.LinkFrom(ctx, "", "")}
	}

	var (
		// Address is text after => till first whitespace
		addr = strings.FieldsFunc(line, func(r rune) bool {
			return r == '|'
		})[0]
		// Display text is what is left
		display = strings.TrimPrefix(strings.TrimPrefix(line, addr), "|")
		rl      = blocks.RocketLink{
			IsEmpty: false,
			Link:    links.LinkFrom(ctx, addr, display),
		}
	)

	return rl
}

func nextCodeBlock(ctx mycocontext.Context) (code blocks.CodeBlock, eof bool) {
	contents := ""
	line, eof := ctxio.NextLine(ctx)
	language := strings.TrimPrefix(line, "```")

	for !eof {
		line, eof = ctxio.NextLine(ctx)
		if strings.HasPrefix(line, "```") {
			break
		}
		contents += "\n" + html.EscapeString(line) // Note: newline added every time
	}
	if len(contents) > 0 {
		contents = contents[1:] // Drop the leading newline
	}
	return blocks.NewCodeBlock(language, contents), eof
}

func linesForQuote(ctx mycocontext.Context) ([]string, bool) {
	var (
		line  string
		lines []string
		done  bool
	)
	for {
		line, done = ctxio.NextLine(ctx)
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
		lines, done   = linesForQuote(ctx)
		innerText     bytes.Buffer
		quoteContents = make([]blocks.Block, 0)
	)

	for i, line := range lines {
		if i > 0 {
			innerText.WriteRune('\n')
		}
		innerText.WriteString(line)
	}

	parseSubdocumentForEachBlock(ctx, &innerText, func(block blocks.Block) {
		quoteContents = append(quoteContents, block)
	})

	return blocks.NewQuote(quoteContents), done
}

func nextLineIsSomething(ctx mycocontext.Context) bool {
	// Maybe sort them by how often they used? How to count?
	prefices := []string{"=>", "<=", "```", "* ", "*. ", "*v ", "*x ", ">", "----", "= ", "== ", "=== ", "==== "}
	for _, prefix := range prefices {
		if isPrefixedBy(ctx, prefix) {
			return true
		}
	}
	return matchesEmptyLine(ctx) || matchesImg(ctx) || matchesTable(ctx)
}

func matchesEmptyLine(ctx mycocontext.Context) bool {
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

// NextToken returns the next found block in the input. The second value is true if eof.
func NextToken(ctx mycocontext.Context) (blocks.Block, bool) {
	switch {
	case matchesEmptyLine(ctx):
		_, done := ctxio.NextLine(ctx)
		return nil, done
	case looksLikeList(ctx):
		return nextList(ctx)
	case isPrefixedBy(ctx, "```"):
		return nextCodeBlock(ctx)
	case isPrefixedBy(ctx, "=>"):
		return nextLaunchPad(ctx)
	case isPrefixedBy(ctx, ">"):
		return nextQuote(ctx)
	case isPrefixedBy(ctx, "<="):
		line, done := ctxio.NextLine(ctx)
		return blocks.MakeTransclusion(ctx, line), done
	case isPrefixedBy(ctx, "----"):
		line, done := ctxio.NextLine(ctx)
		return blocks.NewThematicBreak(line), done

	case isPrefixedBy(ctx, "==== "):
		line, done := ctxio.NextLine(ctx)
		return parseHeading(ctx, line, 4), done
	case isPrefixedBy(ctx, "=== "):
		line, done := ctxio.NextLine(ctx)
		return parseHeading(ctx, line, 3), done
	case isPrefixedBy(ctx, "== "):
		line, done := ctxio.NextLine(ctx)
		return parseHeading(ctx, line, 2), done
	case isPrefixedBy(ctx, "= "):
		line, done := ctxio.NextLine(ctx)
		return parseHeading(ctx, line, 1), done

	case matchesImg(ctx):
		return nextImg(ctx)
	case matchesTable(ctx):
		return nextTable(ctx)
	}
	return nextParagraph(ctx)
}
