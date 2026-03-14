package spreadsheet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JohnPitter/openscribe/common"
)

func TestDonutChart(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Donut")
	chart := s.AddChart(ChartTypeDonut, 1, 1, 8, 12)
	chart.SetTitle("Market Share")
	chart.SetCategories([]string{"Product A", "Product B", "Product C"})
	chart.AddSeries("Share", []float64{45, 35, 20}, common.Blue)

	if chart.Type() != ChartTypeDonut {
		t.Errorf("expected ChartTypeDonut, got %d", chart.Type())
	}

	path := filepath.Join(t.TempDir(), "donut.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestRadarChart(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Radar")
	chart := s.AddChart(ChartTypeRadar, 1, 1, 8, 12)
	chart.SetTitle("Skills Assessment")
	chart.SetCategories([]string{"Go", "Python", "JS", "Rust", "Java"})
	chart.AddSeries("Developer A", []float64{90, 70, 85, 60, 75}, common.Blue)
	chart.AddSeries("Developer B", []float64{70, 90, 80, 80, 65}, common.Red)

	if chart.Type() != ChartTypeRadar {
		t.Errorf("expected ChartTypeRadar, got %d", chart.Type())
	}

	path := filepath.Join(t.TempDir(), "radar.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestBarStackedChart(t *testing.T) {
	wb := New()
	s := wb.AddSheet("Stacked")
	chart := s.AddChart(ChartTypeBarStacked, 1, 1, 10, 15)
	chart.SetTitle("Quarterly Revenue by Product")
	chart.SetCategories([]string{"Q1", "Q2", "Q3", "Q4"})
	chart.AddSeries("Product A", []float64{100, 120, 130, 150}, common.Blue)
	chart.AddSeries("Product B", []float64{80, 90, 110, 130}, common.Green)
	chart.AddSeries("Product C", []float64{60, 70, 80, 90}, common.Red)

	if chart.Type() != ChartTypeBarStacked {
		t.Errorf("expected ChartTypeBarStacked, got %d", chart.Type())
	}

	path := filepath.Join(t.TempDir(), "stacked.xlsx")
	if err := wb.Save(path); err != nil {
		t.Fatalf("save error: %v", err)
	}
	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("file should not be empty")
	}
}

func TestNewChartTypesSaveToBytes(t *testing.T) {
	types := []struct {
		name      string
		chartType ChartType
	}{
		{"donut", ChartTypeDonut},
		{"radar", ChartTypeRadar},
		{"bar_stacked", ChartTypeBarStacked},
	}

	for _, tt := range types {
		t.Run(tt.name, func(t *testing.T) {
			wb := New()
			s := wb.AddSheet("Test")
			chart := s.AddChart(tt.chartType, 1, 1, 8, 12)
			chart.SetTitle("Test")
			chart.AddSeries("Data", []float64{10, 20, 30}, common.Blue)

			data, err := wb.SaveToBytes()
			if err != nil {
				t.Fatalf("save error: %v", err)
			}
			if len(data) == 0 {
				t.Error("bytes should not be empty")
			}
		})
	}
}
