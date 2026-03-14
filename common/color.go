package common

import "fmt"

// Color represents an RGB color
type Color struct {
	R, G, B uint8
	A       uint8 // Alpha (0-255, 255 = opaque)
}

func NewColor(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, A: 255}
}

func NewColorWithAlpha(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

func ColorFromHex(hex string) (Color, error) {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 && len(hex) != 8 {
		return Color{}, fmt.Errorf("invalid hex color: %s", hex)
	}
	var r, g, b, a uint8
	_, err := fmt.Sscanf(hex[:6], "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hex color: %s", hex)
	}
	a = 255
	if len(hex) == 8 {
		_, err = fmt.Sscanf(hex[6:8], "%02x", &a)
		if err != nil {
			return Color{}, fmt.Errorf("invalid hex alpha: %s", hex)
		}
	}
	return Color{R: r, G: g, B: b, A: a}, nil
}

func (c Color) Hex() string {
	if c.A == 255 {
		return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
	}
	return fmt.Sprintf("#%02X%02X%02X%02X", c.R, c.G, c.B, c.A)
}

func (c Color) String() string {
	return c.Hex()
}

// Predefined colors
var (
	Black       = NewColor(0, 0, 0)
	White       = NewColor(255, 255, 255)
	Red         = NewColor(255, 0, 0)
	Green       = NewColor(0, 128, 0)
	Blue        = NewColor(0, 0, 255)
	Yellow      = NewColor(255, 255, 0)
	Orange      = NewColor(255, 165, 0)
	Purple      = NewColor(128, 0, 128)
	Gray        = NewColor(128, 128, 128)
	LightGray   = NewColor(211, 211, 211)
	DarkGray    = NewColor(64, 64, 64)
	Transparent = NewColorWithAlpha(0, 0, 0, 0)
)
