package spreadsheet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/style"
)

func TestNewWorkbook(t *testing.T) {
	wb := New()
	if wb == nil {
		t.Fatal("workbook should not be nil")
	}
	if wb.SheetCount() != 0 {
		t.Error("new workbook should have no sheets")
	}
}

func TestNewWithTheme(t *testing.T) {
	theme := style.PremiumModern()
	wb := NewWithTheme(theme)
	if wb.Theme().Name != "Premium Modern" {
		t.Errorf("expected Premium Modern, got %s", wb.Theme().Name)
	}
}

func TestAddSheet(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Sheet1")
	if s == nil {
		t.Fatal("sheet should not be nil")
	}
	if s.Name() != "Sheet1" {
		t.Errorf("expected Sheet1, got %s", s.Name())
	}
	if wb.SheetCount() != 1 {
		t.Errorf("expected 1 sheet, got %d", wb.SheetCount())
	}
}

func TestSheetByName(t *testing.T) {
	wb := New()
	wb.AddSheet("Data")
	wb.AddSheet("Summary")

	s := wb.SheetByName("Summary")
	if s == nil || s.Name() != "Summary" {
		t.Error("should find sheet by name")
	}

	if wb.SheetByName("NotExist") != nil {
		t.Error("should return nil for non-existent sheet")
	}
}

func TestRemoveSheet(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet1")
	wb.AddSheet("Sheet2")

	err := wb.RemoveSheet(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wb.SheetCount() != 1 {
		t.Errorf("expected 1 sheet, got %d", wb.SheetCount())
	}
	if wb.Sheet(0).Name() != "Sheet2" {
		t.Error("remaining sheet should be Sheet2")
	}

	err = wb.RemoveSheet(10)
	if err == nil {
		t.Error("should error on out of range")
	}
}

func TestCellValues(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")

	// String
	s.SetValue(1, 1, "Hello")
	if s.Value(1, 1) != "Hello" {
		t.Errorf("expected Hello, got %v", s.Value(1, 1))
	}

	// Number
	s.SetValue(2, 1, 42.5)
	if s.Value(2, 1) != 42.5 {
		t.Errorf("expected 42.5, got %v", s.Value(2, 1))
	}

	// Boolean
	s.SetValue(3, 1, true)
	if s.Value(3, 1) != true {
		t.Errorf("expected true, got %v", s.Value(3, 1))
	}

	// Formula
	s.Cell(4, 1).SetFormula("SUM(A1:A3)")
	if s.Cell(4, 1).Type() != CellTypeFormula {
		t.Error("should be formula type")
	}
}

func TestCellRef(t *testing.T) {
	tests := []struct {
		row, col int
		want     string
	}{
		{1, 1, "A1"},
		{1, 2, "B1"},
		{1, 26, "Z1"},
		{1, 27, "AA1"},
		{10, 3, "C10"},
	}
	for _, tt := range tests {
		got := CellRef(tt.row, tt.col)
		if got != tt.want {
			t.Errorf("CellRef(%d,%d) = %s, want %s", tt.row, tt.col, got, tt.want)
		}
	}
}

func TestColName(t *testing.T) {
	tests := []struct {
		col  int
		want string
	}{
		{1, "A"},
		{26, "Z"},
		{27, "AA"},
		{28, "AB"},
		{52, "AZ"},
		{53, "BA"},
	}
	for _, tt := range tests {
		got := colName(tt.col)
		if got != tt.want {
			t.Errorf("colName(%d) = %s, want %s", tt.col, got, tt.want)
		}
	}
}

func TestSaveAndOpen(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")
	s.SetValue(1, 1, "Name")
	s.SetValue(1, 2, "Age")
	s.SetValue(2, 1, "Alice")
	s.SetValue(2, 2, 30.0)
	s.SetValue(3, 1, "Bob")
	s.SetValue(3, 2, 25.0)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.xlsx")

	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}

	// Open and verify
	wb2, err := Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	if wb2.SheetCount() != 1 {
		t.Errorf("expected 1 sheet, got %d", wb2.SheetCount())
	}
	if wb2.Sheet(0).Name() != "Data" {
		t.Errorf("expected sheet name Data, got %s", wb2.Sheet(0).Name())
	}
}

func TestSaveToBytes(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.SetValue(1, 1, "Hello")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestDelete(t *testing.T) {
	wb := New()
	wb.AddSheet("Test")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "delete_me.xlsx")

	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	if err := Delete(path); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestMergeCells(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Merge")
	s.SetValue(1, 1, "Merged Header")
	s.MergeCells(1, 1, 1, 3)

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("should produce output")
	}
}

func TestSheetRename(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Old")
	s.SetName("New")
	if s.Name() != "New" {
		t.Errorf("expected New, got %s", s.Name())
	}
}
