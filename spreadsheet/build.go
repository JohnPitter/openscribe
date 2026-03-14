package spreadsheet

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"

	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/internal/xmlutil"
)

func (wb *Workbook) build() error {
	wb.pkg = packaging.NewPackage()

	// Reset shared strings for fresh build
	wb.sharedStrings = nil
	wb.stringIndex = make(map[string]int)

	// Build each sheet
	wbRels := packaging.NewRelationships()
	for i, sheet := range wb.sheets {
		sheetPath := fmt.Sprintf("xl/worksheets/sheet%d.xml", i+1)
		data, err := wb.buildSheetXML(sheet)
		if err != nil {
			return fmt.Errorf("build sheet %s: %w", sheet.name, err)
		}
		wb.pkg.AddFile(sheetPath, data)
		wbRels.Add(packaging.RelTypeWorksheet, fmt.Sprintf("worksheets/sheet%d.xml", i+1))
	}

	// Shared strings
	if len(wb.sharedStrings) > 0 {
		ssData, err := wb.buildSharedStringsXML()
		if err != nil {
			return fmt.Errorf("build shared strings: %w", err)
		}
		wb.pkg.AddFile("xl/sharedStrings.xml", ssData)
		wbRels.Add(packaging.RelTypeSharedStrings, "sharedStrings.xml")
	}

	// Workbook XML
	wbData, err := wb.buildWorkbookXML()
	if err != nil {
		return fmt.Errorf("build workbook: %w", err)
	}
	wb.pkg.AddFile("xl/workbook.xml", wbData)

	// Workbook relationships
	wbRelsData, err := wbRels.Marshal()
	if err != nil {
		return fmt.Errorf("marshal wb rels: %w", err)
	}
	wb.pkg.AddFile("xl/_rels/workbook.xml.rels", wbRelsData)

	// Top-level relationships
	topRels := packaging.NewRelationships()
	topRels.Add(packaging.RelTypeOfficeDocument, "xl/workbook.xml")
	topRelsData, err := topRels.Marshal()
	if err != nil {
		return fmt.Errorf("marshal top rels: %w", err)
	}
	wb.pkg.AddFile("_rels/.rels", topRelsData)

	// Content types
	ct := packaging.NewContentTypes()
	ct.AddOverride("/xl/workbook.xml", packaging.ContentTypeXlsx)
	for i := range wb.sheets {
		ct.AddOverride(fmt.Sprintf("/xl/worksheets/sheet%d.xml", i+1), packaging.ContentTypeWorksheet)
	}
	if len(wb.sharedStrings) > 0 {
		ct.AddOverride("/xl/sharedStrings.xml", packaging.ContentTypeSharedStrings)
	}
	ctData, err := ct.Marshal()
	if err != nil {
		return fmt.Errorf("marshal content types: %w", err)
	}
	wb.pkg.AddFile("[Content_Types].xml", ctData)

	return nil
}

// XML types for workbook

type xmlWorkbook struct {
	XMLName xml.Name  `xml:"workbook"`
	Xmlns   string    `xml:"xmlns,attr"`
	XmlnsR  string    `xml:"xmlns:r,attr"`
	Sheets  xmlSheets `xml:"sheets"`
}

type xmlSheets struct {
	Sheet []xmlSheetRef `xml:"sheet"`
}

type xmlSheetRef struct {
	Name    string `xml:"name,attr"`
	SheetID string `xml:"sheetId,attr"`
	RID     string `xml:"r:id,attr"`
}

func (wb *Workbook) buildWorkbookXML() ([]byte, error) {
	xwb := xmlWorkbook{
		Xmlns:  "http://schemas.openxmlformats.org/spreadsheetml/2006/main",
		XmlnsR: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
	}
	for i, sheet := range wb.sheets {
		xwb.Sheets.Sheet = append(xwb.Sheets.Sheet, xmlSheetRef{
			Name:    sheet.name,
			SheetID: strconv.Itoa(i + 1),
			RID:     fmt.Sprintf("rId%d", i+1),
		})
	}
	return xmlutil.MarshalXML(xwb)
}

// XML types for worksheets

type xmlWorksheet struct {
	XMLName    xml.Name       `xml:"worksheet"`
	Xmlns      string         `xml:"xmlns,attr"`
	SheetData  xmlSheetData   `xml:"sheetData"`
	MergeCells *xmlMergeCells `xml:"mergeCells,omitempty"`
}

type xmlSheetData struct {
	Rows []xmlRow `xml:"row"`
}

type xmlRow struct {
	R     string    `xml:"r,attr"`
	Cells []xmlCell `xml:"c"`
}

type xmlCell struct {
	R string `xml:"r,attr"`
	T string `xml:"t,attr,omitempty"`
	V string `xml:"v,omitempty"`
	F string `xml:"f,omitempty"`
}

type xmlMergeCells struct {
	Count     string          `xml:"count,attr"`
	MergeCell []xmlMergeCell  `xml:"mergeCell"`
}

type xmlMergeCell struct {
	Ref string `xml:"ref,attr"`
}

