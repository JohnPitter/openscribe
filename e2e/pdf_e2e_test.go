package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/pdf"
	"github.com/JohnPitter/openscribe/style"
)

func TestPdfCreate(t *testing.T) {
	doc := pdf.New()
	p := doc.AddPage()
	p.AddText("Hello PDF World", 72, 72, common.NewFont("Helvetica", 24))

	path := filepath.Join(t.TempDir(), "create.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
	assertFileNotEmpty(t, path)
	assertPDFHeader(t, path)
}

func TestPdfCreateWithAllFeatures(t *testing.T) {
	doc := pdf.New()
	doc.SetMetadata(pdf.Metadata{
		Title:  "Full Feature Test",
		Author: "OpenScribe E2E",
	})

	// Page 1: Text and shapes
	p1 := doc.AddPage()
	p1.SetBackground(common.NewColor(250, 250, 245))

	p1.AddText("OpenScribe PDF Test", 72, 50,
		common.NewFont("Helvetica", 28).Bold())
	p1.AddText("Subtitle text with different styling", 72, 85,
		common.NewFont("Helvetica", 14).WithColor(common.Gray))

	// Lines
	p1.AddLine(72, 100, 500, 100, common.Black, 1)
	p1.AddLine(72, 105, 500, 105, common.Red, 0.5)

	// Rectangles
	p1.AddRectangle(72, 120, 200, 80, common.LightGray, nil)
	stroke := common.Blue
	p1.AddRectangle(300, 120, 200, 80, common.White, &stroke)
	rect := p1.AddRectangle(72, 220, 100, 50, common.Yellow, nil)
	rect.SetCornerRadius(5)

	// Page 2: Table
	p2 := doc.AddPage()
	tbl := p2.AddTable(72, 72, 5, 4)
	tbl.SetCellSize(120, 25)
	tbl.SetHeaderBackground(common.DarkGray)
	tbl.SetBorderColor(common.Gray)
	tbl.SetFont(common.NewFont("Helvetica", 10))

	headers := []string{"Product", "Price", "Qty", "Total"}
	for i, h := range headers {
		tbl.SetCell(0, i, h)
	}
	tbl.SetCell(1, 0, "Widget A")
	tbl.SetCell(1, 1, "$25.00")
	tbl.SetCell(1, 2, "100")
	tbl.SetCell(1, 3, "$2,500")
	tbl.SetCell(2, 0, "Widget B")
	tbl.SetCell(2, 1, "$45.00")
	tbl.SetCell(2, 2, "50")
	tbl.SetCell(2, 3, "$2,250")

	// Page 3: Different page size
	p3 := doc.AddPageWithSize(common.PageLetter, common.NarrowMargins())
	p3.AddText("Letter-sized page with narrow margins", 36, 36,
		common.NewFont("Helvetica", 16))

	path := filepath.Join(t.TempDir(), "full_features.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
	assertPDFHeader(t, path)
}

func TestPdfEdit(t *testing.T) {
	// Create initial
	doc := pdf.New()
	doc.AddPage().AddText("Original", 72, 72, common.NewFont("Helvetica", 12))

	path := filepath.Join(t.TempDir(), "edit.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Open and add pages
	doc2, err := pdf.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}

	p := doc2.AddPage()
	p.AddText("Added page", 72, 72, common.NewFont("Helvetica", 16))
	p.AddRectangle(72, 100, 300, 200, common.LightGray, nil)

	editedPath := filepath.Join(t.TempDir(), "edited.pdf")
	if err := doc2.Save(editedPath); err != nil {
		t.Fatalf("save edited error: %v", err)
	}
	assertFileExists(t, editedPath)
}

func TestPdfEditRemovePage(t *testing.T) {
	doc := pdf.New()
	doc.AddPage().AddText("Page 1", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 2", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 3", 72, 72, common.NewFont("Helvetica", 12))

	if err := doc.RemovePage(1); err != nil {
		t.Fatalf("remove error: %v", err)
	}
	if doc.PageCount() != 2 {
		t.Errorf("expected 2 pages, got %d", doc.PageCount())
	}

	path := filepath.Join(t.TempDir(), "removed.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestPdfDelete(t *testing.T) {
	doc := pdf.New()
	doc.AddPage()

	path := filepath.Join(t.TempDir(), "delete.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)

	if err := pdf.Delete(path); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestPdfMerge(t *testing.T) {
	d1 := pdf.New()
	d1.AddPage().AddText("Doc 1 Page 1", 72, 72, common.NewFont("Helvetica", 16))
	d1.AddPage().AddText("Doc 1 Page 2", 72, 72, common.NewFont("Helvetica", 16))

	d2 := pdf.New()
	d2.AddPage().AddText("Doc 2 Page 1", 72, 72, common.NewFont("Helvetica", 16))

	d3 := pdf.New()
	d3.AddPage().AddText("Doc 3 Page 1", 72, 72, common.NewFont("Helvetica", 16))

	merged := pdf.Merge(d1, d2, d3)
	if merged.PageCount() != 4 {
		t.Errorf("expected 4 pages, got %d", merged.PageCount())
	}

	path := filepath.Join(t.TempDir(), "merged.pdf")
	if err := merged.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
	assertPDFHeader(t, path)
}

func TestPdfWatermark(t *testing.T) {
	doc := pdf.New()
	doc.AddPage().AddText("Page 1", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 2", 72, 72, common.NewFont("Helvetica", 12))

	w := pdf.NewWatermark("CONFIDENTIAL")
	w.SetOpacity(0.4)
	w.SetRotation(-30)
	w.SetFont(common.NewFont("Helvetica", 60).Bold())
	doc.AddWatermark(w)

	path := filepath.Join(t.TempDir(), "watermarked.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)

	// Verify each page has the watermark element
	for i := 0; i < doc.PageCount(); i++ {
		if doc.Page(i).ElementCount() < 2 { // original text + watermark
			t.Errorf("page %d should have watermark", i)
		}
	}
}

func TestPdfSaveToBytes(t *testing.T) {
	doc := pdf.New()
	doc.AddPage().AddText("Bytes", 72, 72, common.NewFont("Helvetica", 12))

	data, err := doc.SaveToBytes()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
	if string(data[:5]) != "%PDF-" {
		t.Error("should start with %PDF-")
	}
}

func TestPdfMetadata(t *testing.T) {
	doc := pdf.New()
	doc.SetMetadata(pdf.Metadata{
		Title:   "Test Document",
		Author:  "Test Author",
		Subject: "Testing",
		Creator: "OpenScribe",
	})

	m := doc.GetMetadata()
	if m.Title != "Test Document" {
		t.Error("title mismatch")
	}
	if m.Author != "Test Author" {
		t.Error("author mismatch")
	}
}

func TestPdfPageSizes(t *testing.T) {
	sizes := []struct {
		name string
		size common.PageSize
	}{
		{"A4", common.PageA4},
		{"A3", common.PageA3},
		{"A5", common.PageA5},
		{"Letter", common.PageLetter},
		{"Legal", common.PageLegal},
		{"Tabloid", common.PageTabloid},
	}

	for _, sz := range sizes {
		t.Run(sz.name, func(t *testing.T) {
			doc := pdf.New()
			doc.AddPageWithSize(sz.size, common.NormalMargins())

			path := filepath.Join(t.TempDir(), sz.name+".pdf")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
			assertFileExists(t, path)
		})
	}
}

func TestPdfSpecialChars(t *testing.T) {
	doc := pdf.New()
	p := doc.AddPage()
	p.AddText("Text with (parentheses) and \\backslash", 72, 72,
		common.NewFont("Helvetica", 12))

	path := filepath.Join(t.TempDir(), "special_chars.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestPdfWithThemes(t *testing.T) {
	themes := style.AllThemes()

	for _, theme := range themes {
		t.Run(theme.Name, func(t *testing.T) {
			doc := pdf.NewWithTheme(theme)
			p := doc.AddPage()
			p.AddText("Themed PDF: "+theme.Name, 72, 72,
				theme.Typography.HeadingFont)
			p.AddText("Body text with theme typography", 72, 110,
				theme.Typography.BodyFont)

			path := filepath.Join(t.TempDir(), "themed.pdf")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error with theme %s: %v", theme.Name, err)
			}
			assertFileExists(t, path)
		})
	}
}
