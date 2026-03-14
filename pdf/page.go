package pdf

import (
	"github.com/JohnPitter/openscribe/common"
)

// Page represents a PDF page
type Page struct {
	size       common.PageSize
	margins    common.Margins
	elements   []PageElement
	rawData    []byte
	background *common.Color
}

// PageElement is the interface for all page elements
type PageElement interface {
	pdfElement()
}

// NewPage creates a new page with given size and margins
func NewPage(size common.PageSize, margins common.Margins) *Page {
	return &Page{
		size:    size,
		margins: margins,
	}
}

// Size returns the page size
func (p *Page) Size() common.PageSize { return p.size }

// SetSize sets the page size
func (p *Page) SetSize(size common.PageSize) { p.size = size }

// Margins returns the margins
func (p *Page) Margins() common.Margins { return p.margins }

// SetMargins sets the margins
func (p *Page) SetMargins(margins common.Margins) { p.margins = margins }

// SetBackground sets the page background color
func (p *Page) SetBackground(color common.Color) { p.background = &color }

// AddText adds text to the page
func (p *Page) AddText(text string, x, y float64, font common.Font) *TextElement {
	t := &TextElement{
		text: text,
		x:    x,
		y:    y,
		font: font,
	}
	p.elements = append(p.elements, t)
	return t
}

// AddLine adds a line to the page
func (p *Page) AddLine(x1, y1, x2, y2 float64, color common.Color, width float64) *LineElement {
	l := &LineElement{
		x1: x1, y1: y1, x2: x2, y2: y2,
		color: color, width: width,
	}
	p.elements = append(p.elements, l)
	return l
}

// AddRectangle adds a rectangle to the page
func (p *Page) AddRectangle(x, y, width, height float64, fill common.Color, stroke *common.Color) *RectElement {
	r := &RectElement{
		x: x, y: y, width: width, height: height,
		fill: fill, stroke: stroke,
	}
	p.elements = append(p.elements, r)
	return r
}

// AddTable adds a table to the page
func (p *Page) AddTable(x, y float64, rows, cols int) *TableElement {
	t := &TableElement{
		x: x, y: y, rows: rows, cols: cols,
		cellWidth:  100,
		cellHeight: 20,
		cells:      make([][]string, rows),
		font:       common.NewFont("Helvetica", 10),
	}
	for i := range t.cells {
		t.cells[i] = make([]string, cols)
	}
	p.elements = append(p.elements, t)
	return t
}

// AddImage adds an image to the page
func (p *Page) AddImage(imgData *common.ImageData, x, y, width, height float64) *ImageElement {
	img := &ImageElement{
		data:   imgData,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
	p.elements = append(p.elements, img)
	return img
}

// Elements returns all elements
func (p *Page) Elements() []PageElement { return p.elements }

// ElementCount returns the number of elements
func (p *Page) ElementCount() int { return len(p.elements) }
