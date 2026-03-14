package presentation

import (
	"encoding/xml"
	"fmt"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/internal/xmlutil"
)

func (p *Presentation) build() error {
	p.pkg = packaging.NewPackage()

	presRels := packaging.NewRelationships()

	// Build each slide
	for i, slide := range p.slides {
		// Handle image relationships for this slide
		slideRels := packaging.NewRelationships()
		imgIdx := 0
		for _, elem := range slide.elements {
			if img, ok := elem.(*SlideImage); ok {
				imgIdx++
				ext := img.data.Format.Extension()
				mediaPath := fmt.Sprintf("ppt/media/slide%d_img%d%s", i+1, imgIdx, ext)
				p.pkg.AddFile(mediaPath, img.data.Data)
				relTarget := fmt.Sprintf("../media/slide%d_img%d%s", i+1, imgIdx, ext)
				img.relID = slideRels.Add(packaging.RelTypeImage, relTarget)
			}
		}

		if imgIdx > 0 {
			slideRelsData, err := slideRels.Marshal()
			if err != nil {
				return fmt.Errorf("marshal slide rels: %w", err)
			}
			p.pkg.AddFile(fmt.Sprintf("ppt/slides/_rels/slide%d.xml.rels", i+1), slideRelsData)
		}

		slidePath := fmt.Sprintf("ppt/slides/slide%d.xml", i+1)
		data, err := p.buildSlideXML(slide)
		if err != nil {
			return fmt.Errorf("build slide %d: %w", i+1, err)
		}
		p.pkg.AddFile(slidePath, data)
		presRels.Add(packaging.RelTypeSlide, fmt.Sprintf("slides/slide%d.xml", i+1))
	}

	// Presentation XML
	presData, err := p.buildPresentationXML()
	if err != nil {
		return fmt.Errorf("build presentation: %w", err)
	}
	p.pkg.AddFile("ppt/presentation.xml", presData)

	// Presentation relationships
	presRelsData, err := presRels.Marshal()
	if err != nil {
		return fmt.Errorf("marshal pres rels: %w", err)
	}
	p.pkg.AddFile("ppt/_rels/presentation.xml.rels", presRelsData)

	// Top-level relationships
	topRels := packaging.NewRelationships()
	topRels.Add(packaging.RelTypeOfficeDocument, "ppt/presentation.xml")
	topRelsData, err := topRels.Marshal()
	if err != nil {
		return fmt.Errorf("marshal top rels: %w", err)
	}
	p.pkg.AddFile("_rels/.rels", topRelsData)

	// Content types
	ct := packaging.NewContentTypes()
	ct.AddOverride("/ppt/presentation.xml", packaging.ContentTypePptx)
	for i := range p.slides {
		ct.AddOverride(fmt.Sprintf("/ppt/slides/slide%d.xml", i+1), packaging.ContentTypeSlide)
	}
	ctData, err := ct.Marshal()
	if err != nil {
		return fmt.Errorf("marshal content types: %w", err)
	}
	p.pkg.AddFile("[Content_Types].xml", ctData)

	return nil
}

// XML types

type xmlPresentation struct {
	XMLName xml.Name       `xml:"p:presentation"`
	P       string         `xml:"xmlns:p,attr"`
	R       string         `xml:"xmlns:r,attr"`
	A       string         `xml:"xmlns:a,attr"`
	SldSz   xmlSlideSize   `xml:"p:sldSz"`
	SldList xmlSlideIDList `xml:"p:sldIdLst"`
}

type xmlSlideSize struct {
	Cx string `xml:"cx,attr"`
	Cy string `xml:"cy,attr"`
}

type xmlSlideIDList struct {
	SldID []xmlSlideID `xml:"p:sldId"`
}

type xmlSlideID struct {
	ID  string `xml:"id,attr"`
	RID string `xml:"r:id,attr"`
}

func (p *Presentation) buildPresentationXML() ([]byte, error) {
	xp := xmlPresentation{
		P: "http://schemas.openxmlformats.org/presentationml/2006/main",
		R: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
		A: "http://schemas.openxmlformats.org/drawingml/2006/main",
		SldSz: xmlSlideSize{
			Cx: fmt.Sprintf("%d", p.width.EMUs()),
			Cy: fmt.Sprintf("%d", p.height.EMUs()),
		},
	}

	for i := range p.slides {
		xp.SldList.SldID = append(xp.SldList.SldID, xmlSlideID{
			ID:  fmt.Sprintf("%d", 256+i),
			RID: fmt.Sprintf("rId%d", i+1),
		})
	}

	return xmlutil.MarshalXML(xp)
}

