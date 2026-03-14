package presentation

import (
	"encoding/xml"
	"fmt"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/internal/packaging"
	"github.com/JohnPitter/openscribe/internal/xmlutil"
)

// RelTypeNotesSlide is the relationship type for notes slides
const relTypeNotesSlide = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/notesSlide"

// ContentTypeNotesSlide is the content type for notes slides
const contentTypeNotesSlide = "application/vnd.openxmlformats-officedocument.presentationml.notesSlide+xml"

func (p *Presentation) build() error {
	p.pkg = packaging.NewPackage()

	presRels := packaging.NewRelationships()

	// Build each slide
	for i, slide := range p.slides {
		// Handle image relationships and notes for this slide
		slideRels := packaging.NewRelationships()
		hasRels := false
		imgIdx := 0
		for _, elem := range slide.elements {
			if img, ok := elem.(*SlideImage); ok {
				imgIdx++
				ext := img.data.Format.Extension()
				mediaPath := fmt.Sprintf("ppt/media/slide%d_img%d%s", i+1, imgIdx, ext)
				p.pkg.AddFile(mediaPath, img.data.Data)
				relTarget := fmt.Sprintf("../media/slide%d_img%d%s", i+1, imgIdx, ext)
				img.relID = slideRels.Add(packaging.RelTypeImage, relTarget)
				hasRels = true
			}
		}

		// Build notes part if the slide has formatted notes or plain notes
		hasNotes := len(slide.formattedNotes) > 0 || slide.notes != ""
		notesRelID := ""
		if hasNotes {
			notesRelID = slideRels.Add(relTypeNotesSlide, fmt.Sprintf("../notesSlides/notesSlide%d.xml", i+1))
			hasRels = true

			notesData, err := buildNotesXML(slide, i+1, notesRelID)
			if err != nil {
				return fmt.Errorf("build notes slide %d: %w", i+1, err)
			}
			p.pkg.AddFile(fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", i+1), notesData)

			// Notes slide rels (points back to slide)
			notesRels := packaging.NewRelationships()
			notesRels.Add(packaging.RelTypeSlide, fmt.Sprintf("../slides/slide%d.xml", i+1))
			notesRelsData, err := notesRels.Marshal()
			if err != nil {
				return fmt.Errorf("marshal notes rels: %w", err)
			}
			p.pkg.AddFile(fmt.Sprintf("ppt/notesSlides/_rels/notesSlide%d.xml.rels", i+1), notesRelsData)
		}

		if hasRels {
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
	for i, slide := range p.slides {
		ct.AddOverride(fmt.Sprintf("/ppt/slides/slide%d.xml", i+1), packaging.ContentTypeSlide)
		hasNotes := len(slide.formattedNotes) > 0 || slide.notes != ""
		if hasNotes {
			ct.AddOverride(fmt.Sprintf("/ppt/notesSlides/notesSlide%d.xml", i+1), contentTypeNotesSlide)
		}
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
	Timing     *xmlTiming     `xml:"p:timing,omitempty"`
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
	NvGrpSpPr     xmlNvGrpSpPr      `xml:"p:nvGrpSpPr"`
	GrpSpPr       xmlGrpSpPr        `xml:"p:grpSpPr"`
	Shapes        []xmlSp           `xml:"p:sp"`
	Pics          []xmlPic          `xml:"p:pic"`
	GraphicFrames []xmlGraphicFrame `xml:"p:graphicFrame"`
	CxnSps        []xmlCxnSp        `xml:"p:cxnSp"`
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
	Ln        *xmlLn        `xml:"a:ln,omitempty"`
}

type xmlLn struct {
	W         string        `xml:"w,attr,omitempty"`
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
	PPr  *xmlAPPr    `xml:"a:pPr,omitempty"`
	Runs []xmlARun   `xml:"a:r"`
	Flds []xmlAFld   `xml:"a:fld"`
	EndR *xmlAEndRPr `xml:"a:endParaRPr,omitempty"`
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

type xmlAFld struct {
	ID   string   `xml:"id,attr"`
	Type string   `xml:"type,attr"`
	RPr  *xmlARPr `xml:"a:rPr,omitempty"`
	T    string   `xml:"a:t"`
}

type xmlAEndRPr struct {
	Lang string `xml:"lang,attr,omitempty"`
}

// Table XML types
type xmlGraphicFrame struct {
	XMLName     xml.Name       `xml:"p:graphicFrame"`
	NvGrFramePr xmlNvGrFramePr `xml:"p:nvGraphicFramePr"`
	Xfrm        xmlXfrm        `xml:"p:xfrm"`
	Graphic     xmlGraphic     `xml:"a:graphic"`
}

type xmlNvGrFramePr struct {
	CNvPr        xmlCNvPr  `xml:"p:cNvPr"`
	CNvGrFramePr xmlEmpty2 `xml:"p:cNvGraphicFramePr"`
	NvPr         xmlEmpty2 `xml:"p:nvPr"`
}

type xmlGraphic struct {
	GraphicData xmlGraphicData `xml:"a:graphicData"`
}

type xmlGraphicData struct {
	URI string `xml:"uri,attr"`
	Tbl xmlTbl `xml:"a:tbl"`
}

type xmlTbl struct {
	TblPr   xmlTblPr    `xml:"a:tblPr"`
	TblGrid xmlTblGrid  `xml:"a:tblGrid"`
	Rows    []xmlTblRow `xml:"a:tr"`
}

type xmlTblPr struct {
	FirstRow string `xml:"firstRow,attr,omitempty"`
	BandRow  string `xml:"bandRow,attr,omitempty"`
}

type xmlTblGrid struct {
	GridCols []xmlTblGridCol `xml:"a:gridCol"`
}

type xmlTblGridCol struct {
	W string `xml:"w,attr"`
}

type xmlTblRow struct {
	H     string       `xml:"h,attr"`
	Cells []xmlTblCell `xml:"a:tc"`
}

type xmlTblCell struct {
	TxBody xmlTblCellTxBody `xml:"a:txBody"`
	TcPr   *xmlTcPr         `xml:"a:tcPr,omitempty"`
}

type xmlTblCellTxBody struct {
	BodyPr xmlEmpty2  `xml:"a:bodyPr"`
	Paras  []xmlAPara `xml:"a:p"`
}

type xmlTcPr struct {
	SolidFill *xmlSolidFill `xml:"a:solidFill,omitempty"`
}

// Connector XML types
type xmlCxnSp struct {
	XMLName   xml.Name     `xml:"p:cxnSp"`
	NvCxnSpPr xmlNvCxnSpPr `xml:"p:nvCxnSpPr"`
	SpPr      xmlSpPr      `xml:"p:spPr"`
}

type xmlNvCxnSpPr struct {
	CNvPr      xmlCNvPr  `xml:"p:cNvPr"`
	CNvCxnSpPr xmlEmpty2 `xml:"p:cNvCxnSpPr"`
	NvPr       xmlEmpty2 `xml:"p:nvPr"`
}

// Animation/Timing XML types
type xmlTiming struct {
	XMLName xml.Name `xml:"p:timing"`
	TnLst   xmlTnLst `xml:"p:tnLst"`
}

type xmlTnLst struct {
	Par xmlPar `xml:"p:par"`
}

type xmlPar struct {
	CTn xmlCTn `xml:"p:cTn"`
}

type xmlCTn struct {
	ID       string       `xml:"id,attr"`
	Dur      string       `xml:"dur,attr,omitempty"`
	Restart  string       `xml:"restart,attr,omitempty"`
	NodeType string       `xml:"nodeType,attr,omitempty"`
	ChildLst *xmlChildLst `xml:"p:childTnLst,omitempty"`
}

type xmlChildLst struct {
	Seq []xmlSeq `xml:"p:seq"`
}

type xmlSeq struct {
	CTn xmlSeqCTn `xml:"p:cTn"`
}

type xmlSeqCTn struct {
	ID       string          `xml:"id,attr"`
	Dur      string          `xml:"dur,attr,omitempty"`
	NodeType string          `xml:"nodeType,attr,omitempty"`
	ChildLst *xmlSeqChildLst `xml:"p:childTnLst,omitempty"`
}

type xmlSeqChildLst struct {
	Pars []xmlAnimPar `xml:"p:par"`
}

type xmlAnimPar struct {
	CTn xmlAnimCTn `xml:"p:cTn"`
}

type xmlAnimCTn struct {
	ID          string        `xml:"id,attr"`
	PresetID    string        `xml:"presetID,attr,omitempty"`
	PresetClass string        `xml:"presetClass,attr,omitempty"`
	Fill        string        `xml:"fill,attr,omitempty"`
	NodeType    string        `xml:"nodeType,attr,omitempty"`
	Dur         string        `xml:"dur,attr,omitempty"`
	Delay       string        `xml:"decel,attr,omitempty"`
	StCondLst   *xmlStCondLst `xml:"p:stCondLst,omitempty"`
}

type xmlStCondLst struct {
	Conds []xmlCond `xml:"p:cond"`
}

type xmlCond struct {
	Delay string `xml:"delay,attr,omitempty"`
	Evt   string `xml:"evt,attr,omitempty"`
}

// Notes XML types
type xmlNotes struct {
	XMLName xml.Name `xml:"p:notes"`
	P       string   `xml:"xmlns:p,attr"`
	A       string   `xml:"xmlns:a,attr"`
	R       string   `xml:"xmlns:r,attr"`
	CSld    xmlCSld  `xml:"p:cSld"`
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
		case *SlideTable:
			gf := buildTableXML(e, shapeID)
			xs.CSld.SpTree.GraphicFrames = append(xs.CSld.SpTree.GraphicFrames, gf)
			shapeID++
		case *Connector:
			cxn := buildConnectorXML(e, shapeID)
			xs.CSld.SpTree.CxnSps = append(xs.CSld.SpTree.CxnSps, cxn)
			shapeID++
		}
	}

	// Slide number shape
	if p.showSlideNumbers {
		sp := buildSlideNumberXML(shapeID, p.width, p.height)
		xs.CSld.SpTree.Shapes = append(xs.CSld.SpTree.Shapes, sp)
		shapeID++
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

	// Animations
	if len(slide.animations) > 0 {
		xs.Timing = buildTimingXML(slide.animations)
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
	prst := shapePresetGeom(sh.shapeType)

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

func shapePresetGeom(st ShapeType) string {
	switch st {
	case ShapeRoundedRectangle:
		return "roundRect"
	case ShapeCircle, ShapeEllipse:
		return "ellipse"
	case ShapeTriangle:
		return "triangle"
	case ShapeArrowRight:
		return "rightArrow"
	case ShapeArrowLeft:
		return "leftArrow"
	case ShapeArrowUp:
		return "upArrow"
	case ShapeArrowDown:
		return "downArrow"
	case ShapeStar:
		return "star5"
	case ShapeDiamond:
		return "diamond"
	case ShapeLine:
		return "line"
	case ShapeCallout:
		return "wedgeRoundRectCallout"
	case ShapeFlowchartProcess:
		return "flowChartProcess"
	case ShapeFlowchartDecision:
		return "flowChartDecision"
	case ShapeFlowchartTerminator:
		return "flowChartTerminator"
	case ShapeBrace:
		return "leftBrace"
	case ShapeBracket:
		return "leftBracket"
	default:
		return "rect"
	}
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

func buildTableXML(tbl *SlideTable, id int) xmlGraphicFrame {
	colW := tbl.width.EMUs() / int64(tbl.cols)
	rowH := tbl.height.EMUs() / int64(tbl.rows)

	grid := xmlTblGrid{}
	for c := 0; c < tbl.cols; c++ {
		grid.GridCols = append(grid.GridCols, xmlTblGridCol{
			W: fmt.Sprintf("%d", colW),
		})
	}

	var rows []xmlTblRow
	for r := 0; r < tbl.rows; r++ {
		row := xmlTblRow{
			H: fmt.Sprintf("%d", rowH),
		}
		for c := 0; c < tbl.cols; c++ {
			cell := tbl.cells[r][c]
			tc := xmlTblCell{
				TxBody: xmlTblCellTxBody{
					Paras: []xmlAPara{
						{
							Runs: []xmlARun{
								{
									RPr: buildCellRunPr(cell),
									T:   cell.text,
								},
							},
						},
					},
				},
			}

			// Cell background
			bg := cell.background
			if bg == nil && r == 0 && tbl.headerBg != nil {
				bg = tbl.headerBg
			}
			if bg != nil {
				tc.TcPr = &xmlTcPr{
					SolidFill: &xmlSolidFill{
						SrgbClr: xmlSrgbClr{Val: colorToHex(*bg)},
					},
				}
			}

			row.Cells = append(row.Cells, tc)
		}
		rows = append(rows, row)
	}

	return xmlGraphicFrame{
		NvGrFramePr: xmlNvGrFramePr{
			CNvPr: xmlCNvPr{ID: fmt.Sprintf("%d", id), Name: fmt.Sprintf("Table %d", id)},
		},
		Xfrm: xmlXfrm{
			Off: xmlOff{X: fmt.Sprintf("%d", tbl.x.EMUs()), Y: fmt.Sprintf("%d", tbl.y.EMUs())},
			Ext: xmlExt{Cx: fmt.Sprintf("%d", tbl.width.EMUs()), Cy: fmt.Sprintf("%d", tbl.height.EMUs())},
		},
		Graphic: xmlGraphic{
			GraphicData: xmlGraphicData{
				URI: "http://schemas.openxmlformats.org/drawingml/2006/table",
				Tbl: xmlTbl{
					TblPr: xmlTblPr{
						FirstRow: "1",
						BandRow:  "1",
					},
					TblGrid: grid,
					Rows:    rows,
				},
			},
		},
	}
}

func buildCellRunPr(cell *SlideTableCell) *xmlARPr {
	rPr := &xmlARPr{Lang: "en-US", Sz: "1200"}
	if cell.font != nil {
		rPr.Sz = fmt.Sprintf("%d", int(cell.font.Size*100))
		if cell.font.Weight >= common.FontWeightBold {
			rPr.B = "1"
		}
		if cell.font.Style == common.FontStyleItalic {
			rPr.I = "1"
		}
	}
	return rPr
}

func buildConnectorXML(conn *Connector, id int) xmlCxnSp {
	// Calculate position and size from endpoints
	x := conn.x1.EMUs()
	y := conn.y1.EMUs()
	cx := conn.x2.EMUs() - conn.x1.EMUs()
	cy := conn.y2.EMUs() - conn.y1.EMUs()
	if cx < 0 {
		x = conn.x2.EMUs()
		cx = -cx
	}
	if cy < 0 {
		y = conn.y2.EMUs()
		cy = -cy
	}
	// Ensure minimum size
	if cx == 0 {
		cx = 1
	}
	if cy == 0 {
		cy = 1
	}

	prst := "line"
	switch conn.connType {
	case ConnectorElbow:
		prst = "bentConnector3"
	case ConnectorCurved:
		prst = "curvedConnector3"
	}

	return xmlCxnSp{
		NvCxnSpPr: xmlNvCxnSpPr{
			CNvPr: xmlCNvPr{ID: fmt.Sprintf("%d", id), Name: fmt.Sprintf("Connector %d", id)},
		},
		SpPr: xmlSpPr{
			Xfrm: &xmlXfrm{
				Off: xmlOff{X: fmt.Sprintf("%d", x), Y: fmt.Sprintf("%d", y)},
				Ext: xmlExt{Cx: fmt.Sprintf("%d", cx), Cy: fmt.Sprintf("%d", cy)},
			},
			PrstGeom: &xmlPrstGeom{Prst: prst},
			Ln: &xmlLn{
				W: fmt.Sprintf("%d", conn.width.EMUs()),
				SolidFill: &xmlSolidFill{
					SrgbClr: xmlSrgbClr{Val: colorToHex(conn.color)},
				},
			},
		},
	}
}

func buildSlideNumberXML(id int, slideWidth, slideHeight common.Measurement) xmlSp {
	// Position slide number in the bottom-right footer area
	numWidth := common.In(1.5)
	numHeight := common.In(0.4)
	x := common.EMU(slideWidth.EMUs() - numWidth.EMUs() - common.In(0.5).EMUs())
	y := common.EMU(slideHeight.EMUs() - numHeight.EMUs() - common.In(0.25).EMUs())

	return xmlSp{
		NvSpPr: xmlNvSpPr{
			CNvPr: xmlCNvPr{ID: fmt.Sprintf("%d", id), Name: fmt.Sprintf("Slide Number %d", id)},
		},
		SpPr: xmlSpPr{
			Xfrm: &xmlXfrm{
				Off: xmlOff{X: fmt.Sprintf("%d", x.EMUs()), Y: fmt.Sprintf("%d", y.EMUs())},
				Ext: xmlExt{Cx: fmt.Sprintf("%d", numWidth.EMUs()), Cy: fmt.Sprintf("%d", numHeight.EMUs())},
			},
			PrstGeom: &xmlPrstGeom{Prst: "rect"},
		},
		TxBody: &xmlTxBody{
			Paras: []xmlAPara{
				{
					PPr: &xmlAPPr{Algn: "r"},
					Flds: []xmlAFld{
						{
							ID:   "{B6F15528-F159-4107-2052-41F4E69A2D00}",
							Type: "slidenum",
							RPr:  &xmlARPr{Lang: "en-US", Sz: "1000"},
							T:    "<#>",
						},
					},
				},
			},
		},
	}
}

func buildTimingXML(animations []*Animation) *xmlTiming {
	var pars []xmlAnimPar
	for i, anim := range animations {
		presetClass := "entr"
		if !animationIsEntrance(anim.Type) {
			presetClass = "exit"
		}

		ctn := xmlAnimCTn{
			ID:          fmt.Sprintf("%d", i+3),
			PresetID:    fmt.Sprintf("%d", animationPresetID(anim.Type)),
			PresetClass: presetClass,
			Fill:        "hold",
			Dur:         fmt.Sprintf("%d", anim.Duration),
		}

		// Trigger conditions
		switch anim.Trigger {
		case TriggerOnClick:
			ctn.StCondLst = &xmlStCondLst{
				Conds: []xmlCond{{Delay: "0", Evt: "onClick"}},
			}
			ctn.NodeType = "clickEffect"
		case TriggerWithPrevious:
			ctn.StCondLst = &xmlStCondLst{
				Conds: []xmlCond{{Delay: fmt.Sprintf("%d", anim.Delay)}},
			}
			ctn.NodeType = "withEffect"
		case TriggerAfterPrevious:
			ctn.StCondLst = &xmlStCondLst{
				Conds: []xmlCond{{Delay: fmt.Sprintf("%d", anim.Delay)}},
			}
			ctn.NodeType = "afterEffect"
		}

		pars = append(pars, xmlAnimPar{CTn: ctn})
	}

	return &xmlTiming{
		TnLst: xmlTnLst{
			Par: xmlPar{
				CTn: xmlCTn{
					ID:       "1",
					Dur:      "indefinite",
					Restart:  "never",
					NodeType: "tmRoot",
					ChildLst: &xmlChildLst{
						Seq: []xmlSeq{
							{
								CTn: xmlSeqCTn{
									ID:       "2",
									Dur:      "indefinite",
									NodeType: "mainSeq",
									ChildLst: &xmlSeqChildLst{
										Pars: pars,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildNotesXML(slide *Slide, slideNum int, _ string) ([]byte, error) {
	notes := xmlNotes{
		P: "http://schemas.openxmlformats.org/presentationml/2006/main",
		A: "http://schemas.openxmlformats.org/drawingml/2006/main",
		R: "http://schemas.openxmlformats.org/officeDocument/2006/relationships",
	}

	notes.CSld.SpTree.NvGrpSpPr.CNvPr = xmlCNvPr{ID: "1", Name: ""}

	// Notes text shape
	sp := xmlSp{
		NvSpPr: xmlNvSpPr{
			CNvPr: xmlCNvPr{ID: "2", Name: "Notes Placeholder"},
		},
		SpPr: xmlSpPr{
			Xfrm: &xmlXfrm{
				Off: xmlOff{X: "0", Y: "0"},
				Ext: xmlExt{Cx: fmt.Sprintf("%d", common.In(6).EMUs()), Cy: fmt.Sprintf("%d", common.In(4).EMUs())},
			},
		},
	}

	txBody := &xmlTxBody{}

	if len(slide.formattedNotes) > 0 {
		for _, np := range slide.formattedNotes {
			para := xmlAPara{}
			rPr := &xmlARPr{Lang: "en-US"}
			if np.FontSize > 0 {
				rPr.Sz = fmt.Sprintf("%d", int(np.FontSize*100))
			} else {
				rPr.Sz = "1200"
			}
			if np.Bold {
				rPr.B = "1"
			}
			if np.Italic {
				rPr.I = "1"
			}
			para.Runs = append(para.Runs, xmlARun{
				RPr: rPr,
				T:   np.Text,
			})
			txBody.Paras = append(txBody.Paras, para)
		}
	} else if slide.notes != "" {
		txBody.Paras = append(txBody.Paras, xmlAPara{
			Runs: []xmlARun{
				{
					RPr: &xmlARPr{Lang: "en-US", Sz: "1200"},
					T:   slide.notes,
				},
			},
		})
	}

	sp.TxBody = txBody
	notes.CSld.SpTree.Shapes = append(notes.CSld.SpTree.Shapes, sp)

	return xmlutil.MarshalXML(notes)
}

func colorToHex(c common.Color) string {
	return fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B)
}
