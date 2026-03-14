package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestAddBulletList(t *testing.T) {
	doc := New()
	doc.AddHeading("Lists Test", 1)

	list := doc.AddList(ListBullet)
	item1 := list.AddItem("First item")
	item1.AddSubItem("Sub-item A")
	item1.AddSubItem("Sub-item B")
	list.AddItem("Second item")
	list.AddItem("Third item")

	if len(doc.Lists()) != 1 {
		t.Errorf("expected 1 list, got %d", len(doc.Lists()))
	}
	if list.Type() != ListBullet {
		t.Errorf("expected ListBullet, got %d", list.Type())
	}
	if len(list.Items()) != 3 {
		t.Errorf("expected 3 items, got %d", len(list.Items()))
	}
	if len(item1.SubItems()) != 2 {
		t.Errorf("expected 2 sub-items, got %d", len(item1.SubItems()))
	}
	if item1.Text() != "First item" {
		t.Errorf("expected 'First item', got '%s'", item1.Text())
	}
	if item1.Level() != 0 {
		t.Errorf("expected level 0, got %d", item1.Level())
	}
	if item1.SubItems()[0].Level() != 1 {
		t.Errorf("expected sub-item level 1, got %d", item1.SubItems()[0].Level())
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bullet_list.docx")
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
	if !pkg.HasFile("word/numbering.xml") {
		t.Error("numbering.xml should exist in package")
	}

	docXML, ok := pkg.GetFile("word/document.xml")
	if !ok {
		t.Fatal("document.xml not found")
	}
	if !containsSubstring(string(docXML), "w:numPr") {
		t.Error("document.xml should contain w:numPr for list items")
	}
}

func TestAddNumberedList(t *testing.T) {
	doc := New()
	list := doc.AddList(ListNumbered)
	list.AddItem("Step 1")
	list.AddItem("Step 2")

	if list.Type() != ListNumbered {
		t.Errorf("expected ListNumbered, got %d", list.Type())
	}
	if list.NumID() != 1 {
		t.Errorf("expected numID 1, got %d", list.NumID())
	}

	path := filepath.Join(t.TempDir(), "numbered_list.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestAddLetteredAndRomanLists(t *testing.T) {
	doc := New()
	l1 := doc.AddList(ListLettered)
	l1.AddItem("Alpha")
	l2 := doc.AddList(ListRoman)
	l2.AddItem("Roman I")

	if len(doc.Lists()) != 2 {
		t.Errorf("expected 2 lists, got %d", len(doc.Lists()))
	}
	if l1.Type() != ListLettered {
		t.Errorf("expected ListLettered")
	}
	if l2.Type() != ListRoman {
		t.Errorf("expected ListRoman")
	}

	path := filepath.Join(t.TempDir(), "multi_lists.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}
