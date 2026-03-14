package presentation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/style"
)

func TestNewPresentation(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("should not be nil")
	}
	if p.SlideCount() != 0 {
		t.Error("new presentation should have no slides")
	}
	w, h := p.SlideSize()
	if w.Inches() < 13 || h.Inches() < 7 {
		t.Error("default should be widescreen 16:9")
	}
}

func TestNewWithTheme(t *testing.T) {
	theme := style.LuxuryAgency()
	p := NewWithTheme(theme)
	if p.Theme().Name != "Luxury Agency" {
		t.Errorf("expected Luxury Agency, got %s", p.Theme().Name)
	}
}

func TestNew4x3(t *testing.T) {
	p := New4x3()
	w, _ := p.SlideSize()
	if w.Inches() != 10 {
		t.Errorf("4:3 width should be 10 inches, got %f", w.Inches())
	}
}

func TestAddSlide(t *testing.T) {
	p := New()
	s := p.AddSlide()
	if s == nil {
		t.Fatal("slide should not be nil")
	}
	if p.SlideCount() != 1 {
		t.Errorf("expected 1 slide, got %d", p.SlideCount())
	}
}

func TestRemoveSlide(t *testing.T) {
	p := New()
	p.AddSlide()
	p.AddSlide()
	p.AddSlide()

	err := p.RemoveSlide(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.SlideCount() != 2 {
		t.Errorf("expected 2 slides, got %d", p.SlideCount())
	}

	err = p.RemoveSlide(10)
	if err == nil {
		t.Error("should error on out of range")
	}
}

func TestMoveSlide(t *testing.T) {
	p := New()
	p.AddSlide().SetNotes("Slide 1")
	p.AddSlide().SetNotes("Slide 2")
	p.AddSlide().SetNotes("Slide 3")

	err := p.MoveSlide(2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Slide(0).Notes() != "Slide 3" {
		t.Error("slide 3 should now be at position 0")
	}
}

func TestSlideTextBox(t *testing.T) {
	p := New()
	s := p.AddSlide()

	tb := s.AddTextBox(common.In(1), common.In(1), common.In(5), common.In(2))
	font := common.NewFont("Arial", 24).Bold()
	tb.SetText("Hello World", font)

	if tb.Text() != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", tb.Text())
	}
	if s.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", s.ElementCount())
	}
}

func TestSlideShape(t *testing.T) {
	p := New()
	s := p.AddSlide()

	sh := s.AddShape(ShapeCircle, common.In(2), common.In(2), common.In(3), common.In(3))
	sh.SetFill(common.Red)
	sh.SetRotation(45)

	if sh.Type() != ShapeCircle {
		t.Error("should be circle")
	}
	if s.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", s.ElementCount())
	}
}

func TestSlideBackground(t *testing.T) {
	p := New()
	s := p.AddSlide()
	s.SetBackground(common.Blue)

	if s.Background() == nil || *s.Background() != common.Blue {
		t.Error("background should be blue")
	}
}

func TestSlideLayout(t *testing.T) {
	p := New()
	s := p.AddSlide()
	s.SetLayout(LayoutTitleContent)
	if s.Layout() != LayoutTitleContent {
		t.Error("should be TitleContent layout")
	}
}

func TestTransition(t *testing.T) {
	tr := NewTransition(TransitionFade, TransitionMedium)
	if tr.Duration != 0.5 {
		t.Errorf("expected 0.5s, got %f", tr.Duration)
	}
}

func TestRemoveElement(t *testing.T) {
	p := New()
	s := p.AddSlide()
	s.AddTextBox(common.In(1), common.In(1), common.In(3), common.In(1))
	s.AddShape(ShapeRectangle, common.In(1), common.In(3), common.In(2), common.In(2))

	err := s.RemoveElement(0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if s.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", s.ElementCount())
	}
}

func TestSaveAndOpen(t *testing.T) {
	p := New()
	s := p.AddSlide()
	s.SetBackground(common.White)
	tb := s.AddTextBox(common.In(1), common.In(1), common.In(8), common.In(2))
	tb.SetText("Test Presentation", common.NewFont("Arial", 32).Bold())

	s.AddShape(ShapeRectangle, common.In(2), common.In(4), common.In(4), common.In(2))

	p.AddSlide().SetNotes("Blank slide")

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.pptx")

	err := p.Save(path)
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

	// Open and verify
	p2, err := Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	if p2.SlideCount() != 2 {
		t.Errorf("expected 2 slides, got %d", p2.SlideCount())
	}
}

func TestSaveToBytes(t *testing.T) {
	p := New()
	p.AddSlide()

	data, err := p.SaveToBytes()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(data) == 0 {
		t.Error("bytes should not be empty")
	}
}

