package document

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/style"
)

func TestAddChart(t *testing.T) {
	doc := New()
	doc.AddHeading("Report with Chart", 1)
	doc.AddText("Below is a chart:")

	chart := doc.AddChart(ChartTypeBar, common.In(5), common.In(3))
	chart.SetTitle("Quarterly Revenue")
	chart.SetCategories([]string{"Q1", "Q2", "Q3", "Q4"})
	chart.AddSeries("2025", []float64{150, 200, 180, 250}, common.Blue)
	chart.AddSeries("2026", []float64{180, 250, 220, 300}, common.Red)

	if chart.Title() != "Quarterly Revenue" {
		t.Error("title mismatch")
	}
	if chart.Type() != ChartTypeBar {
		t.Error("type should be bar")
	}
	if len(chart.Series()) != 2 {
		t.Errorf("expected 2 series, got %d", len(chart.Series()))
	}
	if len(chart.Categories()) != 4 {
		t.Errorf("expected 4 categories, got %d", len(chart.Categories()))
	}

	path := filepath.Join(t.TempDir(), "chart.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}

	// Verify chart XML exists in the package
	pkg, err := packaging.OpenPackage(path)
	if err != nil {
		t.Fatalf("open package error: %v", err)
	}
	chartXML, ok := pkg.GetFile("word/charts/chart1.xml")
	if !ok {
		t.Fatal("word/charts/chart1.xml should exist")
	}
	if !strings.Contains(string(chartXML), "barChart") {
		t.Error("chart XML should contain barChart")
	}
	if !strings.Contains(string(chartXML), "Quarterly Revenue") {
		t.Error("chart XML should contain title")
	}
}

func TestLineChart(t *testing.T) {
	doc := New()
	chart := doc.AddChart(ChartTypeLine, common.In(6), common.In(4))
	chart.SetTitle("Monthly Trend")
	chart.SetCategories([]string{"Jan", "Feb", "Mar", "Apr", "May"})
	chart.AddSeries("Revenue", []float64{100, 120, 115, 140, 160}, common.Blue)
	chart.AddSeries("Costs", []float64{80, 85, 90, 95, 100}, common.Red)

	path := filepath.Join(t.TempDir(), "line_chart.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackage(path)
	chartXML, ok := pkg.GetFile("word/charts/chart1.xml")
	if !ok {
		t.Fatal("chart should exist")
	}
	if !strings.Contains(string(chartXML), "lineChart") {
		t.Error("should be lineChart")
	}
}

func TestPieChart(t *testing.T) {
	doc := New()
	chart := doc.AddChart(ChartTypePie, common.In(4), common.In(4))
	chart.SetTitle("Market Share")
	chart.SetCategories([]string{"Product A", "Product B", "Product C"})
	chart.AddSeries("Share", []float64{45, 35, 20}, common.Blue)

	path := filepath.Join(t.TempDir(), "pie_chart.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackage(path)
	chartXML, _ := pkg.GetFile("word/charts/chart1.xml")
	if !strings.Contains(string(chartXML), "pieChart") {
		t.Error("should be pieChart")
	}
}

func TestAreaChart(t *testing.T) {
	doc := New()
	chart := doc.AddChart(ChartTypeArea, common.In(5), common.In(3))
	chart.SetTitle("Growth")
	chart.AddSeries("Users", []float64{1000, 2500, 5000, 8000}, common.Green)

	path := filepath.Join(t.TempDir(), "area_chart.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackage(path)
	chartXML, _ := pkg.GetFile("word/charts/chart1.xml")
	if !strings.Contains(string(chartXML), "areaChart") {
		t.Error("should be areaChart")
	}
}

func TestDonutChart(t *testing.T) {
	doc := New()
	chart := doc.AddChart(ChartTypeDonut, common.In(4), common.In(4))
	chart.SetTitle("Budget Allocation")
	chart.SetCategories([]string{"Engineering", "Marketing", "Sales", "Operations"})
	chart.AddSeries("Budget", []float64{40, 25, 20, 15}, common.Blue)

	path := filepath.Join(t.TempDir(), "donut_chart.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackage(path)
	chartXML, _ := pkg.GetFile("word/charts/chart1.xml")
	if !strings.Contains(string(chartXML), "doughnutChart") {
		t.Error("should be doughnutChart")
	}
}

func TestMultipleCharts(t *testing.T) {
	doc := New()
	doc.AddHeading("Dashboard Report", 1)

	c1 := doc.AddChart(ChartTypeBar, common.In(5), common.In(3))
	c1.SetTitle("Revenue")
	c1.AddSeries("Data", []float64{100, 200, 300}, common.Blue)

	doc.AddText("Some text between charts.")

	c2 := doc.AddChart(ChartTypePie, common.In(4), common.In(4))
	c2.SetTitle("Distribution")
	c2.AddSeries("Data", []float64{60, 30, 10}, common.Red)

	c3 := doc.AddChart(ChartTypeLine, common.In(6), common.In(3))
	c3.SetTitle("Trend")
	c3.AddSeries("Data", []float64{10, 20, 30, 40}, common.Green)

	path := filepath.Join(t.TempDir(), "multi_chart.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackage(path)
	for i := 1; i <= 3; i++ {
		chartPath := fmt.Sprintf("word/charts/chart%d.xml", i)
		if !pkg.HasFile(chartPath) {
			t.Errorf("missing %s", chartPath)
		}
	}
}

func TestChartNoLegend(t *testing.T) {
	doc := New()
	chart := doc.AddChart(ChartTypeBar, common.In(5), common.In(3))
	chart.SetShowLegend(false)
	chart.SetShowTitle(false)
	chart.AddSeries("Data", []float64{1, 2, 3}, common.Blue)

	path := filepath.Join(t.TempDir(), "no_legend.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackage(path)
	chartXML, _ := pkg.GetFile("word/charts/chart1.xml")
	if strings.Contains(string(chartXML), "legendPos") {
		t.Error("should not have legend")
	}
}

func TestChartWithTheme(t *testing.T) {
	themes := style.AllThemes()
	for _, theme := range themes[:3] { // Test first 3 themes
		t.Run(theme.Name, func(t *testing.T) {
			doc := NewWithTheme(theme)
			doc.AddHeading("Themed Report", 1)
			chart := doc.AddChart(ChartTypeBar, common.In(5), common.In(3))
			chart.SetTitle("Data with " + theme.Name)
			chart.AddSeries("Series", []float64{10, 20, 30}, theme.Palette.Primary)

			path := filepath.Join(t.TempDir(), "themed.docx")
			if err := doc.Save(path); err != nil {
				t.Fatalf("save error: %v", err)
			}
		})
	}
}

func TestChartDrawingInDocumentXML(t *testing.T) {
	doc := New()
	doc.AddText("Before chart")
	chart := doc.AddChart(ChartTypeBar, common.In(5), common.In(3))
	chart.SetTitle("Inline Chart")
	chart.AddSeries("S1", []float64{1, 2, 3}, common.Blue)
	doc.AddText("After chart")

	data, err := doc.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackageFromBytes(data)
	docXML, _ := pkg.GetFile("word/document.xml")
	xmlStr := string(docXML)

	if !strings.Contains(xmlStr, "w:drawing") {
		t.Error("document.xml should contain w:drawing for chart")
	}
	if !strings.Contains(xmlStr, "drawingml/2006/chart") {
		t.Error("document.xml should reference chart namespace")
	}
}

func TestChartEscapeXML(t *testing.T) {
	doc := New()
	chart := doc.AddChart(ChartTypeBar, common.In(5), common.In(3))
	chart.SetTitle("Revenue & Costs <2026>")
	chart.SetCategories([]string{"Q1 & Q2", "Q3 <special>"})
	chart.AddSeries("Data", []float64{100, 200}, common.Blue)

	path := filepath.Join(t.TempDir(), "escape.docx")
	if err := doc.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackage(path)
	chartXML, _ := pkg.GetFile("word/charts/chart1.xml")
	xmlStr := string(chartXML)
	if strings.Contains(xmlStr, "Revenue & Costs") {
		t.Error("& should be escaped")
	}
	if !strings.Contains(xmlStr, "Revenue &amp; Costs") {
		t.Error("should contain escaped &amp;")
	}
}
