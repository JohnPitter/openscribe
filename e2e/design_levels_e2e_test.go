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
)

// TestDesignLevelDocx tests that every theme produces valid DOCX at each design level
func TestDesignLevelDocx(t *testing.T) {
	for _, theme := range style.AllThemes() {
		t.Run(theme.Name, func(t *testing.T) {
			doc := document.NewWithTheme(theme)

			// Rich document structure
			doc.AddHeading("Report: "+theme.Name, 1)
			doc.AddText("Generated with " + theme.Level.String() + " design level.")

			doc.AddHeading("Executive Summary", 2)
			p := doc.AddParagraph()
			p.SetAlignment(common.TextAlignJustify)
			r := p.AddRun()
			r.SetFont(theme.Typography.BodyFont)
			r.SetText("This document demonstrates the visual quality and consistency of the " + theme.Name + " theme applied to a standard business report format.")

			doc.AddHeading("Data Overview", 2)
			tbl := doc.AddTable(4, 3)
			tbl.Cell(0, 0).SetText("Metric")
			tbl.Cell(0, 1).SetText("Current")
			tbl.Cell(0, 2).SetText("Target")
			tbl.Cell(1, 0).SetText("Revenue")
			tbl.Cell(1, 1).SetText("$1.2M")
			tbl.Cell(1, 2).SetText("$1.5M")
			tbl.Cell(2, 0).SetText("Users")
			tbl.Cell(2, 1).SetText("45,000")
			tbl.Cell(2, 2).SetText("60,000")
			tbl.Cell(3, 0).SetText("NPS")
			tbl.Cell(3, 1).SetText("72")
			tbl.Cell(3, 2).SetText("80")

			// Header shading from theme
			for col := 0; col < 3; col++ {
				tbl.Cell(0, col).SetShading(theme.Palette.Primary)
			}

			doc.AddPageBreak()
			doc.AddHeading("Conclusion", 2)
			doc.AddText("The design level ensures professional and consistent output.")

			path := filepath.Join(t.TempDir(), "design_level.docx")
			if err := doc.Save(path); err != nil {
				t.Fatalf("error: %v", err)
			}
			assertFileExists(t, path)
			assertFileNotEmpty(t, path)
		})
	}
}

// TestDesignLevelXlsx tests all themes with XLSX
func TestDesignLevelXlsx(t *testing.T) {
	for _, theme := range style.AllThemes() {
		t.Run(theme.Name, func(t *testing.T) {
			wb := spreadsheet.NewWithTheme(theme)
			s := wb.AddSheet("Dashboard")

			// Headers with theme styling
			headers := []string{"Month", "Revenue", "Costs", "Profit", "Margin"}
			for i, h := range headers {
				s.Cell(1, i+1).SetString(h)
				s.Cell(1, i+1).SetFont(theme.Typography.HeadingFont.WithSize(11))
				s.Cell(1, i+1).SetBackgroundColor(theme.Palette.Primary)
			}

			// Data rows
			months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
			for i, m := range months {
				row := i + 2
				s.SetValue(row, 1, m)
				s.SetValue(row, 2, float64(50000+i*5000))
				s.SetValue(row, 3, float64(30000+i*2000))
				s.Cell(row, 4).SetFormula("B" + string(rune('0'+row)) + "-C" + string(rune('0'+row)))
				s.Cell(row, 5).SetFormula("D" + string(rune('0'+row)) + "/B" + string(rune('0'+row)))
				s.Cell(row, 5).SetNumberFormat("0.0%")
			}

			// Summary sheet
			sum := wb.AddSheet("Summary")
			sum.SetValue(1, 1, "Total Revenue")
			sum.Cell(1, 2).SetFormula("SUM(Dashboard!B2:B7)")

			path := filepath.Join(t.TempDir(), "design_level.xlsx")
			if err := wb.Save(path); err != nil {
				t.Fatalf("error: %v", err)
			}
			assertFileExists(t, path)
		})
	}
}

