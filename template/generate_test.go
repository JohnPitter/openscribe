package template

import (
	"path/filepath"
	"testing"
)

func TestGenerateDocx(t *testing.T) {
	docxTemplates := ByFormat(FormatDOCX)
	for _, tmpl := range docxTemplates {
		t.Run(tmpl.Name, func(t *testing.T) {
			doc, err := tmpl.GenerateDocx()
			if err != nil {
				t.Fatalf("generate error: %v", err)
			}
			if len(doc.Paragraphs()) == 0 {
				t.Error("generated doc should have content")
			}
			path := filepath.Join(t.TempDir(), "generated.docx")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
		})
	}
}

func TestGenerateXlsx(t *testing.T) {
	xlsxTemplates := ByFormat(FormatXLSX)
	for _, tmpl := range xlsxTemplates {
		t.Run(tmpl.Name, func(t *testing.T) {
			wb, err := tmpl.GenerateXlsx()
			if err != nil {
				t.Fatalf("generate error: %v", err)
			}
			if wb.SheetCount() == 0 {
				t.Error("generated workbook should have sheets")
			}
			path := filepath.Join(t.TempDir(), "generated.xlsx")
			if err := wb.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
		})
	}
}

func TestGeneratePptx(t *testing.T) {
	pptxTemplates := ByFormat(FormatPPTX)
	for _, tmpl := range pptxTemplates {
		t.Run(tmpl.Name, func(t *testing.T) {
			pres, err := tmpl.GeneratePptx()
			if err != nil {
				t.Fatalf("generate error: %v", err)
			}
			if pres.SlideCount() == 0 {
				t.Error("generated presentation should have slides")
			}
			path := filepath.Join(t.TempDir(), "generated.pptx")
			if err := pres.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
		})
	}
}

func TestGeneratePdf(t *testing.T) {
	pdfTemplates := ByFormat(FormatPDF)
	for _, tmpl := range pdfTemplates {
		t.Run(tmpl.Name, func(t *testing.T) {
			doc, err := tmpl.GeneratePdf()
			if err != nil {
				t.Fatalf("generate error: %v", err)
			}
			if doc.PageCount() == 0 {
				t.Error("generated PDF should have pages")
			}
			path := filepath.Join(t.TempDir(), "generated.pdf")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
		})
	}
}

func TestGenerateWrongFormat(t *testing.T) {
	tmpl := Find("Basic Report") // DOCX template
	if tmpl == nil {
		t.Fatal("should find template")
	}

	_, err := tmpl.GenerateXlsx()
	if err == nil {
		t.Error("should error when generating XLSX from DOCX template")
	}

	_, err = tmpl.GeneratePptx()
	if err == nil {
		t.Error("should error when generating PPTX from DOCX template")
	}

	_, err = tmpl.GeneratePdf()
	if err == nil {
		t.Error("should error when generating PDF from DOCX template")
	}
}

func TestGenerateAllTemplates(t *testing.T) {
	for _, tmpl := range All() {
		t.Run(tmpl.Name, func(t *testing.T) {
			switch tmpl.Format {
			case FormatDOCX:
				doc, err := tmpl.GenerateDocx()
				if err != nil {
					t.Fatalf("error: %v", err)
				}
				path := filepath.Join(t.TempDir(), "gen.docx")
				if err := doc.Save(path); err != nil {
					t.Fatalf("save error: %v", err)
				}
			case FormatXLSX:
				wb, err := tmpl.GenerateXlsx()
				if err != nil {
					t.Fatalf("error: %v", err)
				}
				path := filepath.Join(t.TempDir(), "gen.xlsx")
				if err := wb.Save(path); err != nil {
					t.Fatalf("save error: %v", err)
				}
			case FormatPPTX:
				pres, err := tmpl.GeneratePptx()
				if err != nil {
					t.Fatalf("error: %v", err)
				}
				path := filepath.Join(t.TempDir(), "gen.pptx")
				if err := pres.Save(path); err != nil {
					t.Fatalf("save error: %v", err)
				}
			case FormatPDF:
				doc, err := tmpl.GeneratePdf()
				if err != nil {
					t.Fatalf("error: %v", err)
				}
				path := filepath.Join(t.TempDir(), "gen.pdf")
				if err := doc.Save(path); err != nil {
					t.Fatalf("save error: %v", err)
				}
			}
		})
	}
}
