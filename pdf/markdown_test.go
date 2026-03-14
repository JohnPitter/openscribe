package pdf

import (
	"testing"
)

func TestFromMarkdownHeadings(t *testing.T) {
	md := "# Heading 1\n## Heading 2\n### Heading 3"
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	if doc.PageCount() < 1 {
		t.Fatal("expected at least 1 page")
	}
	page := doc.Page(0)
	if page.ElementCount() == 0 {
		t.Error("expected elements on the page")
	}
	// Verify we have text elements for each heading
	textCount := 0
	for _, elem := range page.Elements() {
		if _, ok := elem.(*TextElement); ok {
			textCount++
		}
	}
	if textCount < 3 {
		t.Errorf("expected at least 3 text elements for headings, got %d", textCount)
	}
}

func TestFromMarkdownParagraph(t *testing.T) {
	md := "This is a paragraph.\n\nThis is another paragraph."
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	if doc.PageCount() < 1 {
		t.Fatal("expected at least 1 page")
	}
	page := doc.Page(0)
	textCount := 0
	for _, elem := range page.Elements() {
		if _, ok := elem.(*TextElement); ok {
			textCount++
		}
	}
	if textCount < 2 {
		t.Errorf("expected at least 2 text elements, got %d", textCount)
	}
}

func TestFromMarkdownBoldItalic(t *testing.T) {
	md := "This has **bold** and *italic* text."
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	if doc.PageCount() < 1 {
		t.Fatal("expected at least 1 page")
	}
	// Should produce multiple text elements for inline segments
	page := doc.Page(0)
	if page.ElementCount() == 0 {
		t.Error("expected elements for bold/italic text")
	}
}

func TestFromMarkdownLists(t *testing.T) {
	md := "- Item 1\n- Item 2\n- Item 3"
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	page := doc.Page(0)
	textCount := 0
	for _, elem := range page.Elements() {
		if _, ok := elem.(*TextElement); ok {
			textCount++
		}
	}
	// Each list item produces a bullet text + item text = 6 texts
	if textCount < 3 {
		t.Errorf("expected at least 3 text elements for list items, got %d", textCount)
	}
}

func TestFromMarkdownOrderedList(t *testing.T) {
	md := "1. First\n2. Second\n3. Third"
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	page := doc.Page(0)
	if page.ElementCount() == 0 {
		t.Error("expected elements for ordered list")
	}
}

func TestFromMarkdownCodeBlock(t *testing.T) {
	md := "```\nfunc main() {\n  fmt.Println(\"hello\")\n}\n```"
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	page := doc.Page(0)
	// Code block should produce text elements and rectangle backgrounds
	hasRect := false
	hasText := false
	for _, elem := range page.Elements() {
		if _, ok := elem.(*RectElement); ok {
			hasRect = true
		}
		if _, ok := elem.(*TextElement); ok {
			hasText = true
		}
	}
	if !hasRect {
		t.Error("expected rectangle background for code block")
	}
	if !hasText {
		t.Error("expected text elements for code block")
	}
}

func TestFromMarkdownHorizontalRule(t *testing.T) {
	md := "Above\n\n---\n\nBelow"
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	page := doc.Page(0)
	hasLine := false
	for _, elem := range page.Elements() {
		if _, ok := elem.(*LineElement); ok {
			hasLine = true
		}
	}
	if !hasLine {
		t.Error("expected line element for horizontal rule")
	}
}

func TestFromMarkdownBlockquote(t *testing.T) {
	md := "> This is a quote"
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	page := doc.Page(0)
	hasLine := false
	hasText := false
	for _, elem := range page.Elements() {
		if _, ok := elem.(*LineElement); ok {
			hasLine = true
		}
		if _, ok := elem.(*TextElement); ok {
			hasText = true
		}
	}
	if !hasLine {
		t.Error("expected line element for blockquote border")
	}
	if !hasText {
		t.Error("expected text element for blockquote content")
	}
}

func TestFromMarkdownLink(t *testing.T) {
	md := "Visit [OpenScribe](https://example.com) for more."
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	page := doc.Page(0)
	if page.ElementCount() == 0 {
		t.Error("expected elements for link text")
	}
}

func TestFromMarkdownMultiPage(t *testing.T) {
	// Generate enough content to overflow to a second page
	md := ""
	for i := 0; i < 100; i++ {
		md += "This is a long paragraph that should eventually cause a page break. "
		md += "Adding more content to fill the page.\n\n"
	}
	doc, err := FromMarkdown(md, DefaultHTMLOptions())
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	if doc.PageCount() < 2 {
		t.Errorf("expected at least 2 pages for long content, got %d", doc.PageCount())
	}
}

func TestParseInlineMarkdown(t *testing.T) {
	segments := parseInlineMarkdown("Hello **bold** and *italic* world")
	if len(segments) < 4 {
		t.Fatalf("expected at least 4 segments, got %d", len(segments))
	}

	// Check that bold segment exists
	foundBold := false
	foundItalic := false
	for _, seg := range segments {
		if seg.bold && seg.text == "bold" {
			foundBold = true
		}
		if seg.italic && seg.text == "italic" {
			foundItalic = true
		}
	}
	if !foundBold {
		t.Error("expected bold segment")
	}
	if !foundItalic {
		t.Error("expected italic segment")
	}
}

func TestIsHorizontalRule(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"---", true},
		{"***", true},
		{"___", true},
		{"- - -", true},
		{"--", false},
		{"hello", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := isHorizontalRule(tt.input); got != tt.want {
			t.Errorf("isHorizontalRule(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
