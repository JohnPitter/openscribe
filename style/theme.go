// Package style provides design system themes, palettes, and typography for openscribe.
package style

import "github.com/JohnPitter/openscribe/common"

// DesignLevel represents the quality tier of a design
type DesignLevel int

const (
	DesignLevelBasic        DesignLevel = iota // Clean, minimal
	DesignLevelProfessional                    // Business-grade
	DesignLevelPremium                         // High-quality (Behance/Freepik)
	DesignLevelLuxury                          // Ultra-premium (Slidesgo/Agency)
)

func (d DesignLevel) String() string {
	switch d {
	case DesignLevelBasic:
		return "Basic"
	case DesignLevelProfessional:
		return "Professional"
	case DesignLevelPremium:
		return "Premium"
	case DesignLevelLuxury:
		return "Luxury"
	default:
		return "Unknown"
	}
}

// Theme defines a complete visual design system
type Theme struct {
	Name       string
	Level      DesignLevel
	Palette    Palette
	Typography Typography
	Spacing    Spacing
	Borders    BorderStyle
}

// Palette defines the color scheme
type Palette struct {
	Primary    common.Color
	Secondary  common.Color
	Accent     common.Color
	Background common.Color
	Surface    common.Color
	Text       common.Color
	TextLight  common.Color
	Success    common.Color
	Warning    common.Color
	Error      common.Color
	Info       common.Color
}

// Typography defines the font scheme
type Typography struct {
	HeadingFont common.Font
	BodyFont    common.Font
	CaptionFont common.Font
	CodeFont    common.Font
}

// Spacing defines consistent spacing values
type Spacing struct {
	XS common.Measurement // Extra small (4pt)
	SM common.Measurement // Small (8pt)
	MD common.Measurement // Medium (16pt)
	LG common.Measurement // Large (24pt)
	XL common.Measurement // Extra large (32pt)
}

// BorderStyle defines border defaults
type BorderStyle struct {
	Radius common.Measurement
	Width  common.Measurement
	Color  common.Color
}
