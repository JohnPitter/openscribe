package presentation

// AnimationType represents the type of animation effect
type AnimationType int

const (
	AnimFadeIn AnimationType = iota
	AnimFadeOut
	AnimFlyIn
	AnimFlyOut
	AnimZoomIn
	AnimZoomOut
	AnimBounce
	AnimWipe
	AnimSplit
	AnimAppear
	AnimDisappear
)

// AnimationTrigger defines when an animation starts
type AnimationTrigger int

const (
	TriggerOnClick AnimationTrigger = iota
	TriggerWithPrevious
	TriggerAfterPrevious
)

// Animation represents an animation applied to a slide element
type Animation struct {
	Type         AnimationType
	Trigger      AnimationTrigger
	Duration     int // milliseconds
	Delay        int // milliseconds
	ElementIndex int
}

// SetDuration sets the animation duration in milliseconds
func (a *Animation) SetDuration(ms int) { a.Duration = ms }

// SetDelay sets the animation delay in milliseconds
func (a *Animation) SetDelay(ms int) { a.Delay = ms }

// AddAnimation adds an animation to an element on the slide
func (s *Slide) AddAnimation(elementIndex int, animType AnimationType, trigger AnimationTrigger) *Animation {
	a := &Animation{
		Type:         animType,
		Trigger:      trigger,
		Duration:     500,
		Delay:        0,
		ElementIndex: elementIndex,
	}
	s.animations = append(s.animations, a)
	return a
}

// Animations returns all animations on the slide
func (s *Slide) Animations() []*Animation { return s.animations }

// animationPresetID returns the OOXML preset animation ID
func animationPresetID(t AnimationType) int {
	switch t {
	case AnimFadeIn, AnimFadeOut:
		return 10
	case AnimFlyIn, AnimFlyOut:
		return 2
	case AnimZoomIn, AnimZoomOut:
		return 53
	case AnimBounce:
		return 26
	case AnimWipe:
		return 22
	case AnimSplit:
		return 16
	case AnimAppear:
		return 1
	case AnimDisappear:
		return 1
	default:
		return 1
	}
}

// animationIsEntrance returns true if the animation is an entrance effect
func animationIsEntrance(t AnimationType) bool {
	switch t {
	case AnimFadeIn, AnimFlyIn, AnimZoomIn, AnimBounce, AnimWipe, AnimSplit, AnimAppear:
		return true
	default:
		return false
	}
}
