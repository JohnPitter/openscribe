package style

import (
	"testing"
)

func TestAllExtendedThemesHaveRequiredFields(t *testing.T) {
	themes := AllThemes()

	for _, theme := range themes {
		if theme.Name == "" {
			t.Error("theme has empty name")
		}

		// Palette
		if theme.Palette.Primary == theme.Palette.Background && theme.Name != "High Contrast Dark" {
			// Only High Contrast Dark might have white-on-black
		}

		// Typography
		if theme.Typography.HeadingFont.Family == "" {
			t.Errorf("theme %q has empty heading font family", theme.Name)
		}
		if theme.Typography.BodyFont.Family == "" {
			t.Errorf("theme %q has empty body font family", theme.Name)
		}
		if theme.Typography.CaptionFont.Family == "" {
			t.Errorf("theme %q has empty caption font family", theme.Name)
		}
		if theme.Typography.CodeFont.Family == "" {
			t.Errorf("theme %q has empty code font family", theme.Name)
		}

		if theme.Typography.HeadingFont.Size == 0 {
			t.Errorf("theme %q has zero heading font size", theme.Name)
		}
		if theme.Typography.BodyFont.Size == 0 {
			t.Errorf("theme %q has zero body font size", theme.Name)
		}
		if theme.Typography.CaptionFont.Size == 0 {
			t.Errorf("theme %q has zero caption font size", theme.Name)
		}
		if theme.Typography.CodeFont.Size == 0 {
			t.Errorf("theme %q has zero code font size", theme.Name)
		}

		// Spacing
		if theme.Spacing.XS.Points() == 0 {
			t.Errorf("theme %q has zero XS spacing", theme.Name)
		}
		if theme.Spacing.SM.Points() == 0 {
			t.Errorf("theme %q has zero SM spacing", theme.Name)
		}
		if theme.Spacing.MD.Points() == 0 {
			t.Errorf("theme %q has zero MD spacing", theme.Name)
		}
		if theme.Spacing.LG.Points() == 0 {
			t.Errorf("theme %q has zero LG spacing", theme.Name)
		}
		if theme.Spacing.XL.Points() == 0 {
			t.Errorf("theme %q has zero XL spacing", theme.Name)
		}

		// Borders
		if theme.Borders.Width.Points() == 0 {
			t.Errorf("theme %q has zero border width", theme.Name)
		}
	}
}

func TestIndustryThemes(t *testing.T) {
	themes := []struct {
		name  string
		theme Theme
	}{
		{"Healthcare", HealthcareTheme()},
		{"Finance", FinanceTheme()},
		{"Education", EducationTheme()},
		{"Tech Startup", TechStartupTheme()},
		{"Legal", LegalTheme()},
		{"Creative Agency", CreativeAgencyTheme()},
		{"Real Estate", RealEstateTheme()},
		{"Retail", RetailTheme()},
		{"Non-Profit", NonProfitTheme()},
		{"Government", GovernmentTheme()},
	}

	for _, tt := range themes {
		if tt.theme.Name != tt.name {
			t.Errorf("expected name %q, got %q", tt.name, tt.theme.Name)
		}
	}
}

func TestDarkThemes(t *testing.T) {
	themes := []struct {
		name  string
		theme Theme
	}{
		{"Dark Modern", DarkModern()},
		{"Dark Elegant", DarkElegant()},
		{"Dark Minimal", DarkMinimal()},
	}

	for _, tt := range themes {
		if tt.theme.Name != tt.name {
			t.Errorf("expected name %q, got %q", tt.name, tt.theme.Name)
		}
		// Dark themes should have dark backgrounds
		bg := tt.theme.Palette.Background
		brightness := int(bg.R) + int(bg.G) + int(bg.B)
		if brightness > 300 {
			t.Errorf("theme %q background seems too bright for a dark theme (brightness=%d)", tt.name, brightness)
		}
	}
}

func TestHighContrastThemes(t *testing.T) {
	light := HighContrastLight()
	if light.Name != "High Contrast Light" {
		t.Errorf("expected 'High Contrast Light', got %q", light.Name)
	}
	if light.Level != DesignLevelBasic {
		t.Errorf("expected Basic level, got %s", light.Level)
	}

	dark := HighContrastDark()
	if dark.Name != "High Contrast Dark" {
		t.Errorf("expected 'High Contrast Dark', got %q", dark.Name)
	}
	if dark.Level != DesignLevelBasic {
		t.Errorf("expected Basic level, got %s", dark.Level)
	}
}

func TestThemeCount(t *testing.T) {
	all := AllThemes()
	// 6 original + 10 industry + 3 dark + 2 high-contrast = 21
	if len(all) != 21 {
		t.Errorf("expected 21 themes, got %d", len(all))
	}
}