type xmlTransition struct {
	XMLName  xml.Name   `xml:"p:transition"`
	Speed    string     `xml:"spd,attr,omitempty"`
	Fade     *xmlEmpty2 `xml:"p:fade,omitempty"`
	Push     *xmlEmpty2 `xml:"p:push,omitempty"`
	Wipe     *xmlEmpty2 `xml:"p:wipe,omitempty"`
	Split    *xmlEmpty2 `xml:"p:split,omitempty"`
	Cover    *xmlEmpty2 `xml:"p:cover,omitempty"`
	Cut      *xmlEmpty2 `xml:"p:cut,omitempty"`
	Dissolve *xmlEmpty2 `xml:"p:dissolve,omitempty"`
}

type xmlPic struct {
	XMLName  xml.Name    `xml:"p:pic"`
	NvPicPr  xmlNvPicPr  `xml:"p:nvPicPr"`
	BlipFill xmlBlipFill `xml:"p:blipFill"`
	SpPr     xmlSpPr     `xml:"p:spPr"`
}

type xmlNvPicPr struct {
	CNvPr    xmlCNvPr  `xml:"p:cNvPr"`
	CNvPicPr xmlEmpty2 `xml:"p:cNvPicPr"`
	NvPr     xmlEmpty2 `xml:"p:nvPr"`
}

type xmlBlipFill struct {
	Blip    xmlBlip    `xml:"a:blip"`
	Stretch xmlStretch `xml:"a:stretch"`
}

type xmlBlip struct {
	Embed string `xml:"r:embed,attr"`
}

type xmlStretch struct {
	FillRect xmlEmpty2 `xml:"a:fillRect"`
}

type xmlSlide struct {
	XMLName    xml.Name       `xml:"p:sld"`
	P          string         `xml:"xmlns:p,attr"`
	A          string         `xml:"xmlns:a,attr"`
	R          string         `xml:"xmlns:r,attr"`
	CSld       xmlCSld        `xml:"p:cSld"`
	Transition *xmlTransition `xml:"p:transition,omitempty"`
}

type xmlCSld struct {
	Bg     *xmlBg    `xml:"p:bg,omitempty"`
	SpTree xmlSpTree `xml:"p:spTree"`
}

type xmlBg struct {
	BgPr xmlBgPr `xml:"p:bgPr"`
}

type xmlBgPr struct {
	SolidFill *xmlSolidFill `xml:"a:solidFill,omitempty"`
}

type xmlSolidFill struct {
	SrgbClr xmlSrgbClr `xml:"a:srgbClr"`
}

type xmlSrgbClr struct {
	Val string `xml:"val,attr"`
}

type xmlSpTree struct {
	NvGrpSpPr xmlNvGrpSpPr `xml:"p:nvGrpSpPr"`
	GrpSpPr   xmlGrpSpPr   `xml:"p:grpSpPr"`
	Shapes    []xmlSp      `xml:"p:sp"`
	Pics      []xmlPic     `xml:"p:pic"`
}

type xmlNvGrpSpPr struct {
	CNvPr      xmlCNvPr  `xml:"p:cNvPr"`
	CNvGrpSpPr xmlEmpty2 `xml:"p:cNvGrpSpPr"`
	NvPr       xmlEmpty2 `xml:"p:nvPr"`
}

type xmlCNvPr struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type xmlGrpSpPr struct{}
type xmlEmpty2 struct{}

type xmlSp struct {
	NvSpPr xmlNvSpPr  `xml:"p:nvSpPr"`
	SpPr   xmlSpPr    `xml:"p:spPr"`
	TxBody *xmlTxBody `xml:"p:txBody,omitempty"`
}

type xmlNvSpPr struct {
	CNvPr   xmlCNvPr  `xml:"p:cNvPr"`
	CNvSpPr xmlEmpty2 `xml:"p:cNvSpPr"`
	NvPr    xmlEmpty2 `xml:"p:nvPr"`
}

type xmlSpPr struct {
	Xfrm      *xmlXfrm      `xml:"a:xfrm,omitempty"`
	PrstGeom  *xmlPrstGeom  `xml:"a:prstGeom,omitempty"`
	SolidFill *xmlSolidFill `xml:"a:solidFill,omitempty"`
}

type xmlXfrm struct {
	Off xmlOff `xml:"a:off"`
	Ext xmlExt `xml:"a:ext"`
}

