package common

import "testing"

func TestNewFont(t *testing.T) {
	f := NewFont("Arial", 12)
	if f.Family != "Arial" || f.Size != 12 || f.Weight != FontWeightRegular {
		t.Errorf("unexpected font: %+v", f)
	}
}

func TestFontChaining(t *testing.T) {
	f := NewFont("Arial", 12).Bold().Italic().Underline().WithColor(Red).WithSize(16)
	if f.Weight != FontWeightBold {
		t.Error("expected bold")
	}
	if f.Style != FontStyleItalic {
		t.Error("expected italic")
	}
	if f.Decoration != TextDecorationUnderline {
		t.Error("expected underline")
	}
	if f.Color != Red {
		t.Error("expected red color")
	}
	if f.Size != 16 {
		t.Error("expected size 16")
	}
}
