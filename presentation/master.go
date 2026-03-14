package presentation

import "github.com/JohnPitter/openscribe/common"

// SlideMaster represents a slide master with default styling
type SlideMaster struct {
	name       string
	background *common.Color
	titleFont  common.Font
	bodyFont   common.Font
	layouts    []*SlideLayoutDef
}

// SlideLayoutDef represents a predefined slide layout
type SlideLayoutDef struct {
	name         string
	layoutType   SlideLayout
	placeholders []Placeholder
}

// Placeholder represents a content placeholder on a slide layout
type Placeholder struct {
	Type   PlaceholderType
	X, Y   common.Measurement
	Width  common.Measurement
	Height common.Measurement
}

// PlaceholderType represents the type of placeholder
type PlaceholderType int

const (
	PlaceholderTitle PlaceholderType = iota
	PlaceholderSubtitle
	PlaceholderBody
	PlaceholderFooter
	PlaceholderDate
	PlaceholderSlideNumber
	PlaceholderImage
)

// NewSlideMaster creates a new slide master
func NewSlideMaster(name string) *SlideMaster {
	return &SlideMaster{
		name:      name,
		titleFont: common.NewFont("Arial", 36).Bold(),
		bodyFont:  common.NewFont("Arial", 18),
	}
}

// SetBackground sets the master background
func (m *SlideMaster) SetBackground(c common.Color) { m.background = &c }

// SetTitleFont sets the default title font
func (m *SlideMaster) SetTitleFont(f common.Font) { m.titleFont = f }

// SetBodyFont sets the default body font
func (m *SlideMaster) SetBodyFont(f common.Font) { m.bodyFont = f }

// Name returns the master name
func (m *SlideMaster) Name() string { return m.name }

// TitleFont returns the title font
func (m *SlideMaster) TitleFont() common.Font { return m.titleFont }

// BodyFont returns the body font
func (m *SlideMaster) BodyFont() common.Font { return m.bodyFont }

// Background returns the background color
func (m *SlideMaster) Background() *common.Color { return m.background }

// AddLayout adds a layout to the master
func (m *SlideMaster) AddLayout(name string, layoutType SlideLayout) *SlideLayoutDef {
	layout := &SlideLayoutDef{
		name:       name,
		layoutType: layoutType,
	}
	m.layouts = append(m.layouts, layout)
	return layout
}

// Layouts returns all layouts
func (m *SlideMaster) Layouts() []*SlideLayoutDef { return m.layouts }

// AddPlaceholder adds a placeholder to the layout
func (l *SlideLayoutDef) AddPlaceholder(phType PlaceholderType, x, y, width, height common.Measurement) {
	l.placeholders = append(l.placeholders, Placeholder{
		Type:   phType,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	})
}

// Name returns the layout name
func (l *SlideLayoutDef) Name() string { return l.name }

// LayoutType returns the layout type
func (l *SlideLayoutDef) LayoutType() SlideLayout { return l.layoutType }

// Placeholders returns all placeholders
func (l *SlideLayoutDef) Placeholders() []Placeholder { return l.placeholders }

// SetSlideMaster sets the presentation's slide master
func (p *Presentation) SetSlideMaster(master *SlideMaster) {
	p.master = master
}

// SlideMaster returns the current slide master
func (p *Presentation) SlideMaster() *SlideMaster {
	return p.master
}

// AddSlideFromLayout adds a slide with elements from a layout
func (p *Presentation) AddSlideFromLayout(layout *SlideLayoutDef) *Slide {
	s := p.AddSlide()
	s.SetLayout(layout.layoutType)

	// Apply master background if set
	if p.master != nil && p.master.background != nil {
		s.SetBackground(*p.master.background)
	}

	// Add placeholders as text boxes
	for _, ph := range layout.placeholders {
		font := common.NewFont("Arial", 18)
		if p.master != nil {
			switch ph.Type {
			case PlaceholderTitle:
				font = p.master.titleFont
			case PlaceholderBody, PlaceholderSubtitle:
				font = p.master.bodyFont
			}
		}
		tb := s.AddTextBox(ph.X, ph.Y, ph.Width, ph.Height)
		placeholder := ""
		switch ph.Type {
		case PlaceholderTitle:
			placeholder = "Click to add title"
		case PlaceholderSubtitle:
			placeholder = "Click to add subtitle"
		case PlaceholderBody:
			placeholder = "Click to add text"
		case PlaceholderFooter:
			placeholder = ""
		}
		if placeholder != "" {
			tb.SetText(placeholder, font)
		}
	}

	return s
}

// DefaultMaster creates a standard slide master with common layouts
func DefaultMaster() *SlideMaster {
	m := NewSlideMaster("Default")

	// Title slide layout
	titleLayout := m.AddLayout("Title Slide", LayoutTitle)
	titleLayout.AddPlaceholder(PlaceholderTitle, common.In(1.5), common.In(2), common.In(10), common.In(2))
	titleLayout.AddPlaceholder(PlaceholderSubtitle, common.In(2.5), common.In(4.5), common.In(8), common.In(1.5))

	// Title + Content layout
	contentLayout := m.AddLayout("Title and Content", LayoutTitleContent)
	contentLayout.AddPlaceholder(PlaceholderTitle, common.In(0.5), common.In(0.3), common.In(12), common.In(1.2))
	contentLayout.AddPlaceholder(PlaceholderBody, common.In(0.5), common.In(1.8), common.In(12), common.In(5))

	// Two Column layout
	twoColLayout := m.AddLayout("Two Column", LayoutTwoColumn)
	twoColLayout.AddPlaceholder(PlaceholderTitle, common.In(0.5), common.In(0.3), common.In(12), common.In(1.2))
	twoColLayout.AddPlaceholder(PlaceholderBody, common.In(0.5), common.In(1.8), common.In(5.5), common.In(5))
	twoColLayout.AddPlaceholder(PlaceholderBody, common.In(6.5), common.In(1.8), common.In(5.5), common.In(5))

	// Blank layout
	m.AddLayout("Blank", LayoutBlank)

	// Section layout
	sectionLayout := m.AddLayout("Section", LayoutSection)
	sectionLayout.AddPlaceholder(PlaceholderTitle, common.In(1), common.In(2.5), common.In(11), common.In(2))

	// Title Only layout
	titleOnlyLayout := m.AddLayout("Title Only", LayoutTitleOnly)
	titleOnlyLayout.AddPlaceholder(PlaceholderTitle, common.In(0.5), common.In(0.3), common.In(12), common.In(1.2))

	return m
}