type xmlOff struct {
	X string `xml:"x,attr"`
	Y string `xml:"y,attr"`
}

type xmlExt struct {
	Cx string `xml:"cx,attr"`
	Cy string `xml:"cy,attr"`
}

type xmlPrstGeom struct {
	Prst string `xml:"prst,attr"`
}

type xmlTxBody struct {
	BodyPr xmlEmpty2  `xml:"a:bodyPr"`
	Paras  []xmlAPara `xml:"a:p"`
}

type xmlAPara struct {
	PPr  *xmlAPPr  `xml:"a:pPr,omitempty"`
	Runs []xmlARun `xml:"a:r"`
}

type xmlAPPr struct {
	Algn string `xml:"algn,attr,omitempty"`
}

type xmlARun struct {
	RPr *xmlARPr `xml:"a:rPr,omitempty"`
	T   string   `xml:"a:t"`
}

type xmlARPr struct {
	Lang      string        `xml:"lang,attr,omitempty"`
	Sz        string        `xml:"sz,attr,omitempty"`
	B         string        `xml:"b,attr,omitempty"`
	I         string        `xml:"i,attr,omitempty"`
	SolidFill *xmlSolidFill `xml:"a:solidFill,omitempty"`
}

func (p *Presentation) buildSlideXML(slide *Slide) ([]byte, error) {
	xs := xmlSlide{
		P: "http://schemas.openxmlformats.org/presentationml/2006/main",
		A: "http://schemas.openxmlformats.org/drawingml/2006/main",
		R: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
	}

	// Background
	if slide.background != nil {
		hex := colorToHex(*slide.background)
		xs.CSld.Bg = &xmlBg{
			BgPr: xmlBgPr{
				SolidFill: &xmlSolidFill{
					SrgbClr: xmlSrgbClr{Val: hex},
				},
			},
		}
	}

	// Shape tree
	xs.CSld.SpTree.NvGrpSpPr.CNvPr = xmlCNvPr{ID: "1", Name: ""}

	shapeID := 2
	for _, elem := range slide.elements {
		switch e := elem.(type) {
		case *TextBox:
			sp := buildTextBoxXML(e, shapeID)
			xs.CSld.SpTree.Shapes = append(xs.CSld.SpTree.Shapes, sp)
			shapeID++
		case *Shape:
			sp := buildShapeXML(e, shapeID)
			xs.CSld.SpTree.Shapes = append(xs.CSld.SpTree.Shapes, sp)
			shapeID++
		case *SlideImage:
			pic := buildPictureXML(e, shapeID)
			xs.CSld.SpTree.Pics = append(xs.CSld.SpTree.Pics, pic)
			shapeID++
		}
	}

	// Transition
	if slide.transition != nil && slide.transition.Type != TransitionNone {
		speed := "med"
		switch slide.transition.Speed {
		case TransitionSlow:
			speed = "slow"
		case TransitionFast:
			speed = "fast"
		}

		xs.Transition = &xmlTransition{Speed: speed}
		switch slide.transition.Type {
		case TransitionFade:
			xs.Transition.Fade = &xmlEmpty2{}
		case TransitionPush:
			xs.Transition.Push = &xmlEmpty2{}
		case TransitionWipe:
			xs.Transition.Wipe = &xmlEmpty2{}
		case TransitionSplit:
			xs.Transition.Split = &xmlEmpty2{}
		case TransitionCover:
			xs.Transition.Cover = &xmlEmpty2{}
		case TransitionCut:
			xs.Transition.Cut = &xmlEmpty2{}
		case TransitionDissolve:
			xs.Transition.Dissolve = &xmlEmpty2{}
		}
	}

	return xmlutil.MarshalXML(xs)
}

