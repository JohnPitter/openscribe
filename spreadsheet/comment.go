package spreadsheet

import (
	"fmt"
	"strings"
)

// Comment represents a cell comment
type Comment struct {
	row    int
	col    int
	author string
	text   string
}

// SetComment sets a comment on the cell
func (c *Cell) SetComment(author, text string) {
	c.comment = &Comment{
		row:    c.row.index,
		col:    c.col,
		author: author,
		text:   text,
	}
	c.row.sheet.addComment(c.comment)
}

// Comment returns the comment author and text, or empty strings if none
func (c *Cell) Comment() (author, text string) {
	if c.comment == nil {
		return "", ""
	}
	return c.comment.author, c.comment.text
}

// addComment registers a comment on the sheet
func (s *Sheet) addComment(c *Comment) {
	// Replace existing comment for same cell
	for i, existing := range s.comments {
		if existing.row == c.row && existing.col == c.col {
			s.comments[i] = c
			return
		}
	}
	s.comments = append(s.comments, c)
}

// buildCommentsXML generates the xl/comments{n}.xml content
func buildCommentsXML(comments []*Comment) string {
	if len(comments) == 0 {
		return ""
	}

	// Collect unique authors
	authorMap := make(map[string]int)
	var authors []string
	for _, c := range comments {
		if _, ok := authorMap[c.author]; !ok {
			authorMap[c.author] = len(authors)
			authors = append(authors, c.author)
		}
	}

	var buf strings.Builder
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<comments xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">`)

	// Authors
	buf.WriteString(`<authors>`)
	for _, a := range authors {
		fmt.Fprintf(&buf, `<author>%s</author>`, escapeXMLText(a))
	}
	buf.WriteString(`</authors>`)

	// Comment list
	buf.WriteString(`<commentList>`)
	for _, c := range comments {
		ref := CellRef(c.row, c.col)
		authorIdx := authorMap[c.author]
		fmt.Fprintf(&buf, `<comment ref="%s" authorId="%d">`, ref, authorIdx)
		fmt.Fprintf(&buf, `<text><r><t>%s</t></r></text>`, escapeXMLText(c.text))
		buf.WriteString(`</comment>`)
	}
	buf.WriteString(`</commentList>`)

	buf.WriteString(`</comments>`)
	return buf.String()
}

// buildVMLDrawingXML generates the VML drawing XML for comment shapes
func buildVMLDrawingXML(comments []*Comment) string {
	if len(comments) == 0 {
		return ""
	}

	var buf strings.Builder
	buf.WriteString(`<xml xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel">`)

	for _, c := range comments {
		col := c.col - 1 // 0-based for VML
		row := c.row - 1
		fmt.Fprintf(&buf, `<v:shape type="#_x0000_t202" style="position:absolute;margin-left:0;margin-top:0;width:108pt;height:60pt;z-index:1;visibility:hidden" fillcolor="#ffffe1" o:insetmode="auto">`)
		fmt.Fprintf(&buf, `<v:fill color2="#ffffe1"/>`)
		fmt.Fprintf(&buf, `<v:shadow color="black" obscured="t"/>`)
		fmt.Fprintf(&buf, `<v:textbox/>`)
		fmt.Fprintf(&buf, `<x:ClientData ObjectType="Note">`)
		fmt.Fprintf(&buf, `<x:MoveWithCells/>`)
		fmt.Fprintf(&buf, `<x:SizeWithCells/>`)
		fmt.Fprintf(&buf, `<x:Anchor>%d,15,%d,10,%d,31,%d,4</x:Anchor>`, col+1, row, col+3, row+3)
		fmt.Fprintf(&buf, `<x:AutoFill>False</x:AutoFill>`)
		fmt.Fprintf(&buf, `<x:Row>%d</x:Row>`, row)
		fmt.Fprintf(&buf, `<x:Column>%d</x:Column>`, col)
		fmt.Fprintf(&buf, `</x:ClientData>`)
		buf.WriteString(`</v:shape>`)
	}

	buf.WriteString(`</xml>`)
	return buf.String()
}
