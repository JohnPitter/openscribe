package presentation

import (
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestNewSlideMaster(t *testing.T) {
	m := NewSlideMaster("Custom")
	if m.Name() != "Custom" {
		t.Error("name mismatch")
	}
}

func TestSlideMasterBackground(t *testing.T) {
	m := NewSlideMaster("Test")
	m.SetBackground(common.Blue)
	if m.Background() == nil || *m.Background() != common.Blue {
		t.Error("background should be blue")
	}
}

func TestSlideMasterFonts(t *testing.T) {
	m := NewSlideMaster("Test")
	titleFont := common.NewFont("Georgia", 40).Bold()
	bodyFont := common.NewFont("Arial", 16)

	m.SetTitleFont(titleFont)
	m.SetBodyFont(bodyFont)

	if m.TitleFont().Family != "Georgia" {
		t.Error("title font should be Georgia")
	}
	if m.BodyFont().Family != "Arial" {
		t.Error("body font should be Arial")
	}
}

func TestAddLayout(t *testing.T) {
	m := NewSlideMaster("Test")
	layout := m.AddLayout("Custom Layout", LayoutTitleContent)
	layout.AddPlaceholder(PlaceholderTitle, common.In(1), common.In(1), common.In(10), common.In(2))
	layout.AddPlaceholder(PlaceholderBody, common.In(1), common.In(3.5), common.In(10), common.In(4))

	if layout.Name() != "Custom Layout" {
		t.Error("name mismatch")
	}
	if layout.LayoutType() != LayoutTitleContent {
		t.Error("layout type mismatch")
	}
	if len(layout.Placeholders()) != 2 {
		t.Errorf("expected 2 placeholders, got %d", len(layout.Placeholders()))
	}
	if len(m.Layouts()) != 1 {
		t.Error("master should have 1 layout")
	}
}

func TestDefaultMaster(t *testing.T) {
	m := DefaultMaster()
	if m.Name() != "Default" {
		t.Error("should be named Default")
	}
	if len(m.Layouts()) != 6 {
		t.Errorf("expected 6 layouts, got %d", len(m.Layouts()))
	}
}

func TestAddSlideFromLayout(t *testing.T) {
	pres := New()
	m := DefaultMaster()
	m.SetBackground(common.DarkGray)
	pres.SetSlideMaster(m)

	// Find title layout
	var titleLayout *SlideLayoutDef
	for _, l := range m.Layouts() {
		if l.LayoutType() == LayoutTitle {
			titleLayout = l
			break
		}
	}
	if titleLayout == nil {
		t.Fatal("should find title layout")
	}

	slide := pres.AddSlideFromLayout(titleLayout)
	if slide == nil {
		t.Fatal("slide should not be nil")
	}
	if slide.Background() == nil || *slide.Background() != common.DarkGray {
		t.Error("should inherit master background")
	}
	if slide.ElementCount() != 2 { // title + subtitle placeholders
		t.Errorf("expected 2 elements from layout, got %d", slide.ElementCount())
	}

	path := filepath.Join(t.TempDir(), "master_layout.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestPresentationWithMasterTheme(t *testing.T) {
	pres := New()
	m := DefaultMaster()
	m.SetTitleFont(common.NewFont("Georgia", 44).Bold().WithColor(common.White))
	m.SetBodyFont(common.NewFont("Georgia", 18).WithColor(common.LightGray))
	m.SetBackground(common.NewColor(30, 30, 60))
	pres.SetSlideMaster(m)

	// Create slides from layouts
	for _, layout := range m.Layouts() {
		pres.AddSlideFromLayout(layout)
	}

	if pres.SlideCount() != 6 {
		t.Errorf("expected 6 slides, got %d", pres.SlideCount())
	}

	path := filepath.Join(t.TempDir(), "all_layouts.pptx")
	if err := pres.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestSlideMasterNil(t *testing.T) {
	pres := New()
	if pres.SlideMaster() != nil {
		t.Error("master should be nil by default")
	}
}
