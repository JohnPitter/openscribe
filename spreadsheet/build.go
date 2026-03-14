package spreadsheet

import (
	"bytes"
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
	ct := packaging.NewContentTypes()
	chartGlobalIdx := 0

	commentIdx := 0
	for i, sheet := range wb.sheets {
		sheetPath := fmt.Sprintf("xl/worksheets/sheet%d.xml", i+1)
		data, err := wb.buildSheetXML(sheet)
		if err != nil {
			return fmt.Errorf("build sheet %s: %w", sheet.name, err)
		}
		wb.pkg.AddFile(sheetPath, data)
		wbRels.Add(packaging.RelTypeWorksheet, fmt.Sprintf("worksheets/sheet%d.xml", i+1))
		ct.AddOverride(fmt.Sprintf("/xl/worksheets/sheet%d.xml", i+1), packaging.ContentTypeWorksheet)

		// Track if we need sheet rels
		hasCharts := len(sheet.charts) > 0
		hasComments := len(sheet.comments) > 0
		needsSheetRels := hasCharts || hasComments

		var sheetRels *packaging.Relationships
		if needsSheetRels {
			sheetRels = packaging.NewRelationships()
		}

		// Charts for this sheet
		if hasCharts {
			drawingIdx := i + 1

			sheetRels.Add(packaging.RelTypeDrawing, fmt.Sprintf("../drawings/drawing%d.xml", drawingIdx))

			// Drawing relationships: link drawing -> charts
			drawingRels := packaging.NewRelationships()
			for j, chart := range sheet.charts {
				chartGlobalIdx++
				chartPath := fmt.Sprintf("xl/charts/chart%d.xml", chartGlobalIdx)

				chartData, cErr := wb.buildChartXML(chart, chartGlobalIdx)
				if cErr != nil {
					return fmt.Errorf("build chart %d: %w", chartGlobalIdx, cErr)
				}
				wb.pkg.AddFile(chartPath, chartData)
				ct.AddOverride(fmt.Sprintf("/xl/charts/chart%d.xml", chartGlobalIdx), packaging.ContentTypeChart)

				_ = j
				drawingRels.Add(packaging.RelTypeChart, fmt.Sprintf("../charts/chart%d.xml", chartGlobalIdx))
			}

			drawingRelsData, drErr := drawingRels.Marshal()
			if drErr != nil {
				return fmt.Errorf("marshal drawing rels: %w", drErr)
			}
			wb.pkg.AddFile(fmt.Sprintf("xl/drawings/_rels/drawing%d.xml.rels", drawingIdx), drawingRelsData)

			// Drawing XML
			drawingData, dErr := wb.buildDrawingXML(sheet.charts)
			if dErr != nil {
				return fmt.Errorf("build drawing %d: %w", drawingIdx, dErr)
			}
			wb.pkg.AddFile(fmt.Sprintf("xl/drawings/drawing%d.xml", drawingIdx), drawingData)
			ct.AddOverride(fmt.Sprintf("/xl/drawings/drawing%d.xml", drawingIdx), packaging.ContentTypeDrawing)
		}

		// Comments for this sheet
		if hasComments {
			commentIdx++
			commentsXML := buildCommentsXML(sheet.comments)
			wb.pkg.AddFile(fmt.Sprintf("xl/comments%d.xml", commentIdx), []byte(commentsXML))
			ct.AddOverride(fmt.Sprintf("/xl/comments%d.xml", commentIdx), packaging.ContentTypeSpreadsheetComments)

			// VML drawing for comment shapes
			vmlXML := buildVMLDrawingXML(sheet.comments)
			wb.pkg.AddFile(fmt.Sprintf("xl/drawings/vmlDrawing%d.vml", commentIdx), []byte(vmlXML))

			sheetRels.Add(packaging.RelTypeComments, fmt.Sprintf("../comments%d.xml", commentIdx))
			sheetRels.Add(packaging.RelTypeVMLDrawing, fmt.Sprintf("../drawings/vmlDrawing%d.vml", commentIdx))
		}

		// Write sheet rels if needed
		if needsSheetRels {
			sheetRelsData, rErr := sheetRels.Marshal()
			if rErr != nil {
				return fmt.Errorf("marshal sheet rels: %w", rErr)
			}
			wb.pkg.AddFile(fmt.Sprintf("xl/worksheets/_rels/sheet%d.xml.rels", i+1), sheetRelsData)
		}
	}

	// Shared strings
	if len(wb.sharedStrings) > 0 {
		ssData, err := wb.buildSharedStringsXML()
		if err != nil {
			return fmt.Errorf("build shared strings: %w", err)
		}
		wb.pkg.AddFile("xl/sharedStrings.xml", ssData)
		wbRels.Add(packaging.RelTypeSharedStrings, "sharedStrings.xml")
		ct.AddOverride("/xl/sharedStrings.xml", packaging.ContentTypeSharedStrings)
	}

	// Workbook XML
	wbData, err := wb.buildWorkbookXML()
	if err != nil {
		return fmt.Errorf("build workbook: %w", err)
	}
	wb.pkg.AddFile("xl/workbook.xml", wbData)
	ct.AddOverride("/xl/workbook.xml", packaging.ContentTypeXlsx)

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
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)

	// Sheets
	buf.WriteString(`<sheets>`)
	for i, sheet := range wb.sheets {
		fmt.Fprintf(&buf, `<sheet name="%s" sheetId="%d" r:id="rId%d"/>`,
			escapeXMLAttr(sheet.name), i+1, i+1)
	}
	buf.WriteString(`</sheets>`)

	// Defined names: named ranges + print areas + print titles
	hasDefinedNames := len(wb.namedRanges) > 0
	if !hasDefinedNames {
		for _, sheet := range wb.sheets {
			if sheet.printSettings != nil && (sheet.printSettings.printArea != nil || sheet.printSettings.repeatRows != "" || sheet.printSettings.repeatCols != "") {
				hasDefinedNames = true
				break
			}
		}
	}

	if hasDefinedNames {
		buf.WriteString(`<definedNames>`)

		// User-defined named ranges
		for _, nr := range wb.namedRanges {
			ref := fmt.Sprintf("'%s'!%s", nr.sheetName, nr.cellRange)
			if nr.sheetID >= 0 {
				fmt.Fprintf(&buf, `<definedName name="%s" localSheetId="%d">%s</definedName>`,
					escapeXMLAttr(nr.name), nr.sheetID, escapeXMLText(ref))
			} else {
				fmt.Fprintf(&buf, `<definedName name="%s">%s</definedName>`,
					escapeXMLAttr(nr.name), escapeXMLText(ref))
			}
		}

		// Print areas and titles
		for i, sheet := range wb.sheets {
			if sheet.printSettings != nil {
				buf.WriteString(buildPrintAreaDefinedName(sheet.name, i, sheet.printSettings.printArea))
				buf.WriteString(buildPrintTitlesDefinedName(sheet.name, i, sheet.printSettings.repeatRows, sheet.printSettings.repeatCols))
			}
		}

		buf.WriteString(`</definedNames>`)
	}

	buf.WriteString(`</workbook>`)
	return buf.Bytes(), nil
}

