package packaging

import (
	"encoding/xml"

	"github.com/JohnPitter/openscribe/internal/xmlutil"
)

// Common content types
const (
	ContentTypeRelationships = "application/vnd.openxmlformats-package.relationships+xml"
	ContentTypeXML           = "application/xml"
	ContentTypeDocx          = "application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"
	ContentTypeStyles        = "application/vnd.openxmlformats-officedocument.wordprocessingml.styles+xml"
	ContentTypeXlsx          = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"
	ContentTypeWorksheet     = "application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"
	ContentTypeSharedStrings = "application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml"
	ContentTypePptx          = "application/vnd.openxmlformats-officedocument.presentationml.presentation.main+xml"
	ContentTypeSlide         = "application/vnd.openxmlformats-officedocument.presentationml.slide+xml"
	ContentTypeSlideMaster   = "application/vnd.openxmlformats-officedocument.presentationml.slideMaster+xml"
	ContentTypeSlideLayout   = "application/vnd.openxmlformats-officedocument.presentationml.slideLayout+xml"
	ContentTypePNG           = "image/png"
	ContentTypeJPEG          = "image/jpeg"
	ContentTypeCoreProps     = "application/vnd.openxmlformats-package.core-properties+xml"
)

// ContentTypes represents [Content_Types].xml
type ContentTypes struct {
	XMLName   xml.Name   `xml:"Types"`
	Xmlns     string     `xml:"xmlns,attr"`
	Defaults  []Default  `xml:"Default"`
	Overrides []Override `xml:"Override"`
}

// Default maps a file extension to a content type
type Default struct {
	Extension   string `xml:"Extension,attr"`
	ContentType string `xml:"ContentType,attr"`
}

// Override maps a specific part name to a content type
type Override struct {
	PartName    string `xml:"PartName,attr"`
	ContentType string `xml:"ContentType,attr"`
}

// NewContentTypes creates content types with common defaults
func NewContentTypes() *ContentTypes {
	return &ContentTypes{
		Xmlns: "http://schemas.openxmlformats.org/package/2006/content-types",
		Defaults: []Default{
			{Extension: "rels", ContentType: ContentTypeRelationships},
			{Extension: "xml", ContentType: ContentTypeXML},
			{Extension: "png", ContentType: ContentTypePNG},
			{Extension: "jpeg", ContentType: ContentTypeJPEG},
		},
	}
}

// AddOverride adds a part-specific content type override
func (ct *ContentTypes) AddOverride(partName, contentType string) {
	ct.Overrides = append(ct.Overrides, Override{
		PartName:    partName,
		ContentType: contentType,
	})
}

// Marshal serializes the content types to XML
func (ct *ContentTypes) Marshal() ([]byte, error) {
	return xmlutil.MarshalXML(ct)
}