func buildTextBoxXML(tb *TextBox, id int) xmlSp {
	sp := xmlSp{
		NvSpPr: xmlNvSpPr{
			CNvPr: xmlCNvPr{ID: fmt.Sprintf("%d", id), Name: fmt.Sprintf("TextBox %d", id)},
		},
		SpPr: xmlSpPr{
			Xfrm: &xmlXfrm{
				Off: xmlOff{X: fmt.Sprintf("%d", tb.x.EMUs()), Y: fmt.Sprintf("%d", tb.y.EMUs())},
				Ext: xmlExt{Cx: fmt.Sprintf("%d", tb.width.EMUs()), Cy: fmt.Sprintf("%d", tb.height.EMUs())},
			},
			PrstGeom: &xmlPrstGeom{Prst: "rect"},
		},
	}

	if tb.fillColor != nil {
		sp.SpPr.SolidFill = &xmlSolidFill{
			SrgbClr: xmlSrgbClr{Val: colorToHex(*tb.fillColor)},
		}
	}

	if len(tb.paragraphs) > 0 {
		txBody := &xmlTxBody{}
		for _, para := range tb.paragraphs {
			xp := xmlAPara{}
			if para.alignment != common.TextAlignLeft {
				algn := "l"
				switch para.alignment {
				case common.TextAlignCenter:
					algn = "ctr"
				case common.TextAlignRight:
					algn = "r"
				case common.TextAlignJustify:
					algn = "just"
				}
				xp.PPr = &xmlAPPr{Algn: algn}
			}
			for _, run := range para.runs {
				xr := xmlARun{T: run.text}
				rPr := &xmlARPr{
					Lang: "en-US",
					Sz:   fmt.Sprintf("%d", int(run.font.Size*100)),
				}
				if run.bold {
					rPr.B = "1"
				}
				if run.italic {
					rPr.I = "1"
				}
				hex := colorToHex(run.color)
				if hex != "000000" {
					rPr.SolidFill = &xmlSolidFill{
						SrgbClr: xmlSrgbClr{Val: hex},
					}
				}
				xr.RPr = rPr
				xp.Runs = append(xp.Runs, xr)
			}
			txBody.Paras = append(txBody.Paras, xp)
		}
		sp.TxBody = txBody
	}

	return sp
}

func buildShapeXML(sh *Shape, id int) xmlSp {
	prst := "rect"
	switch sh.shapeType {
	case ShapeRoundedRectangle:
		prst = "roundRect"
	case ShapeCircle, ShapeEllipse:
		prst = "ellipse"
	case ShapeTriangle:
		prst = "triangle"
	case ShapeArrowRight:
		prst = "rightArrow"
	case ShapeArrowLeft:
		prst = "leftArrow"
	case ShapeArrowUp:
		prst = "upArrow"
	case ShapeArrowDown:
		prst = "downArrow"
	case ShapeStar:
		prst = "star5"
	case ShapeDiamond:
		prst = "diamond"
	case ShapeLine:
		prst = "line"
	}

	sp := xmlSp{
		NvSpPr: xmlNvSpPr{
			CNvPr: xmlCNvPr{ID: fmt.Sprintf("%d", id), Name: fmt.Sprintf("Shape %d", id)},
		},
		SpPr: xmlSpPr{
			Xfrm: &xmlXfrm{
				Off: xmlOff{X: fmt.Sprintf("%d", sh.x.EMUs()), Y: fmt.Sprintf("%d", sh.y.EMUs())},
				Ext: xmlExt{Cx: fmt.Sprintf("%d", sh.width.EMUs()), Cy: fmt.Sprintf("%d", sh.height.EMUs())},
			},
			PrstGeom:  &xmlPrstGeom{Prst: prst},
			SolidFill: &xmlSolidFill{SrgbClr: xmlSrgbClr{Val: colorToHex(sh.fillColor)}},
		},
	}

	// Text inside shape
	if sh.text != "" && sh.textFont != nil {
		sp.TxBody = &xmlTxBody{
			Paras: []xmlAPara{
				{
					PPr: &xmlAPPr{Algn: "ctr"},
					Runs: []xmlARun{
						{
							RPr: &xmlARPr{
								Lang: "en-US",
								Sz:   fmt.Sprintf("%d", int(sh.textFont.Size*100)),
							},
							T: sh.text,
						},
					},
				},
			},
		}
	}

	return sp
}

func buildPictureXML(img *SlideImage, id int) xmlPic {
	return xmlPic{
		NvPicPr: xmlNvPicPr{
			CNvPr: xmlCNvPr{ID: fmt.Sprintf("%d", id), Name: fmt.Sprintf("Picture %d", id)},
		},
		BlipFill: xmlBlipFill{
			Blip: xmlBlip{Embed: img.relID},
		},
		SpPr: xmlSpPr{
			Xfrm: &xmlXfrm{
				Off: xmlOff{X: fmt.Sprintf("%d", img.x.EMUs()), Y: fmt.Sprintf("%d", img.y.EMUs())},
				Ext: xmlExt{Cx: fmt.Sprintf("%d", img.width.EMUs()), Cy: fmt.Sprintf("%d", img.height.EMUs())},
			},
			PrstGeom: &xmlPrstGeom{Prst: "rect"},
		},
	}
}

func colorToHex(c common.Color) string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}
