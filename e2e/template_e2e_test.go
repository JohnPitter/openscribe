package e2e

import (
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/document"
	"github.com/JohnPitter/openscribe/pdf"
	"github.com/JohnPitter/openscribe/presentation"
	"github.com/JohnPitter/openscribe/spreadsheet"
	"github.com/JohnPitter/openscribe/style"
	"github.com/JohnPitter/openscribe/template"
)

func TestTemplateRegistration(t *testing.T) {
	all := template.All()
	if len(all) != 32 {
		t.Errorf("expected 32 templates, got %d", len(all))
	}
}

func TestAllTemplatesHaveRequiredFields(t *testing.T) {
	for _, tmpl := range template.All() {
		if tmpl.Name == "" {
			t.Error("template name should not be empty")
		}
		if tmpl.Description == "" {
			t.Errorf("template %s description should not be empty", tmpl.Name)
		}
		if tmpl.Theme.Name == "" {
			t.Errorf("template %s should have a theme name", tmpl.Name)
		}
	}
}

func TestTemplatesByLevel(t *testing.T) {
	levels := template.Levels()
	for _, level := range levels {
		templates := template.ByLevel(level)
		if len(templates) == 0 {
			t.Errorf("no templates for level %s", level.String())
		}
		for _, tmpl := range templates {
			if tmpl.Level != level {
				t.Errorf("template %s has wrong level: got %s, want %s",
					tmpl.Name, tmpl.Level.String(), level.String())
			}
		}
	}
}

func TestTemplatesByFormat(t *testing.T) {
	formats := template.Formats()
	for _, format := range formats {
		templates := template.ByFormat(format)
		if len(templates) == 0 {
			t.Errorf("no templates for format %s", format.String())
		}
		for _, tmpl := range templates {
			if tmpl.Format != format {
				t.Errorf("template %s has wrong format", tmpl.Name)
			}
		}
	}
}

func TestTemplatesByCategory(t *testing.T) {
	categories := template.Categories()
	for _, cat := range categories {
		templates := template.ByCategory(cat)
		if len(templates) == 0 {
			t.Errorf("no templates for category %s", cat.String())
		}
	}
}

func TestTemplateSearch(t *testing.T) {
	// Luxury PPTX presentations
	level := style.DesignLevelLuxury
	format := template.FormatPPTX
	results := template.Search(&level, &format, nil)
	if len(results) == 0 {
		t.Error("should find luxury PPTX templates")
	}
}

func TestTemplateFind(t *testing.T) {
	tmpl := template.Find("Agency Pitch Deck")
	if tmpl == nil {
		t.Fatal("should find template")
	}
	if tmpl.Level != style.DesignLevelLuxury {
		t.Error("wrong level")
	}
}

// Test that every DOCX template produces a valid document
func TestAllDocxTemplatesGenerate(t *testing.T) {
	docxTemplates := template.ByFormat(template.FormatDOCX)

	for _, tmpl := range docxTemplates {
		t.Run(tmpl.Name, func(t *testing.T) {
			doc := document.NewWithTheme(tmpl.Theme)
			doc.AddHeading(tmpl.Name, 1)
			doc.AddText(tmpl.Description)
			doc.AddHeading("Section", 2)
			doc.AddText("Body content with theme styling.")

			tbl := doc.AddTable(3, 2)
			tbl.Cell(0, 0).SetText("Key")
			tbl.Cell(0, 1).SetText("Value")
			tbl.Cell(1, 0).SetText("Quality")
			tbl.Cell(1, 1).SetText(tmpl.Level.String())
			tbl.Cell(2, 0).SetText("Theme")
			tbl.Cell(2, 1).SetText(tmpl.Theme.Name)

			path := filepath.Join(t.TempDir(), "template.docx")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
			assertFileExists(t, path)
			assertFileNotEmpty(t, path)
		})
	}
}

// Test that every XLSX template produces a valid spreadsheet
func TestAllXlsxTemplatesGenerate(t *testing.T) {
	xlsxTemplates := template.ByFormat(template.FormatXLSX)

	for _, tmpl := range xlsxTemplates {
		t.Run(tmpl.Name, func(t *testing.T) {
			wb := spreadsheet.NewWithTheme(tmpl.Theme)
			s := wb.AddSheet("Data")
			s.SetValue(1, 1, "Template")
			s.SetValue(1, 2, tmpl.Name)
			s.SetValue(2, 1, "Level")
			s.SetValue(2, 2, tmpl.Level.String())
			s.SetValue(3, 1, "Theme")
			s.SetValue(3, 2, tmpl.Theme.Name)
			s.SetValue(4, 1, "Value")
			s.SetValue(4, 2, 42.0)
			s.Cell(5, 2).SetFormula("SUM(B4:B4)")

			path := filepath.Join(t.TempDir(), "template.xlsx")
			if err := wb.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
			assertFileExists(t, path)
			assertFileNotEmpty(t, path)
		})
	}
}

// Test that every PPTX template produces a valid presentation
func TestAllPptxTemplatesGenerate(t *testing.T) {
	pptxTemplates := template.ByFormat(template.FormatPPTX)

	for _, tmpl := range pptxTemplates {
		t.Run(tmpl.Name, func(t *testing.T) {
			pres := presentation.NewWithTheme(tmpl.Theme)

			// Title slide
			s1 := pres.AddSlide()
			s1.SetBackground(tmpl.Theme.Palette.Primary)
			tb := s1.AddTextBox(common.In(2), common.In(2.5), common.In(9), common.In(2))
			tb.SetText(tmpl.Name, tmpl.Theme.Typography.HeadingFont.WithColor(common.White))

			// Content slide
			s2 := pres.AddSlide()
			content := s2.AddTextBox(common.In(1), common.In(1), common.In(10), common.In(5))
			p := content.AddParagraph()
			p.AddRun("Design Level: "+tmpl.Level.String(), tmpl.Theme.Typography.BodyFont)

			path := filepath.Join(t.TempDir(), "template.pptx")
			if err := pres.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
			assertFileExists(t, path)
			assertFileNotEmpty(t, path)
		})
	}
}

// Test that every PDF template produces a valid PDF
func TestAllPdfTemplatesGenerate(t *testing.T) {
	pdfTemplates := template.ByFormat(template.FormatPDF)

	for _, tmpl := range pdfTemplates {
		t.Run(tmpl.Name, func(t *testing.T) {
			doc := pdf.NewWithTheme(tmpl.Theme)
			p := doc.AddPage()

			p.AddText(tmpl.Name, 72, 72, tmpl.Theme.Typography.HeadingFont)
			p.AddText(tmpl.Description, 72, 110, tmpl.Theme.Typography.BodyFont)
			p.AddText("Level: "+tmpl.Level.String(), 72, 140,
				tmpl.Theme.Typography.CaptionFont)
			p.AddRectangle(72, 160, 400, 2, tmpl.Theme.Palette.Primary, nil)

			tbl := p.AddTable(72, 200, 3, 2)
			tbl.SetCellSize(200, 25)
			tbl.SetCell(0, 0, "Attribute")
			tbl.SetCell(0, 1, "Value")
			tbl.SetCell(1, 0, "Theme")
			tbl.SetCell(1, 1, tmpl.Theme.Name)
			tbl.SetCell(2, 0, "Quality")
			tbl.SetCell(2, 1, tmpl.Level.String())

			path := filepath.Join(t.TempDir(), "template.pdf")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
			assertFileExists(t, path)
			assertFileNotEmpty(t, path)
			assertPDFHeader(t, path)
		})
	}
}
