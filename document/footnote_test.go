package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestAddFootnote(t *testing.T) {
	doc := New()
	p := doc.AddParagraph()
	p.AddText("This has a footnote.")

	fnID := doc.AddFootnote("This is the footnote text.")
	p.AddFootnoteRef(fnID)

	if fnID != 1 {
		t.Errorf("expected footnote ID 1, got %d", fnID)
	}
	if len(doc.Footnotes()) != 1 {
		t.Errorf("expected 1 footnote, got %d", len(doc.Footnotes()))
	}
	if doc.Footnotes()[0].Text() != "This is the footnote text." {
		t.Errorf("unexpected footnote text: %s", doc.Footnotes()[0].Text())
	}
	if doc.Footnotes()[0].ID() != 1 {
		t.Errorf("expected ID 1, got %d", doc.Footnotes()[0].ID())
	}
	if len(p.FootnoteRefs()) != 1 {
		t.Errorf("expected 1 footnote ref, got %d", len(p.FootnoteRefs()))
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "footnote.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}

	pkg, err := packaging.OpenPackage(path)
	if err != nil {
		t.Fatalf("open package error: %v", err)
	}
	if !pkg.HasFile("word/footnotes.xml") {
		t.Error("footnotes.xml should exist in package")
	}

	fnXML, ok := pkg.GetFile("word/footnotes.xml")
	if !ok {
		t.Fatal("footnotes.xml not found")
	}
	fnStr := string(fnXML)
	if !containsSubstring(fnStr, "This is the footnote text.") {
		t.Error("footnotes.xml should contain the footnote text")
	}
}

func TestMultipleFootnotes(t *testing.T) {
	doc := New()
	p := doc.AddParagraph()
	p.AddText("Text with multiple footnotes.")

	fn1 := doc.AddFootnote("First footnote.")
	fn2 := doc.AddFootnote("Second footnote.")
	p.AddFootnoteRef(fn1)
	p.AddFootnoteRef(fn2)

	if fn1 != 1 || fn2 != 2 {
		t.Errorf("expected IDs 1 and 2, got %d and %d", fn1, fn2)
	}
	if len(doc.Footnotes()) != 2 {
		t.Errorf("expected 2 footnotes, got %d", len(doc.Footnotes()))
	}

	path := filepath.Join(t.TempDir(), "multi_footnotes.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}
