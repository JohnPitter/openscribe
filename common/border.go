package common

// BorderStyle represents the style of a border
type BorderStyle int

const (
	BorderStyleNone BorderStyle = iota
	BorderStyleSolid
	BorderStyleDashed
	BorderStyleDotted
	BorderStyleDouble
	BorderStyleGroove
	BorderStyleRidge
)

// Border represents a single border edge
type Border struct {
	Style BorderStyle
	Width Measurement
	Color Color
}

func NewBorder(style BorderStyle, width Measurement, color Color) Border {
	return Border{Style: style, Width: width, Color: color}
}

// Borders represents all four borders
type Borders struct {
	Top    Border
	Right  Border
	Bottom Border
	Left   Border
}

func NewBorders(top, right, bottom, left Border) Borders {
	return Borders{Top: top, Right: right, Bottom: bottom, Left: left}
}

func UniformBorders(b Border) Borders {
	return Borders{Top: b, Right: b, Bottom: b, Left: b}
}

func NoBorders() Borders {
	return Borders{}
}

func ThinBorders(color Color) Borders {
	b := NewBorder(BorderStyleSolid, Pt(0.5), color)
	return UniformBorders(b)
}
