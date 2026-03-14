// Package presentation provides PPTX presentation creation, reading, and editing.
package presentation

import (
	"fmt"
	"os"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/style"
)

// Presentation represents a PPTX presentation
type Presentation struct {
	pkg    *packaging.Package
	slides []*Slide
	master *SlideMaster
	theme  *style.Theme
	width  common.Measurement
	height common.Measurement
}

// New creates a new empty presentation (16:9 widescreen)
func New() *Presentation {
	theme := style.BasicClean()
	return &Presentation{
		pkg:    packaging.NewPackage(),
		theme:  &theme,
		width:  common.In(13.333), // 16:9 widescreen
		height: common.In(7.5),
	}
}

// NewWithTheme creates a presentation with a specific theme
func NewWithTheme(theme style.Theme) *Presentation {
	p := New()
	p.theme = &theme
	return p
}

// New4x3 creates a presentation with 4:3 aspect ratio
func New4x3() *Presentation {
	p := New()
	p.width = common.In(10)
	p.height = common.In(7.5)
	return p
}

// Theme returns the current theme
func (p *Presentation) Theme() *style.Theme { return p.theme }

// SetTheme sets the presentation theme
func (p *Presentation) SetTheme(theme style.Theme) { p.theme = &theme }

// SlideSize returns the slide dimensions
func (p *Presentation) SlideSize() (common.Measurement, common.Measurement) {
	return p.width, p.height
}

// SetSlideSize sets custom slide dimensions
func (p *Presentation) SetSlideSize(width, height common.Measurement) {
	p.width = width
	p.height = height
}

// AddSlide adds a new blank slide
func (p *Presentation) AddSlide() *Slide {
	s := newSlide(p, len(p.slides)+1)
	p.slides = append(p.slides, s)
	return s
}

// Slide returns a slide by index (0-based)
func (p *Presentation) Slide(index int) *Slide {
	if index < 0 || index >= len(p.slides) {
		return nil
	}
	return p.slides[index]
}

// SlideCount returns the number of slides
func (p *Presentation) SlideCount() int { return len(p.slides) }

// RemoveSlide removes a slide by index
func (p *Presentation) RemoveSlide(index int) error {
	if index < 0 || index >= len(p.slides) {
		return fmt.Errorf("slide index %d out of range", index)
	}
	p.slides = append(p.slides[:index], p.slides[index+1:]...)
	return nil
}

// MoveSlide moves a slide from one position to another
func (p *Presentation) MoveSlide(from, to int) error {
	if from < 0 || from >= len(p.slides) || to < 0 || to >= len(p.slides) {
		return fmt.Errorf("slide index out of range")
	}
	slide := p.slides[from]
	p.slides = append(p.slides[:from], p.slides[from+1:]...)
	newSlides := make([]*Slide, 0, len(p.slides)+1)
	newSlides = append(newSlides, p.slides[:to]...)
	newSlides = append(newSlides, slide)
	newSlides = append(newSlides, p.slides[to:]...)
	p.slides = newSlides
	return nil
}

// Save writes the presentation to a file
func (p *Presentation) Save(path string) error {
	if err := p.build(); err != nil {
		return fmt.Errorf("build presentation: %w", err)
	}
	return p.pkg.Save(path)
}

// SaveToBytes returns the presentation as bytes
func (p *Presentation) SaveToBytes() ([]byte, error) {
	if err := p.build(); err != nil {
		return nil, fmt.Errorf("build presentation: %w", err)
	}
	return p.pkg.ToBytes()
}

// Open reads a PPTX file
func Open(path string) (*Presentation, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return OpenFromBytes(data)
}

// OpenFromBytes reads a PPTX from bytes
func OpenFromBytes(data []byte) (*Presentation, error) {
	pkg, err := packaging.OpenPackageFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("open package: %w", err)
	}
	pres := &Presentation{
		pkg:    pkg,
		width:  common.In(13.333),
		height: common.In(7.5),
	}

	// Count slides by checking for slide files
	for i := 1; ; i++ {
		path := fmt.Sprintf("ppt/slides/slide%d.xml", i)
		if !pkg.HasFile(path) {
			break
		}
		s := newSlide(pres, i)
		pres.slides = append(pres.slides, s)
	}

	return pres, nil
}

// Delete removes a PPTX file from disk
func Delete(path string) error {
	return os.Remove(path)
}
