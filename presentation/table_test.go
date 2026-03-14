package presentation

import (
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestAddTable(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	tbl := s.AddTable(3, 4, common.In(1), common.In(1), common.In(8), common.In(4))
	if tbl == nil {
		t.Fatal("table should not be nil")
	}
	if tbl.Rows() != 3 {
		t.Errorf("expected 3 rows, got %d", tbl.Rows())
	}
	if tbl.Cols() != 4 {
		t.Errorf("expected 4 cols, got %d", tbl.Cols())
	}
	if s.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", s.ElementCount())
	}
}

func TestTableCellText(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	tbl := s.AddTable(2, 2, common.In(1), common.In(1), common.In(6), common.In(3))
	tbl.Cell(0, 0).SetText("Header 1")
	tbl.Cell(0, 1).SetText("Header 2")
	tbl.Cell(1, 0).SetText("Value 1")
	tbl.Cell(1, 1).SetText("Value 2")

	if tbl.Cell(0, 0).Text() != "Header 1" {
		t.Errorf("expected 'Header 1', got '%s'", tbl.Cell(0, 0).Text())
	}
	if tbl.Cell(1, 1).Text() != "Value 2" {
		t.Errorf("expected 'Value 2', got '%s'", tbl.Cell(1, 1).Text())
	}
}

func TestTableCellOutOfBounds(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	tbl := s.AddTable(2, 2, common.In(1), common.In(1), common.In(4), common.In(2))
	if tbl.Cell(-1, 0) != nil {
		t.Error("negative row should return nil")
	}
	if tbl.Cell(0, 5) != nil {
		t.Error("out of bounds col should return nil")
	}
	if tbl.Cell(3, 0) != nil {
		t.Error("out of bounds row should return nil")
	}
}

func TestTableCellStyling(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	tbl := s.AddTable(2, 2, common.In(1), common.In(1), common.In(6), common.In(3))
	cell := tbl.Cell(0, 0)
	cell.SetText("Styled")
	cell.SetFont(common.NewFont("Arial", 14).Bold())
	cell.SetBackground(common.Blue)

	if cell.text != "Styled" {
		t.Error("text should be set")
	}
	if cell.font == nil {
		t.Error("font should be set")
	}
	if cell.background == nil || *cell.background != common.Blue {
		t.Error("background should be blue")
	}
}

func TestTableHeaderBackground(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	tbl := s.AddTable(3, 3, common.In(1), common.In(1), common.In(6), common.In(4))
	tbl.SetHeaderBackground(common.DarkGray)
	tbl.SetBorderColor(common.Black)

	if tbl.headerBg == nil || *tbl.headerBg != common.DarkGray {
		t.Error("header background should be dark gray")
	}
	if tbl.borderColor == nil || *tbl.borderColor != common.Black {
		t.Error("border color should be black")
	}
}

func TestTableSerialization(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	tbl := s.AddTable(2, 3, common.In(1), common.In(1), common.In(8), common.In(3))
	tbl.Cell(0, 0).SetText("Name")
	tbl.Cell(0, 1).SetText("Age")
	tbl.Cell(0, 2).SetText("City")
	tbl.Cell(1, 0).SetText("Alice")
	tbl.Cell(1, 1).SetText("30")
	tbl.Cell(1, 2).SetText("NYC")
	tbl.SetHeaderBackground(common.LightGray)

	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, err := packaging.OpenPackageFromBytes(data)
	if err != nil {
		t.Fatalf("open package error: %v", err)
	}

	slideXML, ok := pkg.GetFile("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide1.xml should exist")
	}

	xmlStr := string(slideXML)
	if !strings.Contains(xmlStr, "graphicFrame") {
		t.Error("slide XML should contain graphicFrame element")
	}
	if !strings.Contains(xmlStr, "a:tbl") {
		t.Error("slide XML should contain a:tbl element")
	}
	if !strings.Contains(xmlStr, "Alice") {
		t.Error("slide XML should contain cell text 'Alice'")
	}
	if !strings.Contains(xmlStr, "a:tr") {
		t.Error("slide XML should contain table rows")
	}
	if !strings.Contains(xmlStr, "a:tc") {
		t.Error("slide XML should contain table cells")
	}
}

func TestTableElementType(t *testing.T) {
	tbl := &SlideTable{}
	if tbl.elementType() != "table" {
		t.Errorf("expected 'table', got '%s'", tbl.elementType())
	}
}
