package document

import "github.com/JohnPitter/openscribe/common"

// TableOfContents represents a document table of contents
type TableOfContents struct {
	title           string
	maxLevel        int // maximum heading level to include (1-6)
	font            common.Font
	showPageNumbers bool
	entries         []TOCEntry
}

// TOCEntry represents a single entry in the table of contents
type TOCEntry struct {
	Text  string
	Level int
	Page  int
}

// AddTableOfContents adds a TOC to the document
func (d *Document) AddTableOfContents() *TableOfContents {
	toc := &TableOfContents{
		title:           "Table of Contents",
		maxLevel:        3,
		font:            common.NewFont("Arial", 11),
		showPageNumbers: true,
	}
	d.toc = toc
	return toc
}

// SetTitle sets the TOC title
func (t *TableOfContents) SetTitle(title string) { t.title = title }

// SetMaxLevel sets the maximum heading level to include
func (t *TableOfContents) SetMaxLevel(level int) {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	t.maxLevel = level
}

// SetFont sets the TOC entry font
func (t *TableOfContents) SetFont(f common.Font) { t.font = f }

// SetShowPageNumbers toggles page numbers
func (t *TableOfContents) SetShowPageNumbers(show bool) { t.showPageNumbers = show }

// Title returns the TOC title
func (t *TableOfContents) Title() string { return t.title }

// MaxLevel returns the max heading level
func (t *TableOfContents) MaxLevel() int { return t.maxLevel }

// BuildEntries scans document paragraphs to build TOC entries
func (t *TableOfContents) BuildEntries(paragraphs []*Paragraph) {
	t.entries = nil
	for _, p := range paragraphs {
		style := p.Style()
		level := 0
		switch style {
		case "Heading1":
			level = 1
		case "Heading2":
			level = 2
		case "Heading3":
			level = 3
		case "Heading4":
			level = 4
		case "Heading5":
			level = 5
		case "Heading6":
			level = 6
		}
		if level > 0 && level <= t.maxLevel {
			t.entries = append(t.entries, TOCEntry{
				Text:  p.Text(),
				Level: level,
			})
		}
	}
}

// Entries returns the TOC entries
func (t *TableOfContents) Entries() []TOCEntry { return t.entries }
