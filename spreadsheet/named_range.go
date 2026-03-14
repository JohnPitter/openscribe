package spreadsheet

import (
	"fmt"
	"strings"
)

// NamedRange represents a named range in the workbook
type NamedRange struct {
	name      string
	sheetName string
	cellRange string
	sheetID   int // -1 for workbook-level, 0+ for sheet-level
}

// Name returns the named range name
func (nr *NamedRange) Name() string { return nr.name }

// SheetName returns the sheet name
func (nr *NamedRange) SheetName() string { return nr.sheetName }

// CellRange returns the cell range
func (nr *NamedRange) CellRange() string { return nr.cellRange }

// AddNamedRange adds a named range to the workbook
func (wb *Workbook) AddNamedRange(name, sheetName, cellRange string) *NamedRange {
	// Find sheet index
	sheetID := -1
	for i, s := range wb.sheets {
		if s.name == sheetName {
			sheetID = i
			break
		}
	}

	nr := &NamedRange{
		name:      name,
		sheetName: sheetName,
		cellRange: cellRange,
		sheetID:   sheetID,
	}
	wb.namedRanges = append(wb.namedRanges, nr)
	return nr
}

// NamedRange returns a named range by name, or nil if not found
func (wb *Workbook) NamedRange(name string) *NamedRange {
	for _, nr := range wb.namedRanges {
		if nr.name == name {
			return nr
		}
	}
	return nil
}

// NamedRanges returns all named ranges
func (wb *Workbook) NamedRanges() []*NamedRange {
	return wb.namedRanges
}

// buildDefinedNamesXML generates the <definedNames> XML for the workbook
func buildDefinedNamesXML(namedRanges []*NamedRange) string {
	if len(namedRanges) == 0 {
		return ""
	}

	var buf strings.Builder
	buf.WriteString(`<definedNames>`)
	for _, nr := range namedRanges {
		// Reference format: SheetName!CellRange
		ref := fmt.Sprintf("'%s'!%s", nr.sheetName, nr.cellRange)
		if nr.sheetID >= 0 {
			fmt.Fprintf(&buf, `<definedName name="%s" localSheetId="%d">%s</definedName>`,
				escapeXMLAttr(nr.name), nr.sheetID, escapeXMLText(ref))
		} else {
			fmt.Fprintf(&buf, `<definedName name="%s">%s</definedName>`,
				escapeXMLAttr(nr.name), escapeXMLText(ref))
		}
	}
	buf.WriteString(`</definedNames>`)
	return buf.String()
}
