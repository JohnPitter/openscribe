package template

import "github.com/JohnPitter/openscribe/style"

func registerProfessionalTemplates() {
	theme := style.ProfessionalCorporate()

	register(Template{
		Name:        "Corporate Report",
		Description: "Business-grade report with corporate styling and structured sections",
		Category:    CategoryReport,
		Format:      FormatDOCX,
		Level:       style.DesignLevelProfessional,
		Theme:       theme,
	})

	register(Template{
		Name:        "Corporate Invoice",
		Description: "Professional invoice with company branding areas",
		Category:    CategoryInvoice,
		Format:      FormatDOCX,
		Level:       style.DesignLevelProfessional,
		Theme:       theme,
	})

	register(Template{
		Name:        "Financial Dashboard",
		Description: "Professional spreadsheet for financial data and KPIs",
		Category:    CategoryDashboard,
		Format:      FormatXLSX,
		Level:       style.DesignLevelProfessional,
		Theme:       theme,
	})

	register(Template{
		Name:        "Business Presentation",
		Description: "Corporate slide deck for stakeholder meetings",
		Category:    CategoryPitchDeck,
		Format:      FormatPPTX,
		Level:       style.DesignLevelProfessional,
		Theme:       theme,
	})

	register(Template{
		Name:        "Professional PDF Report",
		Description: "Business-grade PDF with structured layout",
		Category:    CategoryReport,
		Format:      FormatPDF,
		Level:       style.DesignLevelProfessional,
		Theme:       theme,
	})

	register(Template{
		Name:        "Business Letter",
		Description: "Formal business letter with corporate styling",
		Category:    CategoryLetter,
		Format:      FormatDOCX,
		Level:       style.DesignLevelProfessional,
		Theme:       theme,
	})

	register(Template{
		Name:        "Professional Resume",
		Description: "Career-focused resume with corporate aesthetics",
		Category:    CategoryResume,
		Format:      FormatDOCX,
		Level:       style.DesignLevelProfessional,
		Theme:       theme,
	})

	register(Template{
		Name:        "Company Newsletter",
		Description: "Internal newsletter template for organizations",
		Category:    CategoryNewsletter,
		Format:      FormatDOCX,
		Level:       style.DesignLevelProfessional,
		Theme:       theme,
	})
}
