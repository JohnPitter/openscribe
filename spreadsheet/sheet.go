package spreadsheet

import (
	"fmt"

	"github.com/JohnPitter/openscribe/common"
)

// Sheet represents a worksheet
type Sheet struct {
	workbook    *Workbook
	name        string
	index       int
	rows        map[int]*Row
	colWidths   map[int]float64
	columns     map[int]*Column
	mergedCells []MergedCell
	charts             []*Chart
	conditionalFormats []*ConditionalFormat
	maxRow             int
	maxCol             int
}

// MergedCell represents a merged cell range
type MergedCell struct {
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
}

func newSheet(wb *Workbook, name string, index int) *Sheet {
	return &Sheet{
		workbook:  wb,
		name:      name,
		index:     index,
		rows:      make(map[int]*Row),
		colWidths: make(map[int]float64),
		columns:   make(map[int]*Column),
	}
}

// Name returns the sheet name
func (s *Sheet) Name() string { return s.name }

// SetName renames the sheet
func (s *Sheet) SetName(name string) { s.name = name }

// Cell returns the cell at the given row/col (1-based)
func (s *Sheet) Cell(row, col int) *Cell {
	r := s.Row(row)
	return r.Cell(col)
}

// Row returns or creates a row (1-based)
func (s *Sheet) Row(index int) *Row {
	if r, ok := s.rows[index]; ok {
		return r
	}
	r := newRow(s, index)
	s.rows[index] = r
	if index > s.maxRow {
		s.maxRow = index
	}
	return r
}

// MaxRow returns the maximum row index used
func (s *Sheet) MaxRow() int { return s.maxRow }

// MaxCol returns the maximum column index used
func (s *Sheet) MaxCol() int { return s.maxCol }

// SetColumnWidth sets the width of a column (1-based)
func (s *Sheet) SetColumnWidth(col int, width float64) {
	s.colWidths[col] = width
}

// MergeCells merges a range of cells
func (s *Sheet) MergeCells(startRow, startCol, endRow, endCol int) {
	s.mergedCells = append(s.mergedCells, MergedCell{
		StartRow: startRow, StartCol: startCol,
		EndRow: endRow, EndCol: endCol,
	})
}

// SetValue is a convenience method to set a cell value
func (s *Sheet) SetValue(row, col int, value interface{}) {
	s.Cell(row, col).SetValue(value)
}

// Value is a convenience method to get a cell value
func (s *Sheet) Value(row, col int) interface{} {
	return s.Cell(row, col).Value()
}

// updateMaxCol tracks the maximum column
func (s *Sheet) updateMaxCol(col int) {
	if col > s.maxCol {
		s.maxCol = col
	}
}

// CellRef converts a 1-based row/col to an Excel cell reference like "A1"
func CellRef(row, col int) string {
	return fmt.Sprintf("%s%d", colName(col), row)
}

// colName converts a 1-based column index to column letters (1=A, 26=Z, 27=AA)
func colName(col int) string {
	name := ""
	for col > 0 {
		col--
		name = string(rune('A'+col%26)) + name
		col /= 26
	}
	return name
}

// SetCellBorders sets borders on a cell
func (s *Sheet) SetCellBorders(row, col int, borders common.Borders) {
	s.Cell(row, col).SetBorders(borders)
}
