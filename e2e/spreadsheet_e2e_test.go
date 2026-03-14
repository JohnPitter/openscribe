package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/spreadsheet"
	"github.com/JohnPitter/openscribe/style"
)

func TestXlsxCreate(t *testing.T) {
	wb := spreadsheet.New()
	s := wb.AddSheet("Sheet1")
	s.SetValue(1, 1, "Hello")
	s.SetValue(1, 2, "World")

	path := filepath.Join(t.TempDir(), "create.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
	assertFileNotEmpty(t, path)
}

func TestXlsxCreateWithAllFeatures(t *testing.T) {
	wb := spreadsheet.New()

	// Multiple sheets
	s1 := wb.AddSheet("Data")
	s2 := wb.AddSheet("Summary")

	// String cells
	s1.SetValue(1, 1, "Product")
	s1.SetValue(1, 2, "Q1")
	s1.SetValue(1, 3, "Q2")
	s1.SetValue(1, 4, "Total")

	// Number cells
	s1.SetValue(2, 1, "Widget A")
	s1.SetValue(2, 2, 15000.0)
	s1.SetValue(2, 3, 18000.0)

	// Boolean
	s1.SetValue(3, 1, "Active")
	s1.SetValue(3, 2, true)

	// Formula
	s1.Cell(2, 4).SetFormula("SUM(B2:C2)")

	// Cell formatting
	headerFont := common.NewFont("Arial", 11).Bold()
	for col := 1; col <= 4; col++ {
		s1.Cell(1, col).SetFont(headerFont)
		s1.Cell(1, col).SetBackgroundColor(common.LightGray)
		s1.Cell(1, col).SetBorders(common.ThinBorders(common.Black))
		s1.Cell(1, col).SetHorizontalAlignment(common.TextAlignCenter)
	}

	// Column widths
	s1.SetColumnWidth(1, 20)
	s1.SetColumnWidth(2, 15)

	// Merged cells
	s1.MergeCells(5, 1, 5, 4)
	s1.SetValue(5, 1, "Merged Header")

	// Summary sheet
	s2.SetValue(1, 1, "Total Revenue")
	s2.Cell(1, 2).SetFormula("SUM(Data!B2:C2)")
	s2.Cell(1, 2).SetNumberFormat("#,##0.00")

	// Row height
	s1.Row(1).SetHeight(25)

	path := filepath.Join(t.TempDir(), "full_features.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestXlsxEdit(t *testing.T) {
	// Create initial
	wb := spreadsheet.New()
	s := wb.AddSheet("Data")
	s.SetValue(1, 1, "Original")
	s.SetValue(2, 1, "Data")

	path := filepath.Join(t.TempDir(), "edit.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Open and edit
	wb2, err := spreadsheet.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}

	// Add new sheet
	s2 := wb2.AddSheet("Added")
	s2.SetValue(1, 1, "New data")

	// Modify existing data
	if wb2.Sheet(0) != nil {
		wb2.Sheet(0).SetValue(3, 1, "Added row")
	}

	editedPath := filepath.Join(t.TempDir(), "edited.xlsx")
	if err := wb2.Save(editedPath); err != nil {
		t.Fatalf("save edited error: %v", err)
	}
	assertFileExists(t, editedPath)
}

func TestXlsxEditRemoveSheet(t *testing.T) {
	wb := spreadsheet.New()
	wb.AddSheet("Keep")
	wb.AddSheet("Remove")
	wb.AddSheet("AlsoKeep")

	if err := wb.RemoveSheet(1); err != nil {
		t.Fatalf("remove error: %v", err)
	}
	if wb.SheetCount() != 2 {
		t.Errorf("expected 2 sheets, got %d", wb.SheetCount())
	}
	if wb.Sheet(1).Name() != "AlsoKeep" {
		t.Error("wrong sheet after removal")
	}

	path := filepath.Join(t.TempDir(), "removed.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestXlsxDelete(t *testing.T) {
	wb := spreadsheet.New()
	wb.AddSheet("Test")

	path := filepath.Join(t.TempDir(), "delete.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)

	if err := spreadsheet.Delete(path); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestXlsxSaveToBytes(t *testing.T) {
	wb := spreadsheet.New()
	s := wb.AddSheet("Test")
	s.SetValue(1, 1, "Bytes")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestXlsxCellTypes(t *testing.T) {
	wb := spreadsheet.New()
	s := wb.AddSheet("Types")

	// All cell types
	s.Cell(1, 1).SetString("text")
	s.Cell(2, 1).SetNumber(42.5)
	s.Cell(3, 1).SetBool(true)
	s.Cell(4, 1).SetFormula("A1&A2")

	// Auto-detection
	s.SetValue(5, 1, "auto string")
	s.SetValue(6, 1, 100)
	s.SetValue(7, 1, int64(200))
	s.SetValue(8, 1, float32(3.14))
	s.SetValue(9, 1, false)

	if s.Cell(1, 1).Type() != spreadsheet.CellTypeString {
		t.Error("should be string")
	}
	if s.Cell(2, 1).Type() != spreadsheet.CellTypeNumber {
		t.Error("should be number")
	}
	if s.Cell(3, 1).Type() != spreadsheet.CellTypeBoolean {
		t.Error("should be boolean")
	}
	if s.Cell(4, 1).Type() != spreadsheet.CellTypeFormula {
		t.Error("should be formula")
	}

	path := filepath.Join(t.TempDir(), "types.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestXlsxCellRef(t *testing.T) {
	wb := spreadsheet.New()
	s := wb.AddSheet("Refs")
	cell := s.Cell(3, 2)
	if cell.Ref() != "B3" {
		t.Errorf("expected B3, got %s", cell.Ref())
	}
}

func TestXlsxSheetByName(t *testing.T) {
	wb := spreadsheet.New()
	wb.AddSheet("Alpha")
	wb.AddSheet("Beta")

	s := wb.SheetByName("Beta")
	if s == nil || s.Name() != "Beta" {
		t.Error("should find Beta sheet")
	}
	if wb.SheetByName("Gamma") != nil {
		t.Error("should return nil for unknown sheet")
	}
}

func TestXlsxSheetRename(t *testing.T) {
	wb := spreadsheet.New()
	s := wb.AddSheet("Old")
	s.SetName("New")
	if s.Name() != "New" {
		t.Errorf("expected New, got %s", s.Name())
	}
}

func TestXlsxMaxRowCol(t *testing.T) {
	wb := spreadsheet.New()
	s := wb.AddSheet("Test")
	s.Cell(10, 5).SetValue("X")

	if s.MaxRow() != 10 {
		t.Errorf("expected max row 10, got %d", s.MaxRow())
	}
	if s.MaxCol() != 5 {
		t.Errorf("expected max col 5, got %d", s.MaxCol())
	}
}

func TestXlsxWithThemes(t *testing.T) {
	themes := style.AllThemes()

	for _, theme := range themes {
		t.Run(theme.Name, func(t *testing.T) {
			wb := spreadsheet.NewWithTheme(theme)
			s := wb.AddSheet("Data")
			s.SetValue(1, 1, "Themed spreadsheet")
			s.SetValue(2, 1, 100.0)

			path := filepath.Join(t.TempDir(), "themed.xlsx")
			if err := wb.Save(path); err != nil {
				t.Fatalf("save error with theme %s: %v", theme.Name, err)
			}
			assertFileExists(t, path)
		})
	}
}

func TestXlsxRoundTrip(t *testing.T) {
	wb := spreadsheet.New()
	s := wb.AddSheet("RoundTrip")
	s.SetValue(1, 1, "Name")
	s.SetValue(1, 2, "Score")
	s.SetValue(2, 1, "Alice")
	s.SetValue(2, 2, 95.5)
	s.SetValue(3, 1, "Bob")
	s.SetValue(3, 2, 87.0)

	path := filepath.Join(t.TempDir(), "roundtrip.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	wb2, err := spreadsheet.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	if wb2.SheetCount() != 1 {
		t.Errorf("expected 1 sheet, got %d", wb2.SheetCount())
	}
	if wb2.Sheet(0).Name() != "RoundTrip" {
		t.Error("sheet name should be preserved")
	}
}
