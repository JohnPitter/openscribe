package packaging

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewPackage(t *testing.T) {
	pkg := NewPackage()
	if len(pkg.Files) != 0 {
		t.Error("new package should have no files")
	}
}

func TestPackageAddGetRemove(t *testing.T) {
	pkg := NewPackage()

	pkg.AddFile("test.xml", []byte("<test/>"))

	if !pkg.HasFile("test.xml") {
		t.Error("should have test.xml")
	}

	data, ok := pkg.GetFile("test.xml")
	if !ok || string(data) != "<test/>" {
		t.Error("should return correct data")
	}

	pkg.RemoveFile("test.xml")
	if pkg.HasFile("test.xml") {
		t.Error("should not have test.xml after removal")
	}
}

func TestPackageSaveAndOpen(t *testing.T) {
	pkg := NewPackage()
	pkg.AddFile("hello.txt", []byte("world"))
	pkg.AddFile("dir/file.xml", []byte("<root/>"))

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.zip")

	err := pkg.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file should exist: %v", err)
	}

	// Open and verify
	pkg2, err := OpenPackage(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}

	if len(pkg2.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(pkg2.Files))
	}

	data, ok := pkg2.GetFile("hello.txt")
	if !ok || string(data) != "world" {
		t.Error("should have correct hello.txt content")
	}
}

func TestPackageToBytes(t *testing.T) {
	pkg := NewPackage()
	pkg.AddFile("test.txt", []byte("data"))

	data, err := pkg.ToBytes()
	if err != nil {
		t.Fatalf("to bytes error: %v", err)
	}

	pkg2, err := OpenPackageFromBytes(data)
	if err != nil {
		t.Fatalf("open from bytes error: %v", err)
	}

	if !pkg2.HasFile("test.txt") {
		t.Error("should have test.txt")
	}
}

func TestRelationships(t *testing.T) {
	rels := NewRelationships()
	id := rels.Add(RelTypeOfficeDocument, "word/document.xml")
	if id != "rId1" {
		t.Errorf("expected rId1, got %s", id)
	}

	id2 := rels.Add(RelTypeStylesheet, "word/styles.xml")
	if id2 != "rId2" {
		t.Errorf("expected rId2, got %s", id2)
	}

	data, err := rels.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	rels2, err := UnmarshalRelationships(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(rels2.Relationships) != 2 {
		t.Errorf("expected 2 relationships, got %d", len(rels2.Relationships))
	}
}

func TestContentTypes(t *testing.T) {
	ct := NewContentTypes()
	ct.AddOverride("/word/document.xml", ContentTypeDocx)

	if len(ct.Defaults) != 4 {
		t.Errorf("expected 4 defaults, got %d", len(ct.Defaults))
	}
	if len(ct.Overrides) != 1 {
		t.Errorf("expected 1 override, got %d", len(ct.Overrides))
	}

	data, err := ct.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if len(data) == 0 {
		t.Error("marshaled data should not be empty")
	}
}
