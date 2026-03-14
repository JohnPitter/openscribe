package spreadsheet

import (
	"fmt"
	"strconv"

	"github.com/JohnPitter/openscribe/common"
)

// CellType represents the type of a cell value
type CellType int

const (
	CellTypeEmpty CellType = iota
	CellTypeString
	CellTypeNumber
	CellTypeBoolean
	CellTypeFormula
)

// Cell represents a cell in a worksheet
type Cell struct {
	row      *Row
	col      int
	cellType CellType
	strVal   string
	numVal   float64
	boolVal  bool
	formula  string
	font     *common.Font
	bgColor  *common.Color
	borders  *common.Borders
	hAlign   common.TextAlignment
	vAlign   common.VerticalAlignment
	numFmt   string
	comment  *Comment
	locked   *bool
}

func newCell(row *Row, col int) *Cell {
	return &Cell{
		row: row,
		col: col,
	}
}

// SetValue sets the cell value (auto-detects type)
func (c *Cell) SetValue(v interface{}) {
	switch val := v.(type) {
	case string:
		c.cellType = CellTypeString
		c.strVal = val
	case int:
		c.cellType = CellTypeNumber
		c.numVal = float64(val)
	case int64:
		c.cellType = CellTypeNumber
		c.numVal = float64(val)
	case float64:
		c.cellType = CellTypeNumber
		c.numVal = val
	case float32:
		c.cellType = CellTypeNumber
		c.numVal = float64(val)
	case bool:
		c.cellType = CellTypeBoolean
		c.boolVal = val
	default:
		c.cellType = CellTypeString
		c.strVal = fmt.Sprintf("%v", v)
	}
}

// Value returns the cell value
func (c *Cell) Value() interface{} {
	switch c.cellType {
	case CellTypeString:
		return c.strVal
	case CellTypeNumber:
		return c.numVal
	case CellTypeBoolean:
		return c.boolVal
	case CellTypeFormula:
		return c.formula
	default:
		return nil
	}
}

// SetString sets a string value
func (c *Cell) SetString(s string) { c.cellType = CellTypeString; c.strVal = s }

// SetNumber sets a numeric value
func (c *Cell) SetNumber(n float64) { c.cellType = CellTypeNumber; c.numVal = n }

// SetBool sets a boolean value
func (c *Cell) SetBool(b bool) { c.cellType = CellTypeBoolean; c.boolVal = b }

// SetFormula sets a formula
func (c *Cell) SetFormula(f string) { c.cellType = CellTypeFormula; c.formula = f }

// Type returns the cell type
func (c *Cell) Type() CellType { return c.cellType }

// String returns the string representation of the cell value
func (c *Cell) String() string {
	switch c.cellType {
	case CellTypeString:
		return c.strVal
	case CellTypeNumber:
		return strconv.FormatFloat(c.numVal, 'f', -1, 64)
	case CellTypeBoolean:
		if c.boolVal {
			return "TRUE"
		}
		return "FALSE"
	case CellTypeFormula:
		return "=" + c.formula
	default:
		return ""
	}
}

// SetFont sets the cell font
func (c *Cell) SetFont(f common.Font) { c.font = &f }

// SetBackgroundColor sets the cell background
func (c *Cell) SetBackgroundColor(color common.Color) { c.bgColor = &color }

// SetBorders sets cell borders
func (c *Cell) SetBorders(b common.Borders) { c.borders = &b }

// SetHorizontalAlignment sets horizontal alignment
func (c *Cell) SetHorizontalAlignment(a common.TextAlignment) { c.hAlign = a }

// SetVerticalAlignment sets vertical alignment
func (c *Cell) SetVerticalAlignment(a common.VerticalAlignment) { c.vAlign = a }

// SetNumberFormat sets the number format string
func (c *Cell) SetNumberFormat(fmt string) { c.numFmt = fmt }

// Col returns the column index
func (c *Cell) Col() int { return c.col }

// Ref returns the cell reference (e.g., "A1")
func (c *Cell) Ref() string { return CellRef(c.row.index, c.col) }
