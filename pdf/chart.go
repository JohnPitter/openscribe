package pdf

import "github.com/JohnPitter/openscribe/common"

// ChartType represents the type of chart
type ChartType int

const (
	ChartTypeBar ChartType = iota
	ChartTypeLine
	ChartTypePie
	ChartTypeArea
	ChartTypeHorizontalBar
)

// ChartSeries represents a data series in a chart
type ChartSeries struct {
	Name   string
	Values []float64
	Color  common.Color
}

// ChartElement represents a chart on a PDF page
type ChartElement struct {
	chartType  ChartType
	x, y       float64
	width      float64
	height     float64
	title      string
	titleFont  common.Font
	categories []string
	series     []ChartSeries
	showLegend bool
	showValues bool
	bgColor    *common.Color
	gridColor  common.Color
	axisColor  common.Color
}

func (c *ChartElement) pdfElement() {}

// AddChart adds a chart to the page
func (p *Page) AddChart(chartType ChartType, x, y, width, height float64) *ChartElement {
	c := &ChartElement{
		chartType:  chartType,
		x:          x,
		y:          y,
		width:      width,
		height:     height,
		titleFont:  common.NewFont("Helvetica", 12).Bold(),
		showLegend: true,
		gridColor:  common.LightGray,
		axisColor:  common.Black,
	}
	p.elements = append(p.elements, c)
	return c
}

// SetTitle sets the chart title
func (c *ChartElement) SetTitle(title string) { c.title = title }

// SetTitleFont sets the title font
func (c *ChartElement) SetTitleFont(f common.Font) { c.titleFont = f }

// SetCategories sets the X-axis categories (labels)
func (c *ChartElement) SetCategories(cats []string) { c.categories = cats }

// AddSeries adds a data series
func (c *ChartElement) AddSeries(name string, values []float64, color common.Color) {
	c.series = append(c.series, ChartSeries{Name: name, Values: values, Color: color})
}

// SetShowLegend toggles legend visibility
func (c *ChartElement) SetShowLegend(show bool) { c.showLegend = show }

// SetShowValues toggles value labels
func (c *ChartElement) SetShowValues(show bool) { c.showValues = show }

// SetBackground sets the chart background
func (c *ChartElement) SetBackground(color common.Color) { c.bgColor = &color }

// SetGridColor sets the grid line color
func (c *ChartElement) SetGridColor(color common.Color) { c.gridColor = color }

// Series returns all data series
func (c *ChartElement) Series() []ChartSeries { return c.series }

// Title returns the chart title
func (c *ChartElement) Title() string { return c.title }

// Type returns the chart type
func (c *ChartElement) Type() ChartType { return c.chartType }
