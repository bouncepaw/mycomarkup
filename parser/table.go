package parser

import (
	"github.com/bouncepaw/mycomarkup/blocks"
	"github.com/bouncepaw/mycomarkup/mycocontext"
	"regexp"
	"strings"
	"unicode"
)

func nextTable(ctx mycocontext.Context) (t blocks.Table, done bool) {
	line, done := mycocontext.NextLine(ctx)
	t = tableFromFirstLine(line, ctx.HyphaName())
	for {
		line, done = mycocontext.NextLine(ctx)
		if processTableLine(&t, line) {
			break
		}
	}
	return t, done
}

func processTableLine(t *blocks.Table, line string) (done bool) {
	if strings.HasPrefix(strings.TrimLeft(line, " \t"), "}") && !t.InMultiline {
		return true
	}
	if !t.InMultiline {
		pushTableRow(t)
	}
	s := initialTableParserState()
	s.lookingForNonSpace = !t.InMultiline
	for _, r := range line {
		parseTableRune(t, s, r)
	}
	parseTableRune(t, s, '\n')
	return false
}

var tableRe = regexp.MustCompile(`^table\s*{`)

func matchesTable(ctx mycocontext.Context) bool {
	return tableRe.Match(ctx.Input().Bytes())
}

func tableFromFirstLine(line, hyphaName string) blocks.Table {
	return blocks.Table{
		HyphaName: hyphaName,
		Caption:   line[strings.IndexRune(line, '{')+1:],
		Rows:      make([]*blocks.TableRow, 0),
	}
}

type tableParserState struct {
	skipNext           bool
	escaping           bool
	lookingForNonSpace bool
	countingColspan    bool
}

func initialTableParserState() *tableParserState {
	return &tableParserState{
		skipNext:           false,
		escaping:           false,
		lookingForNonSpace: false,
		countingColspan:    false,
	}
}

func parseTableRune(t *blocks.Table, s *tableParserState, r rune) (done bool) {
	switch {
	case s.skipNext:
		s.skipNext = false

	case s.lookingForNonSpace && unicode.IsSpace(r):
	case s.lookingForNonSpace && (r == '!' || r == '|'):
		t.CurrCellMarker = r
		t.CurrColspan = 1
		s.lookingForNonSpace = false
		s.countingColspan = true
	case s.lookingForNonSpace:
		t.CurrCellMarker = '^' // ^ represents implicit |, not part of syntax
		t.CurrColspan = 1
		s.lookingForNonSpace = false
		t.CurrCellBuilder.WriteRune(r)

	case s.escaping:
		t.CurrCellBuilder.WriteRune(r)

	case t.InMultiline && r == '}':
		t.InMultiline = false
	case t.InMultiline && r == '\n':
		t.CurrCellBuilder.WriteRune(r)
		t.CurrCellBuilder.WriteRune('\n')
	case t.InMultiline:
		t.CurrCellBuilder.WriteRune(r)

		// Not in multiline:
	case (r == '|' || r == '!') && !s.countingColspan:
		pushTableCell(t)
		t.CurrCellMarker = r
		t.CurrColspan = 1
		s.countingColspan = true
	case r == t.CurrCellMarker && (r == '|' || r == '!') && s.countingColspan:
		t.CurrColspan++
	case r == '{':
		t.InMultiline = true
		s.countingColspan = false
	case r == '\n':
		pushTableCell(t)
	default:
		t.CurrCellBuilder.WriteRune(r)
		s.countingColspan = false
	}
	return false
}

func pushTableRow(t *blocks.Table) {
	t.Rows = append(t.Rows, &blocks.TableRow{
		HyphaName: t.HyphaName,
		Cells:     []*blocks.TableCell{},
	})
}

func pushTableCell(t *blocks.Table) {
	tc := &blocks.TableCell{
		Contents: blocks.MakeFormatted(t.CurrCellBuilder.String(), t.HyphaName),
		Colspan:  t.CurrColspan,
	}
	switch t.CurrCellMarker {
	case '|', '^':
		tc.IsHeaderCell = false
	case '!':
		tc.IsHeaderCell = true
	}
	// We expect the table to have at least one row ready, so no nil-checking
	tr := t.Rows[len(t.Rows)-1]
	tr.Cells = append(tr.Cells, tc)
	t.CurrCellBuilder = strings.Builder{}
}
