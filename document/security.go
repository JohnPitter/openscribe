package document

import "github.com/JohnPitter/openscribe/common"

// SetSecurity sets document protection metadata on the DOCX document.
// This adds basic document protection settings.
func (d *Document) SetSecurity(opts common.SecurityOptions) {
	d.security = &opts
}

// Security returns the current security options, or nil if none set.
func (d *Document) Security() *common.SecurityOptions {
	return d.security
}
