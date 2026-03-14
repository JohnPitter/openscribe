package document

import "github.com/JohnPitter/openscribe/common"

// Section represents a document section with page settings
type Section struct {
	pageSize    common.PageSize
	orientation common.Orientation
	margins     common.Margins
}

// NewSection creates a new section with default A4 portrait settings
func NewSection() *Section {
	return &Section{
		pageSize:    common.PageA4,
		orientation: common.OrientationPortrait,
		margins:     common.NormalMargins(),
	}
}

// SetPageSize sets the page size
func (s *Section) SetPageSize(size common.PageSize) {
	s.pageSize = size
}

// PageSize returns the page size
func (s *Section) PageSize() common.PageSize {
	return s.pageSize
}

// SetOrientation sets page orientation
func (s *Section) SetOrientation(o common.Orientation) {
	s.orientation = o
	if o == common.OrientationLandscape {
		// Swap width and height
		s.pageSize.Width, s.pageSize.Height = s.pageSize.Height, s.pageSize.Width
	}
}

// Orientation returns the page orientation
func (s *Section) Orientation() common.Orientation {
	return s.orientation
}

// SetMargins sets page margins
func (s *Section) SetMargins(m common.Margins) {
	s.margins = m
}

// Margins returns the margins
func (s *Section) Margins() common.Margins {
	return s.margins
}
