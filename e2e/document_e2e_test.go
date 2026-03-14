package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/document"
	"github.com/JohnPitter/openscribe/style"
)

func TestDocxCreate(t *testing.T) {
	doc := document.New()
	doc.AddHeading("Test Document", 1)
	doc.AddText("This is a test paragraph with body text.")
	doc.AddHeading("Section 2", 2)
	doc.AddText("More content here.")

	path := filepath.Join(t.TempDir(), "create.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
	assertFileNotEmpty(t, path)
}

func TestDocxCreateWithAllFeatures(t *testing.T) {
	doc := document.New()

	// Headings
	for level := 1; level <= 6; level++ {
		doc.AddHeading("Heading Level", level)
	}

	// Paragraphs with formatting
	p := doc.AddParagraph()
	p.AddText("Bold").SetBold(true)
	p.AddText(" Italic").SetItalic(true)
	p.AddText(" Underline").SetUnderline(true)
	p.AddText(" Colored").SetColor(common.Red)
	p.AddText(" Sized").SetSize(24)
	p.AddText(" FontFamily").SetFontFamily("Georgia")

	// Alignment
	centered := doc.AddParagraph()
	centered.SetAlignment(common.TextAlignCenter)
	centered.AddText("Centered text")

	right := doc.AddParagraph()
	right.SetAlignment(common.TextAlignRight)
	right.AddText("Right aligned")

	justified := doc.AddParagraph()
	justified.SetAlignment(common.TextAlignJustify)
	justified.AddText("Justified text that spans the full width of the page.")

	// Table
	tbl := doc.AddTable(4, 3)
	tbl.Cell(0, 0).SetText("Header 1")
	tbl.Cell(0, 1).SetText("Header 2")
	tbl.Cell(0, 2).SetText("Header 3")
	tbl.Cell(1, 0).SetText("Data A")
	tbl.Cell(1, 1).SetText("Data B")
	tbl.Cell(1, 2).SetText("Data C")
	tbl.Cell(0, 0).SetShading(common.LightGray)
	tbl.Cell(0, 1).SetShading(common.LightGray)
	tbl.Cell(0, 2).SetShading(common.LightGray)

	// Page break
	doc.AddPageBreak()
	doc.AddText("Content on page 2")

	// Section settings
	sect := doc.Section()
	sect.SetPageSize(common.PageLetter)
	sect.SetMargins(common.NarrowMargins())

	path := filepath.Join(t.TempDir(), "full_features.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestDocxEdit(t *testing.T) {
	// Create initial document
	doc := document.New()
	doc.AddHeading("Original Title", 1)
	doc.AddText("Original paragraph")
	doc.AddText("To be removed")

	path := filepath.Join(t.TempDir(), "edit.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Open and edit
	doc2, err := document.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}

	// Add new content
	doc2.AddHeading("New Section", 2)
	doc2.AddText("Added during edit")

	// Add table
	tbl := doc2.AddTable(2, 2)
	tbl.Cell(0, 0).SetText("New A")
	tbl.Cell(0, 1).SetText("New B")

	editedPath := filepath.Join(t.TempDir(), "edited.docx")
	if err := doc2.Save(editedPath); err != nil {
		t.Fatalf("save edited error: %v", err)
	}
	assertFileExists(t, editedPath)

	// Verify edited file is different from original
	orig, _ := os.ReadFile(path)
	edited, _ := os.ReadFile(editedPath)
	if len(edited) <= len(orig) {
		t.Error("edited file should be larger")
	}
}

func TestDocxEditRemoveElements(t *testing.T) {
	doc := document.New()
	doc.AddText("Para 1")
	doc.AddText("Para 2")
	doc.AddText("Para 3")
	doc.AddTable(2, 2)
	doc.AddTable(3, 3)

	// Remove paragraph
	if err := doc.RemoveParagraph(1); err != nil {
		t.Fatalf("remove paragraph error: %v", err)
	}
	if len(doc.Paragraphs()) != 2 {
		t.Errorf("expected 2 paragraphs, got %d", len(doc.Paragraphs()))
	}

	// Remove table
	if err := doc.RemoveTable(0); err != nil {
		t.Fatalf("remove table error: %v", err)
	}
	if len(doc.Tables()) != 1 {
		t.Errorf("expected 1 table, got %d", len(doc.Tables()))
	}

	path := filepath.Join(t.TempDir(), "removed.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestDocxDelete(t *testing.T) {
	doc := document.New()
	doc.AddText("To be deleted")

	path := filepath.Join(t.TempDir(), "delete.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)

	if err := document.Delete(path); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestDocxSaveToBytes(t *testing.T) {
	doc := document.New()
	doc.AddText("Bytes test")

	data, err := doc.SaveToBytes()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}

	doc2, err := document.OpenFromBytes(data)
	if err != nil {
		t.Fatalf("open from bytes error: %v", err)
	}
	if len(doc2.Paragraphs()) == 0 {
		t.Error("should have paragraphs after round-trip")
	}
}

func TestDocxTableAddRemoveRows(t *testing.T) {
	doc := document.New()
	tbl := doc.AddTable(2, 3)

	tbl.AddRow()
	if tbl.RowCount() != 3 {
		t.Errorf("expected 3 rows after add, got %d", tbl.RowCount())
	}

	if err := tbl.RemoveRow(1); err != nil {
		t.Fatalf("remove row error: %v", err)
	}
	if tbl.RowCount() != 2 {
		t.Errorf("expected 2 rows after remove, got %d", tbl.RowCount())
	}
}

func TestDocxRunFont(t *testing.T) {
	doc := document.New()
	p := doc.AddParagraph()
	font := common.NewFont("Georgia", 14).Bold().Italic().WithColor(common.Blue)
	r := p.AddRun()
	r.SetFont(font)
	r.SetText("Styled text")

	path := filepath.Join(t.TempDir(), "font.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestDocxWithThemes(t *testing.T) {
	themes := style.AllThemes()

	for _, theme := range themes {
		t.Run(theme.Name, func(t *testing.T) {
			doc := document.NewWithTheme(theme)
			doc.AddHeading("Document with "+theme.Name, 1)
			doc.AddText("Body text with theme typography.")

			tbl := doc.AddTable(2, 2)
			tbl.Cell(0, 0).SetText("A")
			tbl.Cell(0, 1).SetText("B")

			path := filepath.Join(t.TempDir(), "themed.docx")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error with theme %s: %v", theme.Name, err)
			}
			assertFileExists(t, path)
		})
	}
}

func TestDocxSectionOrientation(t *testing.T) {
	doc := document.New()
	doc.AddText("Landscape document")
	doc.Section().SetOrientation(common.OrientationLandscape)

	path := filepath.Join(t.TempDir(), "landscape.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}
