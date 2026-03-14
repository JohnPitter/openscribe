package presentation

import (
	"fmt"

	"github.com/JohnPitter/openscribe/common"
)

// SlideLayout represents predefined slide layouts
type SlideLayout int

const (
	LayoutBlank SlideLayout = iota
	LayoutTitle
	LayoutTitleContent
	LayoutTwoColumn
	LayoutComparison
	LayoutTitleOnly
	LayoutSection
)

// Slide represents a single slide in a presentation
type Slide struct {
	presentation *Presentation
	index        int
	elements     []SlideElement
	layout       SlideLayout
	background   *common.Color
	notes        string
	transition   *Transition
}

// SlideElement is the interface for all slide elements
type SlideElement interface {
	elementType() string
}

func newSlide(pres *Presentation, index int) *Slide {
	return &Slide{
		presentation: pres,
		index:        index,
		layout:       LayoutBlank,
	}
}

// AddTextBox adds a text box to the slide
func (s *Slide) AddTextBox(x, y, width, height common.Measurement) *TextBox {
	tb := &TextBox{
		x: x, y: y, width: width, height: height,
	}
	s.elements = append(s.elements, tb)
	return tb
}

// AddShape adds a shape to the slide
func (s *Slide) AddShape(shapeType ShapeType, x, y, width, height common.Measurement) *Shape {
	sh := &Shape{
		shapeType: shapeType,
		x:         x,
		y:         y,
		width:     width,
		height:    height,
		fillColor: common.Blue,
	}
	s.elements = append(s.elements, sh)
	return sh
}

// AddImage adds an image to the slide
func (s *Slide) AddImage(imgData *common.ImageData, x, y, width, height common.Measurement) *SlideImage {
	img := &SlideImage{
		data:   imgData,
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
	s.elements = append(s.elements, img)
	return img
}

// SetLayout sets the slide layout
func (s *Slide) SetLayout(layout SlideLayout) { s.layout = layout }

// Layout returns the slide layout
func (s *Slide) Layout() SlideLayout { return s.layout }

// SetBackground sets the slide background color
func (s *Slide) SetBackground(color common.Color) { s.background = &color }

// Background returns the background color
func (s *Slide) Background() *common.Color { return s.background }

// SetNotes sets speaker notes
func (s *Slide) SetNotes(notes string) { s.notes = notes }

// Notes returns speaker notes
func (s *Slide) Notes() string { return s.notes }

// SetTransition sets the slide transition
func (s *Slide) SetTransition(t Transition) { s.transition = &t }

// Elements returns all elements on the slide
func (s *Slide) Elements() []SlideElement { return s.elements }

// ElementCount returns the number of elements
func (s *Slide) ElementCount() int { return len(s.elements) }

// RemoveElement removes an element by index
func (s *Slide) RemoveElement(index int) error {
	if index < 0 || index >= len(s.elements) {
		return fmt.Errorf("element index %d out of range", index)
	}
	s.elements = append(s.elements[:index], s.elements[index+1:]...)
	return nil
}
