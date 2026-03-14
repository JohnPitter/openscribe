package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/JohnPitter/openscribe/common"
)

// pdfWriter helps construct PDF content
type pdfWriter struct {
	buf     bytes.Buffer
	objects []int // byte offset of each object
	nextObj int
}

func newPDFWriter() *pdfWriter {
	return &pdfWriter{nextObj: 1}
}

func (w *pdfWriter) startObject(id int) {
	for len(w.objects) < id {
		w.objects = append(w.objects, 0)
	}
	w.objects[id-1] = w.buf.Len()
	fmt.Fprintf(&w.buf, "%d 0 obj\n", id)
}

func (w *pdfWriter) endObject() {
	fmt.Fprintf(&w.buf, "endobj\n\n")
}

func (w *pdfWriter) allocObj() int {
	id := w.nextObj
	w.nextObj++
	return id
}

// imageRef tracks an image XObject for a page
type imageRef struct {
	objID   int
	imgName string
	elem    *ImageElement
}

// formFieldRef tracks a form field widget
type formFieldRef struct {
	objID int
}

// annotRef tracks an annotation object
type annotRef struct {
	objID int
}

func (d *Document) build() ([]byte, error) {
	w := newPDFWriter()

	// Pre-allocate object IDs
	catalogID := w.allocObj() // 1
	pagesID := w.allocObj()   // 2
	fontID := w.allocObj()    // 3

	// Allocate page + content stream IDs
	type pageIDs struct {
		pageID    int
		contentID int
	}
	var pids []pageIDs
	for range d.pages {
		pid := pageIDs{
			pageID:    w.allocObj(),
			contentID: w.allocObj(),
		}
		pids = append(pids, pid)
	}

	// Scan for images, form fields, and annotations to pre-allocate object IDs
	type pageExtras struct {
		images     []imageRef
		formFields []formFieldRef
		annots     []annotRef
	}
	pageExtrasMap := make([]pageExtras, len(d.pages))
	var allFormFieldObjIDs []int

	for i, page := range d.pages {
		imgIndex := 0
		for _, elem := range page.elements {
			switch e := elem.(type) {
			case *ImageElement:
				if e.data != nil && len(e.data.Data) > 0 {
					objID := w.allocObj()
					imgName := fmt.Sprintf("Im%d", imgIndex)
					imgIndex++
					pageExtrasMap[i].images = append(pageExtrasMap[i].images, imageRef{
						objID:   objID,
						imgName: imgName,
						elem:    e,
					})
				}
			case *TextField:
				objID := w.allocObj()
				pageExtrasMap[i].formFields = append(pageExtrasMap[i].formFields, formFieldRef{objID: objID})
				allFormFieldObjIDs = append(allFormFieldObjIDs, objID)
			case *Checkbox:
				objID := w.allocObj()
				pageExtrasMap[i].formFields = append(pageExtrasMap[i].formFields, formFieldRef{objID: objID})
				allFormFieldObjIDs = append(allFormFieldObjIDs, objID)
			case *Dropdown:
				objID := w.allocObj()
				pageExtrasMap[i].formFields = append(pageExtrasMap[i].formFields, formFieldRef{objID: objID})
				allFormFieldObjIDs = append(allFormFieldObjIDs, objID)
			case *Annotation:
				objID := w.allocObj()
				pageExtrasMap[i].annots = append(pageExtrasMap[i].annots, annotRef{objID: objID})
			}
		}
	}

	// Allocate AcroForm object ID if there are form fields
	var acroFormID int
	if len(allFormFieldObjIDs) > 0 {
		acroFormID = w.allocObj()
	}

	// Allocate Info dictionary object ID
	infoID := w.allocObj()

	// Write header
	w.buf.WriteString("%PDF-1.4\n")
	w.buf.WriteString("%\xe2\xe3\xcf\xd3\n\n")

	// Font object (Helvetica - always available)
	w.startObject(fontID)
	w.buf.WriteString("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\n")
	w.endObject()

	// Page objects and content streams
	for i, page := range d.pages {
		pid := pids[i]
		extras := pageExtrasMap[i]

		// Build content stream
		var content bytes.Buffer
		pageH := page.size.Height.Points()

		// Background
		if page.background != nil {
			r := float64(page.background.R) / 255
			g := float64(page.background.G) / 255
			b := float64(page.background.B) / 255
			fmt.Fprintf(&content, "%.3f %.3f %.3f rg\n", r, g, b)
			fmt.Fprintf(&content, "0 0 %.2f %.2f re f\n",
				page.size.Width.Points(), pageH)
		}

		imgIdx := 0
		for _, elem := range page.elements {
			switch e := elem.(type) {
			case *TextElement:
				r := float64(e.font.Color.R) / 255
				g := float64(e.font.Color.G) / 255
				b := float64(e.font.Color.B) / 255
				// PDF y-axis is bottom-up, so we invert
				y := pageH - e.y
				fmt.Fprintf(&content, "BT\n")
				fmt.Fprintf(&content, "/F1 %.1f Tf\n", e.font.Size)
				fmt.Fprintf(&content, "%.3f %.3f %.3f rg\n", r, g, b)
				fmt.Fprintf(&content, "%.2f %.2f Td\n", e.x, y)
				fmt.Fprintf(&content, "(%s) Tj\n", escapePDF(e.text))
				fmt.Fprintf(&content, "ET\n")

			case *LineElement:
				r := float64(e.color.R) / 255
				g := float64(e.color.G) / 255
				b := float64(e.color.B) / 255
				fmt.Fprintf(&content, "%.3f %.3f %.3f RG\n", r, g, b)
				fmt.Fprintf(&content, "%.2f w\n", e.width)
				fmt.Fprintf(&content, "%.2f %.2f m %.2f %.2f l S\n",
					e.x1, pageH-e.y1, e.x2, pageH-e.y2)

			case *RectElement:
				rf := float64(e.fill.R) / 255
				gf := float64(e.fill.G) / 255
				bf := float64(e.fill.B) / 255
				fmt.Fprintf(&content, "%.3f %.3f %.3f rg\n", rf, gf, bf)
				y := pageH - e.y - e.height
				fmt.Fprintf(&content, "%.2f %.2f %.2f %.2f re f\n",
					e.x, y, e.width, e.height)
				if e.stroke != nil {
					rs := float64(e.stroke.R) / 255
					gs := float64(e.stroke.G) / 255
					bs := float64(e.stroke.B) / 255
					fmt.Fprintf(&content, "%.3f %.3f %.3f RG\n", rs, gs, bs)
					fmt.Fprintf(&content, "%.2f %.2f %.2f %.2f re S\n",
						e.x, y, e.width, e.height)
				}

			case *TableElement:
				buildTablePDF(&content, e, pageH)

			case *ChartElement:
				buildChartPDF(&content, e, pageH)

			case *ImageElement:
				if e.data != nil && len(e.data.Data) > 0 && imgIdx < len(extras.images) {
					ref := extras.images[imgIdx]
					imgIdx++
					// Place image using cm (concat matrix) operator
					y := pageH - e.y - e.height
					fmt.Fprintf(&content, "q\n")
					fmt.Fprintf(&content, "%.2f 0 0 %.2f %.2f %.2f cm\n",
						e.width, e.height, e.x, y)
					fmt.Fprintf(&content, "/%s Do\n", ref.imgName)
					fmt.Fprintf(&content, "Q\n")
				} else {
					// Placeholder rectangle for images without data
					y := pageH - e.y - e.height
					fmt.Fprintf(&content, "0.9 0.9 0.9 rg\n")
					fmt.Fprintf(&content, "%.2f %.2f %.2f %.2f re f\n", e.x, y, e.width, e.height)
					fmt.Fprintf(&content, "0.5 0.5 0.5 RG\n0.5 w\n")
					fmt.Fprintf(&content, "%.2f %.2f %.2f %.2f re S\n", e.x, y, e.width, e.height)
				}

			case *TextBlock:
				buildTextBlockPDF(&content, e, pageH)

			// Form fields and annotations are handled as separate objects, not in content stream
			case *TextField, *Checkbox, *Dropdown, *Annotation:
				// These are serialized as annotation dictionaries, not in the content stream
			}
		}

		contentStr := content.String()

		// Content stream object
		w.startObject(pid.contentID)
		fmt.Fprintf(&w.buf, "<< /Length %d >>\n", len(contentStr))
		w.buf.WriteString("stream\n")
		w.buf.WriteString(contentStr)
		w.buf.WriteString("\nendstream\n")
		w.endObject()

		// Write image XObject streams
		for _, ref := range extras.images {
			writeImageXObject(w, ref)
		}

		// Write form field widget annotations
		ffIdx := 0
		for _, elem := range page.elements {
			switch e := elem.(type) {
			case *TextField:
				if ffIdx < len(extras.formFields) {
					writeTextFieldWidget(w, extras.formFields[ffIdx].objID, pid.pageID, fontID, e, pageH)
					ffIdx++
				}
			case *Checkbox:
				if ffIdx < len(extras.formFields) {
					writeCheckboxWidget(w, extras.formFields[ffIdx].objID, pid.pageID, e, pageH)
					ffIdx++
				}
			case *Dropdown:
				if ffIdx < len(extras.formFields) {
					writeDropdownWidget(w, extras.formFields[ffIdx].objID, pid.pageID, fontID, e, pageH)
					ffIdx++
				}
			}
		}

		// Write annotation objects
		annotIdx := 0
		for _, elem := range page.elements {
			if a, ok := elem.(*Annotation); ok {
				if annotIdx < len(extras.annots) {
					writeAnnotation(w, extras.annots[annotIdx].objID, pid.pageID, a, pageH, fontID)
					annotIdx++
				}
			}
		}

		// Build XObject references for resources
		var xobjEntries []string
		for _, ref := range extras.images {
			xobjEntries = append(xobjEntries, fmt.Sprintf("/%s %d 0 R", ref.imgName, ref.objID))
		}

		// Collect annotation/widget references for the /Annots array
		var annotRefs []int
		for _, ff := range extras.formFields {
			annotRefs = append(annotRefs, ff.objID)
		}
		for _, an := range extras.annots {
			annotRefs = append(annotRefs, an.objID)
		}

		// Page object
		w.startObject(pid.pageID)
		var pageBuf bytes.Buffer
		fmt.Fprintf(&pageBuf, "<< /Type /Page /Parent %d 0 R /MediaBox [0 0 %.2f %.2f] /Contents %d 0 R",
			pagesID,
			page.size.Width.Points(),
			page.size.Height.Points(),
			pid.contentID,
		)

		// Resources
		pageBuf.WriteString(" /Resources << /Font << /F1 ")
		fmt.Fprintf(&pageBuf, "%d 0 R >>", fontID)
		if len(xobjEntries) > 0 {
			pageBuf.WriteString(" /XObject << ")
			pageBuf.WriteString(strings.Join(xobjEntries, " "))
			pageBuf.WriteString(" >>")
		}
		pageBuf.WriteString(" >>")

		// Annotations array
		if len(annotRefs) > 0 {
			pageBuf.WriteString(" /Annots [")
			for j, ref := range annotRefs {
				if j > 0 {
					pageBuf.WriteString(" ")
				}
				fmt.Fprintf(&pageBuf, "%d 0 R", ref)
			}
			pageBuf.WriteString("]")
		}

		pageBuf.WriteString(" >>\n")
		w.buf.Write(pageBuf.Bytes())
		w.endObject()
	}

	// Pages object
	w.startObject(pagesID)
	fmt.Fprintf(&w.buf, "<< /Type /Pages /Kids [")
	for i, pid := range pids {
		if i > 0 {
			w.buf.WriteString(" ")
		}
		fmt.Fprintf(&w.buf, "%d 0 R", pid.pageID)
	}
	fmt.Fprintf(&w.buf, "] /Count %d >>\n", len(d.pages))
	w.endObject()

	// AcroForm object (if form fields exist)
	if acroFormID > 0 {
		w.startObject(acroFormID)
		fmt.Fprintf(&w.buf, "<< /Fields [")
		for i, fid := range allFormFieldObjIDs {
			if i > 0 {
				w.buf.WriteString(" ")
			}
			fmt.Fprintf(&w.buf, "%d 0 R", fid)
		}
		fmt.Fprintf(&w.buf, "] /DR << /Font << /F1 %d 0 R >> >> /NeedAppearances true >>\n", fontID)
		w.endObject()
	}

	// Info dictionary
	w.startObject(infoID)
	now := time.Now().UTC()
	dateStr := fmt.Sprintf("D:%04d%02d%02d%02d%02d%02dZ",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	fmt.Fprintf(&w.buf, "<< /Producer (OpenScribe PDF Library)")
	if d.metadata.Title != "" {
		fmt.Fprintf(&w.buf, " /Title (%s)", escapePDF(d.metadata.Title))
	}
	if d.metadata.Author != "" {
		fmt.Fprintf(&w.buf, " /Author (%s)", escapePDF(d.metadata.Author))
	}
	if d.metadata.Subject != "" {
		fmt.Fprintf(&w.buf, " /Subject (%s)", escapePDF(d.metadata.Subject))
	}
	if d.metadata.Creator != "" {
		fmt.Fprintf(&w.buf, " /Creator (%s)", escapePDF(d.metadata.Creator))
	}
	fmt.Fprintf(&w.buf, " /CreationDate (%s) /ModDate (%s) >>\n", dateStr, dateStr)
	w.endObject()

	// Catalog object
	w.startObject(catalogID)
	if acroFormID > 0 {
		fmt.Fprintf(&w.buf, "<< /Type /Catalog /Pages %d 0 R /AcroForm %d 0 R >>\n", pagesID, acroFormID)
	} else {
		fmt.Fprintf(&w.buf, "<< /Type /Catalog /Pages %d 0 R >>\n", pagesID)
	}
	w.endObject()

	// Cross-reference table
	xrefOffset := w.buf.Len()
	fmt.Fprintf(&w.buf, "xref\n")
	fmt.Fprintf(&w.buf, "0 %d\n", w.nextObj)
	fmt.Fprintf(&w.buf, "0000000000 65535 f \n")
	for _, offset := range w.objects {
		fmt.Fprintf(&w.buf, "%010d 00000 n \n", offset)
	}

	// Trailer with Info reference
	fmt.Fprintf(&w.buf, "trailer\n")
	fmt.Fprintf(&w.buf, "<< /Size %d /Root %d 0 R /Info %d 0 R >>\n", w.nextObj, catalogID, infoID)
	fmt.Fprintf(&w.buf, "startxref\n%d\n%%%%EOF\n", xrefOffset)

	return w.buf.Bytes(), nil
}

// writeImageXObject writes an image as a PDF XObject stream
func writeImageXObject(w *pdfWriter, ref imageRef) {
	e := ref.elem
	imgData := e.data

	w.startObject(ref.objID)

	// Determine pixel dimensions from ImageData or use element dimensions
	pixW := int(imgData.Width.Points())
	pixH := int(imgData.Height.Points())
	if pixW <= 0 {
		pixW = int(e.width)
	}
	if pixH <= 0 {
		pixH = int(e.height)
	}
	if pixW <= 0 {
		pixW = 1
	}
	if pixH <= 0 {
		pixH = 1
	}

	switch imgData.Format {
	case common.ImageFormatJPEG:
		// Embed raw JPEG data with DCTDecode
		fmt.Fprintf(&w.buf, "<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /DCTDecode /Length %d >>\n",
			pixW, pixH, len(imgData.Data))
		w.buf.WriteString("stream\n")
		w.buf.Write(imgData.Data)
		w.buf.WriteString("\nendstream\n")

	case common.ImageFormatPNG:
		// Compress with zlib (FlateDecode)
		var compressed bytes.Buffer
		zw := zlib.NewWriter(&compressed)
		zw.Write(imgData.Data)
		zw.Close()

		fmt.Fprintf(&w.buf, "<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /FlateDecode /Length %d >>\n",
			pixW, pixH, compressed.Len())
		w.buf.WriteString("stream\n")
		w.buf.Write(compressed.Bytes())
		w.buf.WriteString("\nendstream\n")

	default:
		// Fallback: embed raw data
		fmt.Fprintf(&w.buf, "<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Length %d >>\n",
			pixW, pixH, len(imgData.Data))
		w.buf.WriteString("stream\n")
		w.buf.Write(imgData.Data)
		w.buf.WriteString("\nendstream\n")
	}

	w.endObject()
}

// writeTextFieldWidget writes a text field widget annotation
func writeTextFieldWidget(w *pdfWriter, objID, pageID, fontID int, tf *TextField, pageH float64) {
	y := pageH - tf.y - tf.height
	w.startObject(objID)

	flags := 0
	if tf.readOnly {
		flags |= 1 // Bit 1: ReadOnly
	}
	if tf.required {
		flags |= 2 // Bit 2: Required
	}

	ffFlags := 0
	if tf.multiline {
		ffFlags |= 1 << 12 // Bit 13: Multiline
	}

	fmt.Fprintf(&w.buf, "<< /Type /Annot /Subtype /Widget /FT /Tx /T (%s) /Rect [%.2f %.2f %.2f %.2f] /P %d 0 R /Ff %d",
		escapePDF(tf.name), tf.x, y, tf.x+tf.width, y+tf.height, pageID, flags|ffFlags)

	if tf.value != "" {
		fmt.Fprintf(&w.buf, " /V (%s)", escapePDF(tf.value))
	}
	if tf.maxLength > 0 {
		fmt.Fprintf(&w.buf, " /MaxLen %d", tf.maxLength)
	}

	// Default appearance
	fmt.Fprintf(&w.buf, " /DA (/F1 12 Tf 0 0 0 rg)")

	// Border style
	fmt.Fprintf(&w.buf, " /BS << /W 1 /S /S >>")
	fmt.Fprintf(&w.buf, " /MK << /BC [0 0 0] >>")

	fmt.Fprintf(&w.buf, " >>\n")
	w.endObject()
}

// writeCheckboxWidget writes a checkbox widget annotation
func writeCheckboxWidget(w *pdfWriter, objID, pageID int, cb *Checkbox, pageH float64) {
	y := pageH - cb.y - cb.height
	w.startObject(objID)

	flags := 0
	if cb.readOnly {
		flags |= 1
	}
	if cb.required {
		flags |= 2
	}

	onValue := "Yes"
	currentVal := "Off"
	if cb.checked {
		currentVal = onValue
	}

	fmt.Fprintf(&w.buf, "<< /Type /Annot /Subtype /Widget /FT /Btn /T (%s) /Rect [%.2f %.2f %.2f %.2f] /P %d 0 R /Ff %d /V /%s /AS /%s",
		escapePDF(cb.name), cb.x, y, cb.x+cb.width, y+cb.height, pageID, flags, currentVal, currentVal)

	// Appearance with checkmark
	fmt.Fprintf(&w.buf, " /MK << /CA (4) >>") // 4 = checkmark in ZapfDingbats
	fmt.Fprintf(&w.buf, " /BS << /W 1 /S /S >>")

	fmt.Fprintf(&w.buf, " >>\n")
	w.endObject()
}

// writeDropdownWidget writes a dropdown/choice widget annotation
func writeDropdownWidget(w *pdfWriter, objID, pageID, fontID int, dd *Dropdown, pageH float64) {
	y := pageH - dd.y - dd.height
	w.startObject(objID)

	flags := 0
	if dd.readOnly {
		flags |= 1
	}
	if dd.required {
		flags |= 2
	}

	// Choice field with combo flag (bit 18)
	ffFlags := 1 << 17 // Combo box

	fmt.Fprintf(&w.buf, "<< /Type /Annot /Subtype /Widget /FT /Ch /T (%s) /Rect [%.2f %.2f %.2f %.2f] /P %d 0 R /Ff %d",
		escapePDF(dd.name), dd.x, y, dd.x+dd.width, y+dd.height, pageID, flags|ffFlags)

	// Options
	if len(dd.options) > 0 {
		w.buf.WriteString(" /Opt [")
		for j, opt := range dd.options {
			if j > 0 {
				w.buf.WriteString(" ")
			}
			fmt.Fprintf(&w.buf, "(%s)", escapePDF(opt))
		}
		w.buf.WriteString("]")
	}

	if dd.value != "" {
		fmt.Fprintf(&w.buf, " /V (%s)", escapePDF(dd.value))
	}

	// Default appearance
	fmt.Fprintf(&w.buf, " /DA (/F1 12 Tf 0 0 0 rg)")
	fmt.Fprintf(&w.buf, " /BS << /W 1 /S /S >>")
	fmt.Fprintf(&w.buf, " /MK << /BC [0 0 0] >>")

	fmt.Fprintf(&w.buf, " >>\n")
	w.endObject()
}

// writeAnnotation writes a PDF annotation object
func writeAnnotation(w *pdfWriter, objID, pageID int, a *Annotation, pageH float64, fontID int) {
	// Convert y coordinates to PDF coordinate space
	py1 := pageH - a.y2
	py2 := pageH - a.y1

	w.startObject(objID)

	r := float64(a.color.R) / 255
	g := float64(a.color.G) / 255
	b := float64(a.color.B) / 255

	switch a.annotType {
	case AnnotHighlight:
		fmt.Fprintf(&w.buf, "<< /Type /Annot /Subtype /Highlight /Rect [%.2f %.2f %.2f %.2f] /C [%.3f %.3f %.3f] /P %d 0 R",
			a.x1, py1, a.x2, py2, r, g, b, pageID)
		// QuadPoints for highlight area
		fmt.Fprintf(&w.buf, " /QuadPoints [%.2f %.2f %.2f %.2f %.2f %.2f %.2f %.2f]",
			a.x1, py2, a.x2, py2, a.x1, py1, a.x2, py1)

	case AnnotUnderline:
		fmt.Fprintf(&w.buf, "<< /Type /Annot /Subtype /Underline /Rect [%.2f %.2f %.2f %.2f] /C [%.3f %.3f %.3f] /P %d 0 R",
			a.x1, py1, a.x2, py2, r, g, b, pageID)
		fmt.Fprintf(&w.buf, " /QuadPoints [%.2f %.2f %.2f %.2f %.2f %.2f %.2f %.2f]",
			a.x1, py2, a.x2, py2, a.x1, py1, a.x2, py1)

	case AnnotStrikeout:
		fmt.Fprintf(&w.buf, "<< /Type /Annot /Subtype /StrikeOut /Rect [%.2f %.2f %.2f %.2f] /C [%.3f %.3f %.3f] /P %d 0 R",
			a.x1, py1, a.x2, py2, r, g, b, pageID)
		fmt.Fprintf(&w.buf, " /QuadPoints [%.2f %.2f %.2f %.2f %.2f %.2f %.2f %.2f]",
			a.x1, py2, a.x2, py2, a.x1, py1, a.x2, py1)

	case AnnotStickyNote:
		fmt.Fprintf(&w.buf, "<< /Type /Annot /Subtype /Text /Rect [%.2f %.2f %.2f %.2f] /C [%.3f %.3f %.3f] /P %d 0 R /Open false",
			a.x1, py1, a.x2, py2, r, g, b, pageID)
		if a.text != "" {
			fmt.Fprintf(&w.buf, " /Contents (%s)", escapePDF(a.text))
		}

	case AnnotFreeText:
		fmt.Fprintf(&w.buf, "<< /Type /Annot /Subtype /FreeText /Rect [%.2f %.2f %.2f %.2f] /C [%.3f %.3f %.3f] /P %d 0 R",
			a.x1, py1, a.x2, py2, r, g, b, pageID)
		if a.text != "" {
			fmt.Fprintf(&w.buf, " /Contents (%s)", escapePDF(a.text))
		}
		fmt.Fprintf(&w.buf, " /DA (/F1 %.1f Tf %.3f %.3f %.3f rg)", a.font.Size, r, g, b)
	}

	if a.author != "" {
		fmt.Fprintf(&w.buf, " /T (%s)", escapePDF(a.author))
	}
	if a.subject != "" {
		fmt.Fprintf(&w.buf, " /Subj (%s)", escapePDF(a.subject))
	}

	fmt.Fprintf(&w.buf, " >>\n")
	w.endObject()
}

// buildTextBlockPDF renders a TextBlock to the content stream
func buildTextBlockPDF(content *bytes.Buffer, tb *TextBlock, pageH float64) {
	lines := tb.WrapLines()
	if len(lines) == 0 {
		return
	}

	fontSize := tb.font.Size
	lineHeight := fontSize * tb.lineSpacing
	colWidth := tb.columnWidth()

	r := float64(tb.font.Color.R) / 255
	g := float64(tb.font.Color.G) / 255
	b := float64(tb.font.Color.B) / 255

	if tb.columns <= 1 {
		// Single column
		for i, line := range lines {
			y := pageH - tb.y - float64(i)*lineHeight
			x := tb.x

			x = alignLineX(line, x, colWidth, fontSize, tb.alignment)

			fmt.Fprintf(content, "BT\n")
			fmt.Fprintf(content, "/F1 %.1f Tf\n", fontSize)
			fmt.Fprintf(content, "%.3f %.3f %.3f rg\n", r, g, b)
			fmt.Fprintf(content, "%.2f %.2f Td\n", x, y)
			fmt.Fprintf(content, "(%s) Tj\n", escapePDF(line))
			fmt.Fprintf(content, "ET\n")
		}
	} else {
		// Multi-column layout: distribute lines across columns
		linesPerCol := (len(lines) + tb.columns - 1) / tb.columns

		for col := 0; col < tb.columns; col++ {
			colX := tb.x + float64(col)*(colWidth+tb.columnGap)
			startLine := col * linesPerCol
			endLine := startLine + linesPerCol
			if endLine > len(lines) {
				endLine = len(lines)
			}

			for i := startLine; i < endLine; i++ {
				line := lines[i]
				lineIdx := i - startLine
				y := pageH - tb.y - float64(lineIdx)*lineHeight
				x := alignLineX(line, colX, colWidth, fontSize, tb.alignment)

				fmt.Fprintf(content, "BT\n")
				fmt.Fprintf(content, "/F1 %.1f Tf\n", fontSize)
				fmt.Fprintf(content, "%.3f %.3f %.3f rg\n", r, g, b)
				fmt.Fprintf(content, "%.2f %.2f Td\n", x, y)
				fmt.Fprintf(content, "(%s) Tj\n", escapePDF(line))
				fmt.Fprintf(content, "ET\n")
			}
		}
	}
}

// alignLineX calculates the x offset for a line based on alignment
func alignLineX(line string, baseX, colWidth, fontSize float64, alignment common.TextAlignment) float64 {
	switch alignment {
	case common.TextAlignCenter:
		lineW := measureStringWidth(line, fontSize)
		return baseX + (colWidth-lineW)/2
	case common.TextAlignRight:
		lineW := measureStringWidth(line, fontSize)
		return baseX + colWidth - lineW
	default:
		return baseX
	}
}

func buildTablePDF(content *bytes.Buffer, t *TableElement, pageH float64) {
	y := pageH - t.y

	// Border color
	r := float64(t.borderColor.R) / 255
	g := float64(t.borderColor.G) / 255
	b := float64(t.borderColor.B) / 255
	fmt.Fprintf(content, "%.3f %.3f %.3f RG\n", r, g, b)
	fmt.Fprintf(content, "0.5 w\n")

	for row := 0; row < t.rows; row++ {
		cellY := y - float64(row+1)*t.cellHeight

		// Header background
		if row == 0 && t.headerBg != nil {
			hr := float64(t.headerBg.R) / 255
			hg := float64(t.headerBg.G) / 255
			hb := float64(t.headerBg.B) / 255
			fmt.Fprintf(content, "%.3f %.3f %.3f rg\n", hr, hg, hb)
			fmt.Fprintf(content, "%.2f %.2f %.2f %.2f re f\n",
				t.x, cellY, float64(t.cols)*t.cellWidth, t.cellHeight)
		}

		for col := 0; col < t.cols; col++ {
			cellX := t.x + float64(col)*t.cellWidth

			// Cell border
			fmt.Fprintf(content, "%.2f %.2f %.2f %.2f re S\n",
				cellX, cellY, t.cellWidth, t.cellHeight)

			// Cell text
			text := t.cells[row][col]
			if text != "" {
				textR := float64(t.font.Color.R) / 255
				textG := float64(t.font.Color.G) / 255
				textB := float64(t.font.Color.B) / 255
				fmt.Fprintf(content, "BT\n")
				fmt.Fprintf(content, "/F1 %.1f Tf\n", t.font.Size)
				fmt.Fprintf(content, "%.3f %.3f %.3f rg\n", textR, textG, textB)
				fmt.Fprintf(content, "%.2f %.2f Td\n", cellX+4, cellY+5)
				fmt.Fprintf(content, "(%s) Tj\n", escapePDF(text))
				fmt.Fprintf(content, "ET\n")
			}
		}
	}
}

func buildChartPDF(content *bytes.Buffer, c *ChartElement, pageH float64) {
	chartX := c.x
	chartY := pageH - c.y - c.height
	chartW := c.width
	chartH := c.height

	titleH := 0.0
	if c.title != "" {
		titleH = 20.0
	}
	legendH := 0.0
	if c.showLegend && len(c.series) > 0 {
		legendH = 20.0
	}

	plotX := chartX + 40
	plotY := chartY + 25 + legendH
	plotW := chartW - 50
	plotH := chartH - titleH - 30 - legendH

	// Background
	if c.bgColor != nil {
		setFillColor(content, *c.bgColor)
		fmt.Fprintf(content, "%.2f %.2f %.2f %.2f re f\n", chartX, chartY, chartW, chartH)
	}

	// Title
	if c.title != "" {
		setFillColor(content, c.titleFont.Color)
		fmt.Fprintf(content, "BT\n/F1 %.1f Tf\n", c.titleFont.Size)
		fmt.Fprintf(content, "%.2f %.2f Td\n", chartX+chartW/2-float64(len(c.title))*3, chartY+chartH-titleH+5)
		fmt.Fprintf(content, "(%s) Tj\nET\n", escapePDF(c.title))
	}

	switch c.chartType {
	case ChartTypeBar:
		buildBarChart(content, c, plotX, plotY, plotW, plotH)
	case ChartTypeHorizontalBar:
		buildHorizontalBarChart(content, c, plotX, plotY, plotW, plotH)
	case ChartTypeLine:
		buildLineChart(content, c, plotX, plotY, plotW, plotH)
	case ChartTypePie:
		buildPieChart(content, c, plotX, plotY, plotW, plotH)
	case ChartTypeArea:
		buildAreaChart(content, c, plotX, plotY, plotW, plotH)
	}

	// Legend
	if c.showLegend && len(c.series) > 0 {
		lx := chartX + 40
		ly := chartY + 5
		for i, s := range c.series {
			setFillColor(content, s.Color)
			fmt.Fprintf(content, "%.2f %.2f 8 8 re f\n", lx+float64(i)*80, ly)
			fmt.Fprintf(content, "BT\n/F1 8 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
				lx+float64(i)*80+12, ly+1, escapePDF(s.Name))
		}
	}
}

func buildBarChart(content *bytes.Buffer, c *ChartElement, plotX, plotY, plotW, plotH float64) {
	maxVal := findMaxValue(c.series)
	if maxVal == 0 {
		maxVal = 1
	}

	numCats := len(c.categories)
	if numCats == 0 {
		for _, s := range c.series {
			if len(s.Values) > numCats {
				numCats = len(s.Values)
			}
		}
	}
	if numCats == 0 {
		return
	}

	numSeries := len(c.series)
	if numSeries == 0 {
		return
	}
	catWidth := plotW / float64(numCats)
	barWidth := catWidth / float64(numSeries+1)

	// Draw axes
	setStrokeColor(content, c.axisColor)
	fmt.Fprintf(content, "1 w\n")
	fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, plotY, plotX, plotY+plotH)
	fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, plotY, plotX+plotW, plotY)

	// Grid lines
	setStrokeColor(content, c.gridColor)
	fmt.Fprintf(content, "0.5 w\n")
	for i := 1; i <= 5; i++ {
		gy := plotY + plotH*float64(i)/5
		fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, gy, plotX+plotW, gy)
		label := fmt.Sprintf("%.0f", maxVal*float64(i)/5)
		fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
			plotX-30, gy-3, label)
	}

	// Bars
	for ci := 0; ci < numCats; ci++ {
		for si, s := range c.series {
			if ci >= len(s.Values) {
				continue
			}
			val := s.Values[ci]
			barH := (val / maxVal) * plotH
			bx := plotX + float64(ci)*catWidth + float64(si)*barWidth + barWidth/2
			by := plotY

			setFillColor(content, s.Color)
			fmt.Fprintf(content, "%.2f %.2f %.2f %.2f re f\n", bx, by, barWidth*0.9, barH)

			if c.showValues {
				fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
					bx+2, by+barH+3, fmt.Sprintf("%.0f", val))
			}
		}

		if ci < len(c.categories) {
			cx := plotX + float64(ci)*catWidth + catWidth/2
			fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
				cx-float64(len(c.categories[ci]))*2, plotY-12, escapePDF(c.categories[ci]))
		}
	}
}

