package spreadsheet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestAddChart(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")
	s.SetValue(1, 1, "Q1")
	s.SetValue(1, 2, "Q2")
	s.SetValue(2, 1, 100.0)
	s.SetValue(2, 2, 200.0)

	chart := s.AddChart(ChartTypeBar, 4, 1, 8, 12)
	chart.SetTitle("Sales")
	chart.SetCategories([]string{"Q1", "Q2"})
	chart.AddSeries("Revenue", []float64{100, 200}, common.Blue)

	if chart.Title() != "Sales" {
		t.Error("title mismatch")
	}
	if chart.Type() != ChartTypeBar {
		t.Error("type should be bar")
	}
	if len(chart.Series()) != 1 {
		t.Errorf("expected 1 series, got %d", len(chart.Series()))
	}

	path := filepath.Join(t.TempDir(), "chart.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestLineChart(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Trends")
	chart := s.AddChart(ChartTypeLine, 1, 5, 10, 15)
	chart.SetTitle("Monthly Trend")
	chart.SetCategories([]string{"Jan", "Feb", "Mar", "Apr"})
	chart.AddSeries("Users", []float64{1000, 1500, 2000, 2800}, common.Blue)
	chart.AddSeries("Sessions", []float64{5000, 6000, 8000, 12000}, common.Green)

	path := filepath.Join(t.TempDir(), "line.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestPieChart(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Market")
	chart := s.AddChart(ChartTypePie, 1, 1, 8, 12)
	chart.SetTitle("Market Share")
	chart.SetCategories([]string{"Product A", "Product B", "Product C"})
	chart.AddSeries("Share", []float64{45, 35, 20}, common.Blue)

	path := filepath.Join(t.TempDir(), "pie.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestMultipleCharts(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Dashboard")

	c1 := s.AddChart(ChartTypeBar, 1, 1, 8, 10)
	c1.SetTitle("Revenue")
	c1.AddSeries("2025", []float64{100, 200, 300}, common.Blue)

	c2 := s.AddChart(ChartTypeLine, 1, 9, 8, 10)
	c2.SetTitle("Users")
	c2.AddSeries("Active", []float64{50, 80, 120}, common.Green)

	path := filepath.Join(t.TempDir(), "multi_chart.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestChartFromRange(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Data")
	chart := s.AddChart(ChartTypeColumn, 5, 1, 8, 12)
	chart.SetTitle("From Range")
	chart.AddSeriesFromRange("Revenue", "Data!B1:B4", common.Red)
	chart.SetCategoryRange("Data!A1:A4")
	chart.SetShowLegend(false)
	chart.SetShowTitle(true)

	if chart.showLegend {
		t.Error("legend should be off")
	}

	path := filepath.Join(t.TempDir(), "range_chart.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}

func TestColumnManagement(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Test")

	s.SetColumnWidth(1, 20)
	if s.ColumnWidth(1) != 20 {
		t.Errorf("expected 20, got %f", s.ColumnWidth(1))
	}

	// Default width
	if s.ColumnWidth(99) != 8.43 {
		t.Errorf("default should be 8.43, got %f", s.ColumnWidth(99))
	}

	// Range
	s.SetColumnWidthRange(3, 6, 15)
	if s.ColumnWidth(4) != 15 {
		t.Error("range width should be 15")
	}

	// Hidden
	s.SetColumnHidden(2, true)
	if !s.columns[2].hidden {
		t.Error("column 2 should be hidden")
	}

	// Best fit
	s.SetColumnBestFit(3, true)
	if !s.columns[3].bestFit {
		t.Error("column 3 should be best fit")
	}

	path := filepath.Join(t.TempDir(), "columns.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
}
