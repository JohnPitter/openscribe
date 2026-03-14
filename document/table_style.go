package document

import "github.com/JohnPitter/openscribe/common"

// TableStylePreset defines preset table styling
type TableStylePreset int

const (
	// TableStylePlain has thin black borders, no shading
	TableStylePlain TableStylePreset = iota
	// TableStyleStriped alternates row background colors
	TableStyleStriped
	// TableStyleBanded uses banded columns
	TableStyleBanded
	// TableStyleGrid has thick borders forming a grid
	TableStyleGrid
	// TableStyleDark has a dark header with light text
	TableStyleDark
	// TableStyleColorful uses a colorful header with alternating rows
	TableStyleColorful
)

// tableStyleConfig holds computed styling for a table style preset
type tableStyleConfig struct {
	headerBG    common.Color
	headerFG    common.Color
	evenRowBG   common.Color
	oddRowBG    common.Color
	borderStyle common.BorderStyle
	borderColor common.Color
	borderWidth common.Measurement
	headerBold  bool
}

func configForStyle(preset TableStylePreset) tableStyleConfig {
	switch preset {
	case TableStyleStriped:
		return tableStyleConfig{
			headerBG:    common.NewColor(68, 114, 196),
			headerFG:    common.White,
			evenRowBG:   common.White,
			oddRowBG:    common.NewColor(217, 226, 243),
			borderStyle: common.BorderStyleSingle,
			borderColor: common.NewColor(68, 114, 196),
			borderWidth: common.Pt(0.5),
			headerBold:  true,
		}
	case TableStyleBanded:
		return tableStyleConfig{
			headerBG:    common.NewColor(91, 155, 213),
			headerFG:    common.White,
			evenRowBG:   common.NewColor(222, 235, 247),
			oddRowBG:    common.White,
			borderStyle: common.BorderStyleSingle,
			borderColor: common.NewColor(91, 155, 213),
			borderWidth: common.Pt(0.5),
			headerBold:  true,
		}
	case TableStyleGrid:
		return tableStyleConfig{
			headerBG:    common.White,
			headerFG:    common.Black,
			evenRowBG:   common.White,
			oddRowBG:    common.White,
			borderStyle: common.BorderStyleSingle,
			borderColor: common.Black,
			borderWidth: common.Pt(1.5),
			headerBold:  true,
		}
	case TableStyleDark:
		return tableStyleConfig{
			headerBG:    common.NewColor(51, 51, 51),
			headerFG:    common.White,
			evenRowBG:   common.White,
			oddRowBG:    common.NewColor(242, 242, 242),
			borderStyle: common.BorderStyleSingle,
			borderColor: common.NewColor(51, 51, 51),
			borderWidth: common.Pt(0.5),
			headerBold:  true,
		}
	case TableStyleColorful:
		return tableStyleConfig{
			headerBG:    common.NewColor(255, 165, 0),
			headerFG:    common.White,
			evenRowBG:   common.White,
			oddRowBG:    common.NewColor(255, 243, 224),
			borderStyle: common.BorderStyleSingle,
			borderColor: common.NewColor(255, 165, 0),
			borderWidth: common.Pt(0.5),
			headerBold:  true,
		}
	default: // TableStylePlain
		return tableStyleConfig{
			headerBG:    common.White,
			headerFG:    common.Black,
			evenRowBG:   common.White,
			oddRowBG:    common.White,
			borderStyle: common.BorderStyleSingle,
			borderColor: common.Black,
			borderWidth: common.Pt(0.5),
			headerBold:  false,
		}
	}
}

// SetStyle applies a preset style to the table
func (t *Table) SetStyle(preset TableStylePreset) {
	cfg := configForStyle(preset)

	// Set borders
	border := common.Border{
		Style: cfg.borderStyle,
		Color: cfg.borderColor,
		Width: cfg.borderWidth,
	}
	t.borders = common.Borders{
		Top:    border,
		Left:   border,
		Bottom: border,
		Right:  border,
	}

	// Apply row styling
	for i, row := range t.rows {
		for _, cell := range row.cells {
			if i == 0 {
				// Header row
				cell.SetShading(cfg.headerBG)
				if cfg.headerBold {
					for _, p := range cell.paragraphs {
						for _, r := range p.runs {
							r.SetBold(true)
							r.SetColor(cfg.headerFG)
						}
					}
				}
			} else if i%2 == 0 {
				cell.SetShading(cfg.evenRowBG)
			} else {
				cell.SetShading(cfg.oddRowBG)
			}
		}
	}
}

// StylePreset returns the name string for a table style preset (for debugging)
func (p TableStylePreset) String() string {
	names := []string{"Plain", "Striped", "Banded", "Grid", "Dark", "Colorful"}
	if int(p) < len(names) {
		return names[p]
	}
	return "Unknown"
}
