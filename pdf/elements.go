package pdf

import "github.com/JohnPitter/openscribe/common"

// TextElement represents text on a PDF page
type TextElement struct {
	text string
	x, y float64
	font common.Font
}

func (t *TextElement) pdfElement() {}

// SetText sets the text content
func (t *TextElement) SetText(text string) { t.text = text }

// Text returns the text content
func (t *TextElement) Text() string { return t.text }

// SetPosition sets the position
func (t *TextElement) SetPosition(x, y float64) { t.x = x; t.y = y }

// SetFont sets the font
func (t *TextElement) SetFont(f common.Font) { t.font = f }

// LineElement represents a line
type LineElement struct {
	x1, y1, x2, y2 float64
	color           common.Color
	width           float64
}

func (l *LineElement) pdfElement() {}

// RectElement represents a rectangle
type RectElement struct {
	x, y, width, height float64
	fill                common.Color
	stroke              *common.Color
	cornerRadius        float64
}

func (r *RectElement) pdfElement() {}

// SetCornerRadius sets rounded corners
func (r *RectElement) SetCornerRadius(radius float64) { r.cornerRadius = radius }

// TableElement represents a table
type TableElement struct {
	x, y        float64
	rows, cols  int
	cellWidth   float64
	cellHeight  float64
	cells       [][]string
	font        common.Font
	headerBg    *common.Color
	borderColor common.Color
}

func (t *TableElement) pdfElement() {}

// SetCell sets a cell value
func (t *TableElement) SetCell(row, col int, value string) {
	if row >= 0 && row < t.rows && col >= 0 && col < t.cols {
		t.cells[row][col] = value
	}
}

// Cell returns a cell value
func (t *TableElement) Cell(row, col int) string {
	if row >= 0 && row < t.rows && col >= 0 && col < t.cols {
		return t.cells[row][col]
	}
	return ""
}

// SetCellSize sets cell dimensions
func (t *TableElement) SetCellSize(width, height float64) {
	t.cellWidth = width
	t.cellHeight = height
}

// SetFont sets the table font
func (t *TableElement) SetFont(f common.Font) { t.font = f }

// SetHeaderBackground sets the header row background
func (t *TableElement) SetHeaderBackground(color common.Color) { t.headerBg = &color }

// SetBorderColor sets the border color
func (t *TableElement) SetBorderColor(color common.Color) { t.borderColor = color }

// Rows returns the number of rows
func (t *TableElement) Rows() int { return t.rows }

// Cols returns the number of columns
func (t *TableElement) Cols() int { return t.cols }

// ImageElement represents an image on a PDF page
type ImageElement struct {
	data   *common.ImageData
	x, y   float64
	width  float64
	height float64
}

func (img *ImageElement) pdfElement() {}

// SetPosition sets the image position
func (img *ImageElement) SetPosition(x, y float64) { img.x = x; img.y = y }

// SetSize sets the image dimensions
func (img *ImageElement) SetSize(w, h float64) { img.width = w; img.height = h }
