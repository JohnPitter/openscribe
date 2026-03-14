package document

import (
	"fmt"
)

// Footnote represents a footnote in the document
type Footnote struct {
	id   int
	text string
}

// ID returns the footnote ID
func (fn *Footnote) ID() int {
	return fn.id
}

// Text returns the footnote text
func (fn *Footnote) Text() string {
	return fn.text
}

// buildFootnotesXML creates the word/footnotes.xml content
func buildFootnotesXML(footnotes []*Footnote) []byte {
	var buf []byte
	buf = append(buf, []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)...)
	buf = append(buf, []byte(fmt.Sprintf(
		`<w:footnotes xmlns:w="%s" xmlns:r="%s">`, nsW, nsR,
	))...)

	// Separator and continuation separator (required by Word)
	buf = append(buf, []byte(
		`<w:footnote w:type="separator" w:id="-1"><w:p><w:r><w:separator/></w:r></w:p></w:footnote>`+
			`<w:footnote w:type="continuationSeparator" w:id="0"><w:p><w:r><w:continuationSeparator/></w:r></w:p></w:footnote>`,
	)...)

	for _, fn := range footnotes {
		buf = append(buf, []byte(fmt.Sprintf(
			`<w:footnote w:id="%d"><w:p><w:pPr><w:pStyle w:val="FootnoteText"/></w:pPr>`+
				`<w:r><w:rPr><w:rStyle w:val="FootnoteReference"/></w:rPr><w:footnoteRef/></w:r>`+
				`<w:r><w:t xml:space="preserve"> %s</w:t></w:r>`+
				`</w:p></w:footnote>`,
			fn.id, fn.text,
		))...)
	}

	buf = append(buf, []byte(`</w:footnotes>`)...)
	return buf
}
