package document

import (
	"path/filepath"
	"testing"
)

func TestTableOfContents(t *testing.T) {
	doc := New()

	toc := doc.AddTableOfContents()
	toc.SetTitle("Contents")
	toc.SetMaxLevel(3)

	if toc.Title() != "Contents" {
		t.Error("title mismatch")
	}
	if toc.MaxLevel() != 3 {
		t.Error("max level mismatch")
	}

	doc.AddHeading("Introduction", 1)
	doc.AddText("Intro text")
	doc.AddHeading("Background", 2)
	doc.AddText("Background text")
	doc.AddHeading("Details", 3)
	doc.AddHeading("Sub-details", 4) // should NOT appear in TOC
	doc.AddHeading("Conclusion", 1)

	toc.BuildEntries(doc.Paragraphs())
	entries := toc.Entries()

	if len(entries) != 4 { // Intro, Background, Details, Conclusion (not Sub-details)
		t.Errorf("expected 4 entries, got %d", len(entries))
	}

	if entries[0].Text != "Introduction" || entries[0].Level != 1 {
		t.Errorf("first entry should be Introduction level 1, got %s level %d", entries[0].Text, entries[0].Level)
	}

	if entries[1].Text != "Background" || entries[1].Level != 2 {
		t.Error("second entry should be Background level 2")
	}

	path := filepath.Join(t.TempDir(), "with_toc.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestTOCMaxLevel(t *testing.T) {
	toc := &TableOfContents{maxLevel: 3}

	toc.SetMaxLevel(0) // should clamp to 1
	if toc.MaxLevel() != 1 {
		t.Error("should clamp to 1")
	}

	toc.SetMaxLevel(7) // should clamp to 6
	if toc.MaxLevel() != 6 {
		t.Error("should clamp to 6")
	}
}

func TestTOCShowPageNumbers(t *testing.T) {
	toc := &TableOfContents{}
	toc.SetShowPageNumbers(true)
	if !toc.showPageNumbers {
		t.Error("should show page numbers")
	}
}

func TestTOCEmpty(t *testing.T) {
	doc := New()
	toc := doc.AddTableOfContents()
	doc.AddText("No headings here")

	toc.BuildEntries(doc.Paragraphs())
	if len(toc.Entries()) != 0 {
		t.Error("should have no entries without headings")
	}

	path := filepath.Join(t.TempDir(), "empty_toc.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}
