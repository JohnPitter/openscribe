package spreadsheet

import "github.com/JohnPitter/openscribe/common"

// SetSecurity sets workbook-level protection metadata.
// This complements sheet-level protection with workbook-level settings.
func (wb *Workbook) SetSecurity(opts common.SecurityOptions) {
	wb.security = &opts
}

// Security returns the current security options, or nil if none set.
func (wb *Workbook) Security() *common.SecurityOptions {
	return wb.security
}
