package presentation

import "github.com/JohnPitter/openscribe/common"

// ShapeType represents predefined shape types
type ShapeType int

const (
	ShapeRectangle ShapeType = iota
	ShapeRoundedRectangle
	ShapeCircle
	ShapeEllipse
	ShapeTriangle
	ShapeArrowRight
	ShapeArrowLeft
	ShapeArrowUp
	ShapeArrowDown
	ShapeStar
	ShapeDiamond
	ShapeLine
)

// Shape represents a shape on a slide
type Shape struct {
	shapeType ShapeType
	x, y      common.Measurement
	width     common.Measurement
	height    common.Measurement
	fillColor common.Color
	lineColor *common.Color
	lineWidth common.Measurement
	rotation  float64
	text      string
	textFont  *common.Font
}

func (s *Shape) elementType() string { return "shape" }

// SetPosition sets the shape position
func (s *Shape) SetPosition(x, y common.Measurement) { s.x = x; s.y = y }

// SetSize sets the shape size
func (s *Shape) SetSize(width, height common.Measurement) { s.width = width; s.height = height }

// SetFill sets the fill color
func (s *Shape) SetFill(color common.Color) { s.fillColor = color }

// SetLine sets the outline
func (s *Shape) SetLine(color common.Color, width common.Measurement) {
	s.lineColor = &color
	s.lineWidth = width
}

// SetRotation sets rotation in degrees
func (s *Shape) SetRotation(degrees float64) { s.rotation = degrees }

// SetText sets text inside the shape
func (s *Shape) SetText(text string, font common.Font) {
	s.text = text
	s.textFont = &font
}

// Type returns the shape type
func (s *Shape) Type() ShapeType { return s.shapeType }

// Position returns x, y
func (s *Shape) Position() (common.Measurement, common.Measurement) { return s.x, s.y }

// Size returns width, height
func (s *Shape) Size() (common.Measurement, common.Measurement) { return s.width, s.height }
