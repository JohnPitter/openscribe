package spreadsheet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFreezePanes(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Freeze")
	s.SetValue(1, 1, "Header 1")
	s.SetValue(1, 2, "Header 2")
	for i := 2; i <= 20; i++ {
		s.SetValue(i, 1, i*10)
		s.SetValue(i, 2, i*20)
	}
	s.FreezePanes(2, 2)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "freeze.xlsx")
	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestFreezeTopRow(t *testing.T) {
	wb := New()
	s := wb.AddSheet("FreezeTop")
	s.SetValue(1, 1, "Header")
	s.FreezeTopRow()

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestFreezeFirstColumn(t *testing.T) {
	wb := New()
	s := wb.AddSheet("FreezeCol")
	s.SetValue(1, 1, "Label")
	s.FreezeFirstColumn()

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestBuildSheetViewsXML(t *testing.T) {
	tests := []struct {
		name     string
		fp       *FreezePane
		contains []string
	}{
		{
			name: "freeze top row",
			fp:   &FreezePane{row: 2, col: 0},
			contains: []string{
				"sheetViews", "sheetView", "pane",
				`ySplit="1"`, `state="frozen"`, `activePane="bottomLeft"`,
			},
		},
		{
			name: "freeze first column",
			fp:   &FreezePane{row: 0, col: 2},
			contains: []string{
				`xSplit="1"`, `activePane="topRight"`,
			},
		},
		{
			name: "freeze both",
			fp:   &FreezePane{row: 3, col: 3},
			contains: []string{
				`xSplit="2"`, `ySplit="2"`, `activePane="bottomRight"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xml := buildSheetViewsXML(tt.fp)
			for _, s := range tt.contains {
				if !strings.Contains(xml, s) {
					t.Errorf("expected XML to contain %q, got: %s", s, xml)
				}
			}
		})
	}
}

func TestBuildSheetViewsXMLNil(t *testing.T) {
	xml := buildSheetViewsXML(nil)
	if xml != "" {
		t.Error("nil freeze pane should produce empty string")
	}
}
