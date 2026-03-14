package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestAddHyperlink(t *testing.T) {
	doc := New()
	p := doc.AddParagraph()
	p.AddText("Visit ")
	h := p.AddHyperlink("OpenScribe", "https://github.com/JohnPitter/openscribe")

	if h.Text() != "OpenScribe" {
		t.Errorf("expected 'OpenScribe', got '%s'", h.Text())
	}
	if h.URL() != "https://github.com/JohnPitter/openscribe" {
		t.Errorf("unexpected URL: %s", h.URL())
	}
	if len(p.Hyperlinks()) != 1 {
		t.Errorf("expected 1 hyperlink, got %d", len(p.Hyperlinks()))
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "hyperlink.docx")
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

	// Check that document.xml contains hyperlink
	docXML, ok := pkg.GetFile("word/document.xml")
	if !ok {
		t.Fatal("document.xml not found")
	}
	docStr := string(docXML)
	if !containsSubstring(docStr, "w:hyperlink") {
		t.Error("document.xml should contain w:hyperlink element")
	}

	// Check that rels contain external hyperlink relationship
	relsData, ok := pkg.GetFile("word/_rels/document.xml.rels")
	if !ok {
		t.Fatal("document.xml.rels not found")
	}
	relsStr := string(relsData)
	if !containsSubstring(relsStr, "External") {
		t.Error("rels should contain External TargetMode for hyperlink")
	}
}

func TestMultipleHyperlinks(t *testing.T) {
	doc := New()
	p := doc.AddParagraph()
	p.AddHyperlink("Google", "https://google.com")
	p.AddHyperlink("GitHub", "https://github.com")

	if len(p.Hyperlinks()) != 2 {
		t.Errorf("expected 2 hyperlinks, got %d", len(p.Hyperlinks()))
	}

	path := filepath.Join(t.TempDir(), "multi_hyperlinks.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}
