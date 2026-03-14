package presentation

// SetSlideNumbers enables or disables slide numbers
func (p *Presentation) SetSlideNumbers(show bool) { p.showSlideNumbers = show }

// SetSlideNumberStart sets the starting slide number
func (p *Presentation) SetSlideNumberStart(start int) { p.slideNumberStart = start }

// SlideNumbersEnabled returns whether slide numbers are enabled
func (p *Presentation) SlideNumbersEnabled() bool { return p.showSlideNumbers }

// SlideNumberStart returns the starting slide number
func (p *Presentation) SlideNumberStart() int { return p.slideNumberStart }
