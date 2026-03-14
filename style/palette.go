package style

import "github.com/JohnPitter/openscribe/common"

// Pre-built palettes

func CorporatePalette() Palette {
	return Palette{
		Primary:    mustHex("#1A365D"),
		Secondary:  mustHex("#2D3748"),
		Accent:     mustHex("#3182CE"),
		Background: common.White,
		Surface:    mustHex("#F7FAFC"),
		Text:       mustHex("#1A202C"),
		TextLight:  mustHex("#718096"),
		Success:    mustHex("#38A169"),
		Warning:    mustHex("#D69E2E"),
		Error:      mustHex("#E53E3E"),
		Info:       mustHex("#3182CE"),
	}
}

func ModernPalette() Palette {
	return Palette{
		Primary:    mustHex("#6C63FF"),
		Secondary:  mustHex("#3F3D56"),
		Accent:     mustHex("#FF6584"),
		Background: common.White,
		Surface:    mustHex("#F8F9FA"),
		Text:       mustHex("#2D2D2D"),
		TextLight:  mustHex("#8D8D8D"),
		Success:    mustHex("#00C48C"),
		Warning:    mustHex("#FFB800"),
		Error:      mustHex("#FF4B4B"),
		Info:       mustHex("#00B4D8"),
	}
}

func ElegantPalette() Palette {
	return Palette{
		Primary:    mustHex("#2C3E50"),
		Secondary:  mustHex("#34495E"),
		Accent:     mustHex("#E74C3C"),
		Background: mustHex("#FDFEFE"),
		Surface:    mustHex("#F2F3F4"),
		Text:       mustHex("#17202A"),
		TextLight:  mustHex("#7F8C8D"),
		Success:    mustHex("#27AE60"),
		Warning:    mustHex("#F39C12"),
		Error:      mustHex("#C0392B"),
		Info:       mustHex("#2980B9"),
	}
}

func MinimalPalette() Palette {
	return Palette{
		Primary:    mustHex("#111111"),
		Secondary:  mustHex("#333333"),
		Accent:     mustHex("#0066FF"),
		Background: common.White,
		Surface:    mustHex("#FAFAFA"),
		Text:       mustHex("#111111"),
		TextLight:  mustHex("#999999"),
		Success:    mustHex("#00AA55"),
		Warning:    mustHex("#FFAA00"),
		Error:      mustHex("#FF3333"),
		Info:       mustHex("#0066FF"),
	}
}

func WarmPalette() Palette {
	return Palette{
		Primary:    mustHex("#D35400"),
		Secondary:  mustHex("#8E4400"),
		Accent:     mustHex("#F39C12"),
		Background: mustHex("#FFFDF7"),
		Surface:    mustHex("#FEF9E7"),
		Text:       mustHex("#1C1917"),
		TextLight:  mustHex("#78716C"),
		Success:    mustHex("#6D9B05"),
		Warning:    mustHex("#F59E0B"),
		Error:      mustHex("#DC2626"),
		Info:       mustHex("#0284C7"),
	}
}

func NeonPalette() Palette {
	return Palette{
		Primary:    mustHex("#00FF87"),
		Secondary:  mustHex("#0A0A0A"),
		Accent:     mustHex("#FF00E5"),
		Background: mustHex("#0A0A0A"),
		Surface:    mustHex("#1A1A1A"),
		Text:       mustHex("#FFFFFF"),
		TextLight:  mustHex("#AAAAAA"),
		Success:    mustHex("#00FF87"),
		Warning:    mustHex("#FFE500"),
		Error:      mustHex("#FF0055"),
		Info:       mustHex("#00D4FF"),
	}
}

func mustHex(hex string) common.Color {
	c, err := common.ColorFromHex(hex)
	if err != nil {
		panic("invalid hex color in preset: " + hex)
	}
	return c
}
