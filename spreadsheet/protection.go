package spreadsheet

import (
	"fmt"
	"strings"
)

// ProtectionOptions configures what operations are allowed on a protected sheet
type ProtectionOptions struct {
	AllowSort        bool
	AllowFilter      bool
	AllowInsertRows  bool
	AllowDeleteRows  bool
	AllowFormatCells bool
}

// SheetProtection represents protection settings for a sheet
type SheetProtection struct {
	password string
	options  ProtectionOptions
}

// Protect password-protects a sheet
func (s *Sheet) Protect(password string) {
	if s.protection == nil {
		s.protection = &SheetProtection{}
	}
	s.protection.password = password
}

// SetProtectionOptions configures what operations are allowed on a protected sheet
func (s *Sheet) SetProtectionOptions(opts ProtectionOptions) {
	if s.protection == nil {
		s.protection = &SheetProtection{}
	}
	s.protection.options = opts
}

// SetCellLocked sets whether a specific cell is locked
func (s *Sheet) SetCellLocked(row, col int, locked bool) {
	cell := s.Cell(row, col)
	cell.locked = &locked
}

// hashPassword creates a simple hash for sheet protection (OOXML legacy hash)
func hashPassword(password string) string {
	if password == "" {
		return ""
	}
	var hash uint16
	for i := len(password) - 1; i >= 0; i-- {
		char := uint16(password[i])
		hash = ((hash >> 14) & 0x01) | ((hash << 1) & 0x7FFF)
		hash ^= char
	}
	hash = ((hash >> 14) & 0x01) | ((hash << 1) & 0x7FFF)
	hash ^= uint16(len(password))
	hash ^= 0xCE4B
	return fmt.Sprintf("%04X", hash)
}

// buildProtectionXML generates the <sheetProtection> XML element
func buildProtectionXML(prot *SheetProtection) string {
	if prot == nil {
		return ""
	}

	var buf strings.Builder
	buf.WriteString(`<sheetProtection sheet="1"`)

	if prot.password != "" {
		fmt.Fprintf(&buf, ` password="%s"`, hashPassword(prot.password))
	}

	// By default, everything is locked. These attributes allow specific operations.
	if prot.options.AllowSort {
		buf.WriteString(` sort="0"`)
	}
	if prot.options.AllowFilter {
		buf.WriteString(` autoFilter="0"`)
	}
	if prot.options.AllowInsertRows {
		buf.WriteString(` insertRows="0"`)
	}
	if prot.options.AllowDeleteRows {
		buf.WriteString(` deleteRows="0"`)
	}
	if prot.options.AllowFormatCells {
		buf.WriteString(` formatCells="0"`)
	}

	buf.WriteString(`/>`)
	return buf.String()
}
