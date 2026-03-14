// Package document provides DOCX document creation, reading, and editing.
package document

import (
	"fmt"
	"os"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/style"
)

// Document represents a DOCX document
type Document struct {
	pkg           *packaging.Package
	paragraphs    []*Paragraph
	tables        []*Table
	sections      []*Section
	images        []*ImageRef
	header        *HeaderFooter
	footer        *HeaderFooter
	toc           *TableOfContents
	theme         *style.Theme
	rels          *packaging.Relationships
	docRels       *packaging.Relationships
	ct            *packaging.ContentTypes
	imageCount    int
	lists         []*List
	listCount     int
	footnotes     []*Footnote
	footnoteCount int
	comments      []*Comment
	commentCount  int
	customStyles  []*CustomStyle
	security      *common.SecurityOptions
	charts        []*Chart
	chartCount    int
}

// New creates a new empty DOCX document
func New() *Document {
	theme := style.BasicClean()
	return &Document{
		pkg:     packaging.NewPackage(),
		theme:   &theme,
		rels:    packaging.NewRelationships(),
		docRels: packaging.NewRelationships(),
		ct:      packaging.NewContentTypes(),
		sections: []*Section{
			NewSection(), // default section
		},
	}
}

// NewWithTheme creates a new document with a specific design theme
func NewWithTheme(theme style.Theme) *Document {
	doc := New()
	doc.theme = &theme
	return doc
}

// Theme returns the current theme
func (d *Document) Theme() *style.Theme {
	return d.theme
}

// SetTheme sets the document theme
func (d *Document) SetTheme(theme style.Theme) {
	d.theme = &theme
}

// AddParagraph adds a new paragraph and returns it
func (d *Document) AddParagraph() *Paragraph {
	p := NewParagraph()
	d.paragraphs = append(d.paragraphs, p)
	return p
}

// AddHeading adds a heading paragraph with the specified level (1-6)
func (d *Document) AddHeading(text string, level int) *Paragraph {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	p := NewParagraph()
	p.SetStyle(fmt.Sprintf("Heading%d", level))
	run := p.AddRun()
	run.SetText(text)
	if d.theme != nil {
		run.SetFont(d.theme.Typography.HeadingFont)
	}
	d.paragraphs = append(d.paragraphs, p)
	return p
}

// AddText adds a simple text paragraph
func (d *Document) AddText(text string) *Paragraph {
	p := NewParagraph()
	run := p.AddRun()
	run.SetText(text)
	if d.theme != nil {
		run.SetFont(d.theme.Typography.BodyFont)
	}
	d.paragraphs = append(d.paragraphs, p)
	return p
}

// AddTable adds a new table with specified rows and columns
func (d *Document) AddTable(rows, cols int) *Table {
	t := NewTable(rows, cols)
	if d.theme != nil {
		t.ApplyTheme(*d.theme)
	}
	d.tables = append(d.tables, t)
	return t
}

// AddImage adds an image to the document
func (d *Document) AddImage(imgData *common.ImageData, width, height common.Measurement) *ImageRef {
	d.imageCount++
	ref := &ImageRef{
		id:     fmt.Sprintf("img%d", d.imageCount),
		data:   imgData,
		width:  width,
		height: height,
	}
	d.images = append(d.images, ref)
	return ref
}

// AddList adds a new list to the document and returns it
func (d *Document) AddList(listType ListType) *List {
	d.listCount++
	l := NewList(listType, d.listCount)
	d.lists = append(d.lists, l)
	return l
}

// Lists returns all lists in the document
func (d *Document) Lists() []*List {
	return d.lists
}

// AddFootnote adds a footnote and returns its reference ID
func (d *Document) AddFootnote(text string) int {
	d.footnoteCount++
	fn := &Footnote{
		id:   d.footnoteCount,
		text: text,
	}
	d.footnotes = append(d.footnotes, fn)
	return fn.id
}

// Footnotes returns all footnotes
func (d *Document) Footnotes() []*Footnote {
	return d.footnotes
}

// AddComment adds a comment and returns it
func (d *Document) AddComment(author, text string) *Comment {
	d.commentCount++
	c := NewComment(d.commentCount, author, text)
	d.comments = append(d.comments, c)
	return c
}

// Comments returns all comments
func (d *Document) Comments() []*Comment {
	return d.comments
}

// AddStyle adds a custom paragraph style and returns it
func (d *Document) AddStyle(name, basedOn string) *CustomStyle {
	cs := NewCustomStyle(name, basedOn)
	d.customStyles = append(d.customStyles, cs)
	return cs
}

// CustomStyles returns all custom styles
func (d *Document) CustomStyles() []*CustomStyle {
	return d.customStyles
}

// AddPageBreak adds a page break
func (d *Document) AddPageBreak() {
	p := NewParagraph()
	p.AddPageBreak()
	d.paragraphs = append(d.paragraphs, p)
}

// Paragraphs returns all paragraphs
func (d *Document) Paragraphs() []*Paragraph {
	return d.paragraphs
}

// Tables returns all tables
func (d *Document) Tables() []*Table {
	return d.tables
}

// RemoveParagraph removes a paragraph by index
func (d *Document) RemoveParagraph(index int) error {
	if index < 0 || index >= len(d.paragraphs) {
		return fmt.Errorf("paragraph index %d out of range", index)
	}
	d.paragraphs = append(d.paragraphs[:index], d.paragraphs[index+1:]...)
	return nil
}

// RemoveTable removes a table by index
func (d *Document) RemoveTable(index int) error {
	if index < 0 || index >= len(d.tables) {
		return fmt.Errorf("table index %d out of range", index)
	}
	d.tables = append(d.tables[:index], d.tables[index+1:]...)
	return nil
}

// Section returns the default section
func (d *Document) Section() *Section {
	if len(d.sections) == 0 {
		d.sections = append(d.sections, NewSection())
	}
	return d.sections[0]
}

// Save writes the document to a file
func (d *Document) Save(path string) error {
	if err := d.build(); err != nil {
		return fmt.Errorf("build document: %w", err)
	}
	return d.pkg.Save(path)
}

// SaveToBytes returns the document as bytes
func (d *Document) SaveToBytes() ([]byte, error) {
	if err := d.build(); err != nil {
		return nil, fmt.Errorf("build document: %w", err)
	}
	return d.pkg.ToBytes()
}

// Open reads a DOCX file (basic implementation -- loads the package)
func Open(path string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return OpenFromBytes(data)
}

// OpenFromBytes reads a DOCX from bytes
func OpenFromBytes(data []byte) (*Document, error) {
	pkg, err := packaging.OpenPackageFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("open package: %w", err)
	}

	doc := &Document{
		pkg:      pkg,
		rels:     packaging.NewRelationships(),
		docRels:  packaging.NewRelationships(),
		ct:       packaging.NewContentTypes(),
		sections: []*Section{NewSection()},
	}

	// Parse existing content if available
	if xmlData, ok := pkg.GetFile("word/document.xml"); ok {
		if err := doc.parseDocument(xmlData); err != nil {
			return nil, fmt.Errorf("parse document: %w", err)
		}
	}

	return doc, nil
}

// Delete removes a DOCX file from disk
func Delete(path string) error {
	return os.Remove(path)
}
