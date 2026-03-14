package document

import (
	"fmt"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/internal/xmlutil"
)

// build constructs the DOCX ZIP package
func (d *Document) build() error {
	d.pkg = packaging.NewPackage()

	// Build document.xml
	docXML, err := d.buildDocumentXML()
	if err != nil {
		return fmt.Errorf("build document.xml: %w", err)
	}
	d.pkg.AddFile("word/document.xml", docXML)

	// Build styles.xml
	stylesXML := d.buildStylesXML()
	d.pkg.AddFile("word/styles.xml", stylesXML)

	// Add images
	docRels := packaging.NewRelationships()
	docRels.Add(packaging.RelTypeStylesheet, "styles.xml")
	for _, img := range d.images {
		mediaPath := fmt.Sprintf("word/media/%s%s", img.id, img.data.Format.Extension())
		d.pkg.AddFile(mediaPath, img.data.Data)
		relTarget := fmt.Sprintf("media/%s%s", img.id, img.data.Format.Extension())
		img.relID = docRels.Add(packaging.RelTypeImage, relTarget)
	}

	// Document relationships
	docRelsData, err := docRels.Marshal()
	if err != nil {
		return fmt.Errorf("marshal doc rels: %w", err)
	}
	d.pkg.AddFile("word/_rels/document.xml.rels", docRelsData)

	// Top-level relationships
	topRels := packaging.NewRelationships()
	topRels.Add(packaging.RelTypeOfficeDocument, "word/document.xml")
	topRelsData, err := topRels.Marshal()
	if err != nil {
		return fmt.Errorf("marshal top rels: %w", err)
	}
	d.pkg.AddFile("_rels/.rels", topRelsData)

	// Content types
	ct := packaging.NewContentTypes()
	ct.AddOverride("/word/document.xml", packaging.ContentTypeDocx)
	ct.AddOverride("/word/styles.xml", packaging.ContentTypeStyles)
	ctData, err := ct.Marshal()
	if err != nil {
		return fmt.Errorf("marshal content types: %w", err)
	}
	d.pkg.AddFile("[Content_Types].xml", ctData)

	return nil
}

func (d *Document) buildDocumentXML() ([]byte, error) {
	doc := xmlDocument{
		W: nsW,
		R: nsR,
	}

	// Add paragraphs
	for _, p := range d.paragraphs {
		doc.Body.Paragraphs = append(doc.Body.Paragraphs, p.MarshalXML())
	}

	// Add tables
	for _, tbl := range d.tables {
		doc.Body.Tables = append(doc.Body.Tables, tbl.marshalXML())
	}

	// Section properties
	sect := d.Section()
	doc.Body.SectPr = &xmlSectionProps{
		PgSz: &xmlPageSize{
			W: fmt.Sprintf("%d", int(sect.pageSize.Width.Points()*20)),
			H: fmt.Sprintf("%d", int(sect.pageSize.Height.Points()*20)),
		},
		PgMar: &xmlPageMargins{
			Top:    fmt.Sprintf("%d", int(sect.margins.Top.Points()*20)),
			Right:  fmt.Sprintf("%d", int(sect.margins.Right.Points()*20)),
			Bottom: fmt.Sprintf("%d", int(sect.margins.Bottom.Points()*20)),
			Left:   fmt.Sprintf("%d", int(sect.margins.Left.Points()*20)),
		},
	}
	if sect.orientation == common.OrientationLandscape {
		doc.Body.SectPr.PgSz.Orient = "landscape"
	}

	return xmlutil.MarshalXML(doc)
}

func (d *Document) buildStylesXML() []byte {
	// Minimal styles.xml
	return []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:styles xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:style w:type="paragraph" w:styleId="Heading1">
    <w:name w:val="heading 1"/>
    <w:pPr><w:outlineLvl w:val="0"/></w:pPr>
    <w:rPr><w:b/><w:sz w:val="48"/></w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading2">
    <w:name w:val="heading 2"/>
    <w:pPr><w:outlineLvl w:val="1"/></w:pPr>
    <w:rPr><w:b/><w:sz w:val="36"/></w:rPr>
  </w:style>
  <w:style w:type="paragraph" w:styleId="Heading3">
    <w:name w:val="heading 3"/>
    <w:pPr><w:outlineLvl w:val="2"/></w:pPr>
    <w:rPr><w:b/><w:sz w:val="28"/></w:rPr>
  </w:style>
</w:styles>`)
}

// parseDocument is a basic parser for existing DOCX content.
// Uses readDocument types which handle namespace-aware unmarshaling correctly.
func (d *Document) parseDocument(data []byte) error {
	var doc readDocument
	if err := xmlutil.UnmarshalXML(data, &doc); err != nil {
		// If parsing fails, just start fresh
		return nil
	}

	for _, xp := range doc.Body.Paragraphs {
		p := NewParagraph()
		if xp.Properties != nil && xp.Properties.Style != nil {
			p.SetStyle(xp.Properties.Style.Val)
		}
		for _, xr := range xp.Runs {
			r := p.AddRun()
			if xr.Text != nil {
				r.SetText(xr.Text.Value)
			}
			if xr.Properties != nil {
				if xr.Properties.Bold != nil {
					r.SetBold(true)
				}
				if xr.Properties.Italic != nil {
					r.SetItalic(true)
				}
			}
		}
		d.paragraphs = append(d.paragraphs, p)
	}

	return nil
}

// marshalXML creates the w:tbl XML element for a table
func (t *Table) marshalXML() xmlTable {
	xt := xmlTable{
		Properties: &xmlTableProperties{
			Borders: &xmlTableBorders{
				Top:     borderDefFromCommon(t.borders.Top),
				Left:    borderDefFromCommon(t.borders.Left),
				Bottom:  borderDefFromCommon(t.borders.Bottom),
				Right:   borderDefFromCommon(t.borders.Right),
				InsideH: borderDefFromCommon(t.borders.Top),
				InsideV: borderDefFromCommon(t.borders.Left),
			},
		},
	}

	for _, row := range t.rows {
		xr := xmlTableRow{}
		for _, cell := range row.cells {
			xc := xmlTableCell{}
			if cell.shading != nil {
				hex := cell.shading.Hex()
				if len(hex) > 0 && hex[0] == '#' {
					hex = hex[1:]
				}
				xc.Properties = &xmlTableCellProperties{
					Shading: &xmlShading{
						Val:   "clear",
						Color: "auto",
						Fill:  hex,
					},
				}
			}
			for _, p := range cell.paragraphs {
				xc.Paragraphs = append(xc.Paragraphs, p.MarshalXML())
			}
			if len(xc.Paragraphs) == 0 {
				xc.Paragraphs = append(xc.Paragraphs, xmlParagraph{})
			}
			xr.Cells = append(xr.Cells, xc)
		}
		xt.Rows = append(xt.Rows, xr)
	}

	return xt
}

func borderDefFromCommon(b common.Border) *xmlBorderDef {
	if b.Style == common.BorderStyleNone {
		return nil
	}
	hex := b.Color.Hex()
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	val := "single"
	switch b.Style {
	case common.BorderStyleDashed:
		val = "dashed"
	case common.BorderStyleDotted:
		val = "dotted"
	case common.BorderStyleDouble:
		val = "double"
	}
	return &xmlBorderDef{
		Val:   val,
		Sz:    fmt.Sprintf("%d", int(b.Width.Points()*8)),
		Space: "0",
		Color: hex,
	}
}
