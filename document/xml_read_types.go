package document

import "encoding/xml"

// Read types use full namespace URIs for proper xml.Unmarshal behavior.
// Go's encoding/xml resolves namespace prefixes to URIs during unmarshaling,
// so struct tags must use "namespaceURI localname" format.

type readDocument struct {
	XMLName xml.Name        `xml:"document"`
	Body    readBody        `xml:"body"`
}

type readBody struct {
	Paragraphs []readParagraph `xml:"p"`
	Tables     []readTable     `xml:"tbl"`
}

type readParagraph struct {
	Properties *readParagraphProperties `xml:"pPr"`
	Runs       []readRun                `xml:"r"`
}

type readParagraphProperties struct {
	Style         *readValue `xml:"pStyle"`
	Justification *readValue `xml:"jc"`
}

type readRun struct {
	Properties *readRunProperties `xml:"rPr"`
	Text       *readText          `xml:"t"`
	Break      *readBreak         `xml:"br"`
}

type readRunProperties struct {
	Bold       *readEmpty      `xml:"b"`
	Italic     *readEmpty      `xml:"i"`
	Underline  *readValue      `xml:"u"`
	Strike     *readEmpty      `xml:"strike"`
	Color      *readValue      `xml:"color"`
	Size       *readValue      `xml:"sz"`
	FontFamily *readFontFamily `xml:"rFonts"`
}

type readFontFamily struct {
	Ascii string `xml:"ascii,attr"`
	HAnsi string `xml:"hAnsi,attr"`
}

type readText struct {
	Space string `xml:"space,attr"`
	Value string `xml:",chardata"`
}

type readBreak struct {
	Type string `xml:"type,attr"`
}

type readValue struct {
	Val string `xml:"val,attr"`
}

type readEmpty struct{}

type readTable struct {
	Rows []readTableRow `xml:"tr"`
}

type readTableRow struct {
	Cells []readTableCell `xml:"tc"`
}

type readTableCell struct {
	Paragraphs []readParagraph `xml:"p"`
}
