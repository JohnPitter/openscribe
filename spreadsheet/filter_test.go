package spreadsheet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetAutoFilter(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Filter")
	s.SetValue(1, 1, "Name")
	s.SetValue(1, 2, "Age")
	s.SetValue(1, 3, "City")
	s.SetValue(2, 1, "Alice")
	s.SetValue(2, 2, 30.0)
	s.SetValue(2, 3, "NYC")

	s.SetAutoFilter(1, 1, 10, 3)

	af := s.AutoFilter()
	if af == nil {
		t.Fatal("auto filter should not be nil")
	}
	if af.Ref() != "A1:C10" {
		t.Errorf("expected A1:C10, got %s", af.Ref())
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "filter.xlsx")
	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestAutoFilterNil(t *testing.T) {
	wb := New()
	s := wb.AddSheet("NoFilter")
	if s.AutoFilter() != nil {
		t.Error("should be nil when no auto filter set")
	}
}

func TestBuildAutoFilterXML(t *testing.T) {
	af := &AutoFilter{
		startRow: 1, startCol: 1,
		endRow: 10, endCol: 4,
	}
	xml := buildAutoFilterXML(af)
	if !strings.Contains(xml, "autoFilter") {
		t.Error("should contain autoFilter element")
	}
	if !strings.Contains(xml, `ref="A1:D10"`) {
		t.Errorf("expected ref A1:D10, got %s", xml)
	}
}

func TestBuildAutoFilterXMLNil(t *testing.T) {
	xml := buildAutoFilterXML(nil)
	if xml != "" {
		t.Error("nil auto filter should produce empty string")
	}
}

func TestAutoFilterSaveToBytes(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")
	s.SetValue(1, 1, "Header")
	s.SetAutoFilter(1, 1, 5, 3)

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}
