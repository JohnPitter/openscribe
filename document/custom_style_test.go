package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestAddCustomStyle(t *testing.T) {
	doc := New()

	cs := doc.AddStyle("MyCustomStyle", "Heading1")
	font := common.NewFont("Georgia", 14).Bold().WithColor(common.NewColor(128, 0, 0))
	cs.SetFont(font)
	cs.SetAlignment(common.TextAlignCenter)
	cs.SetSpacing(12, 6, 1.5)
	cs.SetIndent(36, 0, 18)

	if cs.Name() != "MyCustomStyle" {
		t.Errorf("expected name 'MyCustomStyle', got '%s'", cs.Name())
	}
	if cs.BasedOn() != "Heading1" {
		t.Errorf("expected basedOn 'Heading1', got '%s'", cs.BasedOn())
	}
	if len(doc.CustomStyles()) != 1 {
		t.Errorf("expected 1 custom style, got %d", len(doc.CustomStyles()))
	}

	// Use the custom style in a paragraph
	p := doc.AddParagraph()
	p.SetStyle("MyCustomStyle")
	p.AddText("Styled text")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "custom_style.docx")
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

	stylesXML, ok := pkg.GetFile("word/styles.xml")
	if !ok {
		t.Fatal("styles.xml not found")
	}
	stylesStr := string(stylesXML)
	if !containsSubstring(stylesStr, "MyCustomStyle") {
		t.Error("styles.xml should contain custom style 'MyCustomStyle'")
	}
	if !containsSubstring(stylesStr, "Georgia") {
		t.Error("styles.xml should contain font family 'Georgia'")
	}
}

func TestCustomStyleChaining(t *testing.T) {
	doc := New()
	cs := doc.AddStyle("ChainTest", "")

	result := cs.SetFont(common.NewFont("Arial", 12)).
		SetAlignment(common.TextAlignRight).
		SetSpacing(10, 10, 1.0).
		SetIndent(20, 20, 10)

	if result != cs {
		t.Error("chained methods should return the same CustomStyle")
	}
}

func TestMultipleCustomStyles(t *testing.T) {
	doc := New()
	doc.AddStyle("StyleA", "")
	doc.AddStyle("StyleB", "StyleA")
	doc.AddStyle("StyleC", "Heading2")

	if len(doc.CustomStyles()) != 3 {
		t.Errorf("expected 3 custom styles, got %d", len(doc.CustomStyles()))
	}

	path := filepath.Join(t.TempDir(), "multi_styles.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}
