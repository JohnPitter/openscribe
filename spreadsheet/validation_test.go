package spreadsheet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidationList(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Validation")
	s.SetValue(1, 1, "Status")

	v := s.AddValidation("B1:B100", ValidationList)
	v.SetList([]string{"Open", "In Progress", "Closed"})
	v.SetErrorMessage("Invalid", "Please select from the list")
	v.SetPromptMessage("Select Status", "Choose a status value")

	if v.CellRange() != "B1:B100" {
		t.Errorf("expected B1:B100, got %s", v.CellRange())
	}
	if v.Type() != ValidationList {
		t.Errorf("expected ValidationList, got %d", v.Type())
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "validation_list.xlsx")
	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestValidationWholeNumber(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Numbers")
	s.SetValue(1, 1, "Age")

	v := s.AddValidation("B1:B50", ValidationWholeNumber)
	v.SetRange("1", "120")
	v.SetErrorMessage("Invalid Age", "Age must be between 1 and 120")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestValidationDecimal(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Decimals")

	v := s.AddValidation("A1:A10", ValidationDecimal)
	v.SetRange("0.0", "100.0")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestValidationCustom(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Custom")

	v := s.AddValidation("A1:A10", ValidationCustom)
	v.SetCustomFormula("LEN(A1)>0")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestValidationMultiple(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Multi")

	s.AddValidation("A1:A10", ValidationList).SetList([]string{"Yes", "No"})
	s.AddValidation("B1:B10", ValidationWholeNumber).SetRange("0", "999")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestValidationXMLOutput(t *testing.T) {
	validations := []*Validation{
		{
			cellRange:      "A1:A10",
			validationType: ValidationList,
			listItems:      []string{"Red", "Green", "Blue"},
			errorTitle:     "Error",
			errorMessage:   "Pick a color",
		},
	}
	xml := buildValidationsXML(validations)
	if !strings.Contains(xml, "dataValidations") {
		t.Error("should contain dataValidations element")
	}
	if !strings.Contains(xml, `type="list"`) {
		t.Error("should contain type=list")
	}
	if !strings.Contains(xml, "Red,Green,Blue") {
		t.Error("should contain list items")
	}
}

func TestValidationEmptyList(t *testing.T) {
	xml := buildValidationsXML(nil)
	if xml != "" {
		t.Error("empty validations should produce empty string")
	}
}
