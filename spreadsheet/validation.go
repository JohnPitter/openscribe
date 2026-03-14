package spreadsheet

import (
	"fmt"
	"strings"
)

// ValidationType represents the type of data validation
type ValidationType int

const (
	ValidationList        ValidationType = iota // Dropdown list
	ValidationWholeNumber                       // Whole number
	ValidationDecimal                           // Decimal number
	ValidationDate                              // Date
	ValidationTextLength                        // Text length
	ValidationCustom                            // Custom formula
)

// Validation represents a data validation rule on a cell range
type Validation struct {
	cellRange      string
	validationType ValidationType
	listItems      []string
	minValue       string
	maxValue       string
	customFormula  string
	errorTitle     string
	errorMessage   string
	promptTitle    string
	promptMessage  string
}

// AddValidation adds a data validation rule to a cell range
func (s *Sheet) AddValidation(cellRange string, validationType ValidationType) *Validation {
	v := &Validation{
		cellRange:      cellRange,
		validationType: validationType,
	}
	s.validations = append(s.validations, v)
	return v
}

// SetList sets the list items for a dropdown validation
func (v *Validation) SetList(items []string) *Validation {
	v.listItems = items
	return v
}

// SetRange sets the min/max values for number/date/text length validations
func (v *Validation) SetRange(min, max string) *Validation {
	v.minValue = min
	v.maxValue = max
	return v
}

// SetCustomFormula sets the custom formula for custom validations
func (v *Validation) SetCustomFormula(formula string) *Validation {
	v.customFormula = formula
	return v
}

// SetErrorMessage sets the error alert title and message
func (v *Validation) SetErrorMessage(title, msg string) *Validation {
	v.errorTitle = title
	v.errorMessage = msg
	return v
}

// SetPromptMessage sets the input prompt title and message
func (v *Validation) SetPromptMessage(title, msg string) *Validation {
	v.promptTitle = title
	v.promptMessage = msg
	return v
}

// CellRange returns the cell range
func (v *Validation) CellRange() string { return v.cellRange }

// Type returns the validation type
func (v *Validation) Type() ValidationType { return v.validationType }

// buildValidationsXML generates the <dataValidations> XML string
func buildValidationsXML(validations []*Validation) string {
	if len(validations) == 0 {
		return ""
	}

	var buf strings.Builder
	fmt.Fprintf(&buf, `<dataValidations count="%d">`, len(validations))

	for _, v := range validations {
		buf.WriteString(`<dataValidation`)

		// Type attribute
		switch v.validationType {
		case ValidationList:
			buf.WriteString(` type="list"`)
		case ValidationWholeNumber:
			buf.WriteString(` type="whole"`)
		case ValidationDecimal:
			buf.WriteString(` type="decimal"`)
		case ValidationDate:
			buf.WriteString(` type="date"`)
		case ValidationTextLength:
			buf.WriteString(` type="textLength"`)
		case ValidationCustom:
			buf.WriteString(` type="custom"`)
		}

		// Operator for range-based validations
		if v.minValue != "" && v.maxValue != "" && v.validationType != ValidationList && v.validationType != ValidationCustom {
			buf.WriteString(` operator="between"`)
		}

		// Allow blank
		buf.WriteString(` allowBlank="1"`)

		// Show input message
		if v.promptTitle != "" || v.promptMessage != "" {
			buf.WriteString(` showInputMessage="1"`)
		}

		// Show error message
		if v.errorTitle != "" || v.errorMessage != "" {
			buf.WriteString(` showErrorMessage="1"`)
		}

		// Error title/message
		if v.errorTitle != "" {
			fmt.Fprintf(&buf, ` errorTitle="%s"`, escapeXMLAttr(v.errorTitle))
		}
		if v.errorMessage != "" {
			fmt.Fprintf(&buf, ` error="%s"`, escapeXMLAttr(v.errorMessage))
		}

		// Prompt title/message
		if v.promptTitle != "" {
			fmt.Fprintf(&buf, ` promptTitle="%s"`, escapeXMLAttr(v.promptTitle))
		}
		if v.promptMessage != "" {
			fmt.Fprintf(&buf, ` prompt="%s"`, escapeXMLAttr(v.promptMessage))
		}

		// Cell range
		fmt.Fprintf(&buf, ` sqref="%s"`, v.cellRange)
		buf.WriteString(`>`)

		// Formula
		switch v.validationType {
		case ValidationList:
			if len(v.listItems) > 0 {
				fmt.Fprintf(&buf, `<formula1>"%s"</formula1>`, escapeXMLText(strings.Join(v.listItems, ",")))
			}
		case ValidationCustom:
			if v.customFormula != "" {
				fmt.Fprintf(&buf, `<formula1>%s</formula1>`, escapeXMLText(v.customFormula))
			}
		default:
			if v.minValue != "" {
				fmt.Fprintf(&buf, `<formula1>%s</formula1>`, escapeXMLText(v.minValue))
			}
			if v.maxValue != "" {
				fmt.Fprintf(&buf, `<formula2>%s</formula2>`, escapeXMLText(v.maxValue))
			}
		}

		buf.WriteString(`</dataValidation>`)
	}

	buf.WriteString(`</dataValidations>`)
	return buf.String()
}

// escapeXMLAttr escapes special characters for XML attribute values
func escapeXMLAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}
