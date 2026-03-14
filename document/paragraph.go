package document

import (
	"github.com/JohnPitter/openscribe/common"
)

// Paragraph represents a paragraph in a document
type Paragraph struct {
	runs      []*Run
	style     string
	alignment common.TextAlignment
	spacing   ParagraphSpacing
	indent    ParagraphIndent
	pageBreak bool
}

// ParagraphSpacing controls paragraph spacing
type ParagraphSpacing struct {
	Before common.Measurement
	After  common.Measurement
	Line   float64 // line spacing multiplier (1.0, 1.5, 2.0, etc.)
}

// ParagraphIndent controls paragraph indentation
type ParagraphIndent struct {
	Left      common.Measurement
	Right     common.Measurement
	FirstLine common.Measurement
	Hanging   common.Measurement
}

// NewParagraph creates a new paragraph
func NewParagraph() *Paragraph {
	return &Paragraph{
		spacing: ParagraphSpacing{
			Before: common.Pt(0),
			After:  common.Pt(8),
			Line:   1.15,
		},
	}
}

// AddRun adds a text run to the paragraph
func (p *Paragraph) AddRun() *Run {
	r := NewRun()
	p.runs = append(p.runs, r)
	return r
}

// AddText adds a text run with the given text
func (p *Paragraph) AddText(text string) *Run {
	r := NewRun()
	r.SetText(text)
	p.runs = append(p.runs, r)
	return r
}

// Runs returns all runs in the paragraph
func (p *Paragraph) Runs() []*Run {
	return p.runs
}

// SetStyle sets the paragraph style
func (p *Paragraph) SetStyle(style string) {
	p.style = style
}

// Style returns the paragraph style
func (p *Paragraph) Style() string {
	return p.style
}

// SetAlignment sets text alignment
func (p *Paragraph) SetAlignment(align common.TextAlignment) {
	p.alignment = align
}

// Alignment returns the alignment
func (p *Paragraph) Alignment() common.TextAlignment {
	return p.alignment
}

// SetSpacing sets paragraph spacing
func (p *Paragraph) SetSpacing(before, after common.Measurement, lineSpacing float64) {
	p.spacing = ParagraphSpacing{
		Before: before,
		After:  after,
		Line:   lineSpacing,
	}
}

// SetIndent sets paragraph indentation
func (p *Paragraph) SetIndent(left, right, firstLine common.Measurement) {
	p.indent = ParagraphIndent{
		Left:      left,
		Right:     right,
		FirstLine: firstLine,
	}
}

// AddPageBreak adds a page break within this paragraph
func (p *Paragraph) AddPageBreak() {
	p.pageBreak = true
}

// Text returns the concatenated text of all runs
func (p *Paragraph) Text() string {
	var text string
	for _, r := range p.runs {
		text += r.Text()
	}
	return text
}

// MarshalXML creates the w:p element
func (p *Paragraph) MarshalXML() xmlParagraph {
	xp := xmlParagraph{}

	// Paragraph properties
	pPr := xmlParagraphProperties{}
	hasProps := false

	if p.style != "" {
		pPr.Style = &xmlValue{Val: p.style}
		hasProps = true
	}

	alignStr := alignmentToString(p.alignment)
	if alignStr != "" {
		pPr.Justification = &xmlValue{Val: alignStr}
		hasProps = true
	}

	if hasProps {
		xp.Properties = &pPr
	}

	// Runs
	for _, r := range p.runs {
		xp.Runs = append(xp.Runs, r.MarshalXML())
	}

	// Page break
	if p.pageBreak {
		xp.Runs = append(xp.Runs, xmlRun{
			Break: &xmlBreak{Type: "page"},
		})
	}

	return xp
}

func alignmentToString(a common.TextAlignment) string {
	switch a {
	case common.TextAlignCenter:
		return "center"
	case common.TextAlignRight:
		return "right"
	case common.TextAlignJustify:
		return "both"
	default:
		return ""
	}
}
