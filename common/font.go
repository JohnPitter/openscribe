package common

// FontWeight represents the weight of a font
type FontWeight int

const (
	FontWeightThin      FontWeight = 100
	FontWeightLight     FontWeight = 300
	FontWeightRegular   FontWeight = 400
	FontWeightMedium    FontWeight = 500
	FontWeightSemiBold  FontWeight = 600
	FontWeightBold      FontWeight = 700
	FontWeightExtraBold FontWeight = 800
	FontWeightBlack     FontWeight = 900
)

// FontStyle represents italic/normal
type FontStyle int

const (
	FontStyleNormal FontStyle = iota
	FontStyleItalic
)

// TextDecoration represents underline/strikethrough
type TextDecoration int

const (
	TextDecorationNone          TextDecoration = iota
	TextDecorationUnderline
	TextDecorationStrikethrough
	TextDecorationDouble
)

// TextAlignment horizontal text alignment
type TextAlignment int

const (
	TextAlignLeft TextAlignment = iota
	TextAlignCenter
	TextAlignRight
	TextAlignJustify
)

// VerticalAlignment for cells and text boxes
type VerticalAlignment int

const (
	VerticalAlignTop VerticalAlignment = iota
	VerticalAlignMiddle
	VerticalAlignBottom
)

// Font describes a complete font configuration
type Font struct {
	Family     string
	Size       float64 // in points
	Weight     FontWeight
	Style      FontStyle
	Color      Color
	Decoration TextDecoration
}

func NewFont(family string, size float64) Font {
	return Font{
		Family: family,
		Size:   size,
		Weight: FontWeightRegular,
		Style:  FontStyleNormal,
		Color:  Black,
	}
}

func (f Font) WithWeight(w FontWeight) Font        { f.Weight = w; return f }
func (f Font) WithStyle(s FontStyle) Font           { f.Style = s; return f }
func (f Font) WithColor(c Color) Font               { f.Color = c; return f }
func (f Font) WithSize(s float64) Font              { f.Size = s; return f }
func (f Font) WithDecoration(d TextDecoration) Font { f.Decoration = d; return f }
func (f Font) Bold() Font                           { return f.WithWeight(FontWeightBold) }
func (f Font) Italic() Font                         { return f.WithStyle(FontStyleItalic) }
func (f Font) Underline() Font                      { return f.WithDecoration(TextDecorationUnderline) }
