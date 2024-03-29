package parser

import (
	"bytes"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/blocks"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/mycocontext"
	"git.sr.ht/~bouncepaw/mycomarkup/v5/parser/ctxio"
	"regexp"
	"strings"
	"unicode"
)

func nextTable(ctx mycocontext.Context) (t blocks.Table, eof bool) {
	line, eof := ctxio.NextLine(ctx)
	t, tableDone := tableFromFirstLine(line)
	if tableDone || eof {
		return t, eof
	}
	for {
		row, rowFound, tableDone := nextTableRow(ctx)
		if rowFound {
			t = t.WithNewRow(row)
		}
		if tableDone {
			break
		}
	}
	_, eof = ctxio.NextLine(ctx) // Ignore text after }
	return t, eof
}

func nextTableRow(ctx mycocontext.Context) (row blocks.TableRow, foundRow, tableDone bool) {
	var (
		cleaningLeadingWhitespace = true
		countingColspan           = false

		currColspan    uint = 0
		currCellMarker rune

		cellText *bytes.Buffer
		r        rune
		eof      bool

		cells []blocks.TableCell
	)
runeWalker:
	for {
		r, eof = ctxio.NextRune(ctx)
		if eof {
			tableDone = true
			break
		}
	automaton:
		switch {
		case r == '\n':
			break runeWalker
		case cleaningLeadingWhitespace && unicode.IsSpace(r):
			continue
		case cleaningLeadingWhitespace: // When non-space rune, try again
			cleaningLeadingWhitespace = false
			goto automaton // The next and the nextnext case-clauses might trigger

		case (!cleaningLeadingWhitespace && !countingColspan) && r == '}':
			tableDone = true
			break runeWalker
		case (!cleaningLeadingWhitespace && !countingColspan) && (r == '|' || r == '!'):
			// Proper cell marker, great! Let's start counting colspan then
			currCellMarker = r
			currColspan = 1
			countingColspan = true
		case countingColspan && r == currCellMarker:
			currColspan++

		case !cleaningLeadingWhitespace && !countingColspan, countingColspan:
			ctxio.UnreadRune(ctx)
			var contents []blocks.Block
			cellText, tableDone = nextTableCellContents(ctx)
			parseSubdocumentForEachBlock(ctx, cellText, func(block blocks.Block) {
				contents = append(contents, block)
			})
			cell := blocks.NewTableCell(currCellMarker == '!', currColspan, contents)
			cells = append(cells, cell)
			if tableDone {
				break runeWalker
			}

			// Reset state
			countingColspan = false
			currColspan = 0
			currCellMarker = 0
		case r == '}':
			tableDone = true
			break runeWalker
		}
	}

	if len(cells) == 0 {
		return blocks.TableRow{}, false, tableDone
	}

	return blocks.NewTableRow(cells), true, tableDone
}

func nextTableCellContents(
	ctx mycocontext.Context,
) (
	contents *bytes.Buffer,
	tableDone bool,
) {
	var (
		contentsBuilder bytes.Buffer
		escaping        = false
		inLink          = false
	)
runeWalker:
	for {
		r, eof := ctxio.NextRune(ctx)
		if eof {
			tableDone = true
			break
		}
		switch {
		case r == '\n':
			ctxio.UnreadRune(ctx)
			break runeWalker
		case escaping:
			contentsBuilder.WriteRune(r)
			escaping = false
		case r == '\\':
			contentsBuilder.WriteRune('\\')
			escaping = true
		case r == '[':
			contentsBuilder.WriteRune('[')
			r, eof = ctxio.NextRune(ctx)
			if r == '[' {
				inLink = true
			}
			contentsBuilder.WriteRune(r)
		case inLink && r == ']':
			contentsBuilder.WriteRune(']')
			r, eof = ctxio.NextRune(ctx)
			if r == ']' {
				inLink = false
			}
			contentsBuilder.WriteRune(r)
		case !inLink && r == '|', r == '!': // looks like a new cell
			ctxio.UnreadRune(ctx)
			break runeWalker
		case !inLink && r == '{':
			contentsBuilder.WriteString(nextTableMultiline(ctx))

		case r == '}':
			tableDone = true
		default:
			contentsBuilder.WriteRune(r)
		}
	}
	return &contentsBuilder, tableDone
}

// return text until the next unmatched unescaped } (exclusively).
func nextTableMultiline(ctx mycocontext.Context) string {
	var (
		curlyCount = 1 // 1 is the initial state: multiline open. When it is 0, done.
		escaping   = false
		onNewLine  = true
		r          rune
		eof        bool
		ret        strings.Builder
	)
	for {
		r, eof = ctxio.NextRune(ctx)
		if eof {
			break
		}
	automaton:
		switch {
		case r == '\n':
			onNewLine = true
		case escaping:
			escaping = false
		case onNewLine && unicode.IsSpace(r):
			continue
		case onNewLine:
			onNewLine = false
			goto automaton
		case r == '\\':
			escaping = true
		case r == '{':
			curlyCount++
		case r == '}':
			curlyCount--
		}
		if curlyCount == 0 {
			break
		}
		ret.WriteRune(r)
	}
	return ret.String()
}

var tableRe = regexp.MustCompile(`^table\s*{`)

func matchesTable(ctx mycocontext.Context) bool {
	return tableRe.Match(ctx.Input().Bytes())
}

func tableFromFirstLine(line string) (t blocks.Table, done bool) { // V3
	var (
		caption       = line[strings.IndexRune(line, '{')+1:]
		closeBracePos = strings.IndexRune(caption, '}')
		isClosed      = closeBracePos != -1
	)
	if isClosed {
		caption = caption[:closeBracePos]
	}
	return blocks.NewTable(caption, make([]blocks.TableRow, 0)), isClosed
}
