package spreadsheet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestConditionalFormatGreaterThan(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")
	for i := 1; i <= 10; i++ {
		s.SetValue(i, 1, float64(i*10))
	}

	cf := s.AddConditionalFormat("A1:A10", ConditionGreaterThan)
	cf.SetValue("50").SetBackgroundColor(common.Green).SetBold(true)

	if cf.CellRange() != "A1:A10" {
		t.Error("range mismatch")
	}
	if cf.Type() != ConditionGreaterThan {
		t.Error("type mismatch")
	}
	if cf.Value() != "50" {
		t.Error("value mismatch")
	}

	path := filepath.Join(t.TempDir(), "conditional.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestConditionalFormatBetween(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")
	for i := 1; i <= 10; i++ {
		s.SetValue(i, 1, float64(i*10))
	}

	cf := s.AddConditionalFormat("A1:A10", ConditionBetween)
	cf.SetValue("30").SetValue2("70").SetBackgroundColor(common.Yellow)

	path := filepath.Join(t.TempDir(), "between.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestConditionalFormatColorScale(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")
	for i := 1; i <= 10; i++ {
		s.SetValue(i, 1, float64(i*10))
	}

	cf := s.AddConditionalFormat("A1:A10", ConditionColorScale)
	cf.SetColorScale(common.Red, common.Green)

	path := filepath.Join(t.TempDir(), "colorscale.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestConditionalFormatDataBar(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")
	for i := 1; i <= 10; i++ {
		s.SetValue(i, 1, float64(i*10))
	}

	cf := s.AddConditionalFormat("A1:A10", ConditionDataBar)
	cf.SetBarColor(common.Blue)

	path := filepath.Join(t.TempDir(), "databar.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestConditionalFormatChaining(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")

	cf := s.AddConditionalFormat("B1:B5", ConditionLessThan)
	result := cf.SetValue("0").SetBackgroundColor(common.Red).SetFontColor(common.White).SetBold(true).SetItalic(true)

	if result != cf {
		t.Error("chaining should return same pointer")
	}
}

func TestMultipleConditionalFormats(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")
	for i := 1; i <= 10; i++ {
		s.SetValue(i, 1, float64(i*10))
	}

	s.AddConditionalFormat("A1:A10", ConditionGreaterThan).SetValue("80").SetBackgroundColor(common.Green)
	s.AddConditionalFormat("A1:A10", ConditionLessThan).SetValue("30").SetBackgroundColor(common.Red)

	path := filepath.Join(t.TempDir(), "multi_cond.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}
