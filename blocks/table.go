package blocks

import (
	"fmt"
)

// Table is a table, which consists of several rows and has a caption.
type Table struct {
	caption string // Make it Formatted maybe?
	rows    []TableRow
}

// ID returns table's id which is table- and its number.
func (t Table) ID(counter *IDCounter) string {
	counter.tables++
	return fmt.Sprintf("table-%d", counter.tables)
}

// NewTable returns a new Table with the given rows and caption.
func NewTable(caption string, rows []TableRow) Table {
	return Table{
		caption: caption,
		rows:    rows,
	}
}

// Caption returns Table's caption. It may be empty.
func (t Table) Caption() string {
	return t.caption
}

// Rows returns Table's rows.
func (t Table) Rows() []TableRow {
	return t.rows
}

// WithNewRow returns a new table but with a new row.
func (t Table) WithNewRow(row TableRow) Table {
	t.rows = append(t.rows, row)
	return t
}

// TableRow is a row in a table. Thus, it can only be nested inside a table.
type TableRow struct {
	cells []TableCell
}

// ID returns and empty string because table rows do not have ids.
func (tr TableRow) ID(_ *IDCounter) string {
	return ""
}

// NewTableRow returns a new TableRow. It gets the hypha name from ctx.
func NewTableRow(cells []TableCell) TableRow {
	return TableRow{
		cells: cells,
	}
}

// Cells returns the row's cells.
func (tr TableRow) Cells() []TableCell {
	return tr.cells
}

// LooksLikeThead is true if the table row looks like it might as well be a thead row.
//
// Most likely, rows with more than two header cells are theads. I allow one extra datum cell for tables like this:
// |   ! a ! b
// ! c | d | e
// ! f | g | h
func (tr TableRow) LooksLikeThead() bool {
	var (
		headerAmount = 0
		datumAmount  = 0
	)
	for _, tc := range tr.Cells() {
		if tc.IsHeaderCell() {
			headerAmount++
		} else {
			datumAmount++
		}
	}
	return headerAmount >= 2 && datumAmount <= 1
}
