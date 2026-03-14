package pdf

import (
	"strings"

	"github.com/JohnPitter/openscribe/common"
)

// FromMarkdown creates a PDF document from Markdown content.
// Supported syntax: headings (#), paragraphs, bold (**), italic (*),
// lists (- / 1.), code blocks (```), horizontal rules (---),
// links [text](url), and blockquotes (>).
func FromMarkdown(markdown string, opts HTMLOptions) (*Document, error) {
	doc := New()
	page := doc.AddPageWithSize(opts.PageSize, opts.Margins)

	marginLeft := opts.Margins.Left.Points()
	marginTop := opts.Margins.Top.Points()
	pageWidth := opts.PageSize.Width.Points() - opts.Margins.Left.Points() - opts.Margins.Right.Points()

	p := &mdPDFParser{
		doc:        doc,
		page:       page,
		opts:       opts,
		font:       opts.DefaultFont,
		x:          marginLeft,
		y:          marginTop,
		lineHeight: opts.DefaultFont.Size * 1.4,
		pageWidth:  pageWidth,
		marginLeft: marginLeft,
	}

	lines := strings.Split(markdown, "\n")
	p.parseLines(lines)

	return doc, nil
}

type mdPDFParser struct {
	doc         *Document
	page        *Page
	opts        HTMLOptions
	font        common.Font
	x           float64
	y           float64
	lineHeight  float64
	pageWidth   float64
	marginLeft  float64
	inCodeBlock bool
}

func (p *mdPDFParser) parseLines(lines []string) {
	i := 0
	for i < len(lines) {
		line := lines[i]

		// Code block toggle
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if p.inCodeBlock {
				p.inCodeBlock = false
				p.newLine()
				i++
				continue
			}
			p.inCodeBlock = true
			p.newLine()
			i++
			continue
		}

		if p.inCodeBlock {
			codeFont := common.NewFont("Courier New", p.opts.DefaultFont.Size).WithColor(common.DarkGray)
			// Draw light background for code line
			bgColor := common.NewColor(245, 245, 245)
			p.page.AddRectangle(p.marginLeft, p.y-2, p.pageWidth, p.lineHeight, bgColor, nil)
			p.page.AddText(line, p.marginLeft+4, p.y, codeFont)
			p.y += p.lineHeight
			p.checkPageBreak()
			i++
			continue
		}

		trimmed := strings.TrimSpace(line)

		// Empty line
		if trimmed == "" {
			p.y += p.lineHeight * 0.5
			p.checkPageBreak()
			i++
			continue
		}

		// Horizontal rule
		if isHorizontalRule(trimmed) {
			p.y += p.lineHeight * 0.3
			p.page.AddLine(p.marginLeft, p.y, p.marginLeft+p.pageWidth, p.y, common.Gray, 0.5)
			p.y += p.lineHeight * 0.5
			p.checkPageBreak()
			i++
			continue
		}

		// Headings
		if strings.HasPrefix(trimmed, "#") {
			level, text := parseHeading(trimmed)
			fontSize := headingSize(level)
			headingFont := p.opts.DefaultFont.WithSize(fontSize).Bold()
			p.y += p.lineHeight * 0.3
			p.page.AddText(text, p.marginLeft, p.y, headingFont)
			p.y += fontSize * 1.4
			p.checkPageBreak()
			i++
			continue
		}

		// Blockquote
		if strings.HasPrefix(trimmed, ">") {
			text := strings.TrimSpace(strings.TrimPrefix(trimmed, ">"))
			quoteFont := p.opts.DefaultFont.Italic().WithColor(common.Gray)
			// Draw left border for blockquote
			barColor := common.NewColor(200, 200, 200)
			p.page.AddLine(p.marginLeft+4, p.y-2, p.marginLeft+4, p.y+p.lineHeight-2, barColor, 2)
			p.renderInlineText(text, p.marginLeft+14, quoteFont)
			i++
			continue
		}

		// Unordered list
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			text := strings.TrimSpace(trimmed[2:])
			bullet := "\xe2\x80\xa2 " // bullet character
			p.page.AddText(bullet, p.marginLeft+10, p.y, p.font)
			p.renderInlineText(text, p.marginLeft+25, p.font)
			i++
			continue
		}

		// Ordered list
		if num, rest, ok := parseOrderedListItem(trimmed); ok {
			prefix := num + ". "
			p.page.AddText(prefix, p.marginLeft+10, p.y, p.font)
			p.renderInlineText(rest, p.marginLeft+25, p.font)
			i++
			continue
		}

		// Regular paragraph
		p.renderInlineText(trimmed, p.marginLeft, p.font)
		i++
	}
}

