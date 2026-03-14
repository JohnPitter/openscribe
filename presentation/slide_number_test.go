package presentation

import (
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestSlideNumbersEnable(t *testing.T) {
	pres := New()
	pres.SetSlideNumbers(true)

	if !pres.SlideNumbersEnabled() {
		t.Error("slide numbers should be enabled")
	}
}

func TestSlideNumbersDisable(t *testing.T) {
	pres := New()
	pres.SetSlideNumbers(true)
	pres.SetSlideNumbers(false)

	if pres.SlideNumbersEnabled() {
		t.Error("slide numbers should be disabled")
	}
}

func TestSlideNumberStart(t *testing.T) {
	pres := New()
	pres.SetSlideNumberStart(5)

	if pres.SlideNumberStart() != 5 {
		t.Errorf("expected start number 5, got %d", pres.SlideNumberStart())
	}
}

func TestSlideNumberSerialization(t *testing.T) {
	pres := New()
	pres.SetSlideNumbers(true)
	pres.AddSlide()

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
	if !strings.Contains(xmlStr, "slidenum") {
		t.Error("slide XML should contain slidenum field type")
	}
	if !strings.Contains(xmlStr, "Slide Number") {
		t.Error("slide XML should contain Slide Number shape name")
	}
}

func TestSlideNumberNotPresentWhenDisabled(t *testing.T) {
	pres := New()
	pres.AddSlide()

	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackageFromBytes(data)
	slideXML, _ := pkg.GetFile("ppt/slides/slide1.xml")

	if strings.Contains(string(slideXML), "slidenum") {
		t.Error("slide XML should NOT contain slidenum when disabled")
	}
}

func TestSlideNumberOnMultipleSlides(t *testing.T) {
	pres := New()
	pres.SetSlideNumbers(true)
	pres.AddSlide()
	pres.AddSlide()
	pres.AddSlide()

	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackageFromBytes(data)

	for i := 1; i <= 3; i++ {
		slideXML, ok := pkg.GetFile("ppt/slides/slide" + strings.Repeat("", 0) + string(rune('0'+i)) + ".xml")
		if !ok {
			t.Fatalf("slide%d.xml should exist", i)
		}
		if !strings.Contains(string(slideXML), "slidenum") {
			t.Errorf("slide%d should contain slidenum", i)
		}
	}
}