// XML types for worksheets

type xmlWorksheet struct {
	XMLName               xml.Name                   `xml:"worksheet"`
	Xmlns                 string                     `xml:"xmlns,attr"`
	SheetData             xmlSheetData               `xml:"sheetData"`
	MergeCells            *xmlMergeCells             `xml:"mergeCells,omitempty"`
	ConditionalFormatting []xmlConditionalFormatting `xml:"conditionalFormatting,omitempty"`
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
	Count     string         `xml:"count,attr"`
	MergeCell []xmlMergeCell `xml:"mergeCell"`
}

type xmlMergeCell struct {
	Ref string `xml:"ref,attr"`
}

// Conditional formatting XML types

type xmlConditionalFormatting struct {
	SQRef  string      `xml:"sqref,attr"`
	CFRule []xmlCFRule `xml:"cfRule"`
}

type xmlCFRule struct {
	Type       string           `xml:"type,attr"`
	Operator   string           `xml:"operator,attr,omitempty"`
	Priority   string           `xml:"priority,attr"`
	Formula    []string         `xml:"formula,omitempty"`
	DXF        *xmlCFDXF        `xml:"dxf,omitempty"`
	ColorScale *xmlCFColorScale `xml:"colorScale,omitempty"`
	DataBar    *xmlCFDataBar    `xml:"dataBar,omitempty"`
}

