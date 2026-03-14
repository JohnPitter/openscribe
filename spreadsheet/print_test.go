package spreadsheet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestSetPrintArea(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Print")
	s.SetValue(1, 1, "Data")
	s.SetPrintArea(1, 1, 10, 5)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "print_area.xlsx")
	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestSetPrintTitles(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Titles")
	s.SetValue(1, 1, "Header")
	s.SetPrintTitles("1:2", "A:B")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestSetPageOrientation(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Landscape")
	s.SetValue(1, 1, "Wide data")
	s.SetPageOrientation(common.OrientationLandscape)

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestSetPaperSize(t *testing.T) {
	wb := New()
	s := wb.AddSheet("A4")
	s.SetValue(1, 1, "Data")
	s.SetPaperSize(9) // A4

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestSetFitToPage(t *testing.T) {
	wb := New()
	s := wb.AddSheet("FitToPage")
	s.SetValue(1, 1, "Data")
	s.SetFitToPage(1, 2)

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestBuildPageSetupXML(t *testing.T) {
	ps := &PrintSettings{
		orientation: common.OrientationLandscape,
		paperSize:   9,
		fitToWidth:  1,
		fitToHeight: 1,
	}
	xml := buildPageSetupXML(ps)
	if !strings.Contains(xml, "pageSetup") {
		t.Error("should contain pageSetup element")
	}
	if !strings.Contains(xml, "pageMargins") {
		t.Error("should contain pageMargins element")
	}
	if !strings.Contains(xml, `orientation="landscape"`) {
		t.Error("should have landscape orientation")
	}
	if !strings.Contains(xml, `paperSize="9"`) {
		t.Error("should have paper size 9")
	}
	if !strings.Contains(xml, `fitToWidth="1"`) {
		t.Error("should have fitToWidth")
	}
}

func TestBuildPageSetupXMLNil(t *testing.T) {
	xml := buildPageSetupXML(nil)
	if xml != "" {
		t.Error("nil print settings should produce empty string")
	}
}

func TestBuildPageSetupPortrait(t *testing.T) {
	ps := &PrintSettings{
		orientation: common.OrientationPortrait,
	}
	xml := buildPageSetupXML(ps)
	if !strings.Contains(xml, `orientation="portrait"`) {
		t.Error("should have portrait orientation")
	}
}

func TestPrintAreaDefinedName(t *testing.T) {
	pa := &PrintArea{startRow: 1, startCol: 1, endRow: 10, endCol: 5}
	xml := buildPrintAreaDefinedName("Sheet1", 0, pa)
	if !strings.Contains(xml, "_xlnm.Print_Area") {
		t.Error("should contain _xlnm.Print_Area")
	}
	if !strings.Contains(xml, "Sheet1") {
		t.Error("should contain sheet name")
	}
}

func TestPrintAreaDefinedNameNil(t *testing.T) {
	xml := buildPrintAreaDefinedName("Sheet1", 0, nil)
	if xml != "" {
		t.Error("nil print area should produce empty string")
	}
}

func TestPrintTitlesDefinedName(t *testing.T) {
	xml := buildPrintTitlesDefinedName("Sheet1", 0, "1:2", "A:B")
	if !strings.Contains(xml, "_xlnm.Print_Titles") {
		t.Error("should contain _xlnm.Print_Titles")
	}
}

func TestPrintTitlesDefinedNameEmpty(t *testing.T) {
	xml := buildPrintTitlesDefinedName("Sheet1", 0, "", "")
	if xml != "" {
		t.Error("empty print titles should produce empty string")
	}
}

func TestCombinedPrintSettings(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Combined")
	s.SetValue(1, 1, "Data")
	s.SetPrintArea(1, 1, 50, 10)
	s.SetPrintTitles("1:1", "")
	s.SetPageOrientation(common.OrientationLandscape)
	s.SetPaperSize(1) // Letter
	s.SetFitToPage(1, 0)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "combined_print.xlsx")
	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}
