package document

import (
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestNewTable(t *testing.T) {
	tbl := NewTable(3, 4)
	if tbl.RowCount() != 3 {
		t.Errorf("expected 3 rows, got %d", tbl.RowCount())
	}
	if tbl.ColCount() != 4 {
		t.Errorf("expected 4 cols, got %d", tbl.ColCount())
	}
}

func TestTableCell(t *testing.T) {
	tbl := NewTable(2, 2)

	cell := tbl.Cell(0, 0)
	if cell == nil {
		t.Fatal("cell should not be nil")
	}

	cell.SetText("Hello")
	if cell.Text() != "Hello" {
		t.Errorf("expected Hello, got %s", cell.Text())
	}

	// Out of range
	if tbl.Cell(5, 5) != nil {
		t.Error("out of range cell should be nil")
	}
}

func TestTableAddRemoveRow(t *testing.T) {
	tbl := NewTable(2, 3)
	tbl.AddRow()
	if tbl.RowCount() != 3 {
		t.Errorf("expected 3 rows, got %d", tbl.RowCount())
	}

	err := tbl.RemoveRow(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tbl.RowCount() != 2 {
		t.Errorf("expected 2 rows, got %d", tbl.RowCount())
	}

	err = tbl.RemoveRow(10)
	if err == nil {
		t.Error("should error on out of range")
	}
}

func TestCellShading(t *testing.T) {
	tbl := NewTable(1, 1)
	cell := tbl.Cell(0, 0)
	cell.SetShading(common.Blue)
	if cell.shading == nil || *cell.shading != common.Blue {
		t.Error("shading should be blue")
	}
}