func buildHorizontalBarChart(content *bytes.Buffer, c *ChartElement, plotX, plotY, plotW, plotH float64) {
	maxVal := findMaxValue(c.series)
	if maxVal == 0 {
		maxVal = 1
	}

	numCats := len(c.categories)
	if numCats == 0 {
		for _, s := range c.series {
			if len(s.Values) > numCats {
				numCats = len(s.Values)
			}
		}
	}
	if numCats == 0 {
		return
	}

	numSeries := len(c.series)
	if numSeries == 0 {
		return
	}
	catHeight := plotH / float64(numCats)
	barHeight := catHeight / float64(numSeries+1)

	// Axes
	setStrokeColor(content, c.axisColor)
	fmt.Fprintf(content, "1 w\n")
	fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, plotY, plotX, plotY+plotH)
	fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, plotY, plotX+plotW, plotY)

	// Grid lines
	setStrokeColor(content, c.gridColor)
	fmt.Fprintf(content, "0.5 w\n")
	for i := 1; i <= 5; i++ {
		gx := plotX + plotW*float64(i)/5
		fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", gx, plotY, gx, plotY+plotH)
		label := fmt.Sprintf("%.0f", maxVal*float64(i)/5)
		fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
			gx-8, plotY-12, label)
	}

	// Bars
	for ci := 0; ci < numCats; ci++ {
		for si, s := range c.series {
			if ci >= len(s.Values) {
				continue
			}
			val := s.Values[ci]
			barW := (val / maxVal) * plotW
			bx := plotX
			by := plotY + float64(ci)*catHeight + float64(si)*barHeight + barHeight/2

			setFillColor(content, s.Color)
			fmt.Fprintf(content, "%.2f %.2f %.2f %.2f re f\n", bx, by, barW, barHeight*0.9)

			if c.showValues {
				fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
					bx+barW+3, by+2, fmt.Sprintf("%.0f", val))
			}
		}

		if ci < len(c.categories) {
			cy := plotY + float64(ci)*catHeight + catHeight/2
			fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
				plotX-38, cy-3, escapePDF(c.categories[ci]))
		}
	}
}

