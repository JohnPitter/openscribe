package document

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTableStylePresets(t *testing.T) {
	presets := []struct {
		name   string
		preset TableStylePreset
	}{
		{"Plain", TableStylePlain},
		{"Striped", TableStyleStriped},
		{"Banded", TableStyleBanded},
		{"Grid", TableStyleGrid},
		{"Dark", TableStyleDark},
		{"Colorful", TableStyleColorful},
	}

	for _, tc := range presets {
		t.Run(tc.name, func(t *testing.T) {
			doc := New()
			doc.AddHeading("Table Style: "+tc.name, 1)

			tbl := doc.AddTable(4, 3)
			// Set header row
			tbl.Cell(0, 0).SetText("Name")
			tbl.Cell(0, 1).SetText("Value")
			tbl.Cell(0, 2).SetText("Status")
			// Set data rows
			for i := 1; i < 4; i++ {
				for j := 0; j < 3; j++ {
					tbl.Cell(i, j).SetText("Data")
				}
			}

			// Apply style after setting text so runs exist for bold/color
			tbl.SetStyle(tc.preset)

			if tc.preset.String() != tc.name {
				t.Errorf("expected style name '%s', got '%s'", tc.name, tc.preset.String())
			}

			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "table_style_"+tc.name+".docx")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}

			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("file should exist: %v", err)
			}
			if info.Size() == 0 {
				t.Error("file should not be empty")
			}
		})
	}
}

func TestTableStyleHeaderBold(t *testing.T) {
	doc := New()
	tbl := doc.AddTable(2, 2)
	tbl.Cell(0, 0).SetText("Header")
	tbl.Cell(0, 1).SetText("Header2")
	tbl.Cell(1, 0).SetText("Data")
	tbl.Cell(1, 1).SetText("Data2")

	tbl.SetStyle(TableStyleDark)

	// Verify header cell has bold text (the run should have bold set)
	headerCell := tbl.Cell(0, 0)
	if len(headerCell.paragraphs) == 0 || len(headerCell.paragraphs[0].runs) == 0 {
		t.Fatal("header cell should have a paragraph with a run")
	}
	if !headerCell.paragraphs[0].runs[0].bold {
		t.Error("header row run should be bold with Dark style")
	}
}