func (wb *Workbook) buildSheetXML(sheet *Sheet) ([]byte, error) {
	ws := xmlWorksheet{
		Xmlns: "http://schemas.openxmlformats.org/spreadsheetml/2006/main",
	}

	// Sort rows by index
	rowIndices := make([]int, 0, len(sheet.rows))
	for idx := range sheet.rows {
		rowIndices = append(rowIndices, idx)
	}
	sort.Ints(rowIndices)

	for _, rowIdx := range rowIndices {
		row := sheet.rows[rowIdx]
		xr := xmlRow{R: strconv.Itoa(rowIdx)}

		// Sort cells by column
		colIndices := make([]int, 0, len(row.cells))
		for col := range row.cells {
			colIndices = append(colIndices, col)
		}
		sort.Ints(colIndices)

		for _, colIdx := range colIndices {
			cell := row.cells[colIdx]
			xc := xmlCell{R: CellRef(rowIdx, colIdx)}

			switch cell.cellType {
			case CellTypeString:
				idx := wb.addSharedString(cell.strVal)
				xc.T = "s"
				xc.V = strconv.Itoa(idx)
			case CellTypeNumber:
				xc.V = strconv.FormatFloat(cell.numVal, 'f', -1, 64)
			case CellTypeBoolean:
				xc.T = "b"
				if cell.boolVal {
					xc.V = "1"
				} else {
					xc.V = "0"
				}
			case CellTypeFormula:
				xc.F = cell.formula
			}

			xr.Cells = append(xr.Cells, xc)
		}

		ws.SheetData.Rows = append(ws.SheetData.Rows, xr)
	}

	// Merged cells
	if len(sheet.mergedCells) > 0 {
		mc := &xmlMergeCells{
			Count: strconv.Itoa(len(sheet.mergedCells)),
		}
		for _, m := range sheet.mergedCells {
			mc.MergeCell = append(mc.MergeCell, xmlMergeCell{
				Ref: fmt.Sprintf("%s:%s", CellRef(m.StartRow, m.StartCol), CellRef(m.EndRow, m.EndCol)),
			})
		}
		ws.MergeCells = mc
	}

	return xmlutil.MarshalXML(ws)
}

// Shared strings XML

type xmlSST struct {
	XMLName xml.Name `xml:"sst"`
	Xmlns   string   `xml:"xmlns,attr"`
	Count   string   `xml:"count,attr"`
	SI      []xmlSI  `xml:"si"`
}

type xmlSI struct {
	T string `xml:"t"`
}

func (wb *Workbook) buildSharedStringsXML() ([]byte, error) {
	sst := xmlSST{
		Xmlns: "http://schemas.openxmlformats.org/spreadsheetml/2006/main",
		Count: strconv.Itoa(len(wb.sharedStrings)),
	}
	for _, s := range wb.sharedStrings {
		sst.SI = append(sst.SI, xmlSI{T: s})
	}
	return xmlutil.MarshalXML(sst)
}

// parseSharedStrings parses the shared strings table
func (wb *Workbook) parseSharedStrings(data []byte) {
	var sst xmlSST
	if err := xmlutil.UnmarshalXML(data, &sst); err != nil {
		return
	}
	for _, si := range sst.SI {
		wb.addSharedString(si.T)
	}
}

// parseWorkbook parses workbook.xml to get sheet names
func (wb *Workbook) parseWorkbook(data []byte) {
	var xwb xmlWorkbook
	if err := xmlutil.UnmarshalXML(data, &xwb); err != nil {
		return
	}
	for i, s := range xwb.Sheets.Sheet {
		sheet := newSheet(wb, s.Name, i+1)
		// Try to parse existing sheet data
		sheetPath := fmt.Sprintf("xl/worksheets/sheet%d.xml", i+1)
		if sheetData, ok := wb.pkg.GetFile(sheetPath); ok {
			wb.parseSheetData(sheet, sheetData)
		}
		wb.sheets = append(wb.sheets, sheet)
	}
}

// parseSheetData parses a worksheet
func (wb *Workbook) parseSheetData(sheet *Sheet, data []byte) {
	var ws xmlWorksheet
	if err := xmlutil.UnmarshalXML(data, &ws); err != nil {
		return
	}
	for _, xr := range ws.SheetData.Rows {
		rowIdx, _ := strconv.Atoi(xr.R)
		for _, xc := range xr.Cells {
			// Parse cell reference to get column
			col := parseCellRefCol(xc.R)
			cell := sheet.Cell(rowIdx, col)

			switch xc.T {
			case "s":
				// Shared string
				idx, _ := strconv.Atoi(xc.V)
				if idx < len(wb.sharedStrings) {
					cell.SetString(wb.sharedStrings[idx])
				}
			case "b":
				cell.SetBool(xc.V == "1")
			default:
				if xc.F != "" {
					cell.SetFormula(xc.F)
				} else if xc.V != "" {
					if n, err := strconv.ParseFloat(xc.V, 64); err == nil {
						cell.SetNumber(n)
					} else {
						cell.SetString(xc.V)
					}
				}
			}
		}
	}
}

// parseCellRefCol extracts the column number from a cell reference like "B3"
func parseCellRefCol(ref string) int {
	col := 0
	for _, c := range ref {
		if c >= 'A' && c <= 'Z' {
			col = col*26 + int(c-'A'+1)
		} else {
			break
		}
	}
	return col
}
