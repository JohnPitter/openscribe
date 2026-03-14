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
	if len(themes) != 21 {
		t.Errorf("expected 21 themes, got %d", len(themes))
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
	// BasicClean + HighContrastLight + HighContrastDark = 3
	if len(basic) != 3 {
		t.Errorf("expected 3 basic themes, got %d", len(basic))
	}

	professional := ThemesByLevel(DesignLevelProfessional)
	// ProfessionalCorporate + Healthcare + Finance + Education + Legal + NonProfit + Government = 7
	if len(professional) != 7 {
		t.Errorf("expected 7 professional themes, got %d", len(professional))
	}

	premium := ThemesByLevel(DesignLevelPremium)
	// PremiumModern + PremiumElegant + TechStartup + RealEstate + Retail + DarkModern + DarkMinimal = 7
	if len(premium) != 7 {
		t.Errorf("expected 7 premium themes, got %d", len(premium))
	}

	luxury := ThemesByLevel(DesignLevelLuxury)
	// LuxuryAgency + LuxuryWarm + CreativeAgency + DarkElegant = 4
	if len(luxury) != 4 {
		t.Errorf("expected 4 luxury themes, got %d", len(luxury))
	}
}
