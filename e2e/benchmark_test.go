package e2e

import (
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/document"
	"github.com/JohnPitter/openscribe/pdf"
	"github.com/JohnPitter/openscribe/presentation"
	"github.com/JohnPitter/openscribe/spreadsheet"
)

func BenchmarkDocxCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doc := document.New()
		doc.AddHeading("Benchmark Document", 1)
		for j := 0; j < 50; j++ {
			doc.AddText("This is a paragraph of text for benchmarking purposes.")
		}
		tbl := doc.AddTable(10, 5)
		for r := 0; r < 10; r++ {
			for c := 0; c < 5; c++ {
				tbl.Cell(r, c).SetText("Data")
			}
		}
		doc.SaveToBytes()
	}
}

func BenchmarkDocxSaveToBytes(b *testing.B) {
	doc := document.New()
	doc.AddHeading("Benchmark", 1)
	for j := 0; j < 100; j++ {
		doc.AddText("Paragraph content for benchmark testing.")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		doc.SaveToBytes()
	}
}

func BenchmarkXlsxCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wb := spreadsheet.New()
		s := wb.AddSheet("Data")
		for r := 1; r <= 100; r++ {
			for c := 1; c <= 10; c++ {
				s.Cell(r, c).SetNumber(float64(r * c))
			}
		}
		wb.SaveToBytes()
	}
}

func BenchmarkXlsxLarge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wb := spreadsheet.New()
		s := wb.AddSheet("Large")
		for r := 1; r <= 1000; r++ {
			for c := 1; c <= 20; c++ {
				s.Cell(r, c).SetNumber(float64(r * c))
			}
		}
		wb.SaveToBytes()
	}
}

func BenchmarkPptxCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pres := presentation.New()
		for s := 0; s < 10; s++ {
			slide := pres.AddSlide()
			slide.AddTextBox(common.In(1), common.In(1), common.In(8), common.In(2)).
				SetText("Slide Title", common.NewFont("Arial", 28))
			slide.AddShape(presentation.ShapeRectangle,
				common.In(1), common.In(4), common.In(4), common.In(2))
		}
		pres.SaveToBytes()
	}
}

func BenchmarkPdfCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doc := pdf.New()
		for p := 0; p < 10; p++ {
			page := doc.AddPage()
			page.AddText("Page Title", 72, 72, common.NewFont("Helvetica", 24))
			for l := 0; l < 20; l++ {
				page.AddText("Line of text content", 72, float64(100+l*15), common.NewFont("Helvetica", 11))
			}
		}
		doc.SaveToBytes()
	}
}

func BenchmarkPdfWithChart(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doc := pdf.New()
		page := doc.AddPage()
		chart := page.AddChart(pdf.ChartTypeBar, 72, 72, 400, 300)
		chart.SetCategories([]string{"Q1", "Q2", "Q3", "Q4"})
		chart.AddSeries("Revenue", []float64{100, 200, 150, 300}, common.Blue)
		chart.AddSeries("Costs", []float64{80, 150, 120, 200}, common.Red)
		doc.SaveToBytes()
	}
}

func BenchmarkPdfMerge(b *testing.B) {
	docs := make([]*pdf.Document, 5)
	for i := range docs {
		docs[i] = pdf.New()
		docs[i].AddPage().AddText("Content", 72, 72, common.NewFont("Helvetica", 12))
		docs[i].AddPage().AddText("Content", 72, 72, common.NewFont("Helvetica", 12))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		merged := pdf.Merge(docs...)
		merged.SaveToBytes()
	}
}
