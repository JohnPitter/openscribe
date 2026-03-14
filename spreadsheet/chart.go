package spreadsheet

import (
	"github.com/JohnPitter/openscribe/common"
)

// ChartType represents the type of chart
type ChartType int

const (
	ChartTypeBar ChartType = iota
	ChartTypeLine
	ChartTypePie
	ChartTypeArea
	ChartTypeScatter
	ChartTypeColumn
	ChartTypeDonut
	ChartTypeRadar
	ChartTypeBarStacked
)

// ChartSeries represents a data series in a chart
type ChartSeries struct {
	Name      string
	Values    []float64
	Color     common.Color
	DataRange string // e.g., "Sheet1!B1:B10"
}

// Chart represents a chart in a worksheet
type Chart struct {
	sheet      *Sheet
	chartType  ChartType
	title      string
	series     []ChartSeries
	categories []string
	catRange   string // e.g., "Sheet1!A1:A10"
	x, y       int    // cell position (row, col)
	width      int    // in cells
	height     int    // in cells
	showLegend bool
	showTitle  bool
}

// AddChart adds a chart to the sheet at the specified cell position
func (s *Sheet) AddChart(chartType ChartType, row, col, width, height int) *Chart {
	c := &Chart{
		sheet:      s,
		chartType:  chartType,
		x:          col,
		y:          row,
		width:      width,
		height:     height,
		showLegend: true,
		showTitle:  true,
	}
	s.charts = append(s.charts, c)
	return c
}

// SetTitle sets the chart title
func (c *Chart) SetTitle(title string) { c.title = title; c.showTitle = true }

// Title returns the chart title
func (c *Chart) Title() string { return c.title }

// AddSeries adds a data series to the chart
func (c *Chart) AddSeries(name string, values []float64, color common.Color) {
	c.series = append(c.series, ChartSeries{
		Name:   name,
		Values: values,
		Color:  color,
	})
}

// AddSeriesFromRange adds a data series from a cell range
func (c *Chart) AddSeriesFromRange(name, dataRange string, color common.Color) {
	c.series = append(c.series, ChartSeries{
		Name:      name,
		DataRange: dataRange,
		Color:     color,
	})
}

// SetCategories sets the category labels
func (c *Chart) SetCategories(cats []string) { c.categories = cats }

// SetCategoryRange sets categories from a cell range
func (c *Chart) SetCategoryRange(catRange string) { c.catRange = catRange }

// SetShowLegend toggles legend
func (c *Chart) SetShowLegend(show bool) { c.showLegend = show }

// SetShowTitle toggles title visibility
func (c *Chart) SetShowTitle(show bool) { c.showTitle = show }

// Series returns all series
func (c *Chart) Series() []ChartSeries { return c.series }

// Type returns the chart type
func (c *Chart) Type() ChartType { return c.chartType }
