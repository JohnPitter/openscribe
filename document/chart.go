package document

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/JohnPitter/openscribe/common"
)

// ChartType represents the type of chart
type ChartType int

const (
	ChartTypeBar ChartType = iota
	ChartTypeLine
	ChartTypePie
	ChartTypeArea
	ChartTypeColumn
	ChartTypeDonut
)

// ChartSeries represents a data series in a chart
type ChartSeries struct {
	Name   string
	Values []float64
	Color  common.Color
}

// Chart represents a chart embedded in a document
type Chart struct {
	chartType  ChartType
	title      string
	series     []ChartSeries
	categories []string
	width      common.Measurement
	height     common.Measurement
	showLegend bool
	showTitle  bool
	relID      string
	index      int
}

// AddChart adds an inline chart to the document
func (d *Document) AddChart(chartType ChartType, width, height common.Measurement) *Chart {
	d.chartCount++
	c := &Chart{
		chartType:  chartType,
		width:      width,
		height:     height,
		showLegend: true,
		showTitle:  true,
		index:      d.chartCount,
	}
	d.charts = append(d.charts, c)
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

// SetCategories sets the category labels
func (c *Chart) SetCategories(cats []string) { c.categories = cats }

// SetShowLegend toggles legend
func (c *Chart) SetShowLegend(show bool) { c.showLegend = show }

// SetShowTitle toggles title visibility
func (c *Chart) SetShowTitle(show bool) { c.showTitle = show }

// Series returns all series
func (c *Chart) Series() []ChartSeries { return c.series }

// Type returns the chart type
func (c *Chart) Type() ChartType { return c.chartType }

// Categories returns the categories
func (c *Chart) Categories() []string { return c.categories }

// buildChartXML generates the DrawingML chart XML for word/charts/chartN.xml
func (c *Chart) buildChartXML() []byte {
	chartTypeTag := "c:barChart"
	grouping := "clustered"
	switch c.chartType {
	case ChartTypeLine:
		chartTypeTag = "c:lineChart"
		grouping = "standard"
	case ChartTypePie:
		chartTypeTag = "c:pieChart"
		grouping = ""
	case ChartTypeArea:
		chartTypeTag = "c:areaChart"
		grouping = "standard"
	case ChartTypeColumn:
		chartTypeTag = "c:barChart"
		grouping = "clustered"
	case ChartTypeDonut:
		chartTypeTag = "c:doughnutChart"
		grouping = ""
	}

	var buf bytes.Buffer
	buf.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	buf.WriteString(`<c:chartSpace xmlns:c="http://schemas.openxmlformats.org/drawingml/2006/chart" xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)
	buf.WriteString(`<c:chart>`)

	// Title
	if c.showTitle && c.title != "" {
		fmt.Fprintf(&buf, `<c:title><c:tx><c:rich><a:bodyPr/><a:lstStyle/><a:p><a:r><a:t>%s</a:t></a:r></a:p></c:rich></c:tx><c:overlay val="0"/></c:title>`, escapeXML(c.title))
	}

	buf.WriteString(`<c:plotArea><c:layout/>`)

	// Chart type element
	fmt.Fprintf(&buf, `<%s>`, chartTypeTag)
	if grouping != "" {
		fmt.Fprintf(&buf, `<c:grouping val="%s"/>`, grouping)
	}

	// Series
	for i, s := range c.series {
		fmt.Fprintf(&buf, `<c:ser><c:idx val="%d"/><c:order val="%d"/>`, i, i)
		if s.Name != "" {
			fmt.Fprintf(&buf, `<c:tx><c:strRef><c:f>"%s"</c:f></c:strRef></c:tx>`, escapeXML(s.Name))
		}
		// Color
		fmt.Fprintf(&buf, `<c:spPr><a:solidFill><a:srgbClr val="%02X%02X%02X"/></a:solidFill></c:spPr>`, s.Color.R, s.Color.G, s.Color.B)

		// Category data
		if len(c.categories) > 0 {
			buf.WriteString(`<c:cat><c:strLit>`)
			fmt.Fprintf(&buf, `<c:ptCount val="%d"/>`, len(c.categories))
			for j, cat := range c.categories {
				fmt.Fprintf(&buf, `<c:pt idx="%d"><c:v>%s</c:v></c:pt>`, j, escapeXML(cat))
			}
			buf.WriteString(`</c:strLit></c:cat>`)
		}

		// Values
		if len(s.Values) > 0 {
			buf.WriteString(`<c:val><c:numLit>`)
			fmt.Fprintf(&buf, `<c:ptCount val="%d"/>`, len(s.Values))
			for j, v := range s.Values {
				fmt.Fprintf(&buf, `<c:pt idx="%d"><c:v>%s</c:v></c:pt>`, j, strconv.FormatFloat(v, 'f', -1, 64))
			}
			buf.WriteString(`</c:numLit></c:val>`)
		}

		buf.WriteString(`</c:ser>`)
	}

	fmt.Fprintf(&buf, `</%s>`, chartTypeTag)
	buf.WriteString(`</c:plotArea>`)

	// Legend
	if c.showLegend {
		buf.WriteString(`<c:legend><c:legendPos val="b"/></c:legend>`)
	}

	buf.WriteString(`</c:chart></c:chartSpace>`)
	return buf.Bytes()
}

func escapeXML(s string) string {
	var buf bytes.Buffer
	for _, c := range s {
		switch c {
		case '&':
			buf.WriteString("&amp;")
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '"':
			buf.WriteString("&quot;")
		default:
			buf.WriteRune(c)
		}
	}
	return buf.String()
}
