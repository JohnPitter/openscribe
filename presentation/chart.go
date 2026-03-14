package presentation

import (
	"github.com/JohnPitter/openscribe/common"
)

// ChartType for presentation charts
type ChartType int

const (
	PresentationChartBar ChartType = iota
	PresentationChartLine
	PresentationChartPie
	PresentationChartArea
)

// SlideChart represents a chart on a slide
type SlideChart struct {
	chartType  ChartType
	x, y       common.Measurement
	width      common.Measurement
	height     common.Measurement
	title      string
	series     []ChartSeries
	categories []string
}

// ChartSeries for presentation charts
type ChartSeries struct {
	Name   string
	Values []float64
	Color  common.Color
}

func (c *SlideChart) elementType() string { return "chart" }

// AddChart adds a chart to the slide
func (s *Slide) AddChart(chartType ChartType, x, y, width, height common.Measurement) *SlideChart {
	c := &SlideChart{
		chartType: chartType,
		x:         x, y: y, width: width, height: height,
	}
	s.elements = append(s.elements, c)
	return c
}

// SetTitle sets the chart title
func (c *SlideChart) SetTitle(title string) { c.title = title }

// AddSeries adds a data series
func (c *SlideChart) AddSeries(name string, values []float64, color common.Color) {
	c.series = append(c.series, ChartSeries{Name: name, Values: values, Color: color})
}

// SetCategories sets the category labels
func (c *SlideChart) SetCategories(cats []string) { c.categories = cats }

// Title returns the title
func (c *SlideChart) Title() string { return c.title }

// Type returns the chart type
func (c *SlideChart) Type() ChartType { return c.chartType }

// Series returns all series
func (c *SlideChart) Series() []ChartSeries { return c.series }

// Categories returns the category labels
func (c *SlideChart) Categories() []string { return c.categories }
