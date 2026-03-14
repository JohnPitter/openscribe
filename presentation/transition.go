package presentation

// TransitionType represents a slide transition effect
type TransitionType int

const (
	TransitionNone TransitionType = iota
	TransitionFade
	TransitionPush
	TransitionWipe
	TransitionSplit
	TransitionCover
	TransitionUncover
	TransitionCut
	TransitionDissolve
	TransitionZoom
)

// TransitionSpeed represents transition speed
type TransitionSpeed int

const (
	TransitionSlow TransitionSpeed = iota
	TransitionMedium
	TransitionFast
)

// Transition represents a slide transition configuration
type Transition struct {
	Type     TransitionType
	Speed    TransitionSpeed
	Duration float64 // in seconds
}

// NewTransition creates a new transition
func NewTransition(t TransitionType, speed TransitionSpeed) Transition {
	durations := map[TransitionSpeed]float64{
		TransitionSlow:   1.0,
		TransitionMedium: 0.5,
		TransitionFast:   0.25,
	}
	return Transition{
		Type:     t,
		Speed:    speed,
		Duration: durations[speed],
	}
}
