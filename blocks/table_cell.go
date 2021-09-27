package blocks

import "fmt"

// TableCell is a cell in TableRow.
type TableCell struct {
	isHeaderCell bool
	contents     []Block
	colspan      uint
}

// ID returns and empty string because table cells do not have ids.
func (tc TableCell) ID(_ *IDCounter) string {
	return ""
}

// NewTableCell returns TableCell with the given data. Parameter isHeaderCell should be true when the cell starts with !. Colspans 0 and 1 are the same: they mean that the cell does not span columns.
func NewTableCell(isHeaderCell bool, colspan uint, contents []Block) TableCell {
	return TableCell{
		isHeaderCell: isHeaderCell,
		contents:     contents,
		colspan:      colspan,
	}
}

// IsHeaderCell is true for header cells, i/e cells starting with !.
func (tc TableCell) IsHeaderCell() bool {
	return tc.isHeaderCell
}

// Contents returns the cell's contents, which may be any Mycomarkup blocks.
func (tc TableCell) Contents() []Block {
	return tc.contents
}

// Colspan returns how many columns the cell spans.
func (tc TableCell) Colspan() uint {
	return tc.colspan
}

// ColspanAttributeHTML returns either an empty string (if the cell doesn't have colspan) or a string in this format:
//
//     colspan="<number here>"
func (tc TableCell) ColspanAttributeHTML() string {
	if tc.colspan <= 1 {
		return ""
	}
	return fmt.Sprintf(` colspan="%d"`, tc.colspan)
}

// TagNameHTML returns "th" if the cell is a header cell, "td" otherwise.
func (tc TableCell) TagNameHTML() string {
	if tc.IsHeaderCell() {
		return "th"
	}
	return "td"
}
