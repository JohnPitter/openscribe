package document

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/internal/xmlutil"
)

// build constructs the DOCX ZIP package
func (d *Document) build() error {
	d.pkg = packaging.NewPackage()

	// Build styles.xml
	stylesXML := d.buildStylesXML()
	d.pkg.AddFile("word/styles.xml", stylesXML)

	// Document relationships
	docRels := packaging.NewRelationships()
	docRels.Add(packaging.RelTypeStylesheet, "styles.xml")

	// Add images to media and register relationships
	for _, img := range d.images {
		mediaPath := fmt.Sprintf("word/media/%s%s", img.id, img.data.Format.Extension())
		d.pkg.AddFile(mediaPath, img.data.Data)
		relTarget := fmt.Sprintf("media/%s%s", img.id, img.data.Format.Extension())
		img.relID = docRels.Add(packaging.RelTypeImage, relTarget)
	}

	// Content types
	ct := packaging.NewContentTypes()
	ct.AddOverride("/word/document.xml", packaging.ContentTypeDocx)
	ct.AddOverride("/word/styles.xml", packaging.ContentTypeStyles)

	// Build header/footer XML files and add relationships
	var headerRelID, footerRelID string
	if d.header != nil && !d.header.IsEmpty() {
		headerXML := d.buildHeaderFooterXML(d.header, "hdr")
		d.pkg.AddFile("word/header1.xml", headerXML)
		headerRelID = docRels.Add(packaging.RelTypeHeader, "header1.xml")
		ct.AddOverride("/word/header1.xml", packaging.ContentTypeHeader)
	}
	if d.footer != nil && !d.footer.IsEmpty() {
		footerXML := d.buildHeaderFooterXML(d.footer, "ftr")
		d.pkg.AddFile("word/footer1.xml", footerXML)
		footerRelID = docRels.Add(packaging.RelTypeFooter, "footer1.xml")
		ct.AddOverride("/word/footer1.xml", packaging.ContentTypeFooter)
	}

	// Build document.xml (needs rel IDs for header/footer references)
	docXML, err := d.buildDocumentXML(headerRelID, footerRelID)
	if err != nil {
		return fmt.Errorf("build document.xml: %w", err)
	}
	d.pkg.AddFile("word/document.xml", docXML)

	// Marshal document relationships
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

	// Marshal content types
	ctData, err := ct.Marshal()
	if err != nil {
		return fmt.Errorf("marshal content types: %w", err)
	}
	d.pkg.AddFile("[Content_Types].xml", ctData)

	return nil
}

