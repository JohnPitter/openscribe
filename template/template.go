// Package template provides pre-built document templates at various design levels.
package template

import "github.com/JohnPitter/openscribe/style"

// Format represents the output document format
type Format int

const (
	FormatDOCX Format = iota
	FormatXLSX
	FormatPPTX
	FormatPDF
)

func (f Format) String() string {
	switch f {
	case FormatDOCX:
		return "DOCX"
	case FormatXLSX:
		return "XLSX"
	case FormatPPTX:
		return "PPTX"
	case FormatPDF:
		return "PDF"
	default:
		return "Unknown"
	}
}

// Category represents the template category
type Category int

const (
	CategoryReport Category = iota
	CategoryInvoice
	CategoryResume
	CategoryLetter
	CategoryDashboard
	CategoryPitchDeck
	CategoryNewsletter
	CategoryCertificate
)

func (c Category) String() string {
	switch c {
	case CategoryReport:
		return "Report"
	case CategoryInvoice:
		return "Invoice"
	case CategoryResume:
		return "Resume"
	case CategoryLetter:
		return "Letter"
	case CategoryDashboard:
		return "Dashboard"
	case CategoryPitchDeck:
		return "Pitch Deck"
	case CategoryNewsletter:
		return "Newsletter"
	case CategoryCertificate:
		return "Certificate"
	default:
		return "Unknown"
	}
}

// Template describes a document template
type Template struct {
	Name        string
	Description string
	Category    Category
	Format      Format
	Level       style.DesignLevel
	Theme       style.Theme
}

// registry holds all registered templates
var registry []Template

func init() {
	registerBasicTemplates()
	registerProfessionalTemplates()
	registerPremiumTemplates()
	registerLuxuryTemplates()
}

// All returns all registered templates
func All() []Template {
	return registry
}

// ByLevel returns templates at a given design level
func ByLevel(level style.DesignLevel) []Template {
	var result []Template
	for _, t := range registry {
		if t.Level == level {
			result = append(result, t)
		}
	}
	return result
}

// ByFormat returns templates for a given format
func ByFormat(format Format) []Template {
	var result []Template
	for _, t := range registry {
		if t.Format == format {
			result = append(result, t)
		}
	}
	return result
}

// ByCategory returns templates in a given category
func ByCategory(category Category) []Template {
	var result []Template
	for _, t := range registry {
		if t.Category == category {
			result = append(result, t)
		}
	}
	return result
}

// Find searches for a template by name
func Find(name string) *Template {
	for _, t := range registry {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

func register(t Template) {
	registry = append(registry, t)
}
