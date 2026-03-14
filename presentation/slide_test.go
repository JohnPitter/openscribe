package presentation

import (
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestTextBoxParagraphs(t *testing.T) {
	p := New()
	s := p.AddSlide()

	tb := s.AddTextBox(common.In(1), common.In(1), common.In(5), common.In(3))

	p1 := tb.AddParagraph()
	p1.AddRun("Title", common.NewFont("Arial", 28).Bold())
	p1.SetAlignment(common.TextAlignCenter)

	p2 := tb.AddParagraph()
	p2.AddRun("Subtitle", common.NewFont("Arial", 16))

	if len(tb.Paragraphs()) != 2 {
		t.Errorf("expected 2 paragraphs, got %d", len(tb.Paragraphs()))
	}

	if tb.Text() != "Title\nSubtitle" {
		t.Errorf("unexpected text: %s", tb.Text())
	}
}

func TestShapeText(t *testing.T) {
	sh := &Shape{
		shapeType: ShapeRoundedRectangle,
		fillColor: common.Blue,
	}
	font := common.NewFont("Arial", 14)
	sh.SetText("Click Me", font)

	if sh.text != "Click Me" {
		t.Error("shape text should be set")
	}
}

func TestTextBoxSetFill(t *testing.T) {
	tb := &TextBox{}
	tb.SetFill(common.Yellow)
	if tb.fillColor == nil || *tb.fillColor != common.Yellow {
		t.Error("fill should be yellow")
	}
}

func TestTextBoxSetBorder(t *testing.T) {
	tb := &TextBox{}
	tb.SetBorder(common.Black, common.Pt(2))
	if tb.borderColor == nil || *tb.borderColor != common.Black {
		t.Error("border should be black")
	}
}
