package spreadsheet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSheetProtect(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Protected")
	s.SetValue(1, 1, "Locked Data")
	s.Protect("password123")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "protected.xlsx")
	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestSheetProtectWithOptions(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Options")
	s.SetValue(1, 1, "Data")
	s.Protect("secret")
	s.SetProtectionOptions(ProtectionOptions{
		AllowSort:        true,
		AllowFilter:      true,
		AllowInsertRows:  false,
		AllowDeleteRows:  false,
		AllowFormatCells: true,
	})

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestSetCellLocked(t *testing.T) {
	wb := New()
	s := wb.AddSheet("CellLock")
	s.SetValue(1, 1, "Unlocked")
	s.SetCellLocked(1, 1, false)
	s.Protect("pass")

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestHashPassword(t *testing.T) {
	hash := hashPassword("password")
	if hash == "" {
		t.Error("hash should not be empty")
	}
	if len(hash) != 4 {
		t.Errorf("expected 4-char hex hash, got %s", hash)
	}

	// Empty password should return empty hash
	if hashPassword("") != "" {
		t.Error("empty password should return empty hash")
	}
}

func TestBuildProtectionXML(t *testing.T) {
	prot := &SheetProtection{
		password: "test",
		options: ProtectionOptions{
			AllowSort:   true,
			AllowFilter: true,
		},
	}
	xml := buildProtectionXML(prot)
	if !strings.Contains(xml, "sheetProtection") {
		t.Error("should contain sheetProtection element")
	}
	if !strings.Contains(xml, `sheet="1"`) {
		t.Error("should have sheet=1 attribute")
	}
	if !strings.Contains(xml, `password=`) {
		t.Error("should contain password attribute")
	}
	if !strings.Contains(xml, `sort="0"`) {
		t.Error("should allow sort")
	}
	if !strings.Contains(xml, `autoFilter="0"`) {
		t.Error("should allow autoFilter")
	}
}

func TestBuildProtectionXMLNil(t *testing.T) {
	xml := buildProtectionXML(nil)
	if xml != "" {
		t.Error("nil protection should produce empty string")
	}
}

func TestProtectWithoutPassword(t *testing.T) {
	wb := New()
	s := wb.AddSheet("NoPass")
	s.SetProtectionOptions(ProtectionOptions{AllowSort: true})

	data, err := wb.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}
