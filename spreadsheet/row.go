package spreadsheet

// Row represents a row in a worksheet
type Row struct {
	sheet  *Sheet
	index  int
	cells  map[int]*Cell
	height float64
}

func newRow(sheet *Sheet, index int) *Row {
	return &Row{
		sheet: sheet,
		index: index,
		cells: make(map[int]*Cell),
	}
}

// Cell returns or creates a cell at the given column (1-based)
func (r *Row) Cell(col int) *Cell {
	if c, ok := r.cells[col]; ok {
		return c
	}
	c := newCell(r, col)
	r.cells[col] = c
	r.sheet.updateMaxCol(col)
	return c
}

// Index returns the row index
func (r *Row) Index() int { return r.index }

// SetHeight sets the row height in points
func (r *Row) SetHeight(height float64) { r.height = height }

// Height returns the row height
func (r *Row) Height() float64 { return r.height }
