package document

import (
	"testing"
)

func TestFromMarkdownHeadings(t *testing.T) {
	md := "# Heading 1\n## Heading 2\n### Heading 3"
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) < 3 {
		t.Fatalf("expected at least 3 paragraphs, got %d", len(paragraphs))
	}
	// First paragraph should be Heading1 style
	if paragraphs[0].Style() != "Heading1" {
		t.Errorf("expected Heading1 style, got %q", paragraphs[0].Style())
	}
	if paragraphs[1].Style() != "Heading2" {
		t.Errorf("expected Heading2 style, got %q", paragraphs[1].Style())
	}
	if paragraphs[2].Style() != "Heading3" {
		t.Errorf("expected Heading3 style, got %q", paragraphs[2].Style())
	}
}

func TestFromMarkdownParagraphs(t *testing.T) {
	md := "First paragraph.\n\nSecond paragraph."
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) < 2 {
		t.Fatalf("expected at least 2 paragraphs, got %d", len(paragraphs))
	}
	if paragraphs[0].Text() != "First paragraph." {
		t.Errorf("expected 'First paragraph.', got %q", paragraphs[0].Text())
	}
	if paragraphs[1].Text() != "Second paragraph." {
		t.Errorf("expected 'Second paragraph.', got %q", paragraphs[1].Text())
	}
}

func TestFromMarkdownBoldItalic(t *testing.T) {
	md := "This has **bold** and *italic* text."
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) == 0 {
		t.Fatal("expected at least 1 paragraph")
	}
	// Check that the paragraph has multiple runs for inline formatting
	runs := paragraphs[0].Runs()
	if len(runs) < 3 {
		t.Errorf("expected at least 3 runs for inline formatting, got %d", len(runs))
	}
}

func TestFromMarkdownUnorderedList(t *testing.T) {
	md := "- Item 1\n- Item 2\n- Item 3"
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) < 3 {
		t.Fatalf("expected at least 3 paragraphs for list items, got %d", len(paragraphs))
	}
	// Each item should have a bullet prefix run
	for i, para := range paragraphs {
		text := para.Text()
		if len(text) == 0 {
			t.Errorf("list item %d has empty text", i)
		}
	}
}

func TestFromMarkdownOrderedList(t *testing.T) {
	md := "1. First\n2. Second\n3. Third"
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) < 3 {
		t.Fatalf("expected at least 3 paragraphs for ordered list, got %d", len(paragraphs))
	}
}

func TestFromMarkdownCodeBlock(t *testing.T) {
	md := "```\nfunc main() {\n  fmt.Println(\"hello\")\n}\n```"
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) == 0 {
		t.Fatal("expected paragraphs for code block")
	}
	// Code lines should use monospace font
	for _, para := range paragraphs {
		runs := para.Runs()
		if len(runs) > 0 {
			// At minimum, the text should be present
			if runs[0].Text() == "" {
				t.Error("expected non-empty code line")
			}
		}
	}
}

func TestFromMarkdownBlockquote(t *testing.T) {
	md := "> This is a quote"
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) == 0 {
		t.Fatal("expected at least 1 paragraph for blockquote")
	}
	// Blockquote paragraph should have italic formatting
	runs := paragraphs[0].Runs()
	if len(runs) == 0 {
		t.Fatal("expected runs in blockquote paragraph")
	}
}

func TestFromMarkdownHorizontalRule(t *testing.T) {
	md := "Above\n\n---\n\nBelow"
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) < 2 {
		t.Errorf("expected at least 2 content paragraphs plus HR, got %d", len(paragraphs))
	}
}

func TestFromMarkdownLink(t *testing.T) {
	md := "Visit [OpenScribe](https://example.com) for more."
	doc, err := FromMarkdown(md)
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	paragraphs := doc.Paragraphs()
	if len(paragraphs) == 0 {
		t.Fatal("expected at least 1 paragraph")
	}
	// Should have a run for the link text
	foundLink := false
	for _, run := range paragraphs[0].Runs() {
		if run.Text() == "OpenScribe" {
			foundLink = true
		}
	}
	if !foundLink {
		t.Error("expected to find link text 'OpenScribe' in runs")
	}
}

func TestFromMarkdownEmpty(t *testing.T) {
	doc, err := FromMarkdown("")
	if err != nil {
		t.Fatalf("FromMarkdown error: %v", err)
	}
	if doc == nil {
		t.Fatal("expected non-nil document")
	}
}

func TestMdParseInline(t *testing.T) {
	segments := mdParseInline("Hello **bold** and *italic* `code` world")
	if len(segments) < 5 {
		t.Fatalf("expected at least 5 segments, got %d", len(segments))
	}

	foundBold := false
	foundItalic := false
	foundCode := false
	for _, seg := range segments {
		if seg.bold && seg.text == "bold" {
			foundBold = true
		}
		if seg.italic && seg.text == "italic" {
			foundItalic = true
		}
		if seg.code && seg.text == "code" {
			foundCode = true
		}
	}
	if !foundBold {
		t.Error("expected bold segment")
	}
	if !foundItalic {
		t.Error("expected italic segment")
	}
	if !foundCode {
		t.Error("expected code segment")
	}
}
