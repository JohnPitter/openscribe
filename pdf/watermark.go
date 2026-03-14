package pdf

import "github.com/JohnPitter/openscribe/common"

// Watermark represents a text watermark
type Watermark struct {
	text     string
	font     common.Font
	color    common.Color
	opacity  float64
	rotation float64
}

// NewWatermark creates a new watermark
func NewWatermark(text string) *Watermark {
	return &Watermark{
		text:     text,
		font:     common.NewFont("Helvetica", 48),
		color:    common.NewColorWithAlpha(200, 200, 200, 128),
		opacity:  0.3,
		rotation: -45,
	}
}

// SetFont sets the watermark font
func (w *Watermark) SetFont(f common.Font) { w.font = f }

// SetColor sets the watermark color
func (w *Watermark) SetColor(c common.Color) { w.color = c }

// SetOpacity sets the opacity (0.0 to 1.0)
func (w *Watermark) SetOpacity(o float64) {
	if o < 0 {
		o = 0
	}
	if o > 1 {
		o = 1
	}
	w.opacity = o
}

// SetRotation sets rotation in degrees
func (w *Watermark) SetRotation(deg float64) { w.rotation = deg }

// Text returns the watermark text
func (w *Watermark) Text() string { return w.text }

// Opacity returns the opacity
func (w *Watermark) Opacity() float64 { return w.opacity }

// Rotation returns the rotation in degrees
func (w *Watermark) Rotation() float64 { return w.rotation }

// AddWatermark adds a watermark to all pages
func (d *Document) AddWatermark(w *Watermark) {
	for _, page := range d.pages {
		// Add watermark as a centered text element with the watermark font
		pageW := page.size.Width.Points()
		pageH := page.size.Height.Points()
		page.AddText(w.text, pageW/2-100, pageH/2, w.font.WithColor(w.color))
	}
}