type xmlCFDXF struct {
	Font *xmlCFFont `xml:"font,omitempty"`
	Fill *xmlCFFill `xml:"fill,omitempty"`
}

type xmlCFFont struct {
	Bold   *xmlCFBoolVal `xml:"b,omitempty"`
	Italic *xmlCFBoolVal `xml:"i,omitempty"`
	Color  *xmlCFColor   `xml:"color,omitempty"`
}

type xmlCFBoolVal struct {
	Val string `xml:"val,attr,omitempty"`
}

type xmlCFColor struct {
	RGB string `xml:"rgb,attr"`
}

type xmlCFFill struct {
	PatternFill xmlCFPatternFill `xml:"patternFill"`
}

type xmlCFPatternFill struct {
	BgColor xmlCFColor `xml:"bgColor"`
}

type xmlCFColorScale struct {
	CFVOs  []xmlCFVO    `xml:"cfvo"`
	Colors []xmlCFColor `xml:"color"`
}

type xmlCFVO struct {
	Type string `xml:"type,attr"`
}

type xmlCFDataBar struct {
	CFVOs []xmlCFVO  `xml:"cfvo"`
	Color xmlCFColor `xml:"color"`
}

func (wb *Workbook) buildSheetXML(sheet *Sheet) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">`)

	// Sheet views (freeze panes)
	buf.WriteString(buildSheetViewsXML(sheet.freezePane))

	// Sheet data
	buf.WriteString(`<sheetData>`)

	// Sort rows by index
	rowIndices := make([]int, 0, len(sheet.rows))
	for idx := range sheet.rows {
		rowIndices = append(rowIndices, idx)
	}
	sort.Ints(rowIndices)

	for _, rowIdx := range rowIndices {
		row := sheet.rows[rowIdx]
		fmt.Fprintf(&buf, `<row r="%d">`, rowIdx)

		// Sort cells by column
		colIndices := make([]int, 0, len(row.cells))
		for col := range row.cells {
			colIndices = append(colIndices, col)
		}
		sort.Ints(colIndices)

		for _, colIdx := range colIndices {
			cell := row.cells[colIdx]
			ref := CellRef(rowIdx, colIdx)

			switch cell.cellType {
			case CellTypeString:
				idx := wb.addSharedString(cell.strVal)
				fmt.Fprintf(&buf, `<c r="%s" t="s"><v>%d</v></c>`, ref, idx)
			case CellTypeNumber:
				fmt.Fprintf(&buf, `<c r="%s"><v>%s</v></c>`, ref, strconv.FormatFloat(cell.numVal, 'f', -1, 64))
			case CellTypeBoolean:
				bv := "0"
				if cell.boolVal {
					bv = "1"
				}
				fmt.Fprintf(&buf, `<c r="%s" t="b"><v>%s</v></c>`, ref, bv)
			case CellTypeFormula:
				fmt.Fprintf(&buf, `<c r="%s"><f>%s</f></c>`, ref, escapeXMLText(cell.formula))
			default:
				fmt.Fprintf(&buf, `<c r="%s"/>`, ref)
			}
		}

		buf.WriteString(`</row>`)
	}

	buf.WriteString(`</sheetData>`)

	// Sheet protection
	buf.WriteString(buildProtectionXML(sheet.protection))

	// Auto filter
	buf.WriteString(buildAutoFilterXML(sheet.autoFilter))

	// Merged cells
	if len(sheet.mergedCells) > 0 {
		fmt.Fprintf(&buf, `<mergeCells count="%d">`, len(sheet.mergedCells))
		for _, m := range sheet.mergedCells {
			fmt.Fprintf(&buf, `<mergeCell ref="%s:%s"/>`, CellRef(m.StartRow, m.StartCol), CellRef(m.EndRow, m.EndCol))
		}
		buf.WriteString(`</mergeCells>`)
	}

	// Conditional formatting
	if len(sheet.conditionalFormats) > 0 {
		cfs := buildConditionalFormattingXML(sheet.conditionalFormats)
		for _, cf := range cfs {
			cfData, err := xmlutil.MarshalXMLFragment(cf)
			if err != nil {
				return nil, fmt.Errorf("marshal conditional formatting: %w", err)
			}
			buf.Write(cfData)
		}
	}

	// Data validations
	buf.WriteString(buildValidationsXML(sheet.validations))

	// Page setup / print settings
	buf.WriteString(buildPageSetupXML(sheet.printSettings))

	buf.WriteString(`</worksheet>`)
	return buf.Bytes(), nil
}