func (d *Document) buildDocumentXML(headerRelID, footerRelID string) ([]byte, error) {
	doc := xmlDocument{
		W: nsW,
		R: nsR,
	}

	// Add TOC paragraphs before regular content
	if d.toc != nil {
		d.toc.BuildEntries(d.paragraphs)

		// TOC title paragraph
		tocTitle := NewParagraph()
		tocTitle.SetStyle("Heading1")
		tocTitleRun := tocTitle.AddRun()
		tocTitleRun.SetText(d.toc.Title())
		tocTitleRun.SetFont(d.toc.font.Bold())
		doc.Body.Paragraphs = append(doc.Body.Paragraphs, tocTitle.toXML())

		// TOC entry paragraphs
		for _, entry := range d.toc.Entries() {
			entryP := NewParagraph()
			indentTwips := common.Pt(float64((entry.Level - 1) * 12))
			entryP.SetIndent(indentTwips, common.Pt(0), common.Pt(0))
			entryRun := entryP.AddRun()
			entryRun.SetText(entry.Text)
			entryRun.SetFont(d.toc.font)
			doc.Body.Paragraphs = append(doc.Body.Paragraphs, entryP.toXML())
		}
	}

	// Add paragraphs
	for _, p := range d.paragraphs {
		doc.Body.Paragraphs = append(doc.Body.Paragraphs, p.toXML())
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

	// Add header/footer references to section properties
	if headerRelID != "" {
		doc.Body.SectPr.HeaderRef = &xmlHeaderFooterRef{
			Type: "default",
			RID:  headerRelID,
		}
	}
	if footerRelID != "" {
		doc.Body.SectPr.FooterRef = &xmlHeaderFooterRef{
			Type: "default",
			RID:  footerRelID,
		}
	}

	// Marshal to XML
	xmlData, err := xmlutil.MarshalXML(doc)
	if err != nil {
		return nil, err
	}

	// Embed image paragraphs via string insertion before </w:body>
	if len(d.images) > 0 {
		var imgXML bytes.Buffer
		for i, img := range d.images {
			imgID := i + 1
			emuW := img.width.EMUs()
			emuH := img.height.EMUs()
			imgXML.WriteString(fmt.Sprintf(
				`<w:p><w:r><w:drawing>`+
					`<wp:inline xmlns:wp="%s" distT="0" distB="0" distL="0" distR="0">`+
					`<wp:extent cx="%d" cy="%d"/>`+
					`<wp:docPr id="%d" name="%s"/>`+
					`<a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main">`+
					`<a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">`+
					`<pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">`+
					`<pic:nvPicPr><pic:cNvPr id="%d" name="%s"/><pic:cNvPicPr/></pic:nvPicPr>`+
					`<pic:blipFill><a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>`+
					`<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm><a:prstGeom prst="rect"/></pic:spPr>`+
					`</pic:pic>`+
					`</a:graphicData>`+
					`</a:graphic>`+
					`</wp:inline>`+
					`</w:drawing></w:r></w:p>`,
				nsWP, emuW, emuH, imgID, img.id, imgID, img.id, img.relID, emuW, emuH,
			))
		}
		xmlStr := string(xmlData)
		xmlStr = strings.Replace(xmlStr, "</w:body>", imgXML.String()+"</w:body>", 1)
		xmlData = []byte(xmlStr)
	}

	return xmlData, nil
}

// buildHeaderFooterXML creates the XML for a header or footer part.
// tag is "hdr" for header, "ftr" for footer.
func (d *Document) buildHeaderFooterXML(hf *HeaderFooter, tag string) []byte {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)

	var elem string
	if tag == "hdr" {
		elem = "w:hdr"
	} else {
		elem = "w:ftr"
	}

	buf.WriteString(fmt.Sprintf(`<%s xmlns:w="%s" xmlns:r="%s">`, elem, nsW, nsR))

	// Build paragraphs for left, center, right text
	writeParagraph := func(text, align string) {
		if text == "" {
			return
		}
		jc := ""
		if align != "" {
			jc = fmt.Sprintf(`<w:pPr><w:jc w:val="%s"/></w:pPr>`, align)
		}
		fontRP := ""
		if hf.font.Family != "" || hf.font.Size > 0 {
			fontRP = "<w:rPr>"
			if hf.font.Family != "" {
				fontRP += fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, hf.font.Family, hf.font.Family)
			}
			if hf.font.Size > 0 {
				fontRP += fmt.Sprintf(`<w:sz w:val="%d"/>`, int(hf.font.Size*2))
			}
			fontRP += "</w:rPr>"
		}
		buf.WriteString(fmt.Sprintf(`<w:p>%s<w:r>%s<w:t xml:space="preserve">%s</w:t></w:r></w:p>`, jc, fontRP, text))
	}

	writeParagraph(hf.leftText, "left")
	writeParagraph(hf.centerText, "center")
	writeParagraph(hf.rightText, "right")

	// If all texts empty, write an empty paragraph
	if hf.IsEmpty() {
		buf.WriteString(`<w:p/>`)
	}

	buf.WriteString(fmt.Sprintf(`</%s>`, elem))
	return buf.Bytes()
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
				xc.Paragraphs = append(xc.Paragraphs, p.toXML())
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
