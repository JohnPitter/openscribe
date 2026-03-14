package template

import "github.com/JohnPitter/openscribe/style"

func registerBasicTemplates() {
	theme := style.BasicClean()

	register(Template{
		Name:        "Basic Report",
		Description: "Clean, minimal report template for everyday documentation",
		Category:    CategoryReport,
		Format:      FormatDOCX,
		Level:       style.DesignLevelBasic,
		Theme:       theme,
	})

	register(Template{
		Name:        "Basic Invoice",
		Description: "Simple invoice template with clean layout",
		Category:    CategoryInvoice,
		Format:      FormatDOCX,
		Level:       style.DesignLevelBasic,
		Theme:       theme,
	})

	register(Template{
		Name:        "Basic Spreadsheet",
		Description: "Minimal spreadsheet template for data entry",
		Category:    CategoryDashboard,
		Format:      FormatXLSX,
		Level:       style.DesignLevelBasic,
		Theme:       theme,
	})

	register(Template{
		Name:        "Basic Presentation",
		Description: "Clean slide deck for simple presentations",
		Category:    CategoryPitchDeck,
		Format:      FormatPPTX,
		Level:       style.DesignLevelBasic,
		Theme:       theme,
	})

	register(Template{
		Name:        "Basic PDF Report",
		Description: "Minimal PDF report template",
		Category:    CategoryReport,
		Format:      FormatPDF,
		Level:       style.DesignLevelBasic,
		Theme:       theme,
	})

	register(Template{
		Name:        "Basic Letter",
		Description: "Simple letter template",
		Category:    CategoryLetter,
		Format:      FormatDOCX,
		Level:       style.DesignLevelBasic,
		Theme:       theme,
	})

	register(Template{
		Name:        "Basic Resume",
		Description: "Clean resume template",
		Category:    CategoryResume,
		Format:      FormatDOCX,
		Level:       style.DesignLevelBasic,
		Theme:       theme,
	})

	register(Template{
		Name:        "Basic Certificate",
		Description: "Simple certificate template",
		Category:    CategoryCertificate,
		Format:      FormatPDF,
		Level:       style.DesignLevelBasic,
		Theme:       theme,
	})
}
