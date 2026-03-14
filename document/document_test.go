package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/style"
)

func TestNewDocument(t *testing.T) {
	doc := New()
	if doc == nil {
		t.Fatal("document should not be nil")
	}
	if doc.Theme() == nil {
		t.Error("default theme should not be nil")
	}
	if len(doc.Paragraphs()) != 0 {
		t.Error("new document should have no paragraphs")
	}
}

func TestNewWithTheme(t *testing.T) {
	theme := style.PremiumModern()
	doc := NewWithTheme(theme)
	if doc.Theme().Name != "Premium Modern" {
		t.Errorf("expected Premium Modern theme, got %s", doc.Theme().Name)
	}
}

func TestAddParagraph(t *testing.T) {
	doc := New()
	p := doc.AddParagraph()
	p.AddText("Hello World")

	if len(doc.Paragraphs()) != 1 {
		t.Errorf("expected 1 paragraph, got %d", len(doc.Paragraphs()))
	}
	if doc.Paragraphs()[0].Text() != "Hello World" {
		t.Errorf("unexpected text: %s", doc.Paragraphs()[0].Text())
	}
}

func TestAddHeading(t *testing.T) {
	doc := New()
	doc.AddHeading("Title", 1)
	doc.AddHeading("Subtitle", 2)

	if len(doc.Paragraphs()) != 2 {
		t.Errorf("expected 2 paragraphs, got %d", len(doc.Paragraphs()))
	}
	if doc.Paragraphs()[0].Style() != "Heading1" {
		t.Errorf("expected Heading1, got %s", doc.Paragraphs()[0].Style())
	}
}

func TestAddText(t *testing.T) {
	doc := New()
	doc.AddText("Simple text")

	if doc.Paragraphs()[0].Text() != "Simple text" {
		t.Errorf("unexpected text: %s", doc.Paragraphs()[0].Text())
	}
}

func TestAddTable(t *testing.T) {
	doc := New()
	tbl := doc.AddTable(3, 4)

	if tbl.RowCount() != 3 {
		t.Errorf("expected 3 rows, got %d", tbl.RowCount())
	}
	if tbl.ColCount() != 4 {
		t.Errorf("expected 4 cols, got %d", tbl.ColCount())
	}

	tbl.Cell(0, 0).SetText("Header")
	if tbl.Cell(0, 0).Text() != "Header" {
		t.Error("cell text should be Header")
	}
}

func TestRemoveParagraph(t *testing.T) {
	doc := New()
	doc.AddText("First")
	doc.AddText("Second")
	doc.AddText("Third")

	err := doc.RemoveParagraph(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Paragraphs()) != 2 {
		t.Errorf("expected 2 paragraphs, got %d", len(doc.Paragraphs()))
	}
	if doc.Paragraphs()[1].Text() != "Third" {
		t.Error("second paragraph should now be Third")
	}
}

func TestRemoveParagraphOutOfRange(t *testing.T) {
	doc := New()
	err := doc.RemoveParagraph(0)
	if err == nil {
		t.Error("should return error for out of range")
	}
}

func TestRemoveTable(t *testing.T) {
	doc := New()
	doc.AddTable(2, 2)
	doc.AddTable(3, 3)

	err := doc.RemoveTable(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Tables()) != 1 {
		t.Errorf("expected 1 table, got %d", len(doc.Tables()))
	}
}

func TestSection(t *testing.T) {
	doc := New()
	sect := doc.Section()

	sect.SetPageSize(common.PageLetter)
	if sect.PageSize().Width.Inches() != 8.5 {
		t.Error("page width should be 8.5 inches")
	}

	sect.SetMargins(common.NarrowMargins())
	if sect.Margins().Top.Inches() != 0.5 {
		t.Error("top margin should be 0.5 inches")
	}
}

func TestSaveAndOpen(t *testing.T) {
	doc := New()
	doc.AddHeading("Test Document", 1)
	doc.AddText("This is a test paragraph.")

	tbl := doc.AddTable(2, 2)
	tbl.Cell(0, 0).SetText("A1")
	tbl.Cell(0, 1).SetText("B1")
	tbl.Cell(1, 0).SetText("A2")
	tbl.Cell(1, 1).SetText("B2")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.docx")

	err := doc.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Verify file exists
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}

	// Open and verify
	doc2, err := Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	if len(doc2.Paragraphs()) == 0 {
		t.Error("opened document should have paragraphs")
	}
}

func TestSaveToBytes(t *testing.T) {
	doc := New()
	doc.AddText("Bytes test")

	data, err := doc.SaveToBytes()
	if err != nil {
		t.Fatalf("save to bytes error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}

	doc2, err := OpenFromBytes(data)
	if err != nil {
		t.Fatalf("open from bytes error: %v", err)
	}
	if len(doc2.Paragraphs()) == 0 {
		t.Error("should have paragraphs")
	}
}

func TestDelete(t *testing.T) {
	doc := New()
	doc.AddText("To be deleted")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "delete_me.docx")

	err := doc.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	err = Delete(path)
	if err != nil {
		t.Fatalf("delete error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}
