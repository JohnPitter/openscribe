package pdf

// CleanMetadata strips identifying metadata from the PDF document,
// clearing author, title, subject fields while preserving the creator tag.
func CleanMetadata(doc *Document) {
	doc.metadata = Metadata{
		Creator: "OpenScribe",
	}
}
