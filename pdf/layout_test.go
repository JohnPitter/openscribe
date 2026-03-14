package pdf

import (
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestMeasureStringWidth(t *testing.T) {
	// "Hello" at 12pt
	w := measureStringWidth("Hello", 12)
	if w <= 0 {
		t.Error("width should be positive")
	}

	// Empty string should be 0
	w0 := measureStringWidth("", 12)
	if w0 != 0 {
		t.Errorf("empty string width should be 0, got %.2f", w0)
	}

	// Larger font should produce wider result
	w24 := measureStringWidth("Hello", 24)
	if w24 <= w {
		t.Error("24pt should be wider than 12pt")
	}
}

func TestWrapText(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	lines := wrapText(text, 100, 12)

	if len(lines) < 2 {
		t.Errorf("expected at least 2 lines for wrapping, got %d", len(lines))
	}

	// All words should be present
	joined := strings.Join(lines, " ")
	for _, word := range strings.Fields(text) {
		if !strings.Contains(joined, word) {
			t.Errorf("missing word %q after wrapping", word)
		}
	}
}

func TestWrapTextSingleWord(t *testing.T) {
	lines := wrapText("Superlongword", 50, 12)
	if len(lines) != 1 {
		t.Errorf("single word should stay on one line, got %d lines", len(lines))
	}
	if lines[0] != "Superlongword" {
		t.Errorf("expected 'Superlongword', got %q", lines[0])
	}
}

func TestWrapTextPreservesNewlines(t *testing.T) {
	text := "Line one\nLine two\nLine three"
	lines := wrapText(text, 500, 12)

	if len(lines) != 3 {
		t.Errorf("expected 3 lines from newline-separated text, got %d", len(lines))
	}
}

func TestWrapTextEmptyParagraph(t *testing.T) {
	text := "Before\n\nAfter"
	lines := wrapText(text, 500, 12)

	if len(lines) != 3 {
		t.Errorf("expected 3 lines (with empty line), got %d", len(lines))
	}
	if lines[1] != "" {
		t.Errorf("expected empty line, got %q", lines[1])
	}
}

func TestAddTextBlock(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 12)
	tb := p.AddTextBlock(72, 72, 400, "Hello world", font)

	if tb.Text() != "Hello world" {
		t.Errorf("expected 'Hello world', got %q", tb.Text())
	}
	if p.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", p.ElementCount())
	}
}

func TestTextBlockAlignment(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 12)
	tb := p.AddTextBlock(72, 72, 400, "Test", font)

	if tb.Alignment() != common.TextAlignLeft {
		t.Error("default alignment should be left")
	}

	tb.SetAlignment(common.TextAlignCenter)
	if tb.Alignment() != common.TextAlignCenter {
		t.Error("alignment should be center")
	}

	tb.SetAlignment(common.TextAlignRight)
	if tb.Alignment() != common.TextAlignRight {
		t.Error("alignment should be right")
	}

	tb.SetAlignment(common.TextAlignJustify)
	if tb.Alignment() != common.TextAlignJustify {
		t.Error("alignment should be justify")
	}
}

func TestTextBlockLineSpacing(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 12)
	tb := p.AddTextBlock(72, 72, 400, "Test", font)

	if tb.LineSpacing() != 1.2 {
		t.Errorf("default line spacing should be 1.2, got %.2f", tb.LineSpacing())
	}

	tb.SetLineSpacing(1.5)
	if tb.LineSpacing() != 1.5 {
		t.Errorf("expected 1.5, got %.2f", tb.LineSpacing())
	}

	// Test clamping
	tb.SetLineSpacing(0.1)
	if tb.LineSpacing() != 0.5 {
		t.Errorf("expected clamped to 0.5, got %.2f", tb.LineSpacing())
	}
}

func TestTextBlockColumns(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 12)
	tb := p.AddTextBlock(72, 72, 400, "Test", font)

	if tb.Columns() != 1 {
		t.Errorf("default columns should be 1, got %d", tb.Columns())
	}

	tb.SetColumns(2, 20)
	if tb.Columns() != 2 {
		t.Errorf("expected 2 columns, got %d", tb.Columns())
	}
	if tb.ColumnGap() != 20 {
		t.Errorf("expected gap 20, got %.2f", tb.ColumnGap())
	}

	// Negative columns should clamp to 1
	tb.SetColumns(0, 10)
	if tb.Columns() != 1 {
		t.Errorf("expected clamped to 1, got %d", tb.Columns())
	}
}

func TestTextBlockWrapLines(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 12)
	text := "This is a fairly long paragraph that should wrap across multiple lines when rendered in a narrow column."
	tb := p.AddTextBlock(72, 72, 150, text, font)

	lines := tb.WrapLines()
	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines for narrow width, got %d", len(lines))
	}
}

func TestTextBlockMultiColumnBuild(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 10)
	text := "Word one two three four five six seven eight nine ten eleven twelve thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty."
	tb := p.AddTextBlock(72, 72, 400, text, font)
	tb.SetColumns(2, 20)

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("should produce output")
	}

	content := string(data)
	if !strings.Contains(content, "BT") {
		t.Error("should contain text operators")
	}
}

func TestTextBlockCenterAlignBuild(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 12)
	tb := p.AddTextBlock(72, 72, 400, "Centered text line", font)
	tb.SetAlignment(common.TextAlignCenter)

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("should produce output")
	}
}

func TestTextBlockRightAlignBuild(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 12)
	tb := p.AddTextBlock(72, 72, 400, "Right-aligned text", font)
	tb.SetAlignment(common.TextAlignRight)

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}
	if len(data) == 0 {
		t.Error("should produce output")
	}
}

func TestColumnWidth(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 12)
	tb := p.AddTextBlock(72, 72, 400, "Test", font)

	// Single column: full width
	cw := tb.columnWidth()
	if cw != 400 {
		t.Errorf("single column width should be 400, got %.2f", cw)
	}

	// Two columns with 20 gap: (400-20)/2 = 190
	tb.SetColumns(2, 20)
	cw = tb.columnWidth()
	expected := (400.0 - 20.0) / 2.0
	if cw != expected {
		t.Errorf("expected column width %.2f, got %.2f", expected, cw)
	}

	// Three columns with 10 gap: (400-20)/3
	tb.SetColumns(3, 10)
	cw = tb.columnWidth()
	expected = (400.0 - 20.0) / 3.0
	if cw != expected {
		t.Errorf("expected column width %.2f, got %.2f", expected, cw)
	}
}