// TestDesignLevelPptx tests all themes with PPTX
func TestDesignLevelPptx(t *testing.T) {
	for _, theme := range style.AllThemes() {
		t.Run(theme.Name, func(t *testing.T) {
			pres := presentation.NewWithTheme(theme)

			// Title slide
			s1 := pres.AddSlide()
			s1.SetBackground(theme.Palette.Primary)
			title := s1.AddTextBox(common.In(1.5), common.In(2), common.In(10), common.In(2.5))
			tp := title.AddParagraph()
			tp.AddRun(theme.Name+" Presentation",
				theme.Typography.HeadingFont.WithSize(40).WithColor(common.White))
			tp.SetAlignment(common.TextAlignCenter)

			sub := s1.AddTextBox(common.In(3), common.In(5), common.In(7), common.In(1))
			sp := sub.AddParagraph()
			sp.AddRun(theme.Level.String()+" Design Level",
				theme.Typography.BodyFont.WithColor(common.LightGray))
			sp.SetAlignment(common.TextAlignCenter)

			// Content slide
			s2 := pres.AddSlide()
			s2.AddShape(presentation.ShapeRoundedRectangle,
				common.In(0.5), common.In(0.5), common.In(5.5), common.In(6.5))

			contentBox := s2.AddTextBox(common.In(6.5), common.In(1), common.In(6), common.In(5))
			for i := 0; i < 4; i++ {
				bp := contentBox.AddParagraph()
				bp.AddRun("Content point with "+theme.Name+" styling",
					theme.Typography.BodyFont)
				bp.SetSpacing(1.5)
			}

			// Shapes slide
			s3 := pres.AddSlide()
			s3.AddShape(presentation.ShapeCircle,
				common.In(2), common.In(2), common.In(3), common.In(3)).
				SetFill(theme.Palette.Accent)
			s3.AddShape(presentation.ShapeRoundedRectangle,
				common.In(7), common.In(2), common.In(4), common.In(3)).
				SetFill(theme.Palette.Secondary)

			path := filepath.Join(t.TempDir(), "design_level.pptx")
			if err := pres.Save(path); err != nil {
				t.Fatalf("error: %v", err)
			}
			assertFileExists(t, path)
		})
	}
}

// TestDesignLevelPdf tests all themes with PDF
func TestDesignLevelPdf(t *testing.T) {
	for _, theme := range style.AllThemes() {
		t.Run(theme.Name, func(t *testing.T) {
			doc := pdf.NewWithTheme(theme)

			p := doc.AddPage()

			// Header banner
			p.AddRectangle(0, 0, 595, 80, theme.Palette.Primary, nil)
			p.AddText(theme.Name+" Report", 72, 30,
				theme.Typography.HeadingFont.WithSize(28).WithColor(common.White))

			// Body
			p.AddText("Design Level: "+theme.Level.String(), 72, 120,
				theme.Typography.BodyFont)
			p.AddText("Generated with OpenScribe design system", 72, 145,
				theme.Typography.CaptionFont)

			// Separator
			p.AddLine(72, 165, 523, 165, theme.Palette.Accent, 2)

			// Table
			tbl := p.AddTable(72, 200, 4, 3)
			tbl.SetCellSize(150, 25)
			tbl.SetHeaderBackground(theme.Palette.Primary)
			tbl.SetFont(theme.Typography.BodyFont.WithSize(10))
			tbl.SetCell(0, 0, "Feature")
			tbl.SetCell(0, 1, "Status")
			tbl.SetCell(0, 2, "Score")
			tbl.SetCell(1, 0, "Typography")
			tbl.SetCell(1, 1, "Configured")
			tbl.SetCell(1, 2, "100%")
			tbl.SetCell(2, 0, "Colors")
			tbl.SetCell(2, 1, "Applied")
			tbl.SetCell(2, 2, "100%")
			tbl.SetCell(3, 0, "Spacing")
			tbl.SetCell(3, 1, "Tuned")
			tbl.SetCell(3, 2, "100%")

			path := filepath.Join(t.TempDir(), "design_level.pdf")
			if err := doc.Save(path); err != nil {
				t.Fatalf("error: %v", err)
			}
			assertFileExists(t, path)
			assertPDFHeader(t, path)
		})
	}
}

// TestDesignLevelConsistency verifies design level consistency
func TestDesignLevelConsistency(t *testing.T) {
	for _, theme := range style.AllThemes() {
		t.Run(theme.Name+"_consistency", func(t *testing.T) {
			// Typography hierarchy
			if theme.Typography.HeadingFont.Size <= theme.Typography.BodyFont.Size {
				t.Errorf("heading should be larger than body in %s", theme.Name)
			}

			// Spacing progression
			if theme.Spacing.XS.Points() >= theme.Spacing.SM.Points() {
				t.Errorf("XS should be less than SM in %s", theme.Name)
			}
			if theme.Spacing.SM.Points() >= theme.Spacing.MD.Points() {
				t.Errorf("SM should be less than MD in %s", theme.Name)
			}
			if theme.Spacing.MD.Points() >= theme.Spacing.LG.Points() {
				t.Errorf("MD should be less than LG in %s", theme.Name)
			}
			if theme.Spacing.LG.Points() >= theme.Spacing.XL.Points() {
				t.Errorf("LG should be less than XL in %s", theme.Name)
			}

			// Text/Background contrast
			if theme.Palette.Text == theme.Palette.Background {
				t.Errorf("text and background should differ in %s", theme.Name)
			}
		})
	}
}
