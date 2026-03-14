package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
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

func TestHeaderFooter(t *testing.T) {
	doc := New()
	doc.AddText("Content")

	h := doc.Header()
	h.SetLeft("Company Name")
	h.SetCenter("Confidential")
	h.SetRight("Page 1")

	f := doc.Footer()
	f.SetCenter("© 2026 Company")

	if h.Left() != "Company Name" {
		t.Error("header left mismatch")
	}
	if h.Center() != "Confidential" {
		t.Error("header center mismatch")
	}
	if f.Center() != "© 2026 Company" {
		t.Error("footer center mismatch")
	}
	if f.IsEmpty() {
		t.Error("footer should not be empty")
	}

	empty := &HeaderFooter{}
	if !empty.IsEmpty() {
		t.Error("empty header should be empty")
	}

	path := filepath.Join(t.TempDir(), "header_footer.docx")
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

	// Verify header and footer files exist in the package
	pkg, err := packaging.OpenPackage(path)
	if err != nil {
		t.Fatalf("open package error: %v", err)
	}
	if !pkg.HasFile("word/header1.xml") {
		t.Error("header1.xml should exist in package")
	}
	if !pkg.HasFile("word/footer1.xml") {
		t.Error("footer1.xml should exist in package")
	}
}

func TestImageEmbedding(t *testing.T) {
	doc := New()
	doc.AddHeading("Document with Image", 1)

	// Create test image data (minimal PNG header bytes)
	imgData := &common.ImageData{
		Data:   []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
		Format: common.ImageFormatPNG,
	}
	ref := doc.AddImage(imgData, common.In(4), common.In(3))

	if ref.ID() == "" {
		t.Error("image should have an ID")
	}

	doc.AddText("Text after image")

	path := filepath.Join(t.TempDir(), "with_image.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Verify the image file is in the package
	pkg, err := packaging.OpenPackage(path)
	if err != nil {
		t.Fatalf("open package error: %v", err)
	}

	// Check that media file exists
	found := false
	for name := range pkg.Files {
		if len(name) > 11 && name[:11] == "word/media/" {
			found = true
			break
		}
	}
	if !found {
		t.Error("image media file not found in package (expected word/media/img1.png)")
	}

	// Verify the document.xml contains drawing elements
	docXML, ok := pkg.GetFile("word/document.xml")
	if !ok {
		t.Fatal("document.xml not found in package")
	}
	docStr := string(docXML)
	if !containsSubstring(docStr, "w:drawing") {
		t.Error("document.xml should contain w:drawing element for embedded image")
	}
	if !containsSubstring(docStr, "r:embed") {
		t.Error("document.xml should contain r:embed attribute for image relationship")
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
