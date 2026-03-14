package spreadsheet

import (
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestCellTypes(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")

	// Empty
	c := s.Cell(1, 1)
	if c.Type() != CellTypeEmpty {
		t.Error("new cell should be empty")
	}
	if c.String() != "" {
		t.Error("empty cell string should be empty")
	}

	// String
	c.SetString("Hello")
	if c.Type() != CellTypeString {
		t.Error("should be string type")
	}
	if c.String() != "Hello" {
		t.Errorf("expected Hello, got %s", c.String())
	}

	// Number
	c2 := s.Cell(2, 1)
	c2.SetNumber(3.14)
	if c2.Type() != CellTypeNumber {
		t.Error("should be number type")
	}
	if c2.String() != "3.14" {
		t.Errorf("expected 3.14, got %s", c2.String())
	}

	// Boolean
	c3 := s.Cell(3, 1)
	c3.SetBool(true)
	if c3.Type() != CellTypeBoolean {
		t.Error("should be boolean type")
	}
	if c3.String() != "TRUE" {
		t.Errorf("expected TRUE, got %s", c3.String())
	}

	// Formula
	c4 := s.Cell(4, 1)
	c4.SetFormula("SUM(A1:A3)")
	if c4.Type() != CellTypeFormula {
		t.Error("should be formula type")
	}
	if c4.String() != "=SUM(A1:A3)" {
		t.Errorf("expected =SUM(A1:A3), got %s", c4.String())
	}
}

func TestCellRef2(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	c := s.Cell(3, 2)
	if c.Ref() != "B3" {
		t.Errorf("expected B3, got %s", c.Ref())
	}
}

func TestCellSetValue(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")

	// Int
	s.Cell(1, 1).SetValue(42)
	if s.Cell(1, 1).Type() != CellTypeNumber {
		t.Error("int should be number type")
	}

	// Int64
	s.Cell(2, 1).SetValue(int64(100))
	if s.Cell(2, 1).Type() != CellTypeNumber {
		t.Error("int64 should be number type")
	}

	// Float32
	s.Cell(3, 1).SetValue(float32(1.5))
	if s.Cell(3, 1).Type() != CellTypeNumber {
		t.Error("float32 should be number type")
	}
}

func TestCellFormatting(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	c := s.Cell(1, 1)

	f := common.NewFont("Arial", 12).Bold()
	c.SetFont(f)
	c.SetBackgroundColor(common.Yellow)
	c.SetBorders(common.ThinBorders(common.Black))
	c.SetHorizontalAlignment(common.TextAlignCenter)
	c.SetVerticalAlignment(common.VerticalAlignMiddle)
	c.SetNumberFormat("#,##0.00")

	// Just verify no panics and values are set
	if c.font == nil {
		t.Error("font should be set")
	}
	if c.bgColor == nil {
		t.Error("bg color should be set")
	}
	if c.borders == nil {
		t.Error("borders should be set")
	}
}

func TestMaxRowCol(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")

	s.Cell(5, 10).SetValue("X")
	if s.MaxRow() != 5 {
		t.Errorf("expected max row 5, got %d", s.MaxRow())
	}
	if s.MaxCol() != 10 {
		t.Errorf("expected max col 10, got %d", s.MaxCol())
	}
}
