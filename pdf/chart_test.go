package pdf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestBarChart(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	chart := p.AddChart(ChartTypeBar, 72, 72, 400, 250)
	chart.SetTitle("Sales by Quarter")
	chart.SetCategories([]string{"Q1", "Q2", "Q3", "Q4"})
	chart.AddSeries("2025", []float64{150, 200, 180, 250}, common.Blue)
	chart.AddSeries("2026", []float64{180, 220, 210, 280}, common.Red)
	chart.SetShowValues(true)

	if chart.Title() != "Sales by Quarter" {
		t.Error("title mismatch")
	}
	if chart.Type() != ChartTypeBar {
		t.Error("type mismatch")
	}
	if len(chart.Series()) != 2 {
		t.Errorf("expected 2 series, got %d", len(chart.Series()))
	}

	path := filepath.Join(t.TempDir(), "bar_chart.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestLineChart(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	chart := p.AddChart(ChartTypeLine, 72, 72, 400, 250)
	chart.SetTitle("Revenue Trend")
	chart.SetCategories([]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"})
	chart.AddSeries("Revenue", []float64{100, 120, 115, 140, 160, 180}, common.Blue)
	chart.AddSeries("Costs", []float64{80, 85, 90, 95, 100, 105}, common.Red)
	chart.SetBackground(common.White)

	path := filepath.Join(t.TempDir(), "line_chart.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestPieChart(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	chart := p.AddChart(ChartTypePie, 72, 72, 300, 300)
	chart.SetTitle("Market Share")
	chart.AddSeries("Products", []float64{35, 25, 20, 15, 5},
		common.Blue)
	chart.SetCategories([]string{"Product A", "Product B", "Product C", "Product D", "Other"})

	path := filepath.Join(t.TempDir(), "pie_chart.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestAreaChart(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	chart := p.AddChart(ChartTypeArea, 72, 72, 400, 250)
	chart.SetTitle("Growth Over Time")
	chart.SetCategories([]string{"2020", "2021", "2022", "2023", "2024"})
	chart.AddSeries("Users", []float64{1000, 2500, 5000, 8000, 15000}, common.Blue)
	chart.SetGridColor(common.LightGray)

	path := filepath.Join(t.TempDir(), "area_chart.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestHorizontalBarChart(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	chart := p.AddChart(ChartTypeHorizontalBar, 72, 72, 400, 250)
	chart.SetTitle("Performance Scores")
	chart.SetCategories([]string{"Team A", "Team B", "Team C"})
	chart.AddSeries("Score", []float64{85, 92, 78}, common.Green)
	chart.SetShowValues(true)

	path := filepath.Join(t.TempDir(), "hbar_chart.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestChartNoData(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	p.AddChart(ChartTypeBar, 72, 72, 400, 250)

	// Should not panic with empty data
	path := filepath.Join(t.TempDir(), "empty_chart.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestMultipleChartsOnPage(t *testing.T) {
	doc := New()
	p := doc.AddPage()

	c1 := p.AddChart(ChartTypeBar, 72, 72, 220, 200)
	c1.SetTitle("Bar Chart")
	c1.AddSeries("S1", []float64{10, 20, 30}, common.Blue)

	c2 := p.AddChart(ChartTypePie, 320, 72, 220, 200)
	c2.SetTitle("Pie Chart")
	c2.AddSeries("S1", []float64{40, 30, 20, 10}, common.Red)

	c3 := p.AddChart(ChartTypeLine, 72, 350, 220, 200)
	c3.SetTitle("Line Chart")
	c3.AddSeries("S1", []float64{5, 15, 10, 25}, common.Green)

	path := filepath.Join(t.TempDir(), "multi_chart.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)
}

func TestSplit(t *testing.T) {
	doc := New()
	doc.AddPage().AddText("Page 1", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 2", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 3", 72, 72, common.NewFont("Helvetica", 12))

	d1, d2, err := doc.Split(2)
	if err != nil {
		t.Fatalf("split error: %v", err)
	}
	if d1.PageCount() != 2 {
		t.Errorf("d1 should have 2 pages, got %d", d1.PageCount())
	}
	if d2.PageCount() != 1 {
		t.Errorf("d2 should have 1 page, got %d", d2.PageCount())
	}

	// Invalid split
	_, _, err = doc.Split(0)
	if err == nil {
		t.Error("should error on split at 0")
	}
	_, _, err = doc.Split(3)
	if err == nil {
		t.Error("should error on split at end")
	}
}

func TestExtractPages(t *testing.T) {
	doc := New()
	doc.AddPage().AddText("Page 1", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 2", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 3", 72, 72, common.NewFont("Helvetica", 12))
	doc.AddPage().AddText("Page 4", 72, 72, common.NewFont("Helvetica", 12))

	extracted, err := doc.ExtractPages(0, 2)
	if err != nil {
		t.Fatalf("extract error: %v", err)
	}
	if extracted.PageCount() != 2 {
		t.Errorf("expected 2 pages, got %d", extracted.PageCount())
	}

	path := filepath.Join(t.TempDir(), "extracted.pdf")
	if err := extracted.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	assertPDFValid(t, path)

	// Invalid index
	_, err = doc.ExtractPages(10)
	if err == nil {
		t.Error("should error on invalid index")
	}
}

func TestPdfImage(t *testing.T) {
	doc := New()
	p := doc.AddPage()
	imgData := &common.ImageData{
		Data:   []byte{0x89, 0x50, 0x4E, 0x47},
		Format: common.ImageFormatPNG,
	}
	p.AddImage(imgData, 72, 72, 200, 150)

	if p.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", p.ElementCount())
	}
}

func assertPDFValid(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatal("not a valid PDF")
	}
	if len(data) == 0 {
		t.Fatal("empty file")
	}
}
