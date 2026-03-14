package pdf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/style"
)

func TestNewDocument(t *testing.T) {
	d := New()
	if d == nil {
		t.Fatal("should not be nil")
	}
	if d.PageCount() != 0 {
		t.Error("new doc should have no pages")
	}
	if d.Theme() == nil {
		t.Error("should have default theme")
	}
}

func TestNewWithTheme(t *testing.T) {
	theme := style.PremiumElegant()
	d := NewWithTheme(theme)
	if d.Theme().Name != "Premium Elegant" {
		t.Errorf("expected Premium Elegant, got %s", d.Theme().Name)
	}
}

func TestAddPage(t *testing.T) {
	d := New()
	p := d.AddPage()
	if p == nil {
		t.Fatal("page should not be nil")
	}
	if d.PageCount() != 1 {
		t.Errorf("expected 1 page, got %d", d.PageCount())
	}
}

func TestAddPageWithSize(t *testing.T) {
	d := New()
	p := d.AddPageWithSize(common.PageLetter, common.NarrowMargins())
	if p.Size().Width.Inches() != 8.5 {
		t.Error("should be letter size")
	}
}

func TestRemovePage(t *testing.T) {
	d := New()
	d.AddPage()
	d.AddPage()
	d.AddPage()

	err := d.RemovePage(1)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if d.PageCount() != 2 {
		t.Errorf("expected 2 pages, got %d", d.PageCount())
	}

	err = d.RemovePage(10)
	if err == nil {
		t.Error("should error on out of range")
	}
}

func TestPageElements(t *testing.T) {
	d := New()
	p := d.AddPage()

	p.AddText("Hello PDF", 72, 72, common.NewFont("Helvetica", 24).Bold())
	p.AddLine(72, 100, 500, 100, common.Black, 1)
	p.AddRectangle(72, 120, 200, 100, common.LightGray, nil)

	if p.ElementCount() != 3 {
		t.Errorf("expected 3 elements, got %d", p.ElementCount())
	}
}

func TestPageTable(t *testing.T) {
	d := New()
	p := d.AddPage()

	tbl := p.AddTable(72, 200, 3, 4)
	tbl.SetCell(0, 0, "Name")
	tbl.SetCell(0, 1, "Age")
	tbl.SetCell(1, 0, "Alice")
	tbl.SetCell(1, 1, "30")

	if tbl.Cell(0, 0) != "Name" {
		t.Errorf("expected Name, got %s", tbl.Cell(0, 0))
	}
	if tbl.Rows() != 3 {
		t.Errorf("expected 3 rows, got %d", tbl.Rows())
	}
	if tbl.Cols() != 4 {
		t.Errorf("expected 4 cols, got %d", tbl.Cols())
	}
}

func TestTextElement(t *testing.T) {
	te := &TextElement{text: "Hello", x: 10, y: 20}
	if te.Text() != "Hello" {
		t.Error("text should be Hello")
	}
	te.SetText("World")
	if te.Text() != "World" {
		t.Error("text should be World")
	}
	te.SetPosition(30, 40)
	if te.x != 30 || te.y != 40 {
		t.Error("position should be updated")
	}
}

func TestRectCornerRadius(t *testing.T) {
	r := &RectElement{}
	r.SetCornerRadius(5)
	if r.cornerRadius != 5 {
		t.Error("corner radius should be 5")
	}
}

func TestTableSetCellSize(t *testing.T) {
	d := New()
	p := d.AddPage()
	tbl := p.AddTable(72, 72, 2, 2)
	tbl.SetCellSize(150, 30)
	if tbl.cellWidth != 150 || tbl.cellHeight != 30 {
		t.Error("cell size should be updated")
	}
}

func TestTableHeaderBg(t *testing.T) {
	d := New()
	p := d.AddPage()
	tbl := p.AddTable(72, 72, 2, 2)
	tbl.SetHeaderBackground(common.Blue)
	if tbl.headerBg == nil || *tbl.headerBg != common.Blue {
		t.Error("header bg should be blue")
	}
}

func TestWatermark(t *testing.T) {
	w := NewWatermark("DRAFT")
	if w.Text() != "DRAFT" {
		t.Error("text should be DRAFT")
	}
	if w.Opacity() != 0.3 {
		t.Errorf("default opacity should be 0.3, got %f", w.Opacity())
	}
	if w.Rotation() != -45 {
		t.Error("default rotation should be -45")
	}

	w.SetOpacity(0.5)
	if w.Opacity() != 0.5 {
		t.Error("opacity should be 0.5")
	}

	w.SetOpacity(-1)
	if w.Opacity() != 0 {
		t.Error("negative opacity should be clamped to 0")
	}

	w.SetOpacity(2)
	if w.Opacity() != 1 {
		t.Error("opacity > 1 should be clamped to 1")
	}
}

func TestDocumentWatermark(t *testing.T) {
	d := New()
	d.AddPage()
	d.AddPage()

	w := NewWatermark("CONFIDENTIAL")
	d.AddWatermark(w)

	// Each page should now have a text element
	for i := 0; i < d.PageCount(); i++ {
		if d.Page(i).ElementCount() == 0 {
			t.Errorf("page %d should have watermark element", i)
		}
	}
}

func TestMetadata(t *testing.T) {
	d := New()
	d.SetMetadata(Metadata{
		Title:  "Test PDF",
		Author: "OpenScribe",
	})
	m := d.GetMetadata()
	if m.Title != "Test PDF" {
		t.Error("title should be Test PDF")
	}
}

func TestMerge(t *testing.T) {
	d1 := New()
	d1.AddPage()
	d1.AddPage()

	d2 := New()
	d2.AddPage()

	merged := Merge(d1, d2)
	if merged.PageCount() != 3 {
		t.Errorf("expected 3 pages, got %d", merged.PageCount())
	}
}

func TestSaveAndOpen(t *testing.T) {
	d := New()
	p := d.AddPage()
	p.AddText("Hello World", 72, 72, common.NewFont("Helvetica", 24))
	p.AddRectangle(72, 100, 200, 50, common.LightGray, nil)

	p2 := d.AddPage()
	tbl := p2.AddTable(72, 72, 2, 3)
	tbl.SetCell(0, 0, "Col 1")
	tbl.SetCell(0, 1, "Col 2")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.pdf")

	err := d.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}

	// Verify it starts with %PDF
	data, _ := os.ReadFile(path)
	if string(data[:5]) != "%PDF-" {
		t.Error("should start with %PDF-")
	}

	// Open
	d2Doc, err := Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	if d2Doc.PageCount() == 0 {
		t.Error("opened doc should have pages")
	}
}

func TestSaveToBytes(t *testing.T) {
	d := New()
	d.AddPage().AddText("Bytes test", 72, 72, common.NewFont("Helvetica", 12))

	data, err := d.SaveToBytes()
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

func TestDelete(t *testing.T) {
	d := New()
	d.AddPage()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "delete_me.pdf")

	if err := d.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	if err := Delete(path); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestPageBackground(t *testing.T) {
	d := New()
	p := d.AddPage()
	p.SetBackground(common.LightGray)

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(data) == 0 {
		t.Error("should produce output")
	}
}

func TestEscapePDF(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello", "Hello"},
		{"Hello (World)", "Hello \\(World\\)"},
		{"Back\\slash", "Back\\\\slash"},
	}
	for _, tt := range tests {
		got := escapePDF(tt.input)
		if got != tt.want {
			t.Errorf("escapePDF(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
