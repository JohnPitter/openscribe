package spreadsheet

import (
	"fmt"
	"strings"
)

// FreezePane represents a frozen pane configuration
type FreezePane struct {
	row int // Row split (1-based, rows above this are frozen)
	col int // Column split (1-based, columns left of this are frozen)
}

// FreezePanes freezes rows above and columns to the left of the specified position
func (s *Sheet) FreezePanes(row, col int) {
	s.freezePane = &FreezePane{row: row, col: col}
}

// FreezeTopRow freezes the top row
func (s *Sheet) FreezeTopRow() {
	s.freezePane = &FreezePane{row: 2, col: 0}
}

// FreezeFirstColumn freezes the first column
func (s *Sheet) FreezeFirstColumn() {
	s.freezePane = &FreezePane{row: 0, col: 2}
}

// buildSheetViewsXML generates the <sheetViews> XML element with freeze pane
func buildSheetViewsXML(fp *FreezePane) string {
	if fp == nil {
		return ""
	}

	var buf strings.Builder
	buf.WriteString(`<sheetViews><sheetView tabSelected="1" workbookViewId="0">`)

	// Determine pane state and active pane
	ySplit := 0
	xSplit := 0
	if fp.row > 1 {
		ySplit = fp.row - 1
	}
	if fp.col > 1 {
		xSplit = fp.col - 1
	}

	topLeftCell := CellRef(max(fp.row, 1), max(fp.col, 1))

	activePane := "bottomRight"
	if xSplit > 0 && ySplit == 0 {
		activePane = "topRight"
	} else if xSplit == 0 && ySplit > 0 {
		activePane = "bottomLeft"
	}

	buf.WriteString(`<pane`)
	if xSplit > 0 {
		fmt.Fprintf(&buf, ` xSplit="%d"`, xSplit)
	}
	if ySplit > 0 {
		fmt.Fprintf(&buf, ` ySplit="%d"`, ySplit)
	}
	fmt.Fprintf(&buf, ` topLeftCell="%s" activePane="%s" state="frozen"/>`, topLeftCell, activePane)

	buf.WriteString(`</sheetView></sheetViews>`)
	return buf.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