func buildLineChart(content *bytes.Buffer, c *ChartElement, plotX, plotY, plotW, plotH float64) {
	maxVal := findMaxValue(c.series)
	if maxVal == 0 {
		maxVal = 1
	}

	numCats := len(c.categories)
	if numCats == 0 {
		for _, s := range c.series {
			if len(s.Values) > numCats {
				numCats = len(s.Values)
			}
		}
	}
	if numCats < 2 {
		return
	}

	// Axes
	setStrokeColor(content, c.axisColor)
	fmt.Fprintf(content, "1 w\n")
	fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, plotY, plotX, plotY+plotH)
	fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, plotY, plotX+plotW, plotY)

	// Grid lines
	setStrokeColor(content, c.gridColor)
	fmt.Fprintf(content, "0.5 w\n")
	for i := 1; i <= 5; i++ {
		gy := plotY + plotH*float64(i)/5
		fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, gy, plotX+plotW, gy)
		label := fmt.Sprintf("%.0f", maxVal*float64(i)/5)
		fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
			plotX-30, gy-3, label)
	}

	// Category labels
	for ci := 0; ci < numCats; ci++ {
		cx := plotX + float64(ci)*plotW/float64(numCats-1)
		if ci < len(c.categories) {
			fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
				cx-float64(len(c.categories[ci]))*2, plotY-12, escapePDF(c.categories[ci]))
		}
	}

	// Lines and points
	for _, s := range c.series {
		if len(s.Values) < 2 {
			continue
		}
		setStrokeColor(content, s.Color)
		fmt.Fprintf(content, "1.5 w\n")

		for i := 0; i < len(s.Values) && i < numCats; i++ {
			px := plotX + float64(i)*plotW/float64(numCats-1)
			py := plotY + (s.Values[i]/maxVal)*plotH

			if i == 0 {
				fmt.Fprintf(content, "%.2f %.2f m\n", px, py)
			} else {
				fmt.Fprintf(content, "%.2f %.2f l\n", px, py)
			}
		}
		fmt.Fprintf(content, "S\n")

		// Data points (small circles)
		setFillColor(content, s.Color)
		for i := 0; i < len(s.Values) && i < numCats; i++ {
			px := plotX + float64(i)*plotW/float64(numCats-1)
			py := plotY + (s.Values[i]/maxVal)*plotH
			r := 2.5
			// Approximate circle with 4 bezier curves
			k := r * 0.5523
			fmt.Fprintf(content, "%.2f %.2f m\n", px+r, py)
			fmt.Fprintf(content, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", px+r, py+k, px+k, py+r, px, py+r)
			fmt.Fprintf(content, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", px-k, py+r, px-r, py+k, px-r, py)
			fmt.Fprintf(content, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", px-r, py-k, px-k, py-r, px, py-r)
			fmt.Fprintf(content, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", px+k, py-r, px+r, py-k, px+r, py)
			fmt.Fprintf(content, "f\n")
		}
	}
}

func buildAreaChart(content *bytes.Buffer, c *ChartElement, plotX, plotY, plotW, plotH float64) {
	maxVal := findMaxValue(c.series)
	if maxVal == 0 {
		maxVal = 1
	}

	numCats := len(c.categories)
	if numCats == 0 {
		for _, s := range c.series {
			if len(s.Values) > numCats {
				numCats = len(s.Values)
			}
		}
	}
	if numCats < 2 {
		return
	}

	// Axes
	setStrokeColor(content, c.axisColor)
	fmt.Fprintf(content, "1 w\n")
	fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, plotY, plotX, plotY+plotH)
	fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, plotY, plotX+plotW, plotY)

	// Grid lines
	setStrokeColor(content, c.gridColor)
	fmt.Fprintf(content, "0.5 w\n")
	for i := 1; i <= 5; i++ {
		gy := plotY + plotH*float64(i)/5
		fmt.Fprintf(content, "%.2f %.2f m %.2f %.2f l S\n", plotX, gy, plotX+plotW, gy)
		label := fmt.Sprintf("%.0f", maxVal*float64(i)/5)
		fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
			plotX-30, gy-3, label)
	}

	// Category labels
	for ci := 0; ci < numCats; ci++ {
		cx := plotX + float64(ci)*plotW/float64(numCats-1)
		if ci < len(c.categories) {
			fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
				cx-float64(len(c.categories[ci]))*2, plotY-12, escapePDF(c.categories[ci]))
		}
	}

	// Filled areas and lines
	for _, s := range c.series {
		if len(s.Values) < 2 {
			continue
		}

		// Semi-transparent fill
		fmt.Fprintf(content, "%.3f %.3f %.3f rg\n",
			float64(s.Color.R)/255*0.7+0.3,
			float64(s.Color.G)/255*0.7+0.3,
			float64(s.Color.B)/255*0.7+0.3)

		// Build filled path: start at baseline, go up through points, back to baseline
		firstX := plotX
		fmt.Fprintf(content, "%.2f %.2f m\n", firstX, plotY)
		for i := 0; i < len(s.Values) && i < numCats; i++ {
			px := plotX + float64(i)*plotW/float64(numCats-1)
			py := plotY + (s.Values[i]/maxVal)*plotH
			fmt.Fprintf(content, "%.2f %.2f l\n", px, py)
		}
		lastIdx := len(s.Values) - 1
		if lastIdx >= numCats {
			lastIdx = numCats - 1
		}
		lastX := plotX + float64(lastIdx)*plotW/float64(numCats-1)
		fmt.Fprintf(content, "%.2f %.2f l\n", lastX, plotY)
		fmt.Fprintf(content, "f\n")

		// Stroke the line on top
		setStrokeColor(content, s.Color)
		fmt.Fprintf(content, "1.5 w\n")
		for i := 0; i < len(s.Values) && i < numCats; i++ {
			px := plotX + float64(i)*plotW/float64(numCats-1)
			py := plotY + (s.Values[i]/maxVal)*plotH
			if i == 0 {
				fmt.Fprintf(content, "%.2f %.2f m\n", px, py)
			} else {
				fmt.Fprintf(content, "%.2f %.2f l\n", px, py)
			}
		}
		fmt.Fprintf(content, "S\n")
	}
}

