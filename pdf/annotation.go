package pdf

import "github.com/JohnPitter/openscribe/common"

// AnnotationType represents the type of annotation
type AnnotationType int

const (
	AnnotHighlight AnnotationType = iota
	AnnotUnderline
	AnnotStrikeout
	AnnotStickyNote
	AnnotFreeText
)

// Annotation represents a PDF annotation
type Annotation struct {
	annotType AnnotationType
	x1, y1    float64
	x2, y2    float64
	color     common.Color
	text      string
	font      common.Font
	author    string
	subject   string
}

func (a *Annotation) pdfElement() {}

// Type returns the annotation type
func (a *Annotation) Type() AnnotationType { return a.annotType }

// SetAuthor sets the annotation author
func (a *Annotation) SetAuthor(author string) { a.author = author }

// Author returns the author
func (a *Annotation) Author() string { return a.author }

// SetSubject sets the annotation subject
func (a *Annotation) SetSubject(subject string) { a.subject = subject }

// Subject returns the subject
func (a *Annotation) Subject() string { return a.subject }

// Text returns the annotation text
func (a *Annotation) Text() string { return a.text }

// Color returns the annotation color
func (a *Annotation) Color() common.Color { return a.color }

// AddHighlight adds a highlight annotation to the page
func (p *Page) AddHighlight(x1, y1, x2, y2 float64, color common.Color) *Annotation {
	a := &Annotation{
		annotType: AnnotHighlight,
		x1:        x1,
		y1:        y1,
		x2:        x2,
		y2:        y2,
		color:     color,
	}
	p.elements = append(p.elements, a)
	return a
}

// AddStickyNote adds a sticky note annotation to the page
func (p *Page) AddStickyNote(x, y float64, text string, color common.Color) *Annotation {
	a := &Annotation{
		annotType: AnnotStickyNote,
		x1:        x,
		y1:        y,
		x2:        x + 24,
		y2:        y + 24,
		color:     color,
		text:      text,
	}
	p.elements = append(p.elements, a)
	return a
}

// AddFreeText adds a free text annotation to the page
func (p *Page) AddFreeText(x, y, width, height float64, text string, font common.Font) *Annotation {
	a := &Annotation{
		annotType: AnnotFreeText,
		x1:        x,
		y1:        y,
		x2:        x + width,
		y2:        y + height,
		color:     font.Color,
		text:      text,
		font:      font,
	}
	p.elements = append(p.elements, a)
	return a
}
