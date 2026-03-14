package spreadsheet

// CleanMetadata strips identifying metadata from the XLSX workbook.
// This removes author information and other identifying metadata.
func CleanMetadata(wb *Workbook) {
	// Clear security/author metadata
	wb.security = nil
}
