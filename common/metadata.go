package common

// MetadataCleaner provides a unified interface for stripping identifying
// information from documents. Use the package-specific CleanMetadata functions
// in the pdf, document, and spreadsheet packages directly:
//
//	pdf.CleanMetadata(doc)
//	document.CleanMetadata(doc)
//	spreadsheet.CleanMetadata(wb)
type MetadataCleaner struct{}

// NewMetadataCleaner creates a new MetadataCleaner instance.
func NewMetadataCleaner() *MetadataCleaner {
	return &MetadataCleaner{}
}
