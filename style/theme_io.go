package style

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JohnPitter/openscribe/common"
)

// themeJSON is the JSON-serializable representation of a Theme.
type themeJSON struct {
	Name       string         `json:"name"`
	Level      string         `json:"level"`
	Palette    paletteJSON    `json:"palette"`
	Typography typographyJSON `json:"typography"`
	Spacing    spacingJSON    `json:"spacing"`
	Borders    borderJSON     `json:"borders"`
}

type paletteJSON struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Accent     string `json:"accent"`
	Background string `json:"background"`
	Surface    string `json:"surface"`
	Text       string `json:"text"`
	TextLight  string `json:"textLight"`
	Success    string `json:"success"`
	Warning    string `json:"warning"`
	Error      string `json:"error"`
	Info       string `json:"info"`
}

type fontJSON struct {
	Family     string  `json:"family"`
	Size       float64 `json:"size"`
	Weight     int     `json:"weight"`
	Style      int     `json:"style"`
	Color      string  `json:"color"`
	Decoration int     `json:"decoration"`
}

type typographyJSON struct {
	HeadingFont fontJSON `json:"headingFont"`
	BodyFont    fontJSON `json:"bodyFont"`
	CaptionFont fontJSON `json:"captionFont"`
	CodeFont    fontJSON `json:"codeFont"`
}

type spacingJSON struct {
	XS float64 `json:"xs"`
	SM float64 `json:"sm"`
	MD float64 `json:"md"`
	LG float64 `json:"lg"`
	XL float64 `json:"xl"`
}

type borderJSON struct {
	Radius float64 `json:"radius"`
	Width  float64 `json:"width"`
	Color  string  `json:"color"`
}

// ExportTheme serializes a theme to JSON bytes.
func ExportTheme(theme Theme) ([]byte, error) {
	tj := themeJSON{
		Name:  theme.Name,
		Level: theme.Level.String(),
		Palette: paletteJSON{
			Primary:    theme.Palette.Primary.Hex(),
			Secondary:  theme.Palette.Secondary.Hex(),
			Accent:     theme.Palette.Accent.Hex(),
			Background: theme.Palette.Background.Hex(),
			Surface:    theme.Palette.Surface.Hex(),
			Text:       theme.Palette.Text.Hex(),
			TextLight:  theme.Palette.TextLight.Hex(),
			Success:    theme.Palette.Success.Hex(),
			Warning:    theme.Palette.Warning.Hex(),
			Error:      theme.Palette.Error.Hex(),
			Info:       theme.Palette.Info.Hex(),
		},
		Typography: typographyJSON{
			HeadingFont: fontToJSON(theme.Typography.HeadingFont),
			BodyFont:    fontToJSON(theme.Typography.BodyFont),
			CaptionFont: fontToJSON(theme.Typography.CaptionFont),
			CodeFont:    fontToJSON(theme.Typography.CodeFont),
		},
		Spacing: spacingJSON{
			XS: theme.Spacing.XS.Points(),
			SM: theme.Spacing.SM.Points(),
			MD: theme.Spacing.MD.Points(),
			LG: theme.Spacing.LG.Points(),
			XL: theme.Spacing.XL.Points(),
		},
		Borders: borderJSON{
			Radius: theme.Borders.Radius.Points(),
			Width:  theme.Borders.Width.Points(),
			Color:  theme.Borders.Color.Hex(),
		},
	}

	data, err := json.MarshalIndent(tj, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal theme: %w", err)
	}
	return data, nil
}

// ImportTheme deserializes a theme from JSON bytes.
func ImportTheme(data []byte) (Theme, error) {
	var tj themeJSON
	if err := json.Unmarshal(data, &tj); err != nil {
		return Theme{}, fmt.Errorf("unmarshal theme: %w", err)
	}

	theme := Theme{
		Name:  tj.Name,
		Level: parseDesignLevel(tj.Level),
		Palette: Palette{
			Primary:    mustParseColor(tj.Palette.Primary),
			Secondary:  mustParseColor(tj.Palette.Secondary),
			Accent:     mustParseColor(tj.Palette.Accent),
			Background: mustParseColor(tj.Palette.Background),
			Surface:    mustParseColor(tj.Palette.Surface),
			Text:       mustParseColor(tj.Palette.Text),
			TextLight:  mustParseColor(tj.Palette.TextLight),
			Success:    mustParseColor(tj.Palette.Success),
			Warning:    mustParseColor(tj.Palette.Warning),
			Error:      mustParseColor(tj.Palette.Error),
			Info:       mustParseColor(tj.Palette.Info),
		},
		Typography: Typography{
			HeadingFont: fontFromJSON(tj.Typography.HeadingFont),
			BodyFont:    fontFromJSON(tj.Typography.BodyFont),
			CaptionFont: fontFromJSON(tj.Typography.CaptionFont),
			CodeFont:    fontFromJSON(tj.Typography.CodeFont),
		},
		Spacing: Spacing{
			XS: common.Pt(tj.Spacing.XS),
			SM: common.Pt(tj.Spacing.SM),
			MD: common.Pt(tj.Spacing.MD),
			LG: common.Pt(tj.Spacing.LG),
			XL: common.Pt(tj.Spacing.XL),
		},
		Borders: BorderStyle{
			Radius: common.Pt(tj.Borders.Radius),
			Width:  common.Pt(tj.Borders.Width),
			Color:  mustParseColor(tj.Borders.Color),
		},
	}

	return theme, nil
}

// ExportThemeToFile writes a theme to a JSON file.
func ExportThemeToFile(theme Theme, path string) error {
	data, err := ExportTheme(theme)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ImportThemeFromFile reads a theme from a JSON file.
func ImportThemeFromFile(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Theme{}, fmt.Errorf("read theme file: %w", err)
	}
	return ImportTheme(data)
}

func fontToJSON(f common.Font) fontJSON {
	return fontJSON{
		Family:     f.Family,
		Size:       f.Size,
		Weight:     int(f.Weight),
		Style:      int(f.Style),
		Color:      f.Color.Hex(),
		Decoration: int(f.Decoration),
	}
}

func fontFromJSON(fj fontJSON) common.Font {
	return common.Font{
		Family:     fj.Family,
		Size:       fj.Size,
		Weight:     common.FontWeight(fj.Weight),
		Style:      common.FontStyle(fj.Style),
		Color:      mustParseColor(fj.Color),
		Decoration: common.TextDecoration(fj.Decoration),
	}
}

func mustParseColor(hex string) common.Color {
	if hex == "" {
		return common.Black
	}
	c, err := common.ColorFromHex(hex)
	if err != nil {
		return common.Black
	}
	return c
}

func parseDesignLevel(s string) DesignLevel {
	switch s {
	case "Basic":
		return DesignLevelBasic
	case "Professional":
		return DesignLevelProfessional
	case "Premium":
		return DesignLevelPremium
	case "Luxury":
		return DesignLevelLuxury
	default:
		return DesignLevelBasic
	}
}
