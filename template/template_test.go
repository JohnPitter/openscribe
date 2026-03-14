package template

import (
	"testing"

	"github.com/JohnPitter/openscribe/style"
)

func TestAllTemplates(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Fatal("should have registered templates")
	}
	// 8 basic + 8 professional + 8 premium + 8 luxury = 32
	if len(all) != 32 {
		t.Errorf("expected 32 templates, got %d", len(all))
	}
}

func TestCount(t *testing.T) {
	if Count() != 32 {
		t.Errorf("expected 32, got %d", Count())
	}
}

func TestByLevel(t *testing.T) {
	basic := ByLevel(style.DesignLevelBasic)
	if len(basic) != 8 {
		t.Errorf("expected 8 basic templates, got %d", len(basic))
	}

	pro := ByLevel(style.DesignLevelProfessional)
	if len(pro) != 8 {
		t.Errorf("expected 8 professional templates, got %d", len(pro))
	}

	premium := ByLevel(style.DesignLevelPremium)
	if len(premium) != 8 {
		t.Errorf("expected 8 premium templates, got %d", len(premium))
	}

	luxury := ByLevel(style.DesignLevelLuxury)
	if len(luxury) != 8 {
		t.Errorf("expected 8 luxury templates, got %d", len(luxury))
	}
}

func TestByFormat(t *testing.T) {
	docx := ByFormat(FormatDOCX)
	if len(docx) == 0 {
		t.Error("should have DOCX templates")
	}

	xlsx := ByFormat(FormatXLSX)
	if len(xlsx) == 0 {
		t.Error("should have XLSX templates")
	}

	pptx := ByFormat(FormatPPTX)
	if len(pptx) == 0 {
		t.Error("should have PPTX templates")
	}

	pdf := ByFormat(FormatPDF)
	if len(pdf) == 0 {
		t.Error("should have PDF templates")
	}
}

func TestByCategory(t *testing.T) {
	reports := ByCategory(CategoryReport)
	if len(reports) == 0 {
		t.Error("should have report templates")
	}

	invoices := ByCategory(CategoryInvoice)
	if len(invoices) == 0 {
		t.Error("should have invoice templates")
	}
}

func TestFind(t *testing.T) {
	tmpl := Find("Agency Pitch Deck")
	if tmpl == nil {
		t.Fatal("should find Agency Pitch Deck")
	}
	if tmpl.Level != style.DesignLevelLuxury {
		t.Error("Agency Pitch Deck should be luxury level")
	}
	if tmpl.Format != FormatPPTX {
		t.Error("Agency Pitch Deck should be PPTX")
	}

	if Find("NonExistent") != nil {
		t.Error("should return nil for non-existent template")
	}
}

func TestSearch(t *testing.T) {
	// Search by level only
	level := style.DesignLevelPremium
	results := Search(&level, nil, nil)
	if len(results) != 8 {
		t.Errorf("expected 8 premium results, got %d", len(results))
	}

	// Search by format only
	format := FormatPPTX
	results = Search(nil, &format, nil)
	if len(results) == 0 {
		t.Error("should find PPTX templates")
	}

	// Search by category only
	cat := CategoryReport
	results = Search(nil, nil, &cat)
	if len(results) == 0 {
		t.Error("should find report templates")
	}

	// Combined search
	luxLevel := style.DesignLevelLuxury
	pptxFormat := FormatPPTX
	pitchCat := CategoryPitchDeck
	results = Search(&luxLevel, &pptxFormat, &pitchCat)
	if len(results) != 2 {
		t.Errorf("expected 2 luxury PPTX pitch decks, got %d", len(results))
	}
}

func TestFormats(t *testing.T) {
	f := Formats()
	if len(f) != 4 {
		t.Errorf("expected 4 formats, got %d", len(f))
	}
}

func TestCategories(t *testing.T) {
	c := Categories()
	if len(c) != 8 {
		t.Errorf("expected 8 categories, got %d", len(c))
	}
}

func TestLevels(t *testing.T) {
	l := Levels()
	if len(l) != 4 {
		t.Errorf("expected 4 levels, got %d", len(l))
	}
}

func TestFormatString(t *testing.T) {
	if FormatDOCX.String() != "DOCX" {
		t.Error("expected DOCX")
	}
	if FormatXLSX.String() != "XLSX" {
		t.Error("expected XLSX")
	}
	if FormatPPTX.String() != "PPTX" {
		t.Error("expected PPTX")
	}
	if FormatPDF.String() != "PDF" {
		t.Error("expected PDF")
	}
}

func TestCategoryString(t *testing.T) {
	if CategoryReport.String() != "Report" {
		t.Error("expected Report")
	}
	if CategoryPitchDeck.String() != "Pitch Deck" {
		t.Error("expected Pitch Deck")
	}
}

func TestTemplateFields(t *testing.T) {
	for _, tmpl := range All() {
		if tmpl.Name == "" {
			t.Error("template name should not be empty")
		}
		if tmpl.Description == "" {
			t.Errorf("template %s description should not be empty", tmpl.Name)
		}
		if tmpl.Theme.Name == "" {
			t.Errorf("template %s should have a theme", tmpl.Name)
		}
	}
}
