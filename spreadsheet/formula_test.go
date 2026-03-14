package spreadsheet

import (
	"math"
	"testing"
)

func TestFormulaSum(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.Cell(1, 1).SetNumber(10)
	s.Cell(2, 1).SetNumber(20)
	s.Cell(3, 1).SetNumber(30)

	result, err := s.EvaluateFormula("SUM(A1:A3)")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result != 60 {
		t.Errorf("expected 60, got %f", result)
	}
}

func TestFormulaAverage(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.Cell(1, 1).SetNumber(10)
	s.Cell(2, 1).SetNumber(20)
	s.Cell(3, 1).SetNumber(30)

	result, err := s.EvaluateFormula("AVERAGE(A1:A3)")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result != 20 {
		t.Errorf("expected 20, got %f", result)
	}
}

func TestFormulaMin(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.Cell(1, 1).SetNumber(15)
	s.Cell(2, 1).SetNumber(5)
	s.Cell(3, 1).SetNumber(25)

	result, err := s.EvaluateFormula("MIN(A1:A3)")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result != 5 {
		t.Errorf("expected 5, got %f", result)
	}
}

func TestFormulaMax(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.Cell(1, 1).SetNumber(15)
	s.Cell(2, 1).SetNumber(5)
	s.Cell(3, 1).SetNumber(25)

	result, err := s.EvaluateFormula("MAX(A1:A3)")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result != 25 {
		t.Errorf("expected 25, got %f", result)
	}
}

func TestFormulaCount(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.Cell(1, 1).SetNumber(10)
	s.Cell(2, 1).SetNumber(20)
	s.Cell(3, 1).SetString("text") // not counted
	s.Cell(4, 1).SetNumber(30)

	result, err := s.EvaluateFormula("COUNT(A1:A4)")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result != 3 {
		t.Errorf("expected 3, got %f", result)
	}
}

func TestFormulaAbs(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.Cell(1, 1).SetNumber(-42)

	result, err := s.EvaluateFormula("ABS(A1)")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result != 42 {
		t.Errorf("expected 42, got %f", result)
	}
}

func TestFormulaRound(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.Cell(1, 1).SetNumber(3.14159)

	result, err := s.EvaluateFormula("ROUND(A1, 2)")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if math.Abs(result-3.14) > 0.001 {
		t.Errorf("expected 3.14, got %f", result)
	}
}

func TestFormulaWithLiteral(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")

	result, err := s.EvaluateFormula("ABS(-10)")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result != 10 {
		t.Errorf("expected 10, got %f", result)
	}
}

func TestFormulaErrors(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")

	_, err := s.EvaluateFormula("")
	if err == nil {
		t.Error("should error on empty formula")
	}

	_, err = s.EvaluateFormula("UNKNOWN(A1:A3)")
	if err == nil {
		t.Error("should error on unknown function")
	}
}

func TestParseCellRef(t *testing.T) {
	tests := []struct {
		ref     string
		wantCol int
		wantRow int
	}{
		{"A1", 1, 1},
		{"B3", 2, 3},
		{"Z1", 26, 1},
		{"AA1", 27, 1},
		{"C10", 3, 10},
	}
	for _, tt := range tests {
		col, row := parseCellRef(tt.ref)
		if col != tt.wantCol || row != tt.wantRow {
			t.Errorf("parseCellRef(%s) = (%d, %d), want (%d, %d)", tt.ref, col, row, tt.wantCol, tt.wantRow)
		}
	}
}
