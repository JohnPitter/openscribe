package document

import (
	"fmt"
	"time"
)

// Comment represents a comment annotation in the document
type Comment struct {
	id     int
	author string
	text   string
	date   string
}

// NewComment creates a new comment with auto-generated date
func NewComment(id int, author, text string) *Comment {
	return &Comment{
		id:     id,
		author: author,
		text:   text,
		date:   time.Now().Format("2006-01-02T15:04:00Z"),
	}
}

// ID returns the comment ID
func (c *Comment) ID() int {
	return c.id
}

// Author returns the comment author
func (c *Comment) Author() string {
	return c.author
}

// Text returns the comment text
func (c *Comment) Text() string {
	return c.text
}

// Date returns the comment date string
func (c *Comment) Date() string {
	return c.date
}

// buildCommentsXML creates the word/comments.xml content
func buildCommentsXML(comments []*Comment) []byte {
	var buf []byte
	buf = append(buf, []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)...)
	buf = append(buf, []byte(fmt.Sprintf(
		`<w:comments xmlns:w="%s" xmlns:r="%s">`, nsW, nsR,
	))...)

	for _, c := range comments {
		buf = append(buf, []byte(fmt.Sprintf(
			`<w:comment w:id="%d" w:author="%s" w:date="%s">`+
				`<w:p><w:r><w:t>%s</w:t></w:r></w:p>`+
				`</w:comment>`,
			c.id, c.author, c.date, c.text,
		))...)
	}

	buf = append(buf, []byte(`</w:comments>`)...)
	return buf
}
