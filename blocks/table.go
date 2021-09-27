package blocks

import (
	"fmt"
)

// Table is a table, which consists of several Rows and has a Caption.
type Table struct {
	// data
	HyphaName string
	Caption   string
	Rows      []*TableRow // V3
}

// ID returns table's id which is table- and its number.
func (t Table) ID(counter *IDCounter) string {
	counter.tables++
	return fmt.Sprintf("table-%d", counter.tables)
}

func (t Table) isBlock() {}

// TableRow is a row in a table. Thus, it can only be nested inside a table.
type TableRow struct {
	HyphaName string
	Cells     []*TableCell // V3
}

// ID returns and empty string because table rows do not have ids.
func (tr TableRow) ID(_ *IDCounter) string {
	return ""
}

func (tr TableRow) isBlock() {}

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
	Contents     []Block
	Colspan      uint
}

// ID returns and empty string because table cells do not have ids.
func (tc TableCell) ID(_ *IDCounter) string {
	return ""
}

func (tc TableCell) isBlock() {}

// ColspanAttribute returns either an empty string (if the cell doesn't have colspan) or a string in this format:
//
//     colspan="<number here>"
func (tc *TableCell) ColspanAttribute() string {
	if tc.Colspan <= 1 {
		return ""
	}
	return fmt.Sprintf(` colspan="%d"`, tc.Colspan)
}

// TagName returns "th" if the cell is a header cell, "td" otherise.
func (tc *TableCell) TagName() string {
	if tc.IsHeaderCell {
		return "th"
	}
	return "td"
}
