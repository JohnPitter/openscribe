package template

import "github.com/JohnPitter/openscribe/style"

func registerLuxuryTemplates() {
	agency := style.LuxuryAgency()
	warm := style.LuxuryWarm()

	// Agency Luxury
	register(Template{
		Name:        "Agency Pitch Deck",
		Description: "Slidesgo-quality pitch deck with bold neon accents on dark background",
		Category:    CategoryPitchDeck,
		Format:      FormatPPTX,
		Level:       style.DesignLevelLuxury,
		Theme:       agency,
	})

	register(Template{
		Name:        "Agency Report",
		Description: "High-impact report with tech-forward design language",
		Category:    CategoryReport,
		Format:      FormatDOCX,
		Level:       style.DesignLevelLuxury,
		Theme:       agency,
	})

	register(Template{
		Name:        "Tech Dashboard",
		Description: "Data-rich dashboard with dark theme and neon accents",
		Category:    CategoryDashboard,
		Format:      FormatXLSX,
		Level:       style.DesignLevelLuxury,
		Theme:       agency,
	})

	register(Template{
		Name:        "Agency Invoice",
		Description: "Premium invoice with bold design and striking typography",
		Category:    CategoryInvoice,
		Format:      FormatPDF,
		Level:       style.DesignLevelLuxury,
		Theme:       agency,
	})

	// Warm Luxury
	register(Template{
		Name:        "Executive Presentation",
		Description: "Warm, sophisticated deck for C-suite presentations",
		Category:    CategoryPitchDeck,
		Format:      FormatPPTX,
		Level:       style.DesignLevelLuxury,
		Theme:       warm,
	})

	register(Template{
		Name:        "Executive Report",
		Description: "Premium report with warm palette and elegant typography",
		Category:    CategoryReport,
		Format:      FormatDOCX,
		Level:       style.DesignLevelLuxury,
		Theme:       warm,
	})

	register(Template{
		Name:        "Luxury Certificate",
		Description: "Award-grade certificate with premium warm aesthetics",
		Category:    CategoryCertificate,
		Format:      FormatPDF,
		Level:       style.DesignLevelLuxury,
		Theme:       warm,
	})

	register(Template{
		Name:        "Creative Resume",
		Description: "Unique resume design with warm luxury feel",
		Category:    CategoryResume,
		Format:      FormatDOCX,
		Level:       style.DesignLevelLuxury,
		Theme:       warm,
	})
}
