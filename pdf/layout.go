package pdf

import (
	"strings"

	"github.com/JohnPitter/openscribe/common"
)

// helveticaWidths contains approximate character widths for Helvetica in 1/1000 of font size
var helveticaWidths = map[byte]int{
	' ': 278, '!': 278, '"': 355, '#': 556, '$': 556, '%': 889,
	'&': 667, '\'': 191, '(': 333, ')': 333, '*': 389, '+': 584,
	',': 278, '-': 333, '.': 278, '/': 278,
	'0': 556, '1': 556, '2': 556, '3': 556, '4': 556,
	'5': 556, '6': 556, '7': 556, '8': 556, '9': 556,
	':': 278, ';': 278, '<': 584, '=': 584, '>': 584, '?': 556,
	'@': 1015,
	'A': 667, 'B': 667, 'C': 722, 'D': 722, 'E': 667, 'F': 611,
	'G': 778, 'H': 722, 'I': 278, 'J': 500, 'K': 667, 'L': 556,
	'M': 833, 'N': 722, 'O': 778, 'P': 667, 'Q': 778, 'R': 722,
	'S': 667, 'T': 611, 'U': 722, 'V': 667, 'W': 944, 'X': 667,
	'Y': 667, 'Z': 611,
	'[': 278, '\\': 278, ']': 278, '^': 469, '_': 556, '`': 333,
	'a': 556, 'b': 556, 'c': 500, 'd': 556, 'e': 556, 'f': 278,
	'g': 556, 'h': 556, 'i': 222, 'j': 222, 'k': 500, 'l': 222,
	'm': 833, 'n': 556, 'o': 556, 'p': 556, 'q': 556, 'r': 333,
	's': 500, 't': 278, 'u': 556, 'v': 500, 'w': 722, 'x': 500,
	'y': 500, 'z': 500,
	'{': 334, '|': 260, '}': 334, '~': 584,
}

// measureStringWidth calculates the width of a string in points using Helvetica metrics
func measureStringWidth(s string, fontSize float64) float64 {
	total := 0
	for i := 0; i < len(s); i++ {
		w, ok := helveticaWidths[s[i]]
		if !ok {
			w = 556 // default width for unknown characters
		}
		total += w
	}
	return float64(total) * fontSize / 1000.0
}

// measureWordWidth returns the width of a single word in points
func measureWordWidth(word string, fontSize float64) float64 {
	return measureStringWidth(word, fontSize)
}

// TextBlock represents a block of text with word wrapping and alignment
type TextBlock struct {
	x, y        float64
	width       float64
	text        string
	font        common.Font
	alignment   common.TextAlignment
	lineSpacing float64
	columns     int
	columnGap   float64
	// computed wrapped lines (filled during build)
	lines []string
	// document reference for overflow page creation
	document *Document
	page     *Page
}

func (tb *TextBlock) pdfElement() {}

// SetAlignment sets the text alignment
func (tb *TextBlock) SetAlignment(a common.TextAlignment) { tb.alignment = a }

// Alignment returns the alignment
func (tb *TextBlock) Alignment() common.TextAlignment { return tb.alignment }

// SetLineSpacing sets the line spacing multiplier (1.0 = single, 1.5 = 1.5x, 2.0 = double)
func (tb *TextBlock) SetLineSpacing(multiplier float64) {
	if multiplier < 0.5 {
		multiplier = 0.5
	}
	tb.lineSpacing = multiplier
}

// LineSpacing returns the line spacing multiplier
func (tb *TextBlock) LineSpacing() float64 { return tb.lineSpacing }

// SetColumns sets multi-column layout
func (tb *TextBlock) SetColumns(n int, gap float64) {
	if n < 1 {
		n = 1
	}
	tb.columns = n
	tb.columnGap = gap
}

// Columns returns the number of columns
func (tb *TextBlock) Columns() int { return tb.columns }

// ColumnGap returns the gap between columns
func (tb *TextBlock) ColumnGap() float64 { return tb.columnGap }

// Text returns the text content
func (tb *TextBlock) Text() string { return tb.text }

// WrapLines performs word wrapping and returns the wrapped lines
func (tb *TextBlock) WrapLines() []string {
	colWidth := tb.columnWidth()
	return wrapText(tb.text, colWidth, tb.font.Size)
}

// columnWidth returns the width available for text in a single column
func (tb *TextBlock) columnWidth() float64 {
	if tb.columns <= 1 {
		return tb.width
	}
	totalGaps := float64(tb.columns-1) * tb.columnGap
	return (tb.width - totalGaps) / float64(tb.columns)
}

// wrapText splits text into lines that fit within maxWidth
func wrapText(text string, maxWidth, fontSize float64) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	paragraphs := strings.Split(text, "\n")
	var result []string

	for _, para := range paragraphs {
		if para == "" {
			result = append(result, "")
			continue
		}

		words := strings.Fields(para)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}

		currentLine := words[0]
		currentWidth := measureWordWidth(currentLine, fontSize)

		spaceWidth := measureStringWidth(" ", fontSize)

		for i := 1; i < len(words); i++ {
			wordWidth := measureWordWidth(words[i], fontSize)
			if currentWidth+spaceWidth+wordWidth <= maxWidth {
				currentLine += " " + words[i]
				currentWidth += spaceWidth + wordWidth
			} else {
				result = append(result, currentLine)
				currentLine = words[i]
				currentWidth = wordWidth
			}
		}
		result = append(result, currentLine)
	}

	return result
}

// AddTextBlock adds a text block with word wrapping to the page
func (p *Page) AddTextBlock(x, y, width float64, text string, font common.Font) *TextBlock {
	tb := &TextBlock{
		x:           x,
		y:           y,
		width:       width,
		text:        text,
		font:        font,
		alignment:   common.TextAlignLeft,
		lineSpacing: 1.2,
		columns:     1,
		columnGap:   12,
		page:        p,
	}
	p.elements = append(p.elements, tb)
	return tb
}
