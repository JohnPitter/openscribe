package document

import (
	"fmt"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/style"
)

// Table represents a table in the document
type Table struct {
	rows    []*TableRow
	cols    int
	borders common.Borders
	width   common.Measurement
}

// TableRow represents a table row
type TableRow struct {
	cells  []*TableCell
	height common.Measurement
}

// TableCell represents a table cell
type TableCell struct {
	paragraphs []*Paragraph
	width      common.Measurement
	borders    common.Borders
	shading    *common.Color
	vAlign     common.VerticalAlignment
	colSpan    int
	rowSpan    int
}

// NewTable creates a new table
func NewTable(rows, cols int) *Table {
	t := &Table{
		cols:    cols,
		borders: common.ThinBorders(common.Black),
	}
	for i := 0; i < rows; i++ {
		row := &TableRow{}
		for j := 0; j < cols; j++ {
			cell := &TableCell{
				paragraphs: []*Paragraph{NewParagraph()},
				colSpan:    1,
				rowSpan:    1,
			}
			row.cells = append(row.cells, cell)
		}
		t.rows = append(t.rows, row)
	}
	return t
}

// Rows returns all rows
func (t *Table) Rows() []*TableRow {
	return t.rows
}

// Cell returns a cell at the given position
func (t *Table) Cell(row, col int) *TableCell {
	if row < 0 || row >= len(t.rows) || col < 0 || col >= len(t.rows[row].cells) {
		return nil
	}
	return t.rows[row].cells[col]
}

// AddRow adds a new row
func (t *Table) AddRow() *TableRow {
	row := &TableRow{}
	for i := 0; i < t.cols; i++ {
		cell := &TableCell{
			paragraphs: []*Paragraph{NewParagraph()},
			colSpan:    1,
			rowSpan:    1,
		}
		row.cells = append(row.cells, cell)
	}
	t.rows = append(t.rows, row)
	return row
}

// RemoveRow removes a row by index
func (t *Table) RemoveRow(index int) error {
	if index < 0 || index >= len(t.rows) {
		return fmt.Errorf("row index %d out of range", index)
	}
	t.rows = append(t.rows[:index], t.rows[index+1:]...)
	return nil
}

// SetBorders sets table borders
func (t *Table) SetBorders(borders common.Borders) {
	t.borders = borders
}

// ApplyTheme applies a theme to the table
func (t *Table) ApplyTheme(theme style.Theme) {
	t.borders = common.ThinBorders(theme.Palette.Text)
}

// RowCount returns number of rows
func (t *Table) RowCount() int { return len(t.rows) }

// ColCount returns number of columns
func (t *Table) ColCount() int { return t.cols }

// Cells returns all cells in the row
func (r *TableRow) Cells() []*TableCell {
	return r.cells
}

// SetText sets the text in the first paragraph of a cell
func (c *TableCell) SetText(text string) {
	if len(c.paragraphs) == 0 {
		c.paragraphs = append(c.paragraphs, NewParagraph())
	}
	c.paragraphs[0] = NewParagraph()
	c.paragraphs[0].AddText(text)
}

// Text returns text of the first paragraph
func (c *TableCell) Text() string {
	if len(c.paragraphs) == 0 {
		return ""
	}
	return c.paragraphs[0].Text()
}

// AddParagraph adds a paragraph to the cell
func (c *TableCell) AddParagraph() *Paragraph {
	p := NewParagraph()
	c.paragraphs = append(c.paragraphs, p)
	return p
}

// SetShading sets cell background color
func (c *TableCell) SetShading(color common.Color) {
	c.shading = &color
}

// SetVerticalAlignment sets cell vertical alignment
func (c *TableCell) SetVerticalAlignment(align common.VerticalAlignment) {
	c.vAlign = align
}

// SetColSpan sets column span
func (c *TableCell) SetColSpan(span int) {
	if span < 1 {
		span = 1
	}
	c.colSpan = span
}
