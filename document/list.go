package document

import (
	"fmt"
)

// ListType defines the type of list numbering
type ListType int

const (
	// ListBullet is an unordered bullet list
	ListBullet ListType = iota
	// ListNumbered is an ordered numeric list (1, 2, 3...)
	ListNumbered
	// ListLettered is an ordered letter list (a, b, c...)
	ListLettered
	// ListRoman is an ordered roman numeral list (i, ii, iii...)
	ListRoman
)

// List represents an ordered or unordered list in the document
type List struct {
	listType ListType
	items    []*ListItem
	numID    int
}

// ListItem represents an item in a list
type ListItem struct {
	text     string
	level    int
	subItems []*ListItem
	numID    int
}

// NewList creates a new list with the given type and numbering ID
func NewList(listType ListType, numID int) *List {
	return &List{
		listType: listType,
		numID:    numID,
	}
}

// AddItem adds a top-level item to the list
func (l *List) AddItem(text string) *ListItem {
	item := &ListItem{
		text:  text,
		level: 0,
		numID: l.numID,
	}
	l.items = append(l.items, item)
	return item
}

// Items returns all top-level list items
func (l *List) Items() []*ListItem {
	return l.items
}

// Type returns the list type
func (l *List) Type() ListType {
	return l.listType
}

// NumID returns the numbering definition ID
func (l *List) NumID() int {
	return l.numID
}

// AddSubItem adds a nested sub-item under this list item
func (li *ListItem) AddSubItem(text string) *ListItem {
	sub := &ListItem{
		text:  text,
		level: li.level + 1,
		numID: li.numID,
	}
	li.subItems = append(li.subItems, sub)
	return sub
}

// Text returns the item text
func (li *ListItem) Text() string {
	return li.text
}

// Level returns the nesting level (0 = top)
func (li *ListItem) Level() int {
	return li.level
}

// SubItems returns nested sub-items
func (li *ListItem) SubItems() []*ListItem {
	return li.subItems
}

// toParagraphs converts the list item and its sub-items to paragraphs with numPr
func (li *ListItem) toParagraphs() []xmlParagraph {
	var result []xmlParagraph
	result = append(result, li.toXML())
	for _, sub := range li.subItems {
		result = append(result, sub.toParagraphs()...)
	}
	return result
}

// toXML converts a single list item to an xmlParagraph
func (li *ListItem) toXML() xmlParagraph {
	return xmlParagraph{
		Properties: &xmlParagraphProperties{
			NumPr: &xmlNumPr{
				Ilvl:  &xmlValue{Val: fmt.Sprintf("%d", li.level)},
				NumID: &xmlValue{Val: fmt.Sprintf("%d", li.numID)},
			},
		},
		Runs: []xmlRun{
			{
				Text: &xmlText{
					Space: "preserve",
					Value: li.text,
				},
			},
		},
	}
}

// toParagraphs converts the entire list to XML paragraphs
func (l *List) toParagraphs() []xmlParagraph {
	var result []xmlParagraph
	for _, item := range l.items {
		result = append(result, item.toParagraphs()...)
	}
	return result
}

// numFmtForType returns the OOXML number format string for a list type
func numFmtForType(lt ListType) string {
	switch lt {
	case ListBullet:
		return "bullet"
	case ListNumbered:
		return "decimal"
	case ListLettered:
		return "lowerLetter"
	case ListRoman:
		return "lowerRoman"
	default:
		return "decimal"
	}
}

// bulletCharForLevel returns the bullet character for a given nesting level
func bulletCharForLevel(level int) string {
	chars := []string{"\u2022", "\u25E6", "\u25AA"} // bullet, white bullet, small square
	return chars[level%len(chars)]
}

// buildNumberingXML creates the word/numbering.xml content for all lists
func buildNumberingXML(lists []*List) []byte {
	var buf []byte
	buf = append(buf, []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)...)
	buf = append(buf, []byte(fmt.Sprintf(`<w:numbering xmlns:w="%s">`, nsW))...)

	for _, l := range lists {
		// Abstract numbering definition
		buf = append(buf, []byte(fmt.Sprintf(`<w:abstractNum w:abstractNumId="%d">`, l.numID))...)
		for level := 0; level < 9; level++ {
			numFmt := numFmtForType(l.listType)
			var textVal string
			if l.listType == ListBullet {
				textVal = bulletCharForLevel(level)
			} else {
				textVal = fmt.Sprintf("%%%d.", level+1)
			}
			startVal := "1"
			buf = append(buf, []byte(fmt.Sprintf(
				`<w:lvl w:ilvl="%d"><w:start w:val="%s"/><w:numFmt w:val="%s"/><w:lvlText w:val="%s"/><w:lvlJc w:val="left"/><w:pPr><w:ind w:left="%d" w:hanging="360"/></w:pPr>`,
				level, startVal, numFmt, textVal, (level+1)*720,
			))...)
			if l.listType == ListBullet {
				buf = append(buf, []byte(`<w:rPr><w:rFonts w:ascii="Symbol" w:hAnsi="Symbol"/></w:rPr>`)...)
			}
			buf = append(buf, []byte(`</w:lvl>`)...)
		}
		buf = append(buf, []byte(`</w:abstractNum>`)...)

		// Numbering instance referencing the abstract definition
		buf = append(buf, []byte(fmt.Sprintf(
			`<w:num w:numId="%d"><w:abstractNumId w:val="%d"/></w:num>`,
			l.numID, l.numID,
		))...)
	}

	buf = append(buf, []byte(`</w:numbering>`)...)
	return buf
}