func buildPieChart(content *bytes.Buffer, c *ChartElement, plotX, plotY, plotW, plotH float64) {
	if len(c.series) == 0 || len(c.series[0].Values) == 0 {
		return
	}

	values := c.series[0].Values
	total := 0.0
	for _, v := range values {
		total += v
	}
	if total == 0 {
		return
	}

	// Default slice colors
	sliceColors := []common.Color{
		common.Blue,
		common.Red,
		common.Green,
		common.Orange,
		common.Purple,
		common.Yellow,
		{R: 0, G: 191, B: 255, A: 255},   // deep sky blue
		{R: 255, G: 105, B: 180, A: 255}, // hot pink
		{R: 34, G: 139, B: 34, A: 255},   // forest green
		{R: 255, G: 215, B: 0, A: 255},   // gold
	}

	// Center and radius
	cx := plotX + plotW/2
	cy := plotY + plotH/2
	radius := math.Min(plotW, plotH) / 2 * 0.8

	startAngle := 0.0
	for i, v := range values {
		sweepAngle := (v / total) * 2 * math.Pi
		endAngle := startAngle + sweepAngle

		color := sliceColors[i%len(sliceColors)]
		setFillColor(content, color)

		// Draw pie slice using move to center, line to arc start, bezier arcs, close
		drawPieSlice(content, cx, cy, radius, startAngle, endAngle)

		startAngle = endAngle
	}

	// Legend for pie (category names with colors)
	if c.showLegend && len(c.categories) > 0 {
		lx := plotX + plotW + 5
		for i := 0; i < len(values) && i < len(c.categories); i++ {
			ly := plotY + plotH - float64(i)*14 - 10
			color := sliceColors[i%len(sliceColors)]
			setFillColor(content, color)
			fmt.Fprintf(content, "%.2f %.2f 8 8 re f\n", lx, ly)
			pct := fmt.Sprintf("%s (%.0f%%)", escapePDF(c.categories[i]), pctVal(values[i], total))
			fmt.Fprintf(content, "BT\n/F1 7 Tf\n0 0 0 rg\n%.2f %.2f Td\n(%s) Tj\nET\n",
				lx+12, ly+1, pct)
		}
	}
}

