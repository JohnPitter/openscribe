package pdf

import (
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestExtractTextFromElements(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	p.AddText("Hello World", 72, 72, common.NewFont("Helvetica", 12))
	p.AddText("Second line", 72, 90, common.NewFont("Helvetica", 12))

	text, err := doc.ExtractText()
	if err != nil {
		t.Fatalf("extract error: %v", err)
	}
	if !strings.Contains(text, "Hello World") {
		t.Error("should contain 'Hello World'")
	}
	if !strings.Contains(text, "Second line") {
		t.Error("should contain 'Second line'")
	}
}

func TestExtractTextFromTable(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	tbl := p.AddTable(72, 72, 2, 2)
	tbl.SetCell(0, 0, "A1")
	tbl.SetCell(0, 1, "B1")
	tbl.SetCell(1, 0, "A2")
	tbl.SetCell(1, 1, "B2")

	text, err := doc.ExtractText()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(text, "A1") {
		t.Error("should contain A1")
	}
	if !strings.Contains(text, "B2") {
		t.Error("should contain B2")
	}
}

func TestExtractPageText(t *testing.T) {
	doc := New()
	doc.AddPage().AddText("Page 1", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 2", 72, 72, common.NewFont("Helvetica", 12))

	text1, err := doc.ExtractPageText(0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(text1, "Page 1") {
		t.Error("should contain 'Page 1'")
	}

	text2, err := doc.ExtractPageText(1)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(text2, "Page 2") {
		t.Error("should contain 'Page 2'")
	}

	_, err = doc.ExtractPageText(5)
	if err == nil {
		t.Error("should error on invalid index")
	}
}

func TestExtractTextRoundTrip(t *testing.T) {
	// Create PDF, save, open, extract
	doc := New()
	p := doc.AddPage()
	p.AddText("Round trip test", 72, 72, common.NewFont("Helvetica", 14))
	p.AddText("Second paragraph", 72, 100, common.NewFont("Helvetica", 12))

	data, err := doc.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	doc2, err := OpenFromBytes(data)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}

	text, err := doc2.ExtractText()
	if err != nil {
		t.Fatalf("extract error: %v", err)
	}
	if !strings.Contains(text, "Round trip test") {
		t.Error("should extract 'Round trip test' from round-trip")
	}
}

func TestExtractTextEmpty(t *testing.T) {
	doc := New()
	doc.AddPage()

	text, err := doc.ExtractText()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if text != "" {
		t.Error("empty page should have empty text")
	}
}

func TestExtractMultiplePages(t *testing.T) {
	doc := New()
	doc.AddPage().AddText("First", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Second", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Third", 72, 72, common.NewFont("Helvetica", 12))

	text, err := doc.ExtractText()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(text, "First") || !strings.Contains(text, "Second") || !strings.Contains(text, "Third") {
		t.Error("should contain text from all pages")
	}
	if !strings.Contains(text, "Page Break") {
		t.Error("should have page break markers")
	}
}

func TestUnescapePDFString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello", "Hello"},
		{"Hello\\nWorld", "Hello\nWorld"},
		{"\\(parens\\)", "(parens)"},
		{"back\\\\slash", "back\\slash"},
		{"tab\\there", "tab\there"},
	}
	for _, tt := range tests {
		got := unescapePDFString(tt.input)
		if got != tt.want {
			t.Errorf("unescapePDFString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
