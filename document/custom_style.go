package document

import (
	"fmt"

	"github.com/JohnPitter/openscribe/common"
)

// CustomStyle represents a user-defined paragraph style
type CustomStyle struct {
	name      string
	basedOn   string
	font      *common.Font
	alignment common.TextAlignment
	spacing   *styleSpacing
	indent    *styleIndent
}

// styleSpacing holds custom style spacing values (in twips)
type styleSpacing struct {
	before int
	after  int
	line   int
}

// styleIndent holds custom style indentation values (in twips)
type styleIndent struct {
	left      int
	right     int
	firstLine int
}

// NewCustomStyle creates a new custom style
func NewCustomStyle(name, basedOn string) *CustomStyle {
	return &CustomStyle{
		name:    name,
		basedOn: basedOn,
	}
}

// Name returns the style name
func (cs *CustomStyle) Name() string {
	return cs.name
}

// BasedOn returns the parent style name
func (cs *CustomStyle) BasedOn() string {
	return cs.basedOn
}

// SetFont sets the font for this style
func (cs *CustomStyle) SetFont(f common.Font) *CustomStyle {
	cs.font = &f
	return cs
}

// SetAlignment sets the text alignment for this style
func (cs *CustomStyle) SetAlignment(a common.TextAlignment) *CustomStyle {
	cs.alignment = a
	return cs
}

// SetSpacing sets the before/after/line spacing in points
func (cs *CustomStyle) SetSpacing(before, after, line float64) *CustomStyle {
	cs.spacing = &styleSpacing{
		before: int(before * 20), // convert pt to twips
		after:  int(after * 20),
		line:   int(line * 240), // line spacing: 240 twips = 1.0 spacing
	}
	return cs
}

// SetIndent sets left/right/firstLine indentation in points
func (cs *CustomStyle) SetIndent(left, right, firstLine float64) *CustomStyle {
	cs.indent = &styleIndent{
		left:      int(left * 20),
		right:     int(right * 20),
		firstLine: int(firstLine * 20),
	}
	return cs
}

// toXML serializes the custom style to an OOXML w:style element string
func (cs *CustomStyle) toXML() string {
	var xml string
	xml += fmt.Sprintf(`<w:style w:type="paragraph" w:styleId="%s">`, cs.name)
	xml += fmt.Sprintf(`<w:name w:val="%s"/>`, cs.name)

	if cs.basedOn != "" {
		xml += fmt.Sprintf(`<w:basedOn w:val="%s"/>`, cs.basedOn)
	}

	// Paragraph properties
	hasPPr := cs.alignment != common.TextAlignLeft || cs.spacing != nil || cs.indent != nil
	if hasPPr {
		xml += `<w:pPr>`
		if cs.alignment != common.TextAlignLeft {
			xml += fmt.Sprintf(`<w:jc w:val="%s"/>`, alignmentToString(cs.alignment))
		}
		if cs.spacing != nil {
			xml += fmt.Sprintf(`<w:spacing w:before="%d" w:after="%d" w:line="%d"/>`,
				cs.spacing.before, cs.spacing.after, cs.spacing.line)
		}
		if cs.indent != nil {
			xml += fmt.Sprintf(`<w:ind w:left="%d" w:right="%d" w:firstLine="%d"/>`,
				cs.indent.left, cs.indent.right, cs.indent.firstLine)
		}
		xml += `</w:pPr>`
	}

	// Run properties (font)
	if cs.font != nil {
		xml += `<w:rPr>`
		if cs.font.Family != "" {
			xml += fmt.Sprintf(`<w:rFonts w:ascii="%s" w:hAnsi="%s"/>`, cs.font.Family, cs.font.Family)
		}
		if cs.font.Size > 0 {
			xml += fmt.Sprintf(`<w:sz w:val="%d"/>`, int(cs.font.Size*2))
		}
		if cs.font.Weight >= common.FontWeightBold {
			xml += `<w:b/>`
		}
		if cs.font.Style == common.FontStyleItalic {
			xml += `<w:i/>`
		}
		if cs.font.Color != (common.Color{}) {
			hex := cs.font.Color.Hex()
			if len(hex) > 0 && hex[0] == '#' {
				hex = hex[1:]
			}
			xml += fmt.Sprintf(`<w:color w:val="%s"/>`, hex)
		}
		xml += `</w:rPr>`
	}

	xml += `</w:style>`
	return xml
}