// pctVal calculates percentage
func pctVal(val, total float64) float64 {
	return val / total * 100
}

// drawPieSlice draws a filled pie slice using bezier curve approximation
func drawPieSlice(content *bytes.Buffer, cx, cy, r, startAngle, endAngle float64) {
	// Move to center
	fmt.Fprintf(content, "%.2f %.2f m\n", cx, cy)

	// Line to start of arc
	sx := cx + r*math.Cos(startAngle)
	sy := cy + r*math.Sin(startAngle)
	fmt.Fprintf(content, "%.2f %.2f l\n", sx, sy)

	// Draw arc segments (split into segments of max 90 degrees)
	angle := startAngle
	remaining := endAngle - startAngle
	for remaining > 0 {
		segment := math.Min(remaining, math.Pi/2)
		drawArcSegment(content, cx, cy, r, angle, angle+segment)
		angle += segment
		remaining -= segment
	}

	// Close path and fill
	fmt.Fprintf(content, "f\n")
}

// drawArcSegment draws a single arc segment using a cubic bezier approximation
func drawArcSegment(content *bytes.Buffer, cx, cy, r, startAngle, endAngle float64) {
	halfAngle := (endAngle - startAngle) / 2
	midAngle := startAngle + halfAngle

	// Control point distance for bezier approximation of arc
	alpha := 4.0 / 3.0 * math.Tan(halfAngle) * r

	// Start point
	x0 := cx + r*math.Cos(startAngle)
	y0 := cy + r*math.Sin(startAngle)

	// End point
	x3 := cx + r*math.Cos(endAngle)
	y3 := cy + r*math.Sin(endAngle)

	// Control points
	_ = midAngle
	x1 := x0 - alpha*math.Sin(startAngle)
	y1 := y0 + alpha*math.Cos(startAngle)
	x2 := x3 + alpha*math.Sin(endAngle)
	y2 := y3 - alpha*math.Cos(endAngle)

	fmt.Fprintf(content, "%.2f %.2f %.2f %.2f %.2f %.2f c\n", x1, y1, x2, y2, x3, y3)
}

func findMaxValue(series []ChartSeries) float64 {
	max := 0.0
	for _, s := range series {
		for _, v := range s.Values {
			if v > max {
				max = v
			}
		}
	}
	return max
}

func setFillColor(content *bytes.Buffer, c common.Color) {
	fmt.Fprintf(content, "%.3f %.3f %.3f rg\n", float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
}

func setStrokeColor(content *bytes.Buffer, c common.Color) {
	fmt.Fprintf(content, "%.3f %.3f %.3f RG\n", float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
}

func escapePDF(s string) string {
	var result []byte
	for _, c := range []byte(s) {
		switch c {
		case '(', ')', '\\':
			result = append(result, '\\', c)
		default:
			result = append(result, c)
		}
	}
	return string(result)
}
