package template

import (
	"fmt"
	"strings"

	"github.com/JohnPitter/openscribe/document"
	"github.com/JohnPitter/openscribe/pdf"
	"github.com/JohnPitter/openscribe/spreadsheet"
)

// TemplateEngine provides variable substitution, conditionals, and loops
// for document templates using a {{key}} syntax.
type TemplateEngine struct {
	data map[string]interface{}
}

// NewEngine creates a new TemplateEngine with an empty data map.
func NewEngine() *TemplateEngine {
	return &TemplateEngine{
		data: make(map[string]interface{}),
	}
}

// SetData sets a single variable in the data map.
func (e *TemplateEngine) SetData(key string, value interface{}) {
	e.data[key] = value
}

// SetDataMap sets multiple variables at once.
func (e *TemplateEngine) SetDataMap(data map[string]interface{}) {
	for k, v := range data {
		e.data[k] = v
	}
}

// RenderDocx replaces {{key}} patterns in all paragraph runs of a DOCX document.
// It also handles {{#if key}}...{{/if}} conditionals and {{#each items}}...{{/each}} loops.
func (e *TemplateEngine) RenderDocx(tmplDoc *document.Document) *document.Document {
	for _, para := range tmplDoc.Paragraphs() {
		for _, run := range para.Runs() {
			text := run.Text()
			text = e.processDirectives(text)
			text = e.replaceVars(text)
			run.SetText(text)
		}
	}
	return tmplDoc
}

// RenderPdf replaces {{key}} patterns in all TextElement text on all pages.
func (e *TemplateEngine) RenderPdf(tmplDoc *pdf.Document) *pdf.Document {
	for i := 0; i < tmplDoc.PageCount(); i++ {
		page := tmplDoc.Page(i)
		if page == nil {
			continue
		}
		for _, elem := range page.Elements() {
			if te, ok := elem.(*pdf.TextElement); ok {
				text := te.Text()
				text = e.processDirectives(text)
				text = e.replaceVars(text)
				te.SetText(text)
			}
		}
	}
	return tmplDoc
}

// RenderXlsx replaces {{key}} patterns in all string cell values.
func (e *TemplateEngine) RenderXlsx(tmplWb *spreadsheet.Workbook) *spreadsheet.Workbook {
	for i := 0; i < tmplWb.SheetCount(); i++ {
		sheet := tmplWb.Sheet(i)
		if sheet == nil {
			continue
		}
		for row := 1; row <= sheet.MaxRow(); row++ {
			for col := 1; col <= sheet.MaxCol(); col++ {
				cell := sheet.Cell(row, col)
				if cell.Type() == spreadsheet.CellTypeString {
					text := cell.String()
					text = e.processDirectives(text)
					text = e.replaceVars(text)
					cell.SetString(text)
				}
			}
		}
	}
	return tmplWb
}

// replaceVars replaces all {{key}} patterns with values from the data map.
// Supports dot notation for nested maps: {{user.name}}.
func (e *TemplateEngine) replaceVars(text string) string {
	for {
		start := strings.Index(text, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], "}}")
		if end == -1 {
			break
		}
		end += start

		key := strings.TrimSpace(text[start+2 : end])

		// Skip directive tags
		if strings.HasPrefix(key, "#") || strings.HasPrefix(key, "/") {
			// Move past this tag to avoid infinite loop
			text = text[:start] + text[end+2:]
			continue
		}

		val := e.resolveKey(key)
		replacement := fmt.Sprintf("%v", val)
		if val == nil {
			replacement = ""
		}

		text = text[:start] + replacement + text[end+2:]
	}
	return text
}

// processDirectives handles {{#if key}}...{{/if}} and {{#each items}}...{{/each}}.
func (e *TemplateEngine) processDirectives(text string) string {
	text = e.processConditionals(text)
	text = e.processLoops(text)
	return text
}

// processConditionals handles {{#if key}}content{{/if}}.
func (e *TemplateEngine) processConditionals(text string) string {
	for {
		ifStart := strings.Index(text, "{{#if ")
		if ifStart == -1 {
			break
		}
		ifEnd := strings.Index(text[ifStart:], "}}")
		if ifEnd == -1 {
			break
		}
		ifEnd += ifStart

		key := strings.TrimSpace(text[ifStart+6 : ifEnd])

		endTag := "{{/if}}"
		endIdx := strings.Index(text[ifEnd+2:], endTag)
		if endIdx == -1 {
			break
		}
		endIdx += ifEnd + 2

		content := text[ifEnd+2 : endIdx]

		if e.isTruthy(key) {
			text = text[:ifStart] + content + text[endIdx+len(endTag):]
		} else {
			text = text[:ifStart] + text[endIdx+len(endTag):]
		}
	}
	return text
}

// processLoops handles {{#each items}}content with {{this}} or {{.key}}{{/each}}.
func (e *TemplateEngine) processLoops(text string) string {
	for {
		eachStart := strings.Index(text, "{{#each ")
		if eachStart == -1 {
			break
		}
		eachEnd := strings.Index(text[eachStart:], "}}")
		if eachEnd == -1 {
			break
		}
		eachEnd += eachStart

		key := strings.TrimSpace(text[eachStart+8 : eachEnd])

		endTag := "{{/each}}"
		endIdx := strings.Index(text[eachEnd+2:], endTag)
		if endIdx == -1 {
			break
		}
		endIdx += eachEnd + 2

		content := text[eachEnd+2 : endIdx]

		val := e.resolveKey(key)
		var result strings.Builder

		if items, ok := val.([]interface{}); ok {
			for _, item := range items {
				rendered := content
				// Replace {{this}} with the item itself if it's a string/number
				rendered = strings.ReplaceAll(rendered, "{{this}}", fmt.Sprintf("%v", item))

				// If item is a map, replace {{.key}} references
				if m, ok := item.(map[string]interface{}); ok {
					for k, v := range m {
						rendered = strings.ReplaceAll(rendered, "{{."+k+"}}", fmt.Sprintf("%v", v))
					}
				}
				result.WriteString(rendered)
			}
		}

		text = text[:eachStart] + result.String() + text[endIdx+len(endTag):]
	}
	return text
}

// resolveKey resolves a potentially dot-notated key from the data map.
func (e *TemplateEngine) resolveKey(key string) interface{} {
	parts := strings.Split(key, ".")
	var current interface{} = e.data

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = m[part]
		if !ok {
			return nil
		}
	}

	return current
}

// isTruthy checks if a key's value is truthy (non-nil, non-empty, non-false, non-zero).
func (e *TemplateEngine) isTruthy(key string) bool {
	val := e.resolveKey(key)
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v != ""
	case int:
		return v != 0
	case int64:
		return v != 0
	case float64:
		return v != 0
	case []interface{}:
		return len(v) > 0
	default:
		return true
	}
}
