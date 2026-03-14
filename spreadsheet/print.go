package spreadsheet

import (
	"fmt"
	"strings"

	"github.com/JohnPitter/openscribe/common"
)

// PrintSettings holds print configuration for a sheet
type PrintSettings struct {
	orientation common.Orientation
	paperSize   int // Standard paper size code (1=Letter, 9=A4, etc.)
	fitToWidth  int // Fit to N pages wide (0 = not set)
	fitToHeight int // Fit to N pages tall (0 = not set)
	printArea   *PrintArea
	repeatRows  string // e.g., "1:2" to repeat first two rows
	repeatCols  string // e.g., "A:B" to repeat first two columns
}

// PrintArea represents the print area range
type PrintArea struct {
	startRow int
	startCol int
	endRow   int
	endCol   int
}

// SetPrintArea sets the print area for the sheet
func (s *Sheet) SetPrintArea(startRow, startCol, endRow, endCol int) {
	if s.printSettings == nil {
		s.printSettings = &PrintSettings{}
	}
	s.printSettings.printArea = &PrintArea{
		startRow: startRow,
		startCol: startCol,
		endRow:   endRow,
		endCol:   endCol,
	}
}

// SetPrintTitles sets the rows and columns to repeat on each printed page
func (s *Sheet) SetPrintTitles(repeatRows, repeatCols string) {
	if s.printSettings == nil {
		s.printSettings = &PrintSettings{}
	}
	s.printSettings.repeatRows = repeatRows
	s.printSettings.repeatCols = repeatCols
}

// SetPageOrientation sets the page orientation for printing
func (s *Sheet) SetPageOrientation(orientation common.Orientation) {
	if s.printSettings == nil {
		s.printSettings = &PrintSettings{}
	}
	s.printSettings.orientation = orientation
}

// SetPaperSize sets the paper size code
func (s *Sheet) SetPaperSize(size int) {
	if s.printSettings == nil {
		s.printSettings = &PrintSettings{}
	}
	s.printSettings.paperSize = size
}

// SetFitToPage configures fitting the sheet to a specific number of pages
func (s *Sheet) SetFitToPage(width, height int) {
	if s.printSettings == nil {
		s.printSettings = &PrintSettings{}
	}
	s.printSettings.fitToWidth = width
	s.printSettings.fitToHeight = height
}

// buildPageSetupXML generates the <pageSetup> XML element
func buildPageSetupXML(ps *PrintSettings) string {
	if ps == nil {
		return ""
	}

	var buf strings.Builder

	// Page margins (default reasonable margins)
	buf.WriteString(`<pageMargins left="0.7" right="0.7" top="0.75" bottom="0.75" header="0.3" footer="0.3"/>`)

	// Page setup
	buf.WriteString(`<pageSetup`)

	if ps.paperSize > 0 {
		fmt.Fprintf(&buf, ` paperSize="%d"`, ps.paperSize)
	}

	if ps.orientation == common.OrientationLandscape {
		buf.WriteString(` orientation="landscape"`)
	} else {
		buf.WriteString(` orientation="portrait"`)
	}

	if ps.fitToWidth > 0 {
		fmt.Fprintf(&buf, ` fitToWidth="%d"`, ps.fitToWidth)
	}
	if ps.fitToHeight > 0 {
		fmt.Fprintf(&buf, ` fitToHeight="%d"`, ps.fitToHeight)
	}

	buf.WriteString(`/>`)
	return buf.String()
}

// buildPrintAreaDefinedName generates the print area as a defined name
func buildPrintAreaDefinedName(sheetName string, sheetIndex int, pa *PrintArea) string {
	if pa == nil {
		return ""
	}
	ref := fmt.Sprintf("'%s'!$%s$%d:$%s$%d",
		sheetName,
		colName(pa.startCol), pa.startRow,
		colName(pa.endCol), pa.endRow,
	)
	return fmt.Sprintf(`<definedName name="_xlnm.Print_Area" localSheetId="%d">%s</definedName>`,
		sheetIndex, escapeXMLText(ref))
}

// buildPrintTitlesDefinedName generates the print titles as a defined name
func buildPrintTitlesDefinedName(sheetName string, sheetIndex int, repeatRows, repeatCols string) string {
	if repeatRows == "" && repeatCols == "" {
		return ""
	}

	var parts []string
	if repeatRows != "" {
		parts = append(parts, fmt.Sprintf("'%s'!$%s", sheetName, repeatRows))
	}
	if repeatCols != "" {
		parts = append(parts, fmt.Sprintf("'%s'!$%s", sheetName, repeatCols))
	}

	ref := strings.Join(parts, ",")
	return fmt.Sprintf(`<definedName name="_xlnm.Print_Titles" localSheetId="%d">%s</definedName>`,
		sheetIndex, escapeXMLText(ref))
}
