// Command openscribe provides a CLI for creating office documents from templates.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/JohnPitter/openscribe/style"
	"github.com/JohnPitter/openscribe/template"
)

var (
	listTemplates = flag.Bool("list", false, "List all available templates")
	listThemes    = flag.Bool("themes", false, "List all available themes")
	tmplName      = flag.String("template", "", "Template name to generate")
	themeName     = flag.String("theme", "", "Theme to apply (optional)")
	output        = flag.String("output", "", "Output file path")
	formatFilter  = flag.String("format", "", "Filter templates by format (DOCX, XLSX, PPTX, PDF)")
	levelFilter   = flag.String("level", "", "Filter templates by level (Basic, Professional, Premium, Luxury)")
	version       = flag.Bool("version", false, "Show version")
)

const appVersion = "0.2.0"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "OpenScribe — Pure Go office document library\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  openscribe [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  openscribe -list                           # List all templates\n")
		fmt.Fprintf(os.Stderr, "  openscribe -list -format DOCX              # List DOCX templates\n")
		fmt.Fprintf(os.Stderr, "  openscribe -list -level Premium            # List premium templates\n")
		fmt.Fprintf(os.Stderr, "  openscribe -themes                         # List all themes\n")
		fmt.Fprintf(os.Stderr, "  openscribe -template \"Agency Pitch Deck\" -output pitch.pptx\n")
		fmt.Fprintf(os.Stderr, "  openscribe -template \"Basic Report\" -output report.docx\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Printf("openscribe v%s\n", appVersion)
		return
	}

	if *listThemes {
		printThemes()
		return
	}

	if *listTemplates {
		printTemplates()
		return
	}

	if *tmplName != "" {
		generateFromTemplate()
		return
	}

	flag.Usage()
}

func printThemes() {
	fmt.Println("Available Themes:")
	fmt.Println()
	fmt.Printf("  %-25s %-15s %s\n", "NAME", "LEVEL", "PALETTE")
	fmt.Printf("  %-25s %-15s %s\n", "----", "-----", "-------")
	for _, t := range style.AllThemes() {
		fmt.Printf("  %-25s %-15s Primary: %s\n", t.Name, t.Level.String(), t.Palette.Primary.Hex())
	}
}

func printTemplates() {
	templates := template.All()

	// Apply filters
	if *formatFilter != "" {
		var filtered []template.Template
		for _, t := range templates {
			if strings.EqualFold(t.Format.String(), *formatFilter) {
				filtered = append(filtered, t)
			}
		}
		templates = filtered
	}

	if *levelFilter != "" {
		var filtered []template.Template
		for _, t := range templates {
			if strings.EqualFold(t.Level.String(), *levelFilter) {
				filtered = append(filtered, t)
			}
		}
		templates = filtered
	}

	fmt.Printf("Templates (%d):\n\n", len(templates))
	fmt.Printf("  %-30s %-6s %-15s %s\n", "NAME", "FORMAT", "LEVEL", "DESCRIPTION")
	fmt.Printf("  %-30s %-6s %-15s %s\n", "----", "------", "-----", "-----------")
	for _, t := range templates {
		desc := t.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		fmt.Printf("  %-30s %-6s %-15s %s\n", t.Name, t.Format.String(), t.Level.String(), desc)
	}
}

func generateFromTemplate() {
	tmpl := template.Find(*tmplName)
	if tmpl == nil {
		fmt.Fprintf(os.Stderr, "Error: template %q not found\n", *tmplName)
		fmt.Fprintf(os.Stderr, "Run 'openscribe -list' to see available templates\n")
		os.Exit(1)
	}

	// Apply custom theme if specified
	if *themeName != "" {
		for _, t := range style.AllThemes() {
			if strings.EqualFold(t.Name, *themeName) {
				tmpl.Theme = t
				break
			}
		}
	}

	// Determine output path
	outPath := *output
	if outPath == "" {
		ext := ".docx"
		switch tmpl.Format {
		case template.FormatXLSX:
			ext = ".xlsx"
		case template.FormatPPTX:
			ext = ".pptx"
		case template.FormatPDF:
			ext = ".pdf"
		}
		outPath = "output" + ext
	}

	var err error
	switch tmpl.Format {
	case template.FormatDOCX:
		doc, genErr := tmpl.GenerateDocx()
		if genErr != nil {
			fmt.Fprintf(os.Stderr, "Error generating: %v\n", genErr)
			os.Exit(1)
		}
		err = doc.Save(outPath)
	case template.FormatXLSX:
		wb, genErr := tmpl.GenerateXlsx()
		if genErr != nil {
			fmt.Fprintf(os.Stderr, "Error generating: %v\n", genErr)
			os.Exit(1)
		}
		err = wb.Save(outPath)
	case template.FormatPPTX:
		pres, genErr := tmpl.GeneratePptx()
		if genErr != nil {
			fmt.Fprintf(os.Stderr, "Error generating: %v\n", genErr)
			os.Exit(1)
		}
		err = pres.Save(outPath)
	case template.FormatPDF:
		doc, genErr := tmpl.GeneratePdf()
		if genErr != nil {
			fmt.Fprintf(os.Stderr, "Error generating: %v\n", genErr)
			os.Exit(1)
		}
		err = doc.Save(outPath)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated: %s (template: %s, theme: %s)\n", outPath, tmpl.Name, tmpl.Theme.Name)
}
