package presentation

import (
	"strings"
	"testing"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
)

func TestAddAnimation(t *testing.T) {
	pres := New()
	s := pres.AddSlide()
	s.AddTextBox(common.In(1), common.In(1), common.In(5), common.In(2))

	anim := s.AddAnimation(0, AnimFadeIn, TriggerOnClick)
	if anim == nil {
		t.Fatal("animation should not be nil")
	}
	if anim.Type != AnimFadeIn {
		t.Error("animation type should be AnimFadeIn")
	}
	if anim.Trigger != TriggerOnClick {
		t.Error("trigger should be TriggerOnClick")
	}
	if anim.Duration != 500 {
		t.Errorf("default duration should be 500, got %d", anim.Duration)
	}
	if len(s.Animations()) != 1 {
		t.Errorf("expected 1 animation, got %d", len(s.Animations()))
	}
}

func TestAnimationSetDurationAndDelay(t *testing.T) {
	pres := New()
	s := pres.AddSlide()
	s.AddShape(ShapeRectangle, common.In(1), common.In(1), common.In(2), common.In(2))

	anim := s.AddAnimation(0, AnimFlyIn, TriggerAfterPrevious)
	anim.SetDuration(1000)
	anim.SetDelay(250)

	if anim.Duration != 1000 {
		t.Errorf("expected duration 1000, got %d", anim.Duration)
	}
	if anim.Delay != 250 {
		t.Errorf("expected delay 250, got %d", anim.Delay)
	}
}

func TestAnimationSerialization(t *testing.T) {
	pres := New()
	s := pres.AddSlide()
	s.AddTextBox(common.In(1), common.In(1), common.In(5), common.In(2))
	s.AddAnimation(0, AnimFadeIn, TriggerOnClick)

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
	if !strings.Contains(xmlStr, "timing") {
		t.Error("slide XML should contain timing element")
	}
	if !strings.Contains(xmlStr, "mainSeq") {
		t.Error("slide XML should contain mainSeq node")
	}
	if !strings.Contains(xmlStr, "presetID") {
		t.Error("slide XML should contain presetID attribute")
	}
}

func TestMultipleAnimations(t *testing.T) {
	pres := New()
	s := pres.AddSlide()
	s.AddTextBox(common.In(1), common.In(1), common.In(5), common.In(2))
	s.AddShape(ShapeCircle, common.In(2), common.In(3), common.In(2), common.In(2))

	s.AddAnimation(0, AnimFadeIn, TriggerOnClick)
	s.AddAnimation(1, AnimZoomIn, TriggerWithPrevious)
	s.AddAnimation(0, AnimFadeOut, TriggerAfterPrevious)

	if len(s.Animations()) != 3 {
		t.Errorf("expected 3 animations, got %d", len(s.Animations()))
	}

	data, err := pres.SaveToBytes()
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	pkg, _ := packaging.OpenPackageFromBytes(data)
	slideXML, _ := pkg.GetFile("ppt/slides/slide1.xml")
	xmlStr := string(slideXML)

	if !strings.Contains(xmlStr, "entr") {
		t.Error("should contain entrance animation class")
	}
	if !strings.Contains(xmlStr, "exit") {
		t.Error("should contain exit animation class")
	}
}

func TestAnimationTypes(t *testing.T) {
	types := []struct {
		anim     AnimationType
		entrance bool
	}{
		{AnimFadeIn, true},
		{AnimFadeOut, false},
		{AnimFlyIn, true},
		{AnimFlyOut, false},
		{AnimZoomIn, true},
		{AnimZoomOut, false},
		{AnimBounce, true},
		{AnimWipe, true},
		{AnimSplit, true},
		{AnimAppear, true},
		{AnimDisappear, false},
	}

	for _, tc := range types {
		if animationIsEntrance(tc.anim) != tc.entrance {
			t.Errorf("animation type %d: expected entrance=%v", tc.anim, tc.entrance)
		}
	}
}