// buildConditionalFormattingXML converts conditional formats to XML structs
func buildConditionalFormattingXML(formats []*ConditionalFormat) []xmlConditionalFormatting {
	// Group by cell range
	groups := make(map[string][]*ConditionalFormat)
	var order []string
	for _, cf := range formats {
		if _, exists := groups[cf.cellRange]; !exists {
			order = append(order, cf.cellRange)
		}
		groups[cf.cellRange] = append(groups[cf.cellRange], cf)
	}

	var result []xmlConditionalFormatting
	priority := 1
	for _, rangeRef := range order {
		cfs := groups[rangeRef]
		xcf := xmlConditionalFormatting{SQRef: rangeRef}

		for _, cf := range cfs {
			rule := xmlCFRule{
				Priority: strconv.Itoa(priority),
			}
			priority++

			switch cf.condType {
			case ConditionColorScale:
				rule.Type = "colorScale"
				cs := &xmlCFColorScale{
					CFVOs: []xmlCFVO{
						{Type: "min"},
						{Type: "max"},
					},
				}
				minC := "FF000000"
				maxC := "FFFFFFFF"
				if cf.minColor != nil {
					minC = fmt.Sprintf("FF%02X%02X%02X", cf.minColor.R, cf.minColor.G, cf.minColor.B)
				}
				if cf.maxColor != nil {
					maxC = fmt.Sprintf("FF%02X%02X%02X", cf.maxColor.R, cf.maxColor.G, cf.maxColor.B)
				}
				cs.Colors = []xmlCFColor{
					{RGB: minC},
					{RGB: maxC},
				}
				rule.ColorScale = cs

			case ConditionDataBar:
				rule.Type = "dataBar"
				barC := "FF638EC6"
				if cf.barColor != nil {
					barC = fmt.Sprintf("FF%02X%02X%02X", cf.barColor.R, cf.barColor.G, cf.barColor.B)
				}
				rule.DataBar = &xmlCFDataBar{
					CFVOs: []xmlCFVO{
						{Type: "min"},
						{Type: "max"},
					},
					Color: xmlCFColor{RGB: barC},
				}

			default:
				rule.Type = "cellIs"
				switch cf.condType {
				case ConditionGreaterThan:
					rule.Operator = "greaterThan"
				case ConditionLessThan:
					rule.Operator = "lessThan"
				case ConditionEqual:
					rule.Operator = "equal"
				case ConditionNotEqual:
					rule.Operator = "notEqual"
				case ConditionBetween:
					rule.Operator = "between"
				case ConditionContains:
					rule.Type = "containsText"
					rule.Operator = "containsText"
				case ConditionBeginsWith:
					rule.Type = "beginsWith"
					rule.Operator = "beginsWith"
				case ConditionEndsWith:
					rule.Type = "endsWith"
					rule.Operator = "endsWith"
				case ConditionTop10:
					rule.Type = "top10"
					rule.Operator = ""
				case ConditionAboveAverage:
					rule.Type = "aboveAverage"
					rule.Operator = ""
				case ConditionBelowAverage:
					rule.Type = "belowAverage"
					rule.Operator = ""
				case ConditionDuplicate:
					rule.Type = "duplicateValues"
					rule.Operator = ""
				case ConditionUnique:
					rule.Type = "uniqueValues"
					rule.Operator = ""
				}

				if cf.value != "" {
					rule.Formula = append(rule.Formula, cf.value)
				}
				if cf.condType == ConditionBetween && cf.value2 != "" {
					rule.Formula = append(rule.Formula, cf.value2)
				}

				// Build DXF (differential formatting)
				dxf := &xmlCFDXF{}
				hasDXF := false

				if cf.bgColor != nil {
					rgb := fmt.Sprintf("FF%02X%02X%02X", cf.bgColor.R, cf.bgColor.G, cf.bgColor.B)
					dxf.Fill = &xmlCFFill{
						PatternFill: xmlCFPatternFill{
							BgColor: xmlCFColor{RGB: rgb},
						},
					}
					hasDXF = true
				}

				if cf.fontColor != nil || cf.bold || cf.italic {
					font := &xmlCFFont{}
					if cf.fontColor != nil {
						rgb := fmt.Sprintf("FF%02X%02X%02X", cf.fontColor.R, cf.fontColor.G, cf.fontColor.B)
						font.Color = &xmlCFColor{RGB: rgb}
					}
					if cf.bold {
						font.Bold = &xmlCFBoolVal{}
					}
					if cf.italic {
						font.Italic = &xmlCFBoolVal{}
					}
					dxf.Font = font
					hasDXF = true
				}

				if hasDXF {
					rule.DXF = dxf
				}
			}

			xcf.CFRule = append(xcf.CFRule, rule)
		}

		result = append(result, xcf)
	}
	return result
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

// buildChartXML generates DrawingML chart XML
func (wb *Workbook) buildChartXML(chart *Chart, chartIndex int) ([]byte, error) {
	chartTypeTag := "c:barChart"
	grouping := "clustered"
	switch chart.chartType {
	case ChartTypeLine:
		chartTypeTag = "c:lineChart"
		grouping = "standard"
	case ChartTypePie:
		chartTypeTag = "c:pieChart"
		grouping = ""
	case ChartTypeArea:
		chartTypeTag = "c:areaChart"
		grouping = "standard"
	case ChartTypeScatter:
		chartTypeTag = "c:scatterChart"
		grouping = ""
	case ChartTypeColumn:
		chartTypeTag = "c:barChart"
		grouping = "clustered"
	case ChartTypeDonut:
		chartTypeTag = "c:doughnutChart"
		grouping = ""
	case ChartTypeRadar:
		chartTypeTag = "c:radarChart"
		grouping = ""
	case ChartTypeBarStacked:
		chartTypeTag = "c:barChart"
		grouping = "stacked"
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)
	buf.WriteString(`<c:chart>`)

	// Title
	if chart.showTitle && chart.title != "" {
		fmt.Fprintf(&buf, `<c:title><c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>%s</a:t></a:r></a:p></c:rich></c:tx><c:overlay val="0"/></c:title>`, escapeXMLText(chart.title))
	}

	buf.WriteString(`<c:plotArea><c:layout/>`)

	// Chart type element
	fmt.Fprintf(&buf, `<%s>`, chartTypeTag)
	if grouping != "" {
		fmt.Fprintf(&buf, `<c:grouping val="%s"/>`, grouping)
	}

	// Series
	for i, s := range chart.series {
		fmt.Fprintf(&buf, `<c:ser><c:idx val="%d"/><c:order val="%d"/>`, i, i)
		if s.Name != "" {
			fmt.Fprintf(&buf, `<c:tx><c:strRef><c:f>"%s"</c:f></c:strRef></c:tx>`, escapeXMLText(s.Name))
		}
		// Color
		fmt.Fprintf(&buf, `<c:spPr><a:solidFill><a:srgbClr val="%02X%02X%02X"/></a:solidFill></c:spPr>`, s.Color.R, s.Color.G, s.Color.B)

		// Category data
		if len(chart.categories) > 0 {
			buf.WriteString(`<c:cat><c:strLit>`)
			fmt.Fprintf(&buf, `<c:ptCount val="%d"/>`, len(chart.categories))
			for j, cat := range chart.categories {
				fmt.Fprintf(&buf, `<c:pt idx="%d"><c:v>%s</c:v></c:pt>`, j, escapeXMLText(cat))
			}
			buf.WriteString(`</c:strLit></c:cat>`)
		}

		// Values
		if len(s.Values) > 0 {
			buf.WriteString(`<c:val><c:numLit>`)
			fmt.Fprintf(&buf, `<c:ptCount val="%d"/>`, len(s.Values))
			for j, v := range s.Values {
				fmt.Fprintf(&buf, `<c:pt idx="%d"><c:v>%s</c:v></c:pt>`, j, strconv.FormatFloat(v, 'f', -1, 64))
			}
			buf.WriteString(`</c:numLit></c:val>`)
		}

		buf.WriteString(`</c:ser>`)
	}

	fmt.Fprintf(&buf, `</%s>`, chartTypeTag)
	buf.WriteString(`</c:plotArea>`)

	// Legend
	if chart.showLegend {
		buf.WriteString(`<c:legend><c:legendPos val="b"/></c:legend>`)
	}

	buf.WriteString(`</c:chart></c:chartSpace>`)
	return buf.Bytes(), nil
}

// buildDrawingXML generates the worksheet drawing XML that anchors charts
func (wb *Workbook) buildDrawingXML(charts []*Chart) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<xdr:wsDr xmlns:xdr="http://schemas.openxmlformats.org/drawingml/2006/spreadsheetDrawing" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)

	for i, c := range charts {
		buf.WriteString(`<xdr:twoCellAnchor>`)
		fmt.Fprintf(&buf, `<xdr:from><xdr:col>%d</xdr:col><xdr:colOff>0</xdr:colOff><xdr:row>%d</xdr:row><xdr:rowOff>0</xdr:rowOff></xdr:from>`, c.x-1, c.y-1)
		fmt.Fprintf(&buf, `<xdr:to><xdr:col>%d</xdr:col><xdr:colOff>0</xdr:colOff><xdr:row>%d</xdr:row><xdr:rowOff>0</xdr:rowOff></xdr:to>`, c.x-1+c.width, c.y-1+c.height)
		fmt.Fprintf(&buf, `<xdr:graphicFrame><xdr:nvGraphicFramePr><xdr:cNvPr id="%d" name="Chart %d"/><xdr:cNvGraphicFramePr/></xdr:nvGraphicFramePr>`, i+2, i+1)
		buf.WriteString(`<xdr:xfrm><a:off x="0" y="0"/><a:ext cx="0" cy="0"/></xdr:xfrm>`)
		fmt.Fprintf(&buf, `<a:graphic><a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/chart"><c:chart xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" r:id="rId%d" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/></a:graphicData></a:graphic>`, i+1)
		buf.WriteString(`</xdr:graphicFrame><xdr:clientData/>`)
		buf.WriteString(`</xdr:twoCellAnchor>`)
	}

	buf.WriteString(`</xdr:wsDr>`)
	return buf.Bytes(), nil
}

// escapeXMLText escapes special characters for XML text content
func escapeXMLText(s string) string {
	var buf bytes.Buffer
	for _, c := range s {
		switch c {
		case '&':
			buf.WriteString("&amp;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '"':
			buf.WriteString("&quot;")
		default:
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
