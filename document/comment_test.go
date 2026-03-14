package document

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestAddComment(t *testing.T) {
	doc := New()
	p := doc.AddParagraph()
	run := p.AddText("This text has a comment.")

	comment := doc.AddComment("John", "This needs revision.")
	run.SetComment(comment)

	if comment.ID() != 1 {
		t.Errorf("expected comment ID 1, got %d", comment.ID())
	}
	if comment.Author() != "John" {
		t.Errorf("expected author 'John', got '%s'", comment.Author())
	}
	if comment.Text() != "This needs revision." {
		t.Errorf("unexpected comment text: %s", comment.Text())
	}
	if comment.Date() == "" {
		t.Error("comment date should not be empty")
	}
	if len(doc.Comments()) != 1 {
		t.Errorf("expected 1 comment, got %d", len(doc.Comments()))
	}
	if run.Comment() != comment {
		t.Error("run should reference the comment")
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "comment.docx")
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
	if !pkg.HasFile("word/comments.xml") {
		t.Error("comments.xml should exist in package")
	}

	cmXML, ok := pkg.GetFile("word/comments.xml")
	if !ok {
		t.Fatal("comments.xml not found")
	}
	cmStr := string(cmXML)
	if !containsSubstring(cmStr, "This needs revision.") {
		t.Error("comments.xml should contain the comment text")
	}
	if !containsSubstring(cmStr, "John") {
		t.Error("comments.xml should contain the author name")
	}

	// Check document.xml for comment markers
	docXML, ok := pkg.GetFile("word/document.xml")
	if !ok {
		t.Fatal("document.xml not found")
	}
	docStr := string(docXML)
	if !containsSubstring(docStr, "commentRangeStart") {
		t.Error("document.xml should contain commentRangeStart")
	}
	if !containsSubstring(docStr, "commentReference") {
		t.Error("document.xml should contain commentReference")
	}
}

func TestMultipleComments(t *testing.T) {
	doc := New()
	p := doc.AddParagraph()
	r1 := p.AddText("First commented text.")
	r2 := p.AddText("Second commented text.")

	c1 := doc.AddComment("Alice", "Comment one.")
	c2 := doc.AddComment("Bob", "Comment two.")
	r1.SetComment(c1)
	r2.SetComment(c2)

	if len(doc.Comments()) != 2 {
		t.Errorf("expected 2 comments, got %d", len(doc.Comments()))
	}
	if c1.ID() != 1 || c2.ID() != 2 {
		t.Errorf("expected IDs 1 and 2, got %d and %d", c1.ID(), c2.ID())
	}

	path := filepath.Join(t.TempDir(), "multi_comments.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}
