package common

// Unit represents a unit of measurement
type Unit int

const (
	UnitPoint      Unit = iota // 1/72 inch
	UnitInch                   // 1 inch
	UnitCentimeter             // 1 cm
	UnitMillimeter             // 1 mm
	UnitPixel                  // 1/96 inch
	UnitEMU                    // English Metric Unit (1/914400 inch)
)

// Measurement stores a value in points internally
type Measurement struct {
	points float64
}

func Pt(v float64) Measurement { return Measurement{points: v} }
func In(v float64) Measurement { return Measurement{points: v * 72} }
func Cm(v float64) Measurement { return Measurement{points: v * 28.3465} }
func Mm(v float64) Measurement { return Measurement{points: v * 2.83465} }
func Px(v float64) Measurement { return Measurement{points: v * 0.75} }
func EMU(v int64) Measurement  { return Measurement{points: float64(v) / 12700} }

func (m Measurement) Points() float64      { return m.points }
func (m Measurement) Inches() float64      { return m.points / 72 }
func (m Measurement) Centimeters() float64 { return m.points / 28.3465 }
func (m Measurement) Millimeters() float64 { return m.points / 2.83465 }
func (m Measurement) Pixels() float64      { return m.points / 0.75 }
func (m Measurement) EMUs() int64          { return int64(m.points * 12700) }

// PageSize represents standard page sizes
type PageSize struct {
	Width  Measurement
	Height Measurement
}

var (
	PageA4      = PageSize{Width: Mm(210), Height: Mm(297)}
	PageA3      = PageSize{Width: Mm(297), Height: Mm(420)}
	PageA5      = PageSize{Width: Mm(148), Height: Mm(210)}
	PageLetter  = PageSize{Width: In(8.5), Height: In(11)}
	PageLegal   = PageSize{Width: In(8.5), Height: In(14)}
	PageTabloid = PageSize{Width: In(11), Height: In(17)}
)

// Orientation
type Orientation int

const (
	OrientationPortrait Orientation = iota
	OrientationLandscape
)

// Margins
type Margins struct {
	Top    Measurement
	Right  Measurement
	Bottom Measurement
	Left   Measurement
}

func NewMargins(top, right, bottom, left Measurement) Margins {
	return Margins{Top: top, Right: right, Bottom: bottom, Left: left}
}

func UniformMargins(m Measurement) Margins {
	return Margins{Top: m, Right: m, Bottom: m, Left: m}
}

// NormalMargins returns standard 1-inch margins
func NormalMargins() Margins { return UniformMargins(In(1)) }

// NarrowMargins returns 0.5-inch margins
func NarrowMargins() Margins { return UniformMargins(In(0.5)) }
