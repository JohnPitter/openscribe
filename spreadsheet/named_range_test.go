package spreadsheet

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddNamedRange(t *testing.T) {
	wb := New()
	wb.AddSheet("Data")

	nr := wb.AddNamedRange("SalesData", "Data", "$A$1:$D$100")
	if nr.Name() != "SalesData" {
		t.Errorf("expected SalesData, got %s", nr.Name())
	}
	if nr.SheetName() != "Data" {
		t.Errorf("expected Data, got %s", nr.SheetName())
	}
	if nr.CellRange() != "$A$1:$D$100" {
		t.Errorf("expected $A$1:$D$100, got %s", nr.CellRange())
	}
}

func TestNamedRangeByName(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet1")
	wb.AddNamedRange("Range1", "Sheet1", "$A$1:$A$10")
	wb.AddNamedRange("Range2", "Sheet1", "$B$1:$B$10")

	nr := wb.NamedRange("Range2")
	if nr == nil {
		t.Fatal("should find named range")
	}
	if nr.Name() != "Range2" {
		t.Errorf("expected Range2, got %s", nr.Name())
	}

	if wb.NamedRange("NotExist") != nil {
		t.Error("should return nil for non-existent range")
	}
}

func TestNamedRanges(t *testing.T) {
	wb := New()
	wb.AddSheet("Sheet1")
	wb.AddNamedRange("R1", "Sheet1", "$A$1:$A$5")
	wb.AddNamedRange("R2", "Sheet1", "$B$1:$B$5")

	ranges := wb.NamedRanges()
	if len(ranges) != 2 {
		t.Errorf("expected 2 ranges, got %d", len(ranges))
	}
}

func TestNamedRangeSave(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Sales")
	s.SetValue(1, 1, "Q1")
	s.SetValue(1, 2, 1000.0)
	wb.AddNamedRange("QuarterData", "Sales", "$A$1:$B$4")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "named_range.xlsx")
	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestNamedRangeSaveToBytes(t *testing.T) {
	wb := New()
	wb.AddSheet("Data")
	wb.AddNamedRange("TestRange", "Data", "$A$1:$C$50")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestNamedRangeNoSheet(t *testing.T) {
	wb := New()
	// Named range referencing non-existent sheet
	nr := wb.AddNamedRange("Orphan", "NonExistent", "$A$1:$A$10")
	if nr.sheetID != -1 {
		t.Errorf("expected sheetID -1 for missing sheet, got %d", nr.sheetID)
	}

	// Should still save without error
	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}
