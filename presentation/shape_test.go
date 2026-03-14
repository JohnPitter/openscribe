package presentation

import (
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestShapeTypes(t *testing.T) {
	shapes := []ShapeType{
		ShapeRectangle, ShapeRoundedRectangle, ShapeCircle,
		ShapeEllipse, ShapeTriangle, ShapeArrowRight,
		ShapeStar, ShapeDiamond, ShapeLine,
	}

	p := New()
	s := p.AddSlide()

	for _, st := range shapes {
		sh := s.AddShape(st, common.In(1), common.In(1), common.In(2), common.In(2))
		if sh.Type() != st {
			t.Errorf("shape type mismatch: got %d, want %d", sh.Type(), st)
		}
	}

	// Verify it builds without error
	_, err := p.SaveToBytes()
	if err != nil {
		t.Fatalf("build error with all shapes: %v", err)
	}
}

func TestShapePosition(t *testing.T) {
	sh := &Shape{}
	sh.SetPosition(common.In(5), common.In(3))
	x, y := sh.Position()
	if x.Inches() != 5 || y.Inches() != 3 {
		t.Error("position should be 5,3 inches")
	}
}

func TestShapeSize(t *testing.T) {
	sh := &Shape{}
	sh.SetSize(common.In(4), common.In(2))
	w, h := sh.Size()
	if w.Inches() != 4 || h.Inches() != 2 {
		t.Error("size should be 4x2 inches")
	}
}

func TestShapeLine(t *testing.T) {
	sh := &Shape{shapeType: ShapeRectangle}
	sh.SetLine(common.Red, common.Pt(2))
	if sh.lineColor == nil || *sh.lineColor != common.Red {
		t.Error("line color should be red")
	}
}
