package style

import (
	"testing"
)

func TestAllPresetsValid(t *testing.T) {
	presets := []struct {
		name string
		fn   func() Theme
	}{
		{"BasicClean", BasicClean},
		{"ProfessionalCorporate", ProfessionalCorporate},
		{"PremiumModern", PremiumModern},
		{"PremiumElegant", PremiumElegant},
		{"LuxuryAgency", LuxuryAgency},
		{"LuxuryWarm", LuxuryWarm},
	}

	for _, p := range presets {
		t.Run(p.name, func(t *testing.T) {
			theme := p.fn()
			if theme.Name == "" {
				t.Error("name empty")
			}
			if theme.Palette.Primary == (theme.Palette.Background) {
				t.Error("primary should differ from background")
			}
			if theme.Typography.HeadingFont.Size <= theme.Typography.BodyFont.Size {
				t.Error("heading font should be larger than body font")
			}
			if theme.Spacing.XS.Points() >= theme.Spacing.XL.Points() {
				t.Error("XS spacing should be less than XL")
			}
		})
	}
}

func TestPalettes(t *testing.T) {
	palettes := []struct {
		name string
		fn   func() Palette
	}{
		{"Corporate", CorporatePalette},
		{"Modern", ModernPalette},
		{"Elegant", ElegantPalette},
		{"Minimal", MinimalPalette},
		{"Warm", WarmPalette},
		{"Neon", NeonPalette},
	}
	for _, p := range palettes {
		t.Run(p.name, func(t *testing.T) {
			pal := p.fn()
			if pal.Text == pal.Background {
				t.Error("text and background should be different")
			}
		})
	}
}
