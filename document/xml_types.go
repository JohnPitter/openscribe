package document

import "encoding/xml"

const (
	nsW  = "http://schemas.openxmlformats.org/wordprocessingml/2006/main"
	nsR  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships"
	nsWP = "http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
)

// xmlDocument is the root w:document element
type xmlDocument struct {
	XMLName xml.Name         `xml:"w:document"`
	W       string           `xml:"xmlns:w,attr"`
	R       string           `xml:"xmlns:r,attr"`
	Body    xmlBody          `xml:"w:body"`
}

type xmlBody struct {
	Paragraphs []xmlParagraph   `xml:"w:p"`
	Tables     []xmlTable       `xml:"w:tbl"`
	SectPr     *xmlSectionProps `xml:"w:sectPr,omitempty"`
}

type xmlParagraph struct {
	Properties *xmlParagraphProperties `xml:"w:pPr,omitempty"`
	Runs       []xmlRun                `xml:"w:r"`
}

type xmlParagraphProperties struct {
	Style         *xmlValue  `xml:"w:pStyle,omitempty"`
	Indent        *xmlIndent `xml:"w:ind,omitempty"`
	Justification *xmlValue  `xml:"w:jc,omitempty"`
}

type xmlIndent struct {
	Left      string `xml:"w:left,attr,omitempty"`
	Right     string `xml:"w:right,attr,omitempty"`
	FirstLine string `xml:"w:firstLine,attr,omitempty"`
}

type xmlRun struct {
	Properties *xmlRunProperties `xml:"w:rPr,omitempty"`
	Text       *xmlText          `xml:"w:t,omitempty"`
	Break      *xmlBreak         `xml:"w:br,omitempty"`
}

type xmlRunProperties struct {
	Bold       *xmlEmpty      `xml:"w:b,omitempty"`
	Italic     *xmlEmpty      `xml:"w:i,omitempty"`
	Underline  *xmlValue      `xml:"w:u,omitempty"`
	Strike     *xmlEmpty      `xml:"w:strike,omitempty"`
	Color      *xmlValue      `xml:"w:color,omitempty"`
	Size       *xmlValue      `xml:"w:sz,omitempty"`
	FontFamily *xmlFontFamily `xml:"w:rFonts,omitempty"`
}

type xmlFontFamily struct {
	Ascii string `xml:"w:ascii,attr"`
	HAnsi string `xml:"w:hAnsi,attr"`
}

type xmlText struct {
	Space string `xml:"xml:space,attr,omitempty"`
	Value string `xml:",chardata"`
}

type xmlBreak struct {
	Type string `xml:"w:type,attr,omitempty"`
}

type xmlValue struct {
	Val string `xml:"w:val,attr"`
}

type xmlEmpty struct{}

type xmlTable struct {
	Properties *xmlTableProperties `xml:"w:tblPr,omitempty"`
	Grid       *xmlTableGrid       `xml:"w:tblGrid,omitempty"`
	Rows       []xmlTableRow       `xml:"w:tr"`
}

type xmlTableProperties struct {
	Borders *xmlTableBorders `xml:"w:tblBorders,omitempty"`
}

type xmlTableBorders struct {
	Top     *xmlBorderDef `xml:"w:top,omitempty"`
	Left    *xmlBorderDef `xml:"w:left,omitempty"`
	Bottom  *xmlBorderDef `xml:"w:bottom,omitempty"`
	Right   *xmlBorderDef `xml:"w:right,omitempty"`
	InsideH *xmlBorderDef `xml:"w:insideH,omitempty"`
	InsideV *xmlBorderDef `xml:"w:insideV,omitempty"`
}

type xmlBorderDef struct {
	Val   string `xml:"w:val,attr"`
	Sz    string `xml:"w:sz,attr"`
	Space string `xml:"w:space,attr"`
	Color string `xml:"w:color,attr"`
}

type xmlTableGrid struct {
	Cols []xmlGridCol `xml:"w:gridCol"`
}

type xmlGridCol struct {
	W string `xml:"w:w,attr"`
}

type xmlTableRow struct {
	Cells []xmlTableCell `xml:"w:tc"`
}

type xmlTableCell struct {
	Properties *xmlTableCellProperties `xml:"w:tcPr,omitempty"`
	Paragraphs []xmlParagraph          `xml:"w:p"`
}

type xmlTableCellProperties struct {
	Shading *xmlShading `xml:"w:shd,omitempty"`
	VAlign  *xmlValue   `xml:"w:vAlign,omitempty"`
}

type xmlShading struct {
	Val   string `xml:"w:val,attr"`
	Color string `xml:"w:color,attr"`
	Fill  string `xml:"w:fill,attr"`
}

type xmlSectionProps struct {
	HeaderRef *xmlHeaderFooterRef `xml:"w:headerReference,omitempty"`
	FooterRef *xmlHeaderFooterRef `xml:"w:footerReference,omitempty"`
	PgSz      *xmlPageSize       `xml:"w:pgSz,omitempty"`
	PgMar     *xmlPageMargins    `xml:"w:pgMar,omitempty"`
}

type xmlHeaderFooterRef struct {
	Type string `xml:"w:type,attr"`
	RID  string `xml:"r:id,attr"`
}

type xmlPageSize struct {
	W      string `xml:"w:w,attr"`
	H      string `xml:"w:h,attr"`
	Orient string `xml:"w:orient,attr,omitempty"`
}

type xmlPageMargins struct {
	Top    string `xml:"w:top,attr"`
	Right  string `xml:"w:right,attr"`
	Bottom string `xml:"w:bottom,attr"`
	Left   string `xml:"w:left,attr"`
}
