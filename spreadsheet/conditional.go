package spreadsheet

import "github.com/JohnPitter/openscribe/common"

// ConditionType represents the type of conditional formatting rule
type ConditionType int

const (
	ConditionGreaterThan ConditionType = iota
	ConditionLessThan
	ConditionEqual
	ConditionNotEqual
	ConditionBetween
	ConditionContains
	ConditionBeginsWith
	ConditionEndsWith
	ConditionTop10
	ConditionAboveAverage
	ConditionBelowAverage
	ConditionDuplicate
	ConditionUnique
	ConditionColorScale
	ConditionDataBar
)

// ConditionalFormat represents a conditional formatting rule
type ConditionalFormat struct {
	condType  ConditionType
	cellRange string // e.g., "A1:A10"
	value     string
	value2    string // for "between" condition
	bgColor   *common.Color
	fontColor *common.Color
	bold      bool
	italic    bool
	// Color scale
	minColor *common.Color
	maxColor *common.Color
	// Data bar
	barColor *common.Color
}

// AddConditionalFormat adds a conditional formatting rule to the sheet
func (s *Sheet) AddConditionalFormat(cellRange string, condType ConditionType) *ConditionalFormat {
	cf := &ConditionalFormat{
		condType:  condType,
		cellRange: cellRange,
	}
	s.conditionalFormats = append(s.conditionalFormats, cf)
	return cf
}

// SetValue sets the comparison value
func (cf *ConditionalFormat) SetValue(val string) *ConditionalFormat {
	cf.value = val
	return cf
}

// SetValue2 sets the second value (for "between" conditions)
func (cf *ConditionalFormat) SetValue2(val string) *ConditionalFormat {
	cf.value2 = val
	return cf
}

// SetBackgroundColor sets the background color when condition is met
func (cf *ConditionalFormat) SetBackgroundColor(c common.Color) *ConditionalFormat {
	cf.bgColor = &c
	return cf
}

// SetFontColor sets the font color when condition is met
func (cf *ConditionalFormat) SetFontColor(c common.Color) *ConditionalFormat {
	cf.fontColor = &c
	return cf
}

// SetBold sets bold formatting when condition is met
func (cf *ConditionalFormat) SetBold(b bool) *ConditionalFormat {
	cf.bold = b
	return cf
}

// SetItalic sets italic formatting when condition is met
func (cf *ConditionalFormat) SetItalic(i bool) *ConditionalFormat {
	cf.italic = i
	return cf
}

// SetColorScale sets min/max colors for color scale formatting
func (cf *ConditionalFormat) SetColorScale(minColor, maxColor common.Color) *ConditionalFormat {
	cf.minColor = &minColor
	cf.maxColor = &maxColor
	return cf
}

// SetBarColor sets the data bar color
func (cf *ConditionalFormat) SetBarColor(c common.Color) *ConditionalFormat {
	cf.barColor = &c
	return cf
}

// CellRange returns the cell range
func (cf *ConditionalFormat) CellRange() string { return cf.cellRange }

// Type returns the condition type
func (cf *ConditionalFormat) Type() ConditionType { return cf.condType }

// Value returns the comparison value
func (cf *ConditionalFormat) Value() string { return cf.value }
