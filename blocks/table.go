package blocks

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Table struct {
	// data
	HyphaName string
	Caption   string
	Rows      []*TableRow
	// state
	inMultiline bool
	// tmp
	currCellMarker  rune
	currColspan     uint
	currCellBuilder strings.Builder
}

var tableRe = regexp.MustCompile(`^table\s+{`)

func MatchesTable(line string) bool {
	return tableRe.MatchString(line)
}

func TableFromFirstLine(line, hyphaName string) *Table {
	return &Table{
		HyphaName: hyphaName,
		Caption:   line[strings.IndexRune(line, '{')+1:],
		Rows:      make([]*TableRow, 0),
	}
}

func (t *Table) Process(line string) (shouldGoBackToNormal bool) {
	if strings.TrimSpace(line) == "}" && !t.inMultiline {
		return true
	}
	if !t.inMultiline {
		t.pushRow()
	}
	var (
		inLink             bool
		skipNext           bool
		escaping           bool
		lookingForNonSpace = !t.inMultiline
		countingColspan    bool
	)
	for i, r := range line {
		switch {
		case skipNext:
			skipNext = false
			continue

		case lookingForNonSpace && unicode.IsSpace(r):
		case lookingForNonSpace && (r == '!' || r == '|'):
			t.currCellMarker = r
			t.currColspan = 1
			lookingForNonSpace = false
			countingColspan = true
		case lookingForNonSpace:
			t.currCellMarker = '^' // ^ represents implicit |, not part of syntax
			t.currColspan = 1
			lookingForNonSpace = false
			t.currCellBuilder.WriteRune(r)

		case escaping:
			t.currCellBuilder.WriteRune(r)
		case inLink && r == ']' && len(line)-1 > i && line[i+1] == ']':
			t.currCellBuilder.WriteString("]]")
			inLink = false
			skipNext = true
		case inLink:
			t.currCellBuilder.WriteRune(r)

		case t.inMultiline && r == '}':
			t.inMultiline = false
		case t.inMultiline && i == len(line)-1:
			t.currCellBuilder.WriteRune('\n')
		case t.inMultiline:
			t.currCellBuilder.WriteRune(r)

			// Not in multiline:
		case (r == '|' || r == '!') && !countingColspan:
			t.pushCell()
			t.currCellMarker = r
			t.currColspan = 1
			countingColspan = true
		case r == t.currCellMarker && (r == '|' || r == '!') && countingColspan:
			t.currColspan++
		case r == '{':
			t.inMultiline = true
			countingColspan = false
		case r == '[' && len(line)-1 > i && line[i+1] == '[':
			t.currCellBuilder.WriteString("[[")
			inLink = true
			skipNext = true
		case i == len(line)-1:
			t.pushCell()
		default:
			t.currCellBuilder.WriteRune(r)
			countingColspan = false
		}
	}
	return false
}

func (t *Table) pushRow() {
	t.Rows = append(t.Rows, &TableRow{
		cells: make([]*TableCell, 0),
	})
}

func (t *Table) pushCell() {
	tc := &TableCell{
		content: t.currCellBuilder.String(),
		colspan: t.currColspan,
	}
	switch t.currCellMarker {
	case '|', '^':
		tc.kind = tableCellDatum
	case '!':
		tc.kind = tableCellHeader
	}
	// We expect the table to have at least one row ready, so no nil-checking
	tr := t.Rows[len(t.Rows)-1]
	tr.cells = append(tr.cells, tc)
	t.currCellBuilder = strings.Builder{}
}

type TableRow struct {
	cells []*TableCell
}

func (tr *TableRow) AsHtml(hyphaName string) (html string) {
	for _, tc := range tr.cells {
		html += tc.asHtml(hyphaName)
	}
	return fmt.Sprintf("<tr>%s</tr>\n", html)
}

// Most likely, rows with more than two header cells are theads. I allow one extra datum cell for tables like this:
// |   ! a ! b
// ! c | d | e
// ! f | g | h
func (tr *TableRow) LooksLikeThead() bool {
	var (
		headerAmount = 0
		datumAmount  = 0
	)
	for _, tc := range tr.cells {
		switch tc.kind {
		case tableCellHeader:
			headerAmount++
		case tableCellDatum:
			datumAmount++
		}
	}
	return headerAmount >= 2 && datumAmount <= 1
}

type TableCell struct {
	kind    TableCellKind
	colspan uint
	content string
}

func (tc *TableCell) asHtml(hyphaName string) string {
	return fmt.Sprintf(
		"<%[1]s %[2]s>%[3]s</%[1]s>\n",
		tc.kind.tagName(),
		tc.colspanAttribute(),
		tc.contentAsHtml(hyphaName),
	)
}

func (tc *TableCell) colspanAttribute() string {
	if tc.colspan <= 1 {
		return ""
	}
	return fmt.Sprintf(`colspan="%d"`, tc.colspan)
}

func (tc *TableCell) contentAsHtml(hyphaName string) (html string) {
	for _, line := range strings.Split(tc.content, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			if html != "" {
				html += `<br>`
			}
			html += ParagraphToHtml(hyphaName, line)
		}
	}
	return html
}

type TableCellKind int

const (
	tableCellUnknown TableCellKind = iota
	tableCellHeader
	tableCellDatum
)

func (tck TableCellKind) tagName() string {
	switch tck {
	case tableCellHeader:
		return "th"
	case tableCellDatum:
		return "td"
	default:
		return "p"
	}
}
