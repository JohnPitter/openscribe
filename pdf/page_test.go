package pdf

import (
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestPageSize(t *testing.T) {
	p := NewPage(common.PageA4, common.NormalMargins())
	if p.Size().Width.Millimeters() < 209 {
		t.Error("should be A4")
	}

	p.SetSize(common.PageLetter)
	if p.Size().Width.Inches() != 8.5 {
		t.Error("should be letter after set")
	}
}

func TestPageMargins(t *testing.T) {
	p := NewPage(common.PageA4, common.NormalMargins())
	if p.Margins().Top.Inches() != 1 {
		t.Error("should be 1 inch margins")
	}

	p.SetMargins(common.NarrowMargins())
	if p.Margins().Top.Inches() != 0.5 {
		t.Error("should be 0.5 inch margins")
	}
}

func TestOutOfRangeCellAccess(t *testing.T) {
	d := New()
	p := d.AddPage()
	tbl := p.AddTable(72, 72, 2, 2)

	// Should not panic
	tbl.SetCell(-1, 0, "bad")
	tbl.SetCell(0, -1, "bad")
	tbl.SetCell(10, 0, "bad")
	val := tbl.Cell(10, 0)
	if val != "" {
		t.Error("out of range should return empty string")
	}
}

func TestTableBorderColor(t *testing.T) {
	d := New()
	p := d.AddPage()
	tbl := p.AddTable(72, 72, 2, 2)
	tbl.SetBorderColor(common.Red)
	if tbl.borderColor != common.Red {
		t.Error("border color should be red")
	}
}
