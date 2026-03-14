package pdf

import (
	"bytes"
	"fmt"
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

		// Page object
		w.startObject(pid.pageID)
		fmt.Fprintf(&w.buf, "<< /Type /Page /Parent %d 0 R /MediaBox [0 0 %.2f %.2f] /Contents %d 0 R /Resources << /Font << /F1 %d 0 R >> >> >>\n",
			pagesID,
			page.size.Width.Points(),
			page.size.Height.Points(),
			pid.contentID,
			fontID,
		)
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

	// Catalog object
	w.startObject(catalogID)
	fmt.Fprintf(&w.buf, "<< /Type /Catalog /Pages %d 0 R >>\n", pagesID)
	w.endObject()

	// Cross-reference table
	xrefOffset := w.buf.Len()
	fmt.Fprintf(&w.buf, "xref\n")
	fmt.Fprintf(&w.buf, "0 %d\n", w.nextObj)
	fmt.Fprintf(&w.buf, "0000000000 65535 f \n")
	for _, offset := range w.objects {
		fmt.Fprintf(&w.buf, "%010d 00000 n \n", offset)
	}

	// Trailer
	fmt.Fprintf(&w.buf, "trailer\n")
	fmt.Fprintf(&w.buf, "<< /Size %d /Root %d 0 R >>\n", w.nextObj, catalogID)
	fmt.Fprintf(&w.buf, "startxref\n%d\n%%%%EOF\n", xrefOffset)

	return w.buf.Bytes(), nil
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
