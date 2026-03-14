package presentation

import "github.com/JohnPitter/openscribe/common"

// SlideTable represents a table on a slide
type SlideTable struct {
	rows        int
	cols        int
	x, y        common.Measurement
	width       common.Measurement
	height      common.Measurement
	cells       [][]*SlideTableCell
	headerBg    *common.Color
	borderColor *common.Color
}

// SlideTableCell represents a single cell in a slide table
type SlideTableCell struct {
	text       string
	font       *common.Font
	background *common.Color
}

func (t *SlideTable) elementType() string { return "table" }

// AddTable adds a table to the slide
func (s *Slide) AddTable(rows, cols int, x, y, width, height common.Measurement) *SlideTable {
	cells := make([][]*SlideTableCell, rows)
	for r := range cells {
		cells[r] = make([]*SlideTableCell, cols)
		for c := range cells[r] {
			cells[r][c] = &SlideTableCell{}
		}
	}
	t := &SlideTable{
		rows:   rows,
		cols:   cols,
		x:      x,
		y:      y,
		width:  width,
		height: height,
		cells:  cells,
	}
	s.elements = append(s.elements, t)
	return t
}

// Cell returns the cell at the given row and column (0-based)
func (t *SlideTable) Cell(row, col int) *SlideTableCell {
	if row < 0 || row >= t.rows || col < 0 || col >= t.cols {
		return nil
	}
	return t.cells[row][col]
}

// Rows returns the number of rows
func (t *SlideTable) Rows() int { return t.rows }

// Cols returns the number of columns
func (t *SlideTable) Cols() int { return t.cols }

// SetHeaderBackground sets the background color for the first row
func (t *SlideTable) SetHeaderBackground(c common.Color) { t.headerBg = &c }

// SetBorderColor sets the border color for the table
func (t *SlideTable) SetBorderColor(c common.Color) { t.borderColor = &c }

// SetText sets the cell text
func (c *SlideTableCell) SetText(text string) { c.text = text }

// SetFont sets the cell font
func (c *SlideTableCell) SetFont(f common.Font) { c.font = &f }

// SetBackground sets the cell background color
func (c *SlideTableCell) SetBackground(color common.Color) { c.background = &color }

// Text returns the cell text
func (c *SlideTableCell) Text() string { return c.text }
