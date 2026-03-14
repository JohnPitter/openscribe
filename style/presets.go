package style

import "github.com/JohnPitter/openscribe/common"

// Pre-built themes for different design levels

// BasicClean returns a minimal, clean theme suitable for everyday documents
func BasicClean() Theme {
	return Theme{
		Name:       "Basic Clean",
		Level:      DesignLevelBasic,
		Palette:    MinimalPalette(),
		Typography: MinimalTypography(),
		Spacing: Spacing{
			XS: common.Pt(4),
			SM: common.Pt(8),
			MD: common.Pt(16),
			LG: common.Pt(24),
			XL: common.Pt(32),
		},
		Borders: BorderStyle{
			Radius: common.Pt(0),
			Width:  common.Pt(0.5),
			Color:  common.LightGray,
		},
	}
}

// ProfessionalCorporate returns a business-grade theme
func ProfessionalCorporate() Theme {
	return Theme{
		Name:       "Professional Corporate",
		Level:      DesignLevelProfessional,
		Palette:    CorporatePalette(),
		Typography: ClassicTypography(),
		Spacing: Spacing{
			XS: common.Pt(4),
			SM: common.Pt(8),
			MD: common.Pt(16),
			LG: common.Pt(24),
			XL: common.Pt(40),
		},
		Borders: BorderStyle{
			Radius: common.Pt(2),
			Width:  common.Pt(1),
			Color:  common.LightGray,
		},
	}
}

// PremiumModern returns a Behance/Freepik quality theme
func PremiumModern() Theme {
	return Theme{
		Name:       "Premium Modern",
		Level:      DesignLevelPremium,
		Palette:    ModernPalette(),
		Typography: ModernTypography(),
		Spacing: Spacing{
			XS: common.Pt(6),
			SM: common.Pt(12),
			MD: common.Pt(20),
			LG: common.Pt(32),
			XL: common.Pt(48),
		},
		Borders: BorderStyle{
			Radius: common.Pt(8),
			Width:  common.Pt(1),
			Color:  common.LightGray,
		},
	}
}

// PremiumElegant returns an elegant premium theme
func PremiumElegant() Theme {
	return Theme{
		Name:       "Premium Elegant",
		Level:      DesignLevelPremium,
		Palette:    ElegantPalette(),
		Typography: ElegantTypography(),
		Spacing: Spacing{
			XS: common.Pt(6),
			SM: common.Pt(12),
			MD: common.Pt(20),
			LG: common.Pt(32),
			XL: common.Pt(48),
		},
		Borders: BorderStyle{
			Radius: common.Pt(4),
			Width:  common.Pt(0.75),
			Color:  common.LightGray,
		},
	}
}

// LuxuryAgency returns a Slidesgo/Agency quality theme
func LuxuryAgency() Theme {
	return Theme{
		Name:       "Luxury Agency",
		Level:      DesignLevelLuxury,
		Palette:    NeonPalette(),
		Typography: TechTypography(),
		Spacing: Spacing{
			XS: common.Pt(8),
			SM: common.Pt(16),
			MD: common.Pt(24),
			LG: common.Pt(40),
			XL: common.Pt(56),
		},
		Borders: BorderStyle{
			Radius: common.Pt(12),
			Width:  common.Pt(1.5),
			Color:  common.DarkGray,
		},
	}
}

// LuxuryWarm returns a warm luxury theme
func LuxuryWarm() Theme {
	return Theme{
		Name:       "Luxury Warm",
		Level:      DesignLevelLuxury,
		Palette:    WarmPalette(),
		Typography: ElegantTypography(),
		Spacing: Spacing{
			XS: common.Pt(8),
			SM: common.Pt(16),
			MD: common.Pt(24),
			LG: common.Pt(40),
			XL: common.Pt(56),
		},
		Borders: BorderStyle{
			Radius: common.Pt(16),
			Width:  common.Pt(1),
			Color:  common.LightGray,
		},
	}
}

// AllThemes returns all pre-built themes
func AllThemes() []Theme {
	return []Theme{
		BasicClean(),
		ProfessionalCorporate(),
		PremiumModern(),
		PremiumElegant(),
		LuxuryAgency(),
		LuxuryWarm(),
	}
}

// ThemesByLevel returns all themes at a given design level
func ThemesByLevel(level DesignLevel) []Theme {
	var result []Theme
	for _, t := range AllThemes() {
		if t.Level == level {
			result = append(result, t)
		}
	}
	return result
}
