package common

import (
	"math"
	"testing"
)

func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func TestMeasurementConversions(t *testing.T) {
	// 1 inch = 72 points
	m := In(1)
	if m.Points() != 72 {
		t.Errorf("1 inch should be 72 points, got %f", m.Points())
	}

	// Points round-trip
	m2 := Pt(36)
	if !almostEqual(m2.Inches(), 0.5, 0.001) {
		t.Errorf("36 points should be 0.5 inches, got %f", m2.Inches())
	}

	// Centimeters
	m3 := Cm(2.54) // 1 inch
	if !almostEqual(m3.Inches(), 1.0, 0.01) {
		t.Errorf("2.54 cm should be ~1 inch, got %f", m3.Inches())
	}

	// Millimeters
	m4 := Mm(25.4) // 1 inch
	if !almostEqual(m4.Inches(), 1.0, 0.01) {
		t.Errorf("25.4 mm should be ~1 inch, got %f", m4.Inches())
	}

	// EMU
	m5 := EMU(914400) // 1 inch
	if !almostEqual(m5.Inches(), 1.0, 0.01) {
		t.Errorf("914400 EMU should be ~1 inch, got %f", m5.Inches())
	}
}

func TestPageSizes(t *testing.T) {
	// A4 should be 210mm x 297mm
	if !almostEqual(PageA4.Width.Millimeters(), 210, 0.1) {
		t.Errorf("A4 width should be 210mm, got %f", PageA4.Width.Millimeters())
	}
	if !almostEqual(PageA4.Height.Millimeters(), 297, 0.1) {
		t.Errorf("A4 height should be 297mm, got %f", PageA4.Height.Millimeters())
	}

	// Letter should be 8.5 x 11 inches
	if PageLetter.Width.Inches() != 8.5 {
		t.Errorf("Letter width should be 8.5 inches, got %f", PageLetter.Width.Inches())
	}
}

func TestMargins(t *testing.T) {
	m := NormalMargins()
	if m.Top.Inches() != 1 || m.Right.Inches() != 1 {
		t.Error("Normal margins should be 1 inch")
	}

	m2 := NarrowMargins()
	if m2.Top.Inches() != 0.5 {
		t.Error("Narrow margins should be 0.5 inches")
	}
}
