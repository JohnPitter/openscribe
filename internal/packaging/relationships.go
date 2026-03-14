package packaging

import (
	"encoding/xml"
	"fmt"

	"github.com/JohnPitter/openscribe/internal/xmlutil"
)

// Common relationship types
const (
	RelTypeOfficeDocument = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
	RelTypeStylesheet     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles"
	RelTypeImage          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	RelTypeHyperlink      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink"
	RelTypeTheme          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/theme"
	RelTypeSharedStrings  = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings"
	RelTypeWorksheet      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet"
	RelTypeSlide          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slide"
	RelTypeSlideMaster    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideMaster"
	RelTypeSlideLayout    = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout"
	RelTypeCoreProperties = "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties"
	RelTypeDrawing        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/drawing"
	RelTypeChart          = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/chart"
	RelTypeHeader         = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/header"
	RelTypeFooter         = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer"
	RelTypeNumbering      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering"
	RelTypeFootnotes      = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/footnotes"
	RelTypeComments       = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments"
	RelTypeVMLDrawing     = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/vmlDrawing"
)

// Relationships represents the _rels/.rels file
type Relationships struct {
	XMLName       xml.Name       `xml:"Relationships"`
	Xmlns         string         `xml:"xmlns,attr"`
	Relationships []Relationship `xml:"Relationship"`
}

// Relationship represents a single relationship
type Relationship struct {
	XMLName    xml.Name `xml:"Relationship"`
	ID         string   `xml:"Id,attr"`
	Type       string   `xml:"Type,attr"`
	Target     string   `xml:"Target,attr"`
	TargetMode string   `xml:"TargetMode,attr,omitempty"`
}

// NewRelationships creates a new relationships container
func NewRelationships() *Relationships {
	return &Relationships{
		Xmlns: "http://schemas.openxmlformats.org/package/2006/relationships",
	}
}

// Add adds a relationship and returns its ID
func (r *Relationships) Add(relType, target string) string {
	id := fmt.Sprintf("rId%d", len(r.Relationships)+1)
	r.Relationships = append(r.Relationships, Relationship{
		ID:     id,
		Type:   relType,
		Target: target,
	})
	return id
}

// AddExternal adds an external relationship (e.g., hyperlink) and returns its ID
func (r *Relationships) AddExternal(relType, target string) string {
	id := fmt.Sprintf("rId%d", len(r.Relationships)+1)
	r.Relationships = append(r.Relationships, Relationship{
		ID:         id,
		Type:       relType,
		Target:     target,
		TargetMode: "External",
	})
	return id
}

// Marshal serializes the relationships to XML
func (r *Relationships) Marshal() ([]byte, error) {
	return xmlutil.MarshalXML(r)
}

// UnmarshalRelationships deserializes relationships from XML
func UnmarshalRelationships(data []byte) (*Relationships, error) {
	var rels Relationships
	if err := xmlutil.UnmarshalXML(data, &rels); err != nil {
		return nil, err
	}
	return &rels, nil
}
