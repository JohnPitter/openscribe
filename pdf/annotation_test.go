package pdf

import (
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestAddHighlight(t *testing.T) {
	d := New()
	p := d.AddPage()
	a := p.AddHighlight(72, 100, 300, 120, common.Yellow)

	if a.Type() != AnnotHighlight {
		t.Error("expected AnnotHighlight type")
	}
	if a.Color() != common.Yellow {
		t.Error("expected yellow color")
	}
	if p.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", p.ElementCount())
	}
}

func TestAddStickyNote(t *testing.T) {
	d := New()
	p := d.AddPage()
	a := p.AddStickyNote(100, 200, "Review this section", common.Orange)

	if a.Type() != AnnotStickyNote {
		t.Error("expected AnnotStickyNote type")
	}
	if a.Text() != "Review this section" {
		t.Errorf("expected 'Review this section', got %q", a.Text())
	}
	if a.Color() != common.Orange {
		t.Error("expected orange color")
	}
}

func TestAddFreeText(t *testing.T) {
	d := New()
	p := d.AddPage()
	font := common.NewFont("Helvetica", 14).WithColor(common.Red)
	a := p.AddFreeText(72, 300, 200, 50, "Important note", font)

	if a.Type() != AnnotFreeText {
		t.Error("expected AnnotFreeText type")
	}
	if a.Text() != "Important note" {
		t.Errorf("expected 'Important note', got %q", a.Text())
	}
}

func TestAnnotationAuthorSubject(t *testing.T) {
	d := New()
	p := d.AddPage()
	a := p.AddHighlight(72, 100, 300, 120, common.Yellow)

	a.SetAuthor("Jane Doe")
	a.SetSubject("Code Review")

	if a.Author() != "Jane Doe" {
		t.Errorf("expected author 'Jane Doe', got %q", a.Author())
	}
	if a.Subject() != "Code Review" {
		t.Errorf("expected subject 'Code Review', got %q", a.Subject())
	}
}

func TestAnnotationsBuildPDF(t *testing.T) {
	d := New()
	p := d.AddPage()

	h := p.AddHighlight(72, 100, 300, 120, common.Yellow)
	h.SetAuthor("Author1")

	p.AddStickyNote(100, 200, "Note text", common.Green)

	font := common.NewFont("Helvetica", 12)
	p.AddFreeText(72, 400, 200, 50, "Free text content", font)

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "/Subtype /Highlight") {
		t.Error("should contain highlight annotation")
	}
	if !strings.Contains(content, "/Subtype /Text") {
		t.Error("should contain sticky note annotation (Text subtype)")
	}
	if !strings.Contains(content, "/Subtype /FreeText") {
		t.Error("should contain free text annotation")
	}
	if !strings.Contains(content, "/T (Author1)") {
		t.Error("should contain author")
	}
	if !strings.Contains(content, "/Contents (Note text)") {
		t.Error("should contain sticky note contents")
	}
	if !strings.Contains(content, "/Annots") {
		t.Error("page should have /Annots array")
	}
}

func TestAnnotationsNoAnnotsWithout(t *testing.T) {
	d := New()
	p := d.AddPage()
	p.AddText("No annotations", 72, 72, common.NewFont("Helvetica", 12))

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	content := string(data)
	if strings.Contains(content, "/Annots") {
		t.Error("should NOT contain /Annots when no annotations exist")
	}
}

func TestHighlightWithSubject(t *testing.T) {
	d := New()
	p := d.AddPage()
	a := p.AddHighlight(72, 100, 300, 120, common.Yellow)
	a.SetSubject("Important")

	data, err := d.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "/Subj (Important)") {
		t.Error("should contain /Subj in annotation")
	}
}
