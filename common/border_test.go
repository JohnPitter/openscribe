package common

import "testing"

func TestNewBorder(t *testing.T) {
	b := NewBorder(BorderStyleSolid, Pt(1), Red)
	if b.Style != BorderStyleSolid {
		t.Error("style should be solid")
	}
	if b.Width.Points() != 1 {
		t.Error("width should be 1pt")
	}
	if b.Color != Red {
		t.Error("color should be red")
	}
}

func TestNewBorders(t *testing.T) {
	top := NewBorder(BorderStyleSolid, Pt(1), Red)
	right := NewBorder(BorderStyleDashed, Pt(2), Blue)
	bottom := NewBorder(BorderStyleDotted, Pt(1), Green)
	left := NewBorder(BorderStyleDouble, Pt(3), Black)

	borders := NewBorders(top, right, bottom, left)
	if borders.Top.Style != BorderStyleSolid {
		t.Error("top should be solid")
	}
	if borders.Right.Style != BorderStyleDashed {
		t.Error("right should be dashed")
	}
	if borders.Bottom.Style != BorderStyleDotted {
		t.Error("bottom should be dotted")
	}
	if borders.Left.Style != BorderStyleDouble {
		t.Error("left should be double")
	}
}

func TestUniformBorders(t *testing.T) {
	b := NewBorder(BorderStyleSolid, Pt(1), Black)
	borders := UniformBorders(b)
	if borders.Top != borders.Bottom || borders.Left != borders.Right {
		t.Error("all borders should be equal")
	}
}

func TestNoBorders(t *testing.T) {
	borders := NoBorders()
	if borders.Top.Style != BorderStyleNone {
		t.Error("should be no borders")
	}
}

func TestThinBorders(t *testing.T) {
	borders := ThinBorders(Red)
	if borders.Top.Style != BorderStyleSolid {
		t.Error("should be solid")
	}
	if borders.Top.Color != Red {
		t.Error("should be red")
	}
	if borders.Top.Width.Points() != 0.5 {
		t.Error("should be 0.5pt")
	}
}

func TestBorderStyles(t *testing.T) {
	styles := []BorderStyle{
		BorderStyleNone, BorderStyleSolid, BorderStyleDashed,
		BorderStyleDotted, BorderStyleDouble, BorderStyleGroove, BorderStyleRidge,
	}
	for _, s := range styles {
		b := NewBorder(s, Pt(1), Black)
		if b.Style != s {
			t.Errorf("style mismatch: got %d, want %d", b.Style, s)
		}
	}
}
