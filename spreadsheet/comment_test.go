package spreadsheet

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCellComment(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Comments")
	s.SetValue(1, 1, "Hello")

	cell := s.Cell(1, 1)
	cell.SetComment("Author1", "This is a comment")

	author, text := cell.Comment()
	if author != "Author1" {
		t.Errorf("expected Author1, got %s", author)
	}
	if text != "This is a comment" {
		t.Errorf("expected 'This is a comment', got %s", text)
	}
}

func TestCellNoComment(t *testing.T) {
	wb := New()
	s := wb.AddSheet("NoComments")
	cell := s.Cell(1, 1)

	author, text := cell.Comment()
	if author != "" || text != "" {
		t.Error("should return empty strings for cell with no comment")
	}
}

func TestCommentSaveToFile(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Comments")
	s.SetValue(1, 1, "Data")
	s.Cell(1, 1).SetComment("John", "Review this value")
	s.SetValue(2, 1, "More Data")
	s.Cell(2, 1).SetComment("Jane", "Check calculation")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "comments.xlsx")
	err := wb.Save(path)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestCommentReplacement(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Replace")
	cell := s.Cell(1, 1)
	cell.SetComment("Author1", "First comment")
	cell.SetComment("Author2", "Updated comment")

	author, text := cell.Comment()
	if author != "Author2" {
		t.Errorf("expected Author2, got %s", author)
	}
	if text != "Updated comment" {
		t.Errorf("expected 'Updated comment', got %s", text)
	}

	// Should only have one comment for the cell, not two
	if len(s.comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(s.comments))
	}
}

func TestBuildCommentsXML(t *testing.T) {
	comments := []*Comment{
		{row: 1, col: 1, author: "Alice", text: "Hello"},
		{row: 2, col: 3, author: "Bob", text: "World"},
	}
	xml := buildCommentsXML(comments)
	if !strings.Contains(xml, "<authors>") {
		t.Error("should contain authors element")
	}
	if !strings.Contains(xml, "Alice") {
		t.Error("should contain author Alice")
	}
	if !strings.Contains(xml, "Hello") {
		t.Error("should contain comment text Hello")
	}
	if !strings.Contains(xml, `ref="A1"`) {
		t.Error("should contain cell ref A1")
	}
}

func TestBuildCommentsXMLEmpty(t *testing.T) {
	xml := buildCommentsXML(nil)
	if xml != "" {
		t.Error("empty comments should produce empty string")
	}
}

func TestBuildVMLDrawingXML(t *testing.T) {
	comments := []*Comment{
		{row: 1, col: 1, author: "Test", text: "Comment"},
	}
	xml := buildVMLDrawingXML(comments)
	if !strings.Contains(xml, "v:shape") {
		t.Error("should contain v:shape element")
	}
	if !strings.Contains(xml, "ClientData") {
		t.Error("should contain ClientData element")
	}
}
