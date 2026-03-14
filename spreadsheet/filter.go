package spreadsheet

import "fmt"

// AutoFilter represents an auto-filter on a sheet
type AutoFilter struct {
	startRow int
	startCol int
	endRow   int
	endCol   int
}

// SetAutoFilter sets an auto-filter on the specified range
func (s *Sheet) SetAutoFilter(startRow, startCol, endRow, endCol int) {
	s.autoFilter = &AutoFilter{
		startRow: startRow,
		startCol: startCol,
		endRow:   endRow,
		endCol:   endCol,
	}
}

// AutoFilter returns the auto-filter, or nil if none is set
func (s *Sheet) AutoFilter() *AutoFilter {
	return s.autoFilter
}

// Ref returns the cell range reference for the auto-filter (e.g., "A1:D10")
func (af *AutoFilter) Ref() string {
	return fmt.Sprintf("%s:%s", CellRef(af.startRow, af.startCol), CellRef(af.endRow, af.endCol))
}

// buildAutoFilterXML generates the <autoFilter> XML element
func buildAutoFilterXML(af *AutoFilter) string {
	if af == nil {
		return ""
	}
	return fmt.Sprintf(`<autoFilter ref="%s"/>`, af.Ref())
}
