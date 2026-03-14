package pdf

import (
	"strings"

	"github.com/JohnPitter/openscribe/common"
)

// HTMLOptions configures HTML to PDF conversion
type HTMLOptions struct {
	PageSize    common.PageSize
	Margins     common.Margins
	DefaultFont common.Font
	BaseURL     string
}

// DefaultHTMLOptions returns default conversion options
func DefaultHTMLOptions() HTMLOptions {
	return HTMLOptions{
		PageSize:    common.PageA4,
		Margins:     common.NormalMargins(),
		DefaultFont: common.NewFont("Helvetica", 11),
	}
}

// FromHTML creates a PDF from HTML content
// This is a basic implementation that handles common HTML elements
func FromHTML(html string, opts HTMLOptions) (*Document, error) {
	doc := New()
	page := doc.AddPageWithSize(opts.PageSize, opts.Margins)

	marginLeft := opts.Margins.Left.Points()
	marginTop := opts.Margins.Top.Points()
	pageWidth := opts.PageSize.Width.Points() - opts.Margins.Left.Points() - opts.Margins.Right.Points()

	y := marginTop
	lineHeight := opts.DefaultFont.Size * 1.4

	// Simple HTML parser
	parser := &htmlParser{
		page:       page,
		font:       opts.DefaultFont,
		x:          marginLeft,
		y:          y,
		lineHeight: lineHeight,
		pageWidth:  pageWidth,
		marginLeft: marginLeft,
		doc:        doc,
		opts:       opts,
	}

	parser.parse(html)

	return doc, nil
}

type htmlParser struct {
	page       *Page
	font       common.Font
	x          float64
	y          float64
	lineHeight float64
	pageWidth  float64
	marginLeft float64
	doc        *Document
	opts       HTMLOptions
	bold       bool
	italic     bool
	inList     bool
	listItem   int
}

func (p *htmlParser) parse(html string) {
	// Strip HTML comments
	for {
		start := strings.Index(html, "<!--")
		if start == -1 {
			break
		}
		end := strings.Index(html[start:], "-->")
		if end == -1 {
			break
		}
		html = html[:start] + html[start+end+3:]
	}

	i := 0
	for i < len(html) {
		if html[i] == '<' {
			// Find tag end
			tagEnd := strings.Index(html[i:], ">")
			if tagEnd == -1 {
				break
			}
			tag := strings.ToLower(html[i+1 : i+tagEnd])
			tag = strings.TrimSpace(tag)

			// Handle tag
			p.handleTag(tag)

			i += tagEnd + 1
		} else {
			// Text content
			textEnd := strings.Index(html[i:], "<")
			var text string
			if textEnd == -1 {
				text = html[i:]
				i = len(html)
			} else {
				text = html[i : i+textEnd]
				i += textEnd
			}

			text = strings.TrimSpace(text)
			text = decodeHTMLEntities(text)
			if text != "" {
				p.addText(text)
			}
		}
	}
}

func (p *htmlParser) handleTag(tag string) {
	// Remove attributes
	spaceIdx := strings.IndexAny(tag, " \t\n")
	tagName := tag
	if spaceIdx != -1 {
		tagName = tag[:spaceIdx]
	}
	tagName = strings.TrimPrefix(tagName, "/")
	isClosing := strings.HasPrefix(tag, "/")

	switch tagName {
	case "h1":
		if !isClosing {
			p.newLine()
			p.font = p.font.WithSize(28).Bold()
		} else {
			p.font = p.opts.DefaultFont
			p.newLine()
		}
	case "h2":
		if !isClosing {
			p.newLine()
			p.font = p.font.WithSize(22).Bold()
		} else {
			p.font = p.opts.DefaultFont
			p.newLine()
		}
	case "h3":
		if !isClosing {
			p.newLine()
			p.font = p.font.WithSize(18).Bold()
		} else {
			p.font = p.opts.DefaultFont
			p.newLine()
		}
	case "h4":
		if !isClosing {
			p.newLine()
			p.font = p.font.WithSize(15).Bold()
		} else {
			p.font = p.opts.DefaultFont
			p.newLine()
		}
	case "h5", "h6":
		if !isClosing {
			p.newLine()
			p.font = p.font.WithSize(13).Bold()
		} else {
			p.font = p.opts.DefaultFont
			p.newLine()
		}
	case "p":
		p.newLine()
	case "br":
		p.newLine()
	case "b", "strong":
		if !isClosing {
			p.bold = true
			p.font = p.font.Bold()
		} else {
			p.bold = false
			p.font = p.font.WithWeight(common.FontWeightRegular)
		}
	case "i", "em":
		if !isClosing {
			p.italic = true
			p.font = p.font.Italic()
		} else {
			p.italic = false
			p.font = p.font.WithStyle(common.FontStyleNormal)
		}
	case "ul", "ol":
		if !isClosing {
			p.inList = true
			p.listItem = 0
		} else {
			p.inList = false
			p.listItem = 0
		}
	case "li":
		if !isClosing {
			p.listItem++
			p.newLine()
			prefix := "\xe2\x80\xa2 " // bullet character
			p.page.AddText(prefix, p.marginLeft+10, p.y, p.font)
			p.x = p.marginLeft + 25
		}
	case "hr":
		p.newLine()
		p.page.AddLine(p.marginLeft, p.y, p.marginLeft+p.pageWidth, p.y, common.Gray, 0.5)
		p.y += 10
	case "div", "section", "article", "header", "footer", "main", "nav":
		if !isClosing {
			p.newLine()
		}
	case "span":
		// no-op for basic rendering
	case "html", "head", "body", "title", "meta", "link", "style", "script":
		// skip structural/invisible elements
	}
}

func (p *htmlParser) addText(text string) {
	if text == "" {
		return
	}
	p.page.AddText(text, p.x, p.y, p.font)
	p.y += p.lineHeight
	p.x = p.marginLeft
}

func (p *htmlParser) newLine() {
	p.y += p.lineHeight
	p.x = p.marginLeft

	// Check if we need a new page
	pageHeight := p.opts.PageSize.Height.Points() - p.opts.Margins.Bottom.Points()
	if p.y > pageHeight {
		p.page = p.doc.AddPageWithSize(p.opts.PageSize, p.opts.Margins)
		p.y = p.opts.Margins.Top.Points()
		p.x = p.marginLeft
	}
}

func decodeHTMLEntities(s string) string {
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", "\"",
		"&#39;", "'",
		"&apos;", "'",
		"&nbsp;", " ",
		"&mdash;", "\xe2\x80\x94",
		"&ndash;", "\xe2\x80\x93",
		"&copy;", "\xc2\xa9",
		"&reg;", "\xc2\xae",
		"&trade;", "\xe2\x84\xa2",
		"&bull;", "\xe2\x80\xa2",
		"&hellip;", "\xe2\x80\xa6",
	)
	return replacer.Replace(s)
}
