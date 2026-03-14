package style

import (
	"testing"
)

func TestDesignLevelString(t *testing.T) {
	tests := []struct {
		level DesignLevel
		want  string
	}{
		{DesignLevelBasic, "Basic"},
		{DesignLevelProfessional, "Professional"},
		{DesignLevelPremium, "Premium"},
		{DesignLevelLuxury, "Luxury"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("DesignLevel(%d).String() = %q, want %q", tt.level, got, tt.want)
		}
	}
}

func TestAllThemes(t *testing.T) {
	themes := AllThemes()
	if len(themes) != 6 {
		t.Errorf("expected 6 themes, got %d", len(themes))
	}

	// Verify each theme has required fields
	for _, theme := range themes {
		if theme.Name == "" {
			t.Error("theme name should not be empty")
		}
		if theme.Typography.HeadingFont.Family == "" {
			t.Errorf("theme %q heading font should not be empty", theme.Name)
		}
		if theme.Typography.BodyFont.Size == 0 {
			t.Errorf("theme %q body font size should not be zero", theme.Name)
		}
	}
}

func TestThemesByLevel(t *testing.T) {
	basic := ThemesByLevel(DesignLevelBasic)
	if len(basic) != 1 {
		t.Errorf("expected 1 basic theme, got %d", len(basic))
	}

	premium := ThemesByLevel(DesignLevelPremium)
	if len(premium) != 2 {
		t.Errorf("expected 2 premium themes, got %d", len(premium))
	}

	luxury := ThemesByLevel(DesignLevelLuxury)
	if len(luxury) != 2 {
		t.Errorf("expected 2 luxury themes, got %d", len(luxury))
	}
}
