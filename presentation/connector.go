package presentation

import "github.com/JohnPitter/openscribe/common"

// ConnectorType represents the type of connector
type ConnectorType int

const (
	ConnectorStraight ConnectorType = iota
	ConnectorElbow
	ConnectorCurved
)

// Additional shape types for smart shapes
const (
	ShapeCallout             ShapeType = iota + 100
	ShapeFlowchartProcess    ShapeType = 101
	ShapeFlowchartDecision   ShapeType = 102
	ShapeFlowchartTerminator ShapeType = 103
	ShapeBrace               ShapeType = 104
	ShapeBracket             ShapeType = 105
)

// Connector represents a connector line between points on a slide
type Connector struct {
	connType ConnectorType
	x1, y1   common.Measurement
	x2, y2   common.Measurement
	color    common.Color
	width    common.Measurement
}

func (c *Connector) elementType() string { return "connector" }

// AddConnector adds a connector to the slide
func (s *Slide) AddConnector(connType ConnectorType, x1, y1, x2, y2 common.Measurement) *Connector {
	c := &Connector{
		connType: connType,
		x1:       x1,
		y1:       y1,
		x2:       x2,
		y2:       y2,
		color:    common.Black,
		width:    common.Pt(1),
	}
	s.elements = append(s.elements, c)
	return c
}

// SetColor sets the connector line color
func (c *Connector) SetColor(color common.Color) { c.color = color }

// SetWidth sets the connector line width
func (c *Connector) SetWidth(w common.Measurement) { c.width = w }

// Type returns the connector type
func (c *Connector) Type() ConnectorType { return c.connType }

// Color returns the connector color
func (c *Connector) Color() common.Color { return c.color }

// Width returns the connector width
func (c *Connector) Width() common.Measurement { return c.width }
