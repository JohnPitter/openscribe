package document

import "github.com/JohnPitter/openscribe/common"

// HeaderFooter represents document headers and footers
type HeaderFooter struct {
	leftText   string
	centerText string
	rightText  string
	font       common.Font
}

// Header returns the document header (creates if needed)
func (d *Document) Header() *HeaderFooter {
	if d.header == nil {
		d.header = &HeaderFooter{
			font: common.NewFont("Arial", 10),
		}
	}
	return d.header
}

// Footer returns the document footer (creates if needed)
func (d *Document) Footer() *HeaderFooter {
	if d.footer == nil {
		d.footer = &HeaderFooter{
			font: common.NewFont("Arial", 10),
		}
	}
	return d.footer
}

// SetLeft sets the left-aligned text
func (hf *HeaderFooter) SetLeft(text string) { hf.leftText = text }

// SetCenter sets the center-aligned text
func (hf *HeaderFooter) SetCenter(text string) { hf.centerText = text }

// SetRight sets the right-aligned text
func (hf *HeaderFooter) SetRight(text string) { hf.rightText = text }

// SetFont sets the header/footer font
func (hf *HeaderFooter) SetFont(f common.Font) { hf.font = f }

// Left returns left text
func (hf *HeaderFooter) Left() string { return hf.leftText }

// Center returns center text
func (hf *HeaderFooter) Center() string { return hf.centerText }

// Right returns right text
func (hf *HeaderFooter) Right() string { return hf.rightText }

// IsEmpty returns true if no text is set
func (hf *HeaderFooter) IsEmpty() bool {
	return hf.leftText == "" && hf.centerText == "" && hf.rightText == ""
}
