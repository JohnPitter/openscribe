package style

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExportImportRoundTrip(t *testing.T) {
	original := ProfessionalCorporate()

	data, err := ExportTheme(original)
	if err != nil {
		t.Fatalf("ExportTheme error: %v", err)
	}

	imported, err := ImportTheme(data)
	if err != nil {
		t.Fatalf("ImportTheme error: %v", err)
	}

	// Verify key fields match
	if imported.Name != original.Name {
		t.Errorf("name: expected %q, got %q", original.Name, imported.Name)
	}
	if imported.Level != original.Level {
		t.Errorf("level: expected %s, got %s", original.Level, imported.Level)
	}
	if imported.Palette.Primary.Hex() != original.Palette.Primary.Hex() {
		t.Errorf("primary: expected %s, got %s", original.Palette.Primary.Hex(), imported.Palette.Primary.Hex())
	}
	if imported.Palette.Secondary.Hex() != original.Palette.Secondary.Hex() {
		t.Errorf("secondary: expected %s, got %s", original.Palette.Secondary.Hex(), imported.Palette.Secondary.Hex())
	}
	if imported.Typography.HeadingFont.Family != original.Typography.HeadingFont.Family {
		t.Errorf("heading font: expected %q, got %q", original.Typography.HeadingFont.Family, imported.Typography.HeadingFont.Family)
	}
	if imported.Typography.BodyFont.Size != original.Typography.BodyFont.Size {
		t.Errorf("body size: expected %.1f, got %.1f", original.Typography.BodyFont.Size, imported.Typography.BodyFont.Size)
	}
	if imported.Spacing.MD.Points() != original.Spacing.MD.Points() {
		t.Errorf("spacing MD: expected %.1f, got %.1f", original.Spacing.MD.Points(), imported.Spacing.MD.Points())
	}
	if imported.Borders.Width.Points() != original.Borders.Width.Points() {
		t.Errorf("border width: expected %.1f, got %.1f", original.Borders.Width.Points(), imported.Borders.Width.Points())
	}
}

func TestExportImportAllThemes(t *testing.T) {
	for _, theme := range AllThemes() {
		data, err := ExportTheme(theme)
		if err != nil {
			t.Fatalf("ExportTheme(%q) error: %v", theme.Name, err)
		}
		imported, err := ImportTheme(data)
		if err != nil {
			t.Fatalf("ImportTheme(%q) error: %v", theme.Name, err)
		}
		if imported.Name != theme.Name {
			t.Errorf("round-trip name mismatch: expected %q, got %q", theme.Name, imported.Name)
		}
	}
}

func TestExportImportFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "theme.json")

	original := HealthcareTheme()

	err := ExportThemeToFile(original, path)
	if err != nil {
		t.Fatalf("ExportThemeToFile error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("exported file does not exist")
	}

	imported, err := ImportThemeFromFile(path)
	if err != nil {
		t.Fatalf("ImportThemeFromFile error: %v", err)
	}

	if imported.Name != original.Name {
		t.Errorf("expected name %q, got %q", original.Name, imported.Name)
	}
	if imported.Palette.Primary.Hex() != original.Palette.Primary.Hex() {
		t.Errorf("primary: expected %s, got %s", original.Palette.Primary.Hex(), imported.Palette.Primary.Hex())
	}
}

func TestImportInvalidJSON(t *testing.T) {
	_, err := ImportTheme([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestImportFileNotFound(t *testing.T) {
	_, err := ImportThemeFromFile("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestExportProducesValidJSON(t *testing.T) {
	theme := DarkModern()
	data, err := ExportTheme(theme)
	if err != nil {
		t.Fatalf("ExportTheme error: %v", err)
	}

	// Should start with { and end with }
	if len(data) < 2 || data[0] != '{' || data[len(data)-1] != '}' {
		t.Error("exported data is not valid JSON structure")
	}
}
