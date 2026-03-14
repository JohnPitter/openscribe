package common

import (
	"testing"
)

func TestNewColor(t *testing.T) {
	c := NewColor(255, 128, 0)
	if c.R != 255 || c.G != 128 || c.B != 0 || c.A != 255 {
		t.Errorf("expected RGBA(255,128,0,255), got RGBA(%d,%d,%d,%d)", c.R, c.G, c.B, c.A)
	}
}

func TestNewColorWithAlpha(t *testing.T) {
	c := NewColorWithAlpha(255, 128, 0, 128)
	if c.A != 128 {
		t.Errorf("expected alpha 128, got %d", c.A)
	}
}

func TestColorFromHex(t *testing.T) {
	tests := []struct {
		hex     string
		want    Color
		wantErr bool
	}{
		{"#FF8000", NewColor(255, 128, 0), false},
		{"FF8000", NewColor(255, 128, 0), false},
		{"#ff8000", NewColor(255, 128, 0), false},
		{"#FF800080", NewColorWithAlpha(255, 128, 0, 128), false},
		{"#FFF", Color{}, true},
		{"ZZZZZZ", Color{}, true},
		{"", Color{}, true},
	}

	for _, tt := range tests {
		got, err := ColorFromHex(tt.hex)
		if (err != nil) != tt.wantErr {
			t.Errorf("ColorFromHex(%q) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ColorFromHex(%q) = %v, want %v", tt.hex, got, tt.want)
		}
	}
}

func TestColorHex(t *testing.T) {
	c := NewColor(255, 128, 0)
	if c.Hex() != "#FF8000" {
		t.Errorf("expected #FF8000, got %s", c.Hex())
	}

	ca := NewColorWithAlpha(255, 128, 0, 128)
	if ca.Hex() != "#FF800080" {
		t.Errorf("expected #FF800080, got %s", ca.Hex())
	}
}

func TestPredefinedColors(t *testing.T) {
	if Black.R != 0 || Black.G != 0 || Black.B != 0 {
		t.Error("Black should be (0,0,0)")
	}
	if White.R != 255 || White.G != 255 || White.B != 255 {
		t.Error("White should be (255,255,255)")
	}
	if Transparent.A != 0 {
		t.Error("Transparent alpha should be 0")
	}
}
