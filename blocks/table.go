package blocks

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Table is a table, which consists of several Rows and has a Caption.
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

func (t Table) ID(counter *IDCounter) string {
	counter.tables++
	return fmt.Sprintf("table-%d", counter.tables)
}

func (t Table) IsBlock() {}

var tableRe = regexp.MustCompile(`^table\s*{`)

func MatchesTable(line string) bool {
	return tableRe.MatchString(line)
}

func TableFromFirstLine(line, hyphaName string) Table {
	return Table{
		HyphaName: hyphaName,
		Caption:   line[strings.IndexRune(line, '{')+1:],
		Rows:      make([]*TableRow, 0),
	}
}

func (t *Table) ProcessLine(line string) (done bool) {
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
		case r == '\r':
			continue
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
			t.currCellBuilder.WriteRune(r)
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
			t.currCellBuilder.WriteRune(r)
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
		HyphaName: t.HyphaName,
		Cells:     []*TableCell{},
	})
}

func (t *Table) pushCell() {
	tc := &TableCell{
		Contents: MakeFormatted(t.currCellBuilder.String(), t.HyphaName),
		colspan:  t.currColspan,
	}
	switch t.currCellMarker {
	case '|', '^':
		tc.IsHeaderCell = false
	case '!':
		tc.IsHeaderCell = true
	}
	// We expect the table to have at least one row ready, so no nil-checking
	tr := t.Rows[len(t.Rows)-1]
	tr.Cells = append(tr.Cells, tc)
	t.currCellBuilder = strings.Builder{}
}

// TableRow is a row in a table. Thus, it can only be nested inside a table.
type TableRow struct {
	HyphaName string
	Cells     []*TableCell
}

func (tr TableRow) ID(_ *IDCounter) string {
	return ""
}

func (tr TableRow) IsBlock() {}

// LooksLikeThead is true if the table row looks like it might as well be a thead row.
//
// Most likely, rows with more than two header cells are theads. I allow one extra datum cell for tables like this:
// |   ! a ! b
// ! c | d | e
// ! f | g | h
func (tr *TableRow) LooksLikeThead() bool {
	var (
		headerAmount = 0
		datumAmount  = 0
	)
	for _, tc := range tr.Cells {
		if tc.IsHeaderCell {
			headerAmount++
		} else {
			datumAmount++
		}
	}
	return headerAmount >= 2 && datumAmount <= 1
}

// TableCell is a cell in TableRow.
type TableCell struct {
	IsHeaderCell bool
	Contents     Formatted
	colspan      uint
}

func (tc TableCell) ID(_ *IDCounter) string {
	return ""
}

func (tc TableCell) IsBlock() {}

// ColspanAttribute returns either an empty string (if the cell doesn't have colspan) or a string in this format:
//
//     colspan="<number here>"
func (tc *TableCell) ColspanAttribute() string {
	if tc.colspan <= 1 {
		return ""
	}
	return fmt.Sprintf(` colspan="%d"`, tc.colspan)
}

// TagName returns "th" if the cell is a header cell, "td" elsewise.
func (tc *TableCell) TagName() string {
	if tc.IsHeaderCell {
		return "th"
	}
	return "td"
}
