package pdf

import (
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestFromHTMLBasic(t *testing.T) {
	html := `<html><body><h1>Hello World</h1><p>This is a paragraph.</p></body></html>`

	doc, err := FromHTML(html, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if doc.PageCount() == 0 {
		t.Error("should have at least 1 page")
	}

	path := filepath.Join(t.TempDir(), "basic.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestFromHTMLHeadings(t *testing.T) {
	html := `
		<h1>Heading 1</h1>
		<h2>Heading 2</h2>
		<h3>Heading 3</h3>
		<h4>Heading 4</h4>
		<h5>Heading 5</h5>
		<h6>Heading 6</h6>
	`

	doc, err := FromHTML(html, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "headings.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestFromHTMLFormatting(t *testing.T) {
	html := `
		<p>Normal text</p>
		<p><b>Bold text</b></p>
		<p><i>Italic text</i></p>
		<p><strong>Strong text</strong></p>
		<p><em>Emphasized text</em></p>
	`

	doc, err := FromHTML(html, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "formatting.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestFromHTMLList(t *testing.T) {
	html := `
		<h2>Shopping List</h2>
		<ul>
			<li>Apples</li>
			<li>Bananas</li>
			<li>Oranges</li>
		</ul>
	`

	doc, err := FromHTML(html, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "list.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestFromHTMLEntities(t *testing.T) {
	html := `<p>Copyright &copy; 2026 &amp; Trademark&trade;</p>`

	doc, err := FromHTML(html, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	text, _ := doc.ExtractText()
	if text == "" {
		t.Error("should have text content")
	}
}

func TestFromHTMLHR(t *testing.T) {
	html := `<p>Above</p><hr><p>Below</p>`

	doc, err := FromHTML(html, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "hr.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestFromHTMLCustomOptions(t *testing.T) {
	html := `<h1>Custom</h1><p>Letter-sized document</p>`

	opts := HTMLOptions{
		PageSize:    common.PageLetter,
		Margins:     common.NarrowMargins(),
		DefaultFont: common.NewFont("Helvetica", 14),
	}

	doc, err := FromHTML(html, opts)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "custom.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestFromHTMLMultiPage(t *testing.T) {
	// Generate enough content to span multiple pages
	html := "<h1>Long Document</h1>"
	for i := 0; i < 80; i++ {
		html += "<p>This is paragraph number that should fill the page with content.</p>"
	}

	doc, err := FromHTML(html, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if doc.PageCount() < 2 {
		t.Errorf("should span multiple pages, got %d", doc.PageCount())
	}

	path := filepath.Join(t.TempDir(), "multipage.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestFromHTMLEmpty(t *testing.T) {
	doc, err := FromHTML("", DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if doc.PageCount() != 1 {
		t.Error("empty HTML should produce 1 page")
	}
}

func TestDecodeHTMLEntities(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"&amp;", "&"},
		{"&lt;", "<"},
		{"&gt;", ">"},
		{"&quot;", "\""},
		{"&copy;", "\xc2\xa9"},
		{"&trade;", "\xe2\x84\xa2"},
		{"&bull;", "\xe2\x80\xa2"},
		{"no entities", "no entities"},
	}
	for _, tt := range tests {
		got := decodeHTMLEntities(tt.input)
		if got != tt.want {
			t.Errorf("decode(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
