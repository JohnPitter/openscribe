package pdf

import "github.com/JohnPitter/openscribe/common"

// SetSecurity sets document protection metadata on the PDF.
// Note: Full PDF encryption (RC4/AES) requires complex implementation.
// This sets protection flags and metadata that will be included in the trailer.
func (d *Document) SetSecurity(opts common.SecurityOptions) {
	d.security = &opts
}

// Security returns the current security options, or nil if none set.
func (d *Document) Security() *common.SecurityOptions {
	return d.security
}
