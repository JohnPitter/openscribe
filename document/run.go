package document

import (
	"fmt"

	"github.com/JohnPitter/openscribe/common"
)

// Run represents a text run within a paragraph
type Run struct {
	text       string
	font       *common.Font
	bold       bool
	italic     bool
	underline  bool
	strike     bool
	color      *common.Color
	size       float64
	fontFamily string
	highlight  string
}

// NewRun creates a new text run
func NewRun() *Run {
	return &Run{}
}

// SetText sets the text content
func (r *Run) SetText(text string) *Run {
	r.text = text
	return r
}

// Text returns the text content
func (r *Run) Text() string {
	return r.text
}

// SetFont applies a font configuration
func (r *Run) SetFont(f common.Font) *Run {
	r.font = &f
	r.fontFamily = f.Family
	r.size = f.Size
	r.bold = f.Weight >= common.FontWeightBold
	r.italic = f.Style == common.FontStyleItalic
	r.underline = f.Decoration == common.TextDecorationUnderline
	r.color = &f.Color
	return r
}

// SetBold sets bold formatting
func (r *Run) SetBold(b bool) *Run { r.bold = b; return r }

// SetItalic sets italic formatting
func (r *Run) SetItalic(i bool) *Run { r.italic = i; return r }

// SetUnderline sets underline formatting
func (r *Run) SetUnderline(u bool) *Run { r.underline = u; return r }

// SetStrikethrough sets strikethrough formatting
func (r *Run) SetStrikethrough(s bool) *Run { r.strike = s; return r }

// SetColor sets the text color
func (r *Run) SetColor(c common.Color) *Run { r.color = &c; return r }

// SetSize sets the font size in points
func (r *Run) SetSize(size float64) *Run { r.size = size; return r }

// SetFontFamily sets the font family name
func (r *Run) SetFontFamily(family string) *Run { r.fontFamily = family; return r }

// SetHighlight sets text highlight color
func (r *Run) SetHighlight(color string) *Run { r.highlight = color; return r }

// toXML creates the w:r XML element
func (r *Run) toXML() xmlRun {
	xr := xmlRun{}

	// Run properties
	rPr := xmlRunProperties{}
	hasProps := false

	if r.bold {
		rPr.Bold = &xmlEmpty{}
		hasProps = true
	}
	if r.italic {
		rPr.Italic = &xmlEmpty{}
		hasProps = true
	}
	if r.underline {
		rPr.Underline = &xmlValue{Val: "single"}
		hasProps = true
	}
	if r.strike {
		rPr.Strike = &xmlEmpty{}
		hasProps = true
	}
	if r.color != nil {
		hex := r.color.Hex()
		if len(hex) > 0 && hex[0] == '#' {
			hex = hex[1:]
		}
		rPr.Color = &xmlValue{Val: hex}
		hasProps = true
	}
	if r.size > 0 {
		// OOXML uses half-points
		rPr.Size = &xmlValue{Val: fmt.Sprintf("%d", int(r.size*2))}
		hasProps = true
	}
	if r.fontFamily != "" {
		rPr.FontFamily = &xmlFontFamily{Ascii: r.fontFamily, HAnsi: r.fontFamily}
		hasProps = true
	}

	if hasProps {
		xr.Properties = &rPr
	}

	if r.text != "" {
		xr.Text = &xmlText{
			Space: "preserve",
			Value: r.text,
		}
	}

	return xr
}