// renderInlineText renders text with inline formatting (bold, italic, links)
func (p *mdPDFParser) renderInlineText(text string, startX float64, baseFont common.Font) {
	segments := parseInlineMarkdown(text)
	x := startX
	for _, seg := range segments {
		f := baseFont
		if seg.bold {
			f = f.Bold()
		}
		if seg.italic {
			f = f.Italic()
		}
		if seg.link != "" {
			f = f.WithColor(common.Blue)
		}
		p.page.AddText(seg.text, x, p.y, f)
		// Approximate character width
		x += float64(len(seg.text)) * f.Size * 0.5
	}
	p.y += p.lineHeight
	p.checkPageBreak()
}

func (p *mdPDFParser) newLine() {
	p.y += p.lineHeight
	p.checkPageBreak()
}

func (p *mdPDFParser) checkPageBreak() {
	pageHeight := p.opts.PageSize.Height.Points() - p.opts.Margins.Bottom.Points()
	if p.y > pageHeight {
		p.page = p.doc.AddPageWithSize(p.opts.PageSize, p.opts.Margins)
		p.y = p.opts.Margins.Top.Points()
		p.x = p.marginLeft
	}
}

// inlineSegment represents a piece of inline-formatted text
type inlineSegment struct {
	text   string
	bold   bool
	italic bool
	link   string
}

// parseInlineMarkdown splits text into segments with bold/italic/link formatting.
func parseInlineMarkdown(text string) []inlineSegment {
	var segments []inlineSegment
	i := 0
	var current strings.Builder

	flush := func(bold, italic bool, link string) {
		if current.Len() > 0 {
			segments = append(segments, inlineSegment{
				text:   current.String(),
				bold:   bold,
				italic: italic,
				link:   link,
			})
			current.Reset()
		}
	}

	for i < len(text) {
		// Bold: **text**
		if i+1 < len(text) && text[i] == '*' && text[i+1] == '*' {
			flush(false, false, "")
			end := strings.Index(text[i+2:], "**")
			if end >= 0 {
				segments = append(segments, inlineSegment{
					text: text[i+2 : i+2+end],
					bold: true,
				})
				i = i + 2 + end + 2
				continue
			}
		}

		// Italic: *text*
		if text[i] == '*' && (i+1 >= len(text) || text[i+1] != '*') {
			flush(false, false, "")
			end := strings.Index(text[i+1:], "*")
			if end >= 0 {
				segments = append(segments, inlineSegment{
					text:   text[i+1 : i+1+end],
					italic: true,
				})
				i = i + 1 + end + 1
				continue
			}
		}

		// Link: [text](url)
		if text[i] == '[' {
			bracketEnd := strings.Index(text[i:], "](")
			if bracketEnd >= 0 {
				parenEnd := strings.Index(text[i+bracketEnd+2:], ")")
				if parenEnd >= 0 {
					flush(false, false, "")
					linkText := text[i+1 : i+bracketEnd]
					linkURL := text[i+bracketEnd+2 : i+bracketEnd+2+parenEnd]
					segments = append(segments, inlineSegment{
						text: linkText,
						link: linkURL,
					})
					i = i + bracketEnd + 2 + parenEnd + 1
					continue
				}
			}
		}

		// Inline code: `text`
		if text[i] == '`' {
			flush(false, false, "")
			end := strings.Index(text[i+1:], "`")
			if end >= 0 {
				segments = append(segments, inlineSegment{
					text: text[i+1 : i+1+end],
				})
				i = i + 1 + end + 1
				continue
			}
		}

		current.WriteByte(text[i])
		i++
	}

	flush(false, false, "")
	return segments
}

func parseHeading(line string) (int, string) {
	level := 0
	for _, c := range line {
		if c == '#' {
			level++
		} else {
			break
		}
	}
	if level > 6 {
		level = 6
	}
	text := strings.TrimSpace(line[level:])
	return level, text
}

func headingSize(level int) float64 {
	sizes := map[int]float64{
		1: 28,
		2: 22,
		3: 18,
		4: 15,
		5: 13,
		6: 12,
	}
	if s, ok := sizes[level]; ok {
		return s
	}
	return 12
}

func isHorizontalRule(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 {
		return false
	}
	allDash := true
	allStar := true
	allUnder := true
	for _, c := range trimmed {
		if c != '-' && c != ' ' {
			allDash = false
		}
		if c != '*' && c != ' ' {
			allStar = false
		}
		if c != '_' && c != ' ' {
			allUnder = false
		}
	}
	return allDash || allStar || allUnder
}

func parseOrderedListItem(line string) (string, string, bool) {
	for i, c := range line {
		if c >= '0' && c <= '9' {
			continue
		}
		if c == '.' && i > 0 && i+1 < len(line) && line[i+1] == ' ' {
			return line[:i], strings.TrimSpace(line[i+2:]), true
		}
		break
	}
	return "", "", false
}