func TestDelete(t *testing.T) {
	p := New()
	p.AddSlide()

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "delete_me.pptx")

	if err := p.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	if err := Delete(path); err != nil {
		t.Fatalf("delete error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}
}

func TestTransitionSerialization(t *testing.T) {
	pres := New()

	s := pres.AddSlide()
	s.SetTransition(NewTransition(TransitionFade, TransitionFast))

	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Verify the transition XML is in the output
	pkg, pkgErr := packaging.OpenPackageFromBytes(data)
	if pkgErr != nil {
		t.Fatalf("open package error: %v", pkgErr)
	}

	slideXML, ok := pkg.GetFile("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide1.xml should exist")
	}

	if !strings.Contains(string(slideXML), "transition") {
		t.Error("slide XML should contain transition element")
	}
	if !strings.Contains(string(slideXML), "fade") {
		t.Error("slide XML should contain fade transition type")
	}
	if !strings.Contains(string(slideXML), "fast") {
		t.Error("slide XML should contain fast speed attribute")
	}
}

func TestTransitionTypes(t *testing.T) {
	types := []struct {
		tt   TransitionType
		name string
	}{
		{TransitionFade, "fade"},
		{TransitionPush, "push"},
		{TransitionWipe, "wipe"},
		{TransitionSplit, "split"},
		{TransitionCover, "cover"},
		{TransitionCut, "cut"},
		{TransitionDissolve, "dissolve"},
	}

	for _, tc := range types {
		t.Run(tc.name, func(t *testing.T) {
			pres := New()
			s := pres.AddSlide()
			s.SetTransition(NewTransition(tc.tt, TransitionMedium))

			data, err := pres.SaveToBytes()
			if err != nil {
				t.Fatalf("save error: %v", err)
			}

			pkg, _ := packaging.OpenPackageFromBytes(data)
			slideXML, _ := pkg.GetFile("ppt/slides/slide1.xml")

			if !strings.Contains(string(slideXML), tc.name) {
				t.Errorf("slide XML should contain %s transition type", tc.name)
			}
		})
	}
}

func TestNoTransition(t *testing.T) {
	pres := New()
	pres.AddSlide() // no transition set

	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackageFromBytes(data)
	slideXML, _ := pkg.GetFile("ppt/slides/slide1.xml")

	if strings.Contains(string(slideXML), "transition") {
		t.Error("slide XML should NOT contain transition when none is set")
	}
}

func TestSlideImageElement(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	// Create fake image data
	imgData := &common.ImageData{
		Data:   []byte{0x89, 0x50, 0x4E, 0x47}, // PNG magic bytes
		Format: common.ImageFormatPNG,
	}

	img := s.AddImage(imgData, common.In(1), common.In(1), common.In(4), common.In(3))
	if img == nil {
		t.Fatal("image should not be nil")
	}
	if s.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", s.ElementCount())
	}

	path := filepath.Join(t.TempDir(), "img_pres.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	// Verify media file and rels exist in the package
	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("save to bytes error: %v", err)
	}
	pkg, _ := packaging.OpenPackageFromBytes(data)
	if !pkg.HasFile("ppt/media/slide1_img1.png") {
		t.Error("media file should exist in package")
	}
	if !pkg.HasFile("ppt/slides/_rels/slide1.xml.rels") {
		t.Error("slide rels file should exist in package")
	}

	slideXML, _ := pkg.GetFile("ppt/slides/slide1.xml")
	if !strings.Contains(string(slideXML), "p:pic") {
		t.Error("slide XML should contain p:pic element")
	}
}

func TestSlideChart(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	chart := s.AddChart(PresentationChartBar,
		common.In(1), common.In(1), common.In(8), common.In(5))
	chart.SetTitle("Sales Data")
	chart.SetCategories([]string{"Q1", "Q2", "Q3", "Q4"})
	chart.AddSeries("Revenue", []float64{100, 200, 150, 300}, common.Blue)

	if chart.Title() != "Sales Data" {
		t.Error("title mismatch")
	}
	if chart.Type() != PresentationChartBar {
		t.Error("type should be bar")
	}
	if len(chart.Series()) != 1 {
		t.Error("should have 1 series")
	}
	if len(chart.Categories()) != 4 {
		t.Error("should have 4 categories")
	}
	if s.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", s.ElementCount())
	}

	path := filepath.Join(t.TempDir(), "chart_pres.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestSlideChartMultipleSeries(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	chart := s.AddChart(PresentationChartLine,
		common.In(1), common.In(1), common.In(8), common.In(5))
	chart.SetTitle("Comparison")
	chart.AddSeries("Series A", []float64{10, 20, 30}, common.Red)
	chart.AddSeries("Series B", []float64{15, 25, 35}, common.Blue)

	if len(chart.Series()) != 2 {
		t.Errorf("expected 2 series, got %d", len(chart.Series()))
	}
}
