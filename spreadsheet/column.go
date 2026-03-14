package spreadsheet

// Column represents column properties
type Column struct {
	index   int
	width   float64
	hidden  bool
	bestFit bool
	style   int
}

// SetColumnHidden hides/shows the column
func (s *Sheet) SetColumnHidden(col int, hidden bool) {
	s.ensureColumn(col)
	s.columns[col].hidden = hidden
}

// SetColumnBestFit marks column for auto-fit
func (s *Sheet) SetColumnBestFit(col int, bestFit bool) {
	s.ensureColumn(col)
	s.columns[col].bestFit = bestFit
}

// ColumnWidth returns the width of a column
func (s *Sheet) ColumnWidth(col int) float64 {
	if w, ok := s.colWidths[col]; ok {
		return w
	}
	return 8.43 // default Excel column width
}

// SetColumnWidthRange sets width for a range of columns
func (s *Sheet) SetColumnWidthRange(startCol, endCol int, width float64) {
	for col := startCol; col <= endCol; col++ {
		s.colWidths[col] = width
	}
}

func (s *Sheet) ensureColumn(col int) {
	if s.columns == nil {
		s.columns = make(map[int]*Column)
	}
	if _, ok := s.columns[col]; !ok {
		s.columns[col] = &Column{index: col, width: 8.43}
	}
}
