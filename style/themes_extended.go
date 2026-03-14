package style

import "github.com/JohnPitter/openscribe/common"

// --- Industry-specific themes ---

// HealthcareTheme returns a blue/teal professional theme for healthcare.
func HealthcareTheme() Theme {
	return Theme{
		Name:  "Healthcare",
		Level: DesignLevelProfessional,
		Palette: Palette{
			Primary:    mustHex("#0077B6"),
			Secondary:  mustHex("#00B4D8"),
			Accent:     mustHex("#48CAE4"),
			Background: common.White,
			Surface:    mustHex("#F0F9FF"),
			Text:       mustHex("#023E8A"),
			TextLight:  mustHex("#6B7280"),
			Success:    mustHex("#10B981"),
			Warning:    mustHex("#F59E0B"),
			Error:      mustHex("#EF4444"),
			Info:       mustHex("#0077B6"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Segoe UI", 24).Bold(),
			BodyFont:    common.NewFont("Segoe UI", 11),
			CaptionFont: common.NewFont("Segoe UI", 9).WithColor(common.Gray),
			CodeFont:    common.NewFont("Consolas", 10),
		},
		Spacing: Spacing{
			XS: common.Pt(4),
			SM: common.Pt(8),
			MD: common.Pt(16),
			LG: common.Pt(24),
			XL: common.Pt(40),
		},
		Borders: BorderStyle{
			Radius: common.Pt(4),
			Width:  common.Pt(1),
			Color:  mustHex("#B0D4E8"),
		},
	}
}

// FinanceTheme returns a navy/gold corporate theme for finance.
func FinanceTheme() Theme {
	return Theme{
		Name:  "Finance",
		Level: DesignLevelProfessional,
		Palette: Palette{
			Primary:    mustHex("#1B2A4A"),
			Secondary:  mustHex("#2C3E6B"),
			Accent:     mustHex("#C9A227"),
			Background: common.White,
			Surface:    mustHex("#F8F6F0"),
			Text:       mustHex("#1B2A4A"),
			TextLight:  mustHex("#6B7280"),
			Success:    mustHex("#059669"),
			Warning:    mustHex("#D97706"),
			Error:      mustHex("#DC2626"),
			Info:       mustHex("#2563EB"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Georgia", 24).Bold(),
			BodyFont:    common.NewFont("Georgia", 11),
			CaptionFont: common.NewFont("Georgia", 9).WithColor(common.Gray),
			CodeFont:    common.NewFont("Courier New", 10),
		},
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
			Color:  mustHex("#D1C7A3"),
		},
	}
}

// EducationTheme returns a green/warm friendly theme for education.
func EducationTheme() Theme {
	return Theme{
		Name:  "Education",
		Level: DesignLevelProfessional,
		Palette: Palette{
			Primary:    mustHex("#16A34A"),
			Secondary:  mustHex("#166534"),
			Accent:     mustHex("#F59E0B"),
			Background: mustHex("#FFFDF7"),
			Surface:    mustHex("#F0FDF4"),
			Text:       mustHex("#1C1917"),
			TextLight:  mustHex("#78716C"),
			Success:    mustHex("#22C55E"),
			Warning:    mustHex("#FBBF24"),
			Error:      mustHex("#EF4444"),
			Info:       mustHex("#3B82F6"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Trebuchet MS", 24).Bold(),
			BodyFont:    common.NewFont("Trebuchet MS", 11),
			CaptionFont: common.NewFont("Trebuchet MS", 9).WithColor(common.Gray),
			CodeFont:    common.NewFont("Consolas", 10),
		},
		Spacing: Spacing{
			XS: common.Pt(4),
			SM: common.Pt(10),
			MD: common.Pt(18),
			LG: common.Pt(28),
			XL: common.Pt(44),
		},
		Borders: BorderStyle{
			Radius: common.Pt(6),
			Width:  common.Pt(1),
			Color:  mustHex("#BBF7D0"),
		},
	}
}

// TechStartupTheme returns a purple/electric modern theme for tech startups.
func TechStartupTheme() Theme {
	return Theme{
		Name:  "Tech Startup",
		Level: DesignLevelPremium,
		Palette: Palette{
			Primary:    mustHex("#7C3AED"),
			Secondary:  mustHex("#4C1D95"),
			Accent:     mustHex("#06B6D4"),
			Background: common.White,
			Surface:    mustHex("#F5F3FF"),
			Text:       mustHex("#1E1B4B"),
			TextLight:  mustHex("#6B7280"),
			Success:    mustHex("#10B981"),
			Warning:    mustHex("#F59E0B"),
			Error:      mustHex("#EF4444"),
			Info:       mustHex("#6366F1"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Segoe UI", 28).Bold(),
			BodyFont:    common.NewFont("Segoe UI", 11),
			CaptionFont: common.NewFont("Segoe UI", 9).WithColor(common.Gray),
			CodeFont:    common.NewFont("Cascadia Code", 10),
		},
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
			Color:  mustHex("#C4B5FD"),
		},
	}
}

// LegalTheme returns a dark gray/maroon formal theme for legal documents.
func LegalTheme() Theme {
	return Theme{
		Name:  "Legal",
		Level: DesignLevelProfessional,
		Palette: Palette{
			Primary:    mustHex("#4A4A4A"),
			Secondary:  mustHex("#6B2737"),
			Accent:     mustHex("#8B3A4A"),
			Background: common.White,
			Surface:    mustHex("#FAFAFA"),
			Text:       mustHex("#1A1A1A"),
			TextLight:  mustHex("#737373"),
			Success:    mustHex("#16A34A"),
			Warning:    mustHex("#CA8A04"),
			Error:      mustHex("#DC2626"),
			Info:       mustHex("#2563EB"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Times New Roman", 24).Bold(),
			BodyFont:    common.NewFont("Times New Roman", 12),
			CaptionFont: common.NewFont("Times New Roman", 10).WithColor(common.Gray),
			CodeFont:    common.NewFont("Courier New", 10),
		},
		Spacing: Spacing{
			XS: common.Pt(4),
			SM: common.Pt(8),
			MD: common.Pt(16),
			LG: common.Pt(24),
			XL: common.Pt(36),
		},
		Borders: BorderStyle{
			Radius: common.Pt(0),
			Width:  common.Pt(1),
			Color:  mustHex("#D4D4D4"),
		},
	}
}

// CreativeAgencyTheme returns a vibrant multi-color theme for creative agencies.
func CreativeAgencyTheme() Theme {
	return Theme{
		Name:  "Creative Agency",
		Level: DesignLevelLuxury,
		Palette: Palette{
			Primary:    mustHex("#FF6B6B"),
			Secondary:  mustHex("#4ECDC4"),
			Accent:     mustHex("#FFE66D"),
			Background: common.White,
			Surface:    mustHex("#FFF8F0"),
			Text:       mustHex("#2D3436"),
			TextLight:  mustHex("#636E72"),
			Success:    mustHex("#00B894"),
			Warning:    mustHex("#FDCB6E"),
			Error:      mustHex("#E17055"),
			Info:       mustHex("#74B9FF"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Helvetica", 28).Bold(),
			BodyFont:    common.NewFont("Helvetica", 11),
			CaptionFont: common.NewFont("Helvetica", 9).WithColor(common.Gray),
			CodeFont:    common.NewFont("Consolas", 10),
		},
		Spacing: Spacing{
			XS: common.Pt(8),
			SM: common.Pt(16),
			MD: common.Pt(24),
			LG: common.Pt(40),
			XL: common.Pt(56),
		},
		Borders: BorderStyle{
			Radius: common.Pt(12),
			Width:  common.Pt(2),
			Color:  mustHex("#DFE6E9"),
		},
	}
}

// RealEstateTheme returns an earth-tone luxury theme for real estate.
func RealEstateTheme() Theme {
	return Theme{
		Name:  "Real Estate",
		Level: DesignLevelPremium,
		Palette: Palette{
			Primary:    mustHex("#5D4037"),
			Secondary:  mustHex("#795548"),
			Accent:     mustHex("#A1887F"),
			Background: mustHex("#FFFDF7"),
			Surface:    mustHex("#EFEBE9"),
			Text:       mustHex("#3E2723"),
			TextLight:  mustHex("#8D6E63"),
			Success:    mustHex("#4CAF50"),
			Warning:    mustHex("#FF9800"),
			Error:      mustHex("#F44336"),
			Info:       mustHex("#607D8B"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Garamond", 26).WithWeight(common.FontWeightSemiBold),
			BodyFont:    common.NewFont("Garamond", 12),
			CaptionFont: common.NewFont("Garamond", 10).Italic().WithColor(common.Gray),
			CodeFont:    common.NewFont("Courier New", 10),
		},
		Spacing: Spacing{
			XS: common.Pt(6),
			SM: common.Pt(12),
			MD: common.Pt(20),
			LG: common.Pt(32),
			XL: common.Pt(48),
		},
		Borders: BorderStyle{
			Radius: common.Pt(4),
			Width:  common.Pt(1),
			Color:  mustHex("#BCAAA4"),
		},
	}
}

// RetailTheme returns a bright, energetic theme for retail.
func RetailTheme() Theme {
	return Theme{
		Name:  "Retail",
		Level: DesignLevelPremium,
		Palette: Palette{
			Primary:    mustHex("#E91E63"),
			Secondary:  mustHex("#FF5722"),
			Accent:     mustHex("#FFC107"),
			Background: common.White,
			Surface:    mustHex("#FFF3E0"),
			Text:       mustHex("#212121"),
			TextLight:  mustHex("#757575"),
			Success:    mustHex("#4CAF50"),
			Warning:    mustHex("#FF9800"),
			Error:      mustHex("#F44336"),
			Info:       mustHex("#2196F3"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Arial", 26).Bold(),
			BodyFont:    common.NewFont("Arial", 11),
			CaptionFont: common.NewFont("Arial", 9).WithColor(common.Gray),
			CodeFont:    common.NewFont("Consolas", 10),
		},
		Spacing: Spacing{
			XS: common.Pt(6),
			SM: common.Pt(12),
			MD: common.Pt(20),
			LG: common.Pt(32),
			XL: common.Pt(48),
		},
		Borders: BorderStyle{
			Radius: common.Pt(8),
			Width:  common.Pt(1.5),
			Color:  mustHex("#F8BBD0"),
		},
	}
}

// NonProfitTheme returns a warm green/orange approachable theme for non-profits.
func NonProfitTheme() Theme {
	return Theme{
		Name:  "Non-Profit",
		Level: DesignLevelProfessional,
		Palette: Palette{
			Primary:    mustHex("#2E7D32"),
			Secondary:  mustHex("#E65100"),
			Accent:     mustHex("#FF8F00"),
			Background: mustHex("#FFFDE7"),
			Surface:    mustHex("#F1F8E9"),
			Text:       mustHex("#1B5E20"),
			TextLight:  mustHex("#689F38"),
			Success:    mustHex("#43A047"),
			Warning:    mustHex("#FFA000"),
			Error:      mustHex("#E53935"),
			Info:       mustHex("#1E88E5"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Verdana", 24).Bold(),
			BodyFont:    common.NewFont("Verdana", 11),
			CaptionFont: common.NewFont("Verdana", 9).WithColor(common.Gray),
			CodeFont:    common.NewFont("Courier New", 10),
		},
		Spacing: Spacing{
			XS: common.Pt(4),
			SM: common.Pt(10),
			MD: common.Pt(18),
			LG: common.Pt(28),
			XL: common.Pt(44),
		},
		Borders: BorderStyle{
			Radius: common.Pt(6),
			Width:  common.Pt(1),
			Color:  mustHex("#A5D6A7"),
		},
	}
}

// GovernmentTheme returns a blue/red official theme for government documents.
func GovernmentTheme() Theme {
	return Theme{
		Name:  "Government",
		Level: DesignLevelProfessional,
		Palette: Palette{
			Primary:    mustHex("#002868"),
			Secondary:  mustHex("#BF0A30"),
			Accent:     mustHex("#3C3B6E"),
			Background: common.White,
			Surface:    mustHex("#F0F4F8"),
			Text:       mustHex("#1A1A2E"),
			TextLight:  mustHex("#6B7280"),
			Success:    mustHex("#16A34A"),
			Warning:    mustHex("#D97706"),
			Error:      mustHex("#DC2626"),
			Info:       mustHex("#2563EB"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Times New Roman", 24).Bold(),
			BodyFont:    common.NewFont("Times New Roman", 12),
			CaptionFont: common.NewFont("Times New Roman", 10).WithColor(common.Gray),
			CodeFont:    common.NewFont("Courier New", 10),
		},
		Spacing: Spacing{
			XS: common.Pt(4),
			SM: common.Pt(8),
			MD: common.Pt(16),
			LG: common.Pt(24),
			XL: common.Pt(36),
		},
		Borders: BorderStyle{
			Radius: common.Pt(0),
			Width:  common.Pt(1),
			Color:  mustHex("#9CA3AF"),
		},
	}
}

// --- Dark mode themes ---

// DarkModern returns a dark theme with neon accents.
func DarkModern() Theme {
	return Theme{
		Name:  "Dark Modern",
		Level: DesignLevelPremium,
		Palette: Palette{
			Primary:    mustHex("#818CF8"),
			Secondary:  mustHex("#6366F1"),
			Accent:     mustHex("#34D399"),
			Background: mustHex("#111827"),
			Surface:    mustHex("#1F2937"),
			Text:       mustHex("#F9FAFB"),
			TextLight:  mustHex("#9CA3AF"),
			Success:    mustHex("#34D399"),
			Warning:    mustHex("#FBBF24"),
			Error:      mustHex("#F87171"),
			Info:       mustHex("#60A5FA"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Segoe UI", 26).Bold().WithColor(mustHex("#F9FAFB")),
			BodyFont:    common.NewFont("Segoe UI", 11).WithColor(mustHex("#E5E7EB")),
			CaptionFont: common.NewFont("Segoe UI", 9).WithColor(mustHex("#9CA3AF")),
			CodeFont:    common.NewFont("Cascadia Code", 10).WithColor(mustHex("#34D399")),
		},
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
			Color:  mustHex("#374151"),
		},
	}
}

// DarkElegant returns a dark theme with gold accents.
func DarkElegant() Theme {
	return Theme{
		Name:  "Dark Elegant",
		Level: DesignLevelLuxury,
		Palette: Palette{
			Primary:    mustHex("#D4AF37"),
			Secondary:  mustHex("#B8860B"),
			Accent:     mustHex("#FFD700"),
			Background: mustHex("#1A1A1A"),
			Surface:    mustHex("#2D2D2D"),
			Text:       mustHex("#F5F5F5"),
			TextLight:  mustHex("#A3A3A3"),
			Success:    mustHex("#4ADE80"),
			Warning:    mustHex("#FCD34D"),
			Error:      mustHex("#F87171"),
			Info:       mustHex("#93C5FD"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Garamond", 26).WithWeight(common.FontWeightSemiBold).WithColor(mustHex("#D4AF37")),
			BodyFont:    common.NewFont("Garamond", 12).WithColor(mustHex("#E5E5E5")),
			CaptionFont: common.NewFont("Garamond", 10).Italic().WithColor(mustHex("#A3A3A3")),
			CodeFont:    common.NewFont("Courier New", 10).WithColor(mustHex("#FCD34D")),
		},
		Spacing: Spacing{
			XS: common.Pt(8),
			SM: common.Pt(16),
			MD: common.Pt(24),
			LG: common.Pt(40),
			XL: common.Pt(56),
		},
		Borders: BorderStyle{
			Radius: common.Pt(4),
			Width:  common.Pt(1),
			Color:  mustHex("#404040"),
		},
	}
}

// DarkMinimal returns a near-black theme with subtle contrast.
func DarkMinimal() Theme {
	return Theme{
		Name:  "Dark Minimal",
		Level: DesignLevelPremium,
		Palette: Palette{
			Primary:    mustHex("#A3A3A3"),
			Secondary:  mustHex("#737373"),
			Accent:     mustHex("#F5F5F5"),
			Background: mustHex("#0A0A0A"),
			Surface:    mustHex("#171717"),
			Text:       mustHex("#E5E5E5"),
			TextLight:  mustHex("#737373"),
			Success:    mustHex("#4ADE80"),
			Warning:    mustHex("#FDE68A"),
			Error:      mustHex("#FCA5A5"),
			Info:       mustHex("#93C5FD"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Helvetica", 24).WithWeight(common.FontWeightLight).WithColor(mustHex("#E5E5E5")),
			BodyFont:    common.NewFont("Helvetica", 11).WithColor(mustHex("#D4D4D4")),
			CaptionFont: common.NewFont("Helvetica", 9).WithColor(mustHex("#737373")),
			CodeFont:    common.NewFont("Consolas", 10).WithColor(mustHex("#A3A3A3")),
		},
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
			Color:  mustHex("#262626"),
		},
	}
}

// --- High-contrast accessible themes ---

// HighContrastLight returns a maximum contrast light theme for accessibility.
func HighContrastLight() Theme {
	return Theme{
		Name:  "High Contrast Light",
		Level: DesignLevelBasic,
		Palette: Palette{
			Primary:    common.Black,
			Secondary:  mustHex("#1A1A1A"),
			Accent:     mustHex("#0000CC"),
			Background: common.White,
			Surface:    mustHex("#F0F0F0"),
			Text:       common.Black,
			TextLight:  mustHex("#333333"),
			Success:    mustHex("#006600"),
			Warning:    mustHex("#CC6600"),
			Error:      mustHex("#CC0000"),
			Info:       mustHex("#0000CC"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Arial", 26).Bold().WithColor(common.Black),
			BodyFont:    common.NewFont("Arial", 13).WithColor(common.Black),
			CaptionFont: common.NewFont("Arial", 11).WithColor(mustHex("#333333")),
			CodeFont:    common.NewFont("Courier New", 12).WithColor(common.Black),
		},
		Spacing: Spacing{
			XS: common.Pt(6),
			SM: common.Pt(12),
			MD: common.Pt(20),
			LG: common.Pt(28),
			XL: common.Pt(40),
		},
		Borders: BorderStyle{
			Radius: common.Pt(0),
			Width:  common.Pt(2),
			Color:  common.Black,
		},
	}
}

// HighContrastDark returns a maximum contrast dark theme for accessibility.
func HighContrastDark() Theme {
	return Theme{
		Name:  "High Contrast Dark",
		Level: DesignLevelBasic,
		Palette: Palette{
			Primary:    common.White,
			Secondary:  mustHex("#E5E5E5"),
			Accent:     mustHex("#FFFF00"),
			Background: common.Black,
			Surface:    mustHex("#1A1A1A"),
			Text:       common.White,
			TextLight:  mustHex("#CCCCCC"),
			Success:    mustHex("#00FF00"),
			Warning:    mustHex("#FFFF00"),
			Error:      mustHex("#FF3333"),
			Info:       mustHex("#33CCFF"),
		},
		Typography: Typography{
			HeadingFont: common.NewFont("Arial", 26).Bold().WithColor(common.White),
			BodyFont:    common.NewFont("Arial", 13).WithColor(common.White),
			CaptionFont: common.NewFont("Arial", 11).WithColor(mustHex("#CCCCCC")),
			CodeFont:    common.NewFont("Courier New", 12).WithColor(mustHex("#FFFF00")),
		},
		Spacing: Spacing{
			XS: common.Pt(6),
			SM: common.Pt(12),
			MD: common.Pt(20),
			LG: common.Pt(28),
			XL: common.Pt(40),
		},
		Borders: BorderStyle{
			Radius: common.Pt(0),
			Width:  common.Pt(2),
			Color:  common.White,
		},
	}
}
