package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/presentation"
	"github.com/JohnPitter/openscribe/style"
)

func TestPptxCreate(t *testing.T) {
	pres := presentation.New()
	s := pres.AddSlide()
	tb := s.AddTextBox(common.In(1), common.In(1), common.In(8), common.In(2))
	tb.SetText("Hello World", common.NewFont("Arial", 32).Bold())

	path := filepath.Join(t.TempDir(), "create.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
	assertFileNotEmpty(t, path)
}

func TestPptxCreateWithAllFeatures(t *testing.T) {
	pres := presentation.New()

	// Title slide
	s1 := pres.AddSlide()
	s1.SetBackground(common.NewColor(10, 10, 30))
	s1.SetLayout(presentation.LayoutTitle)
	title := s1.AddTextBox(common.In(1), common.In(2), common.In(10), common.In(2))
	p1 := title.AddParagraph()
	p1.AddRun("Product Launch 2026", common.NewFont("Helvetica", 44).Bold().WithColor(common.White))
	p1.SetAlignment(common.TextAlignCenter)

	subtitle := s1.AddTextBox(common.In(2), common.In(4.5), common.In(8), common.In(1))
	sp := subtitle.AddParagraph()
	sp.AddRun("Revolutionizing the Industry", common.NewFont("Helvetica", 20).WithColor(common.LightGray))
	sp.SetAlignment(common.TextAlignCenter)

	s1.SetNotes("Welcome everyone to the product launch presentation")

	// Content slide with shapes
	s2 := pres.AddSlide()
	s2.SetLayout(presentation.LayoutTitleContent)

	// All shape types
	shapes := []presentation.ShapeType{
		presentation.ShapeRectangle,
		presentation.ShapeRoundedRectangle,
		presentation.ShapeCircle,
		presentation.ShapeEllipse,
		presentation.ShapeTriangle,
		presentation.ShapeArrowRight,
		presentation.ShapeStar,
		presentation.ShapeDiamond,
	}
	for i, st := range shapes {
		x := float64(1+i%4) * 3
		y := float64(1+i/4) * 3
		sh := s2.AddShape(st, common.In(x), common.In(y), common.In(2), common.In(2))
		sh.SetFill(common.NewColor(uint8(i*30), uint8(100+i*20), uint8(200-i*20)))
	}

	// Shape with text
	textShape := s2.AddShape(presentation.ShapeRoundedRectangle,
		common.In(1), common.In(6), common.In(4), common.In(1))
	textShape.SetFill(common.Blue)
	textShape.SetText("Click Here", common.NewFont("Arial", 16).Bold().WithColor(common.White))
	textShape.SetLine(common.White, common.Pt(2))

	// Slide with transition
	s3 := pres.AddSlide()
	s3.SetTransition(presentation.NewTransition(presentation.TransitionFade, presentation.TransitionMedium))
	tb := s3.AddTextBox(common.In(2), common.In(3), common.In(8), common.In(2))
	tb.SetText("Fade In Slide", common.NewFont("Arial", 28))
	tb.SetFill(common.LightGray)
	tb.SetBorder(common.DarkGray, common.Pt(2))

	// Text box with multiple paragraphs
	s4 := pres.AddSlide()
	multiTb := s4.AddTextBox(common.In(1), common.In(1), common.In(10), common.In(5))
	for i := 0; i < 5; i++ {
		mp := multiTb.AddParagraph()
		mp.AddRun("Bullet point text paragraph", common.NewFont("Arial", 14))
		mp.SetSpacing(1.5)
	}

	// Custom slide size
	pres.SetSlideSize(common.In(16), common.In(9))

	path := filepath.Join(t.TempDir(), "full_features.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestPptxEdit(t *testing.T) {
	// Create initial
	pres := presentation.New()
	pres.AddSlide().SetNotes("Slide 1")
	pres.AddSlide().SetNotes("Slide 2")

	path := filepath.Join(t.TempDir(), "edit.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Open and edit
	pres2, err := presentation.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}

	// Add slides
	s := pres2.AddSlide()
	s.AddTextBox(common.In(1), common.In(1), common.In(5), common.In(2))

	editedPath := filepath.Join(t.TempDir(), "edited.pptx")
	if err := pres2.Save(editedPath); err != nil {
		t.Fatalf("save edited error: %v", err)
	}
	assertFileExists(t, editedPath)
}

func TestPptxEditRemoveSlide(t *testing.T) {
	pres := presentation.New()
	pres.AddSlide().SetNotes("Keep 1")
	pres.AddSlide().SetNotes("Remove")
	pres.AddSlide().SetNotes("Keep 2")

	if err := pres.RemoveSlide(1); err != nil {
		t.Fatalf("remove error: %v", err)
	}
	if pres.SlideCount() != 2 {
		t.Errorf("expected 2 slides, got %d", pres.SlideCount())
	}

	path := filepath.Join(t.TempDir(), "removed.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestPptxEditMoveSlide(t *testing.T) {
	pres := presentation.New()
	pres.AddSlide().SetNotes("A")
	pres.AddSlide().SetNotes("B")
	pres.AddSlide().SetNotes("C")

	if err := pres.MoveSlide(2, 0); err != nil {
		t.Fatalf("move error: %v", err)
	}
	if pres.Slide(0).Notes() != "C" {
		t.Error("C should now be at position 0")
	}

	path := filepath.Join(t.TempDir(), "moved.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestPptxEditRemoveElement(t *testing.T) {
	pres := presentation.New()
	s := pres.AddSlide()
	s.AddTextBox(common.In(1), common.In(1), common.In(3), common.In(1))
	s.AddShape(presentation.ShapeCircle, common.In(5), common.In(1), common.In(2), common.In(2))

	if err := s.RemoveElement(0); err != nil {
		t.Fatalf("remove element error: %v", err)
	}
	if s.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", s.ElementCount())
	}

	path := filepath.Join(t.TempDir(), "element_removed.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestPptxDelete(t *testing.T) {
	pres := presentation.New()
	pres.AddSlide()

	path := filepath.Join(t.TempDir(), "delete.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)

	if err := presentation.Delete(path); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestPptxSaveToBytes(t *testing.T) {
	pres := presentation.New()
	pres.AddSlide()

	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestPptx4x3(t *testing.T) {
	pres := presentation.New4x3()
	pres.AddSlide()

	w, _ := pres.SlideSize()
	if w.Inches() != 10 {
		t.Errorf("expected 10 inches width, got %f", w.Inches())
	}

	path := filepath.Join(t.TempDir(), "4x3.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestPptxAllTransitions(t *testing.T) {
	pres := presentation.New()

	transitions := []presentation.TransitionType{
		presentation.TransitionNone,
		presentation.TransitionFade,
		presentation.TransitionPush,
		presentation.TransitionWipe,
		presentation.TransitionCut,
		presentation.TransitionDissolve,
		presentation.TransitionZoom,
	}

	for _, tr := range transitions {
		s := pres.AddSlide()
		s.SetTransition(presentation.NewTransition(tr, presentation.TransitionFast))
	}

	path := filepath.Join(t.TempDir(), "transitions.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertFileExists(t, path)
}

func TestPptxWithThemes(t *testing.T) {
	themes := style.AllThemes()

	for _, theme := range themes {
		t.Run(theme.Name, func(t *testing.T) {
			pres := presentation.NewWithTheme(theme)
			s := pres.AddSlide()
			tb := s.AddTextBox(common.In(2), common.In(2), common.In(8), common.In(3))
			tb.SetText("Themed presentation: "+theme.Name,
				common.NewFont("Arial", 28).Bold())

			path := filepath.Join(t.TempDir(), "themed.pptx")
			if err := pres.Save(path); err != nil {
				t.Fatalf("save error with theme %s: %v", theme.Name, err)
			}
			assertFileExists(t, path)
		})
	}
}

func TestPptxRoundTrip(t *testing.T) {
	pres := presentation.New()
	pres.AddSlide()
	pres.AddSlide()
	pres.AddSlide()

	path := filepath.Join(t.TempDir(), "roundtrip.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	pres2, err := presentation.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	if pres2.SlideCount() != 3 {
		t.Errorf("expected 3 slides, got %d", pres2.SlideCount())
	}
}
