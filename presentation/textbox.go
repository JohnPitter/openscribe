package presentation

import (
	"github.com/JohnPitter/openscribe/common"
)

// TextBox represents a text box on a slide
type TextBox struct {
	x, y, width, height common.Measurement
	paragraphs          []*TextParagraph
	fillColor           *common.Color
	borderColor         *common.Color
	borderWidth         common.Measurement
}

// TextParagraph represents a paragraph within a text box
type TextParagraph struct {
	runs      []*TextRun
	alignment common.TextAlignment
	spacing   float64
}

// TextRun represents a text run within a paragraph
type TextRun struct {
	text      string
	font      common.Font
	bold      bool
	italic    bool
	underline bool
	color     common.Color
}

func (tb *TextBox) elementType() string { return "textbox" }

// AddParagraph adds a paragraph to the text box
func (tb *TextBox) AddParagraph() *TextParagraph {
	p := &TextParagraph{spacing: 1.0}
	tb.paragraphs = append(tb.paragraphs, p)
	return p
}

// SetText sets simple text (single paragraph, single run)
func (tb *TextBox) SetText(text string, font common.Font) {
	tb.paragraphs = nil
	p := tb.AddParagraph()
	p.AddRun(text, font)
}

// Paragraphs returns all paragraphs
func (tb *TextBox) Paragraphs() []*TextParagraph { return tb.paragraphs }

// SetPosition sets the text box position
func (tb *TextBox) SetPosition(x, y common.Measurement) { tb.x = x; tb.y = y }

// SetSize sets the text box size
func (tb *TextBox) SetSize(width, height common.Measurement) { tb.width = width; tb.height = height }

// SetFill sets the fill color
func (tb *TextBox) SetFill(color common.Color) { tb.fillColor = &color }

// SetBorder sets the border
func (tb *TextBox) SetBorder(color common.Color, width common.Measurement) {
	tb.borderColor = &color
	tb.borderWidth = width
}

// Text returns the concatenated text
func (tb *TextBox) Text() string {
	var result string
	for i, p := range tb.paragraphs {
		if i > 0 {
			result += "\n"
		}
		for _, r := range p.runs {
			result += r.text
		}
	}
	return result
}

// AddRun adds a text run to the paragraph
func (p *TextParagraph) AddRun(text string, font common.Font) *TextRun {
	r := &TextRun{
		text:   text,
		font:   font,
		color:  font.Color,
		bold:   font.Weight >= common.FontWeightBold,
		italic: font.Style == common.FontStyleItalic,
	}
	p.runs = append(p.runs, r)
	return r
}

// SetAlignment sets paragraph alignment
func (p *TextParagraph) SetAlignment(a common.TextAlignment) { p.alignment = a }

// SetSpacing sets line spacing multiplier
func (p *TextParagraph) SetSpacing(s float64) { p.spacing = s }

// Runs returns all runs
func (p *TextParagraph) Runs() []*TextRun { return p.runs }

// SetText sets the run text
func (r *TextRun) SetText(text string) { r.text = text }

// SetBold sets bold
func (r *TextRun) SetBold(b bool) { r.bold = b }

// SetItalic sets italic
func (r *TextRun) SetItalic(i bool) { r.italic = i }

// SetColor sets text color
func (r *TextRun) SetColor(c common.Color) { r.color = c }
