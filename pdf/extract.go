package pdf

import (
	"bytes"
	"fmt"
	"strings"
)

// ExtractText extracts text content from a PDF document
// This works for PDFs created by openscribe. For complex PDFs,
// a full PDF parser would be needed.
func (d *Document) ExtractText() (string, error) {
	var allText strings.Builder

	for i, page := range d.pages {
		if i > 0 {
			allText.WriteString("\n\n--- Page Break ---\n\n")
		}

		// If page has raw data (loaded from file), try to extract from raw PDF
		if len(page.rawData) > 0 {
			text, err := extractTextFromRaw(page.rawData)
			if err != nil {
				// If raw extraction fails, fall back to elements
				text = extractTextFromElements(page)
			}
			allText.WriteString(text)
			continue
		}

		// Extract from page elements
		allText.WriteString(extractTextFromElements(page))
	}

	return allText.String(), nil
}

// ExtractPageText extracts text from a specific page
func (d *Document) ExtractPageText(pageIndex int) (string, error) {
	if pageIndex < 0 || pageIndex >= len(d.pages) {
		return "", fmt.Errorf("page index %d out of range", pageIndex)
	}

	page := d.pages[pageIndex]
	if len(page.rawData) > 0 {
		text, err := extractTextFromRaw(page.rawData)
		if err != nil {
			return extractTextFromElements(page), nil
		}
		return text, nil
	}

	return extractTextFromElements(page), nil
}

// extractTextFromElements extracts text from page elements
func extractTextFromElements(page *Page) string {
	var text strings.Builder

	for _, elem := range page.elements {
		switch e := elem.(type) {
		case *TextElement:
			if text.Len() > 0 {
				text.WriteString("\n")
			}
			text.WriteString(e.text)
		case *TableElement:
			if text.Len() > 0 {
				text.WriteString("\n")
			}
			for row := 0; row < e.rows; row++ {
				if row > 0 {
					text.WriteString("\n")
				}
				for col := 0; col < e.cols; col++ {
					if col > 0 {
						text.WriteString("\t")
					}
					text.WriteString(e.cells[row][col])
				}
			}
		}
	}

	return text.String()
}

// extractTextFromRaw extracts text from raw PDF data
// This is a basic implementation that finds text between BT/ET markers
func extractTextFromRaw(data []byte) (string, error) {
	var text strings.Builder

	// Find stream content between "stream" and "endstream"
	content := string(data)

	idx := 0
	for {
		streamStart := strings.Index(content[idx:], "stream\n")
		if streamStart == -1 {
			streamStart = strings.Index(content[idx:], "stream\r\n")
			if streamStart == -1 {
				break
			}
			streamStart += idx + len("stream\r\n")
		} else {
			streamStart += idx + len("stream\n")
		}

		streamEnd := strings.Index(content[streamStart:], "\nendstream")
		if streamEnd == -1 {
			streamEnd = strings.Index(content[streamStart:], "\r\nendstream")
			if streamEnd == -1 {
				break
			}
		}

		streamContent := content[streamStart : streamStart+streamEnd]
		extractedText := extractTextFromStream(streamContent)
		if extractedText != "" {
			if text.Len() > 0 {
				text.WriteString("\n")
			}
			text.WriteString(extractedText)
		}

		idx = streamStart + streamEnd + len("\nendstream")
		if idx >= len(content) {
			break
		}
	}

	if text.Len() == 0 {
		return "", fmt.Errorf("no text found in PDF")
	}

	return text.String(), nil
}

// extractTextFromStream extracts text from a PDF content stream
func extractTextFromStream(stream string) string {
	var text strings.Builder

	// Find text between ( ) in Tj/TJ operators
	i := 0
	for i < len(stream) {
		// Look for text show operators
		if stream[i] == '(' {
			// Find matching close paren
			j := i + 1
			depth := 1
			for j < len(stream) && depth > 0 {
				if stream[j] == '(' && (j == 0 || stream[j-1] != '\\') {
					depth++
				} else if stream[j] == ')' && (j == 0 || stream[j-1] != '\\') {
					depth--
				}
				j++
			}

			if depth == 0 {
				extracted := stream[i+1 : j-1]
				// Unescape PDF string
				extracted = unescapePDFString(extracted)
				if text.Len() > 0 && extracted != "" {
					text.WriteString("\n")
				}
				text.WriteString(extracted)
			}
			i = j
		} else {
			i++
		}
	}

	return text.String()
}

// unescapePDFString unescapes a PDF string
func unescapePDFString(s string) string {
	var result bytes.Buffer
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result.WriteByte('\n')
			case 'r':
				result.WriteByte('\r')
			case 't':
				result.WriteByte('\t')
			case '(':
				result.WriteByte('(')
			case ')':
				result.WriteByte(')')
			case '\\':
				result.WriteByte('\\')
			default:
				result.WriteByte(s[i+1])
			}
			i += 2
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}
