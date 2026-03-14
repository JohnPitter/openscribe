package document

import (
	"strings"

	"github.com/JohnPitter/openscribe/common"
)

// FromMarkdown creates a DOCX document from Markdown content.
// Supported syntax: headings (#), paragraphs, bold (**), italic (*),
// lists (- / 1.), code blocks (```), horizontal rules (---),
// links [text](url), and blockquotes (>).
func FromMarkdown(markdown string) (*Document, error) {
	doc := New()

	lines := strings.Split(markdown, "\n")
	p := &mdDocxParser{
		doc: doc,
	}
	p.parseLines(lines)

	return doc, nil
}

type mdDocxParser struct {
	doc         *Document
	inCodeBlock bool
}

func (p *mdDocxParser) parseLines(lines []string) {
	i := 0
	for i < len(lines) {
		line := lines[i]

		// Code block toggle
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if p.inCodeBlock {
				p.inCodeBlock = false
				i++
				continue
			}
			p.inCodeBlock = true
			i++
			continue
		}

		if p.inCodeBlock {
			para := p.doc.AddParagraph()
			bgColor := common.NewColor(245, 245, 245)
			run := para.AddRun()
			run.SetText(line)
			run.SetFontFamily("Courier New")
			run.SetSize(10)
			run.SetHighlight("lightGray")
			_ = bgColor
			i++
			continue
		}

		trimmed := strings.TrimSpace(line)

		// Empty line
		if trimmed == "" {
			i++
			continue
		}

		// Horizontal rule
		if mdIsHorizontalRule(trimmed) {
			para := p.doc.AddParagraph()
			run := para.AddRun()
			run.SetText("") // empty paragraph as separator
			// Add a bottom border as horizontal rule via spacing
			para.SetSpacing(common.Pt(6), common.Pt(6), 1.0)
			i++
			continue
		}

		// Headings
		if strings.HasPrefix(trimmed, "#") {
			level, text := mdParseHeading(trimmed)
			p.doc.AddHeading(text, level)
			i++
			continue
		}

		// Blockquote
		if strings.HasPrefix(trimmed, ">") {
			text := strings.TrimSpace(strings.TrimPrefix(trimmed, ">"))
			para := p.doc.AddParagraph()
			para.SetIndent(common.Pt(24), common.Pt(0), common.Pt(0))
			p.addInlineRuns(para, text, true)
			i++
			continue
		}

		// Unordered list
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
			text := strings.TrimSpace(trimmed[2:])
			para := p.doc.AddParagraph()
			para.SetIndent(common.Pt(18), common.Pt(0), common.Pt(0))
			bulletRun := para.AddRun()
			bulletRun.SetText("\xe2\x80\xa2 ")
			p.addInlineRuns(para, text, false)
			i++
			continue
		}

		// Ordered list
		if num, rest, ok := mdParseOrderedListItem(trimmed); ok {
			para := p.doc.AddParagraph()
			para.SetIndent(common.Pt(18), common.Pt(0), common.Pt(0))
			numRun := para.AddRun()
			numRun.SetText(num + ". ")
			p.addInlineRuns(para, rest, false)
			i++
			continue
		}

		// Regular paragraph
		para := p.doc.AddParagraph()
		p.addInlineRuns(para, trimmed, false)
		i++
	}
}

// addInlineRuns parses inline markdown (bold, italic, links) and adds runs to a paragraph.
func (p *mdDocxParser) addInlineRuns(para *Paragraph, text string, italic bool) {
	segments := mdParseInline(text)
	for _, seg := range segments {
		run := para.AddRun()
		run.SetText(seg.text)
		if seg.bold || seg.italic || italic {
			if seg.bold {
				run.SetBold(true)
			}
			if seg.italic || italic {
				run.SetItalic(true)
			}
		}
		if seg.link != "" {
			run.SetColor(common.Blue)
			run.SetUnderline(true)
		}
		if seg.code {
			run.SetFontFamily("Courier New")
			run.SetSize(10)
		}
	}
}

type mdInlineSegment struct {
	text   string
	bold   bool
	italic bool
	link   string
	code   bool
}

func mdParseInline(text string) []mdInlineSegment {
	var segments []mdInlineSegment
	i := 0
	var current strings.Builder

	flush := func() {
		if current.Len() > 0 {
			segments = append(segments, mdInlineSegment{text: current.String()})
			current.Reset()
		}
	}

	for i < len(text) {
		// Bold: **text**
		if i+1 < len(text) && text[i] == '*' && text[i+1] == '*' {
			flush()
			end := strings.Index(text[i+2:], "**")
			if end >= 0 {
				segments = append(segments, mdInlineSegment{
					text: text[i+2 : i+2+end],
					bold: true,
				})
				i = i + 2 + end + 2
				continue
			}
		}

		// Italic: *text*
		if text[i] == '*' && (i+1 >= len(text) || text[i+1] != '*') {
			flush()
			end := strings.Index(text[i+1:], "*")
			if end >= 0 {
				segments = append(segments, mdInlineSegment{
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
					flush()
					linkText := text[i+1 : i+bracketEnd]
					linkURL := text[i+bracketEnd+2 : i+bracketEnd+2+parenEnd]
					segments = append(segments, mdInlineSegment{
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
			flush()
			end := strings.Index(text[i+1:], "`")
			if end >= 0 {
				segments = append(segments, mdInlineSegment{
					text: text[i+1 : i+1+end],
					code: true,
				})
				i = i + 1 + end + 1
				continue
			}
		}

		current.WriteByte(text[i])
		i++
	}

	flush()
	return segments
}

func mdParseHeading(line string) (int, string) {
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

func mdIsHorizontalRule(line string) bool {
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

func mdParseOrderedListItem(line string) (string, string, bool) {
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
