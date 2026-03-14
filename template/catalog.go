package template

import "github.com/JohnPitter/openscribe/style"

// Count returns total number of registered templates
func Count() int {
	return len(registry)
}

// Formats returns all available formats
func Formats() []Format {
	return []Format{FormatDOCX, FormatXLSX, FormatPPTX, FormatPDF}
}

// Categories returns all available categories
func Categories() []Category {
	return []Category{
		CategoryReport, CategoryInvoice, CategoryResume,
		CategoryLetter, CategoryDashboard, CategoryPitchDeck,
		CategoryNewsletter, CategoryCertificate,
	}
}

// Levels returns all design levels
func Levels() []style.DesignLevel {
	return []style.DesignLevel{
		style.DesignLevelBasic,
		style.DesignLevelProfessional,
		style.DesignLevelPremium,
		style.DesignLevelLuxury,
	}
}

// Search finds templates matching all given criteria
func Search(level *style.DesignLevel, format *Format, category *Category) []Template {
	var result []Template
	for _, t := range registry {
		if level != nil && t.Level != *level {
			continue
		}
		if format != nil && t.Format != *format {
			continue
		}
		if category != nil && t.Category != *category {
			continue
		}
		result = append(result, t)
	}
	return result
}
