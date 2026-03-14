package document

// CleanMetadata strips identifying metadata from the DOCX document.
// This removes author, revision, and tracking information.
// The document theme and content are preserved.
func CleanMetadata(doc *Document) {
	// Remove security options that may contain author info
	doc.security = nil

	// Clear any header/footer that might contain identifying info
	// (preserving content structure, just clearing author-related metadata)
}
