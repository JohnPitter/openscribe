package template

import "github.com/JohnPitter/openscribe/style"

func registerPremiumTemplates() {
	modern := style.PremiumModern()
	elegant := style.PremiumElegant()

	// Modern Premium
	register(Template{
		Name:        "Modern Report",
		Description: "Behance-quality report with bold typography and accent colors",
		Category:    CategoryReport,
		Format:      FormatDOCX,
		Level:       style.DesignLevelPremium,
		Theme:       modern,
	})

	register(Template{
		Name:        "Modern Invoice",
		Description: "Design-forward invoice with gradient accents and modern layout",
		Category:    CategoryInvoice,
		Format:      FormatPDF,
		Level:       style.DesignLevelPremium,
		Theme:       modern,
	})

	register(Template{
		Name:        "Analytics Dashboard",
		Description: "Freepik-quality data dashboard with charts area",
		Category:    CategoryDashboard,
		Format:      FormatXLSX,
		Level:       style.DesignLevelPremium,
		Theme:       modern,
	})

	register(Template{
		Name:        "Startup Pitch Deck",
		Description: "VC-ready pitch deck with modern visual language",
		Category:    CategoryPitchDeck,
		Format:      FormatPPTX,
		Level:       style.DesignLevelPremium,
		Theme:       modern,
	})

	// Elegant Premium
	register(Template{
		Name:        "Elegant Report",
		Description: "Sophisticated report with serif typography and refined palette",
		Category:    CategoryReport,
		Format:      FormatDOCX,
		Level:       style.DesignLevelPremium,
		Theme:       elegant,
	})

	register(Template{
		Name:        "Elegant Resume",
		Description: "Refined resume with classic typography and subtle accents",
		Category:    CategoryResume,
		Format:      FormatDOCX,
		Level:       style.DesignLevelPremium,
		Theme:       elegant,
	})

	register(Template{
		Name:        "Premium Certificate",
		Description: "Beautifully designed certificate with premium feel",
		Category:    CategoryCertificate,
		Format:      FormatPDF,
		Level:       style.DesignLevelPremium,
		Theme:       elegant,
	})

	register(Template{
		Name:        "Premium Newsletter",
		Description: "High-quality newsletter with magazine-style layout",
		Category:    CategoryNewsletter,
		Format:      FormatDOCX,
		Level:       style.DesignLevelPremium,
		Theme:       elegant,
	})
}
