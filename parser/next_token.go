package parser

import (
	"bytes"
	"github.com/bouncepaw/mycomarkup/v3/blocks"
	"github.com/bouncepaw/mycomarkup/v3/mycocontext"
	"strings"
)

func isPrefixedBy(ctx mycocontext.Context, s string) bool { // This function has bugs in it!
	return bytes.HasPrefix(ctx.Input().Bytes(), []byte(s))
}

func nextLaunchPad(ctx mycocontext.Context) (blocks.LaunchPad, bool) {
	var (
		hyphaName   = ctx.HyphaName()
		line        string
		done        bool
		rocketLinks = make([]blocks.RocketLink, 0)
	)
	for isPrefixedBy(ctx, "=>") {
		line, done = mycocontext.NextLine(ctx)
		rocketLinks = append(rocketLinks, blocks.ParseRocketLink(line, hyphaName))
	}
	return blocks.NewLaunchPad(rocketLinks), done
}

func nextCodeBlock(ctx mycocontext.Context) (code blocks.CodeBlock, eof bool) {
	contents := ""
	line, eof := mycocontext.NextLine(ctx)
	language := strings.TrimPrefix(line, "```")

	for !eof {
		line, eof = mycocontext.NextLine(ctx)
		if strings.HasPrefix(line, "```") {
			break
		}
		contents += "\n" + line // Note: newline added every time
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
	prefices := []string{"=>", "<=", "```", "* ", "*. ", "*v ", "*x ", "# ", "## ", "### ", "#### ", "##### ", "###### ", ">", "----"}
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

func NextToken(ctx mycocontext.Context) (blocks.Block, bool) {
	switch {
	case matchesEmptyLine(ctx):
		_, done := mycocontext.NextLine(ctx)
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
		line, done := mycocontext.NextLine(ctx)
		return blocks.MakeTransclusion(ctx, line), done
	case isPrefixedBy(ctx, "----"):
		line, done := mycocontext.NextLine(ctx)
		return blocks.NewHorizontalLine(line), done

	case isPrefixedBy(ctx, "==== "):
		line, done := mycocontext.NextLine(ctx)
		return parseHeading(line, ctx.HyphaName(), 5), done
	case isPrefixedBy(ctx, "=== "):
		line, done := mycocontext.NextLine(ctx)
		return parseHeading(line, ctx.HyphaName(), 4), done
	case isPrefixedBy(ctx, "== "):
		line, done := mycocontext.NextLine(ctx)
		return parseHeading(line, ctx.HyphaName(), 3), done
	case isPrefixedBy(ctx, "= "):
		line, done := mycocontext.NextLine(ctx)
		return parseHeading(line, ctx.HyphaName(), 2), done

	case isPrefixedBy(ctx, "###### "):
		line, done := mycocontext.NextLine(ctx)
		return parseLegacyHeading(line, ctx.HyphaName(), 6), done
	case isPrefixedBy(ctx, "##### "):
		line, done := mycocontext.NextLine(ctx)
		return parseLegacyHeading(line, ctx.HyphaName(), 5), done
	case isPrefixedBy(ctx, "#### "):
		line, done := mycocontext.NextLine(ctx)
		return parseLegacyHeading(line, ctx.HyphaName(), 4), done
	case isPrefixedBy(ctx, "### "):
		line, done := mycocontext.NextLine(ctx)
		return parseLegacyHeading(line, ctx.HyphaName(), 3), done
	case isPrefixedBy(ctx, "## "):
		line, done := mycocontext.NextLine(ctx)
		return parseLegacyHeading(line, ctx.HyphaName(), 2), done
	case isPrefixedBy(ctx, "# "):
		line, done := mycocontext.NextLine(ctx)
		return parseLegacyHeading(line, ctx.HyphaName(), 1), done

	case matchesImg(ctx):
		return nextImg(ctx)
	case matchesTable(ctx):
		return nextTable(ctx)
	}
	return nextParagraph(ctx)
}
