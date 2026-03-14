// Package pdf provides PDF document creation and manipulation.
package pdf

import (
	"fmt"
	"os"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/style"
)

// Document represents a PDF document
type Document struct {
	pages    []*Page
	theme    *style.Theme
	metadata Metadata
}

// Metadata holds document metadata
type Metadata struct {
	Title   string
	Author  string
	Subject string
	Creator string
}

// New creates a new empty PDF document
func New() *Document {
	theme := style.BasicClean()
	return &Document{
		theme: &theme,
		metadata: Metadata{
			Creator: "OpenScribe",
		},
	}
}

// NewWithTheme creates a PDF with a specific theme
func NewWithTheme(theme style.Theme) *Document {
	d := New()
	d.theme = &theme
	return d
}

// Theme returns the current theme
func (d *Document) Theme() *style.Theme { return d.theme }

// SetTheme sets the document theme
func (d *Document) SetTheme(theme style.Theme) { d.theme = &theme }

// SetMetadata sets document metadata
func (d *Document) SetMetadata(m Metadata) { d.metadata = m }

// GetMetadata returns document metadata
func (d *Document) GetMetadata() Metadata { return d.metadata }

// AddPage adds a new page with default A4 size
func (d *Document) AddPage() *Page {
	p := NewPage(common.PageA4, common.NormalMargins())
	d.pages = append(d.pages, p)
	return p
}

// AddPageWithSize adds a page with a specific size
func (d *Document) AddPageWithSize(size common.PageSize, margins common.Margins) *Page {
	p := NewPage(size, margins)
	d.pages = append(d.pages, p)
	return p
}

// Page returns a page by index (0-based)
func (d *Document) Page(index int) *Page {
	if index < 0 || index >= len(d.pages) {
		return nil
	}
	return d.pages[index]
}

// PageCount returns the number of pages
func (d *Document) PageCount() int { return len(d.pages) }

// RemovePage removes a page by index
func (d *Document) RemovePage(index int) error {
	if index < 0 || index >= len(d.pages) {
		return fmt.Errorf("page index %d out of range", index)
	}
	d.pages = append(d.pages[:index], d.pages[index+1:]...)
	return nil
}

// Save writes the PDF to a file
func (d *Document) Save(path string) error {
	data, err := d.build()
	if err != nil {
		return fmt.Errorf("build pdf: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// SaveToBytes returns the PDF as bytes
func (d *Document) SaveToBytes() ([]byte, error) {
	return d.build()
}

// Open reads a PDF file (basic: loads raw data for merge/split)
func Open(path string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return OpenFromBytes(data)
}

// OpenFromBytes reads a PDF from bytes
func OpenFromBytes(data []byte) (*Document, error) {
	doc := New()
	doc.pages = append(doc.pages, &Page{
		rawData: data,
		size:    common.PageA4,
		margins: common.NormalMargins(),
	})
	return doc, nil
}

// Delete removes a PDF file from disk
func Delete(path string) error {
	return os.Remove(path)
}

// Split splits the document at the given page index, returning two documents
func (d *Document) Split(atPage int) (*Document, *Document, error) {
	if atPage < 1 || atPage >= len(d.pages) {
		return nil, nil, fmt.Errorf("split index %d out of range (1 to %d)", atPage, len(d.pages)-1)
	}

	d1 := New()
	d1.pages = make([]*Page, atPage)
	copy(d1.pages, d.pages[:atPage])
	d1.metadata = d.metadata

	d2 := New()
	d2.pages = make([]*Page, len(d.pages)-atPage)
	copy(d2.pages, d.pages[atPage:])
	d2.metadata = d.metadata

	return d1, d2, nil
}

// ExtractPages extracts specific pages into a new document
func (d *Document) ExtractPages(pageIndices ...int) (*Document, error) {
	result := New()
	result.metadata = d.metadata
	for _, idx := range pageIndices {
		if idx < 0 || idx >= len(d.pages) {
			return nil, fmt.Errorf("page index %d out of range", idx)
		}
		result.pages = append(result.pages, d.pages[idx])
	}
	return result, nil
}

// Merge combines multiple PDF documents (simplified: concatenates pages)
func Merge(docs ...*Document) *Document {
	merged := New()
	for _, doc := range docs {
		for _, page := range doc.pages {
			merged.pages = append(merged.pages, page)
		}
	}
	return merged
}
