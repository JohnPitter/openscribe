package document

import (
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestParagraphText(t *testing.T) {
	p := NewParagraph()
	p.AddText("Hello ")
	p.AddText("World")

	if p.Text() != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", p.Text())
	}
}

func TestParagraphStyle(t *testing.T) {
	p := NewParagraph()
	p.SetStyle("Heading1")
	if p.Style() != "Heading1" {
		t.Errorf("expected Heading1, got %s", p.Style())
	}
}

func TestParagraphAlignment(t *testing.T) {
	p := NewParagraph()
	p.SetAlignment(common.TextAlignCenter)
	if p.Alignment() != common.TextAlignCenter {
		t.Error("expected center alignment")
	}
}

func TestRunFormatting(t *testing.T) {
	r := NewRun()
	r.SetText("Bold text").SetBold(true).SetItalic(true).SetColor(common.Red)

	if r.Text() != "Bold text" {
		t.Error("unexpected text")
	}
	if !r.bold {
		t.Error("should be bold")
	}
	if !r.italic {
		t.Error("should be italic")
	}
}
