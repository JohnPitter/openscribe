package presentation

import (
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestAddConnector(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	conn := s.AddConnector(ConnectorStraight,
		common.In(1), common.In(1), common.In(5), common.In(3))
	if conn == nil {
		t.Fatal("connector should not be nil")
	}
	if conn.Type() != ConnectorStraight {
		t.Error("connector type should be straight")
	}
	if s.ElementCount() != 1 {
		t.Errorf("expected 1 element, got %d", s.ElementCount())
	}
}

func TestConnectorSetColorAndWidth(t *testing.T) {
	pres := New()
	s := pres.AddSlide()

	conn := s.AddConnector(ConnectorElbow,
		common.In(1), common.In(1), common.In(5), common.In(3))
	conn.SetColor(common.Red)
	conn.SetWidth(common.Pt(3))

	if conn.Color() != common.Red {
		t.Error("color should be red")
	}
	if conn.Width().Points() != 3 {
		t.Errorf("expected width 3pt, got %f", conn.Width().Points())
	}
}

func TestConnectorTypes(t *testing.T) {
	types := []struct {
		ct   ConnectorType
		name string
	}{
		{ConnectorStraight, "line"},
		{ConnectorElbow, "bentConnector3"},
		{ConnectorCurved, "curvedConnector3"},
	}

	for _, tc := range types {
		t.Run(tc.name, func(t *testing.T) {
			pres := New()
			s := pres.AddSlide()
			s.AddConnector(tc.ct,
				common.In(1), common.In(1), common.In(4), common.In(3))

			data, err := pres.SaveToBytes()
			if err != nil {
				t.Fatalf("save error: %v", err)
			}

			pkg, _ := packaging.OpenPackageFromBytes(data)
			slideXML, _ := pkg.GetFile("ppt/slides/slide1.xml")
			xmlStr := string(slideXML)

			if !strings.Contains(xmlStr, "cxnSp") {
				t.Error("slide XML should contain cxnSp element")
			}
			if !strings.Contains(xmlStr, tc.name) {
				t.Errorf("slide XML should contain preset geometry %s", tc.name)
			}
		})
	}
}

func TestConnectorElementType(t *testing.T) {
	conn := &Connector{}
	if conn.elementType() != "connector" {
		t.Errorf("expected 'connector', got '%s'", conn.elementType())
	}
}

func TestConnectorSerialization(t *testing.T) {
	pres := New()
	s := pres.AddSlide()
	conn := s.AddConnector(ConnectorStraight,
		common.In(1), common.In(2), common.In(5), common.In(4))
	conn.SetColor(common.Blue)
	conn.SetWidth(common.Pt(2))

	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, err := packaging.OpenPackageFromBytes(data)
	if err != nil {
		t.Fatalf("open package error: %v", err)
	}

	slideXML, ok := pkg.GetFile("ppt/slides/slide1.xml")
	if !ok {
		t.Fatal("slide1.xml should exist")
	}

	xmlStr := string(slideXML)
	if !strings.Contains(xmlStr, "p:cxnSp") {
		t.Error("slide XML should contain p:cxnSp element")
	}
	if !strings.Contains(xmlStr, "Connector") {
		t.Error("slide XML should contain Connector name")
	}
}

func TestNewShapeTypes(t *testing.T) {
	shapes := []struct {
		st   ShapeType
		prst string
	}{
		{ShapeCallout, "wedgeRoundRectCallout"},
		{ShapeFlowchartProcess, "flowChartProcess"},
		{ShapeFlowchartDecision, "flowChartDecision"},
		{ShapeFlowchartTerminator, "flowChartTerminator"},
		{ShapeBrace, "leftBrace"},
		{ShapeBracket, "leftBracket"},
	}

	for _, tc := range shapes {
		t.Run(tc.prst, func(t *testing.T) {
			pres := New()
			s := pres.AddSlide()
			s.AddShape(tc.st, common.In(1), common.In(1), common.In(3), common.In(2))

			data, err := pres.SaveToBytes()
			if err != nil {
				t.Fatalf("save error: %v", err)
			}

			pkg, _ := packaging.OpenPackageFromBytes(data)
			slideXML, _ := pkg.GetFile("ppt/slides/slide1.xml")

			if !strings.Contains(string(slideXML), tc.prst) {
				t.Errorf("slide XML should contain preset geometry %s", tc.prst)
			}
		})
	}
}
