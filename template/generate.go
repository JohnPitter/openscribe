package template

import (
	"fmt"

	"github.com/JohnPitter/openscribe/common"
	"github.com/JohnPitter/openscribe/document"
	"github.com/JohnPitter/openscribe/pdf"
	"github.com/JohnPitter/openscribe/presentation"
	"github.com/JohnPitter/openscribe/spreadsheet"
)

// GenerateDocx creates a DOCX document from this template with placeholder content
func (t *Template) GenerateDocx() (*document.Document, error) {
	if t.Format != FormatDOCX {
		return nil, fmt.Errorf("template %q is format %s, not DOCX", t.Name, t.Format)
	}

	doc := document.NewWithTheme(t.Theme)

	switch t.Category {
	case CategoryReport:
		generateReportDocx(doc, t)
	case CategoryInvoice:
		generateInvoiceDocx(doc, t)
	case CategoryResume:
		generateResumeDocx(doc, t)
	case CategoryLetter:
		generateLetterDocx(doc, t)
	case CategoryNewsletter:
		generateNewsletterDocx(doc, t)
	default:
		generateReportDocx(doc, t)
	}

	return doc, nil
}

// GenerateXlsx creates an XLSX workbook from this template
func (t *Template) GenerateXlsx() (*spreadsheet.Workbook, error) {
	if t.Format != FormatXLSX {
		return nil, fmt.Errorf("template %q is format %s, not XLSX", t.Name, t.Format)
	}

	wb := spreadsheet.NewWithTheme(t.Theme)

	switch t.Category {
	case CategoryDashboard:
		generateDashboardXlsx(wb, t)
	default:
		generateDashboardXlsx(wb, t)
	}

	return wb, nil
}

// GeneratePptx creates a PPTX presentation from this template
func (t *Template) GeneratePptx() (*presentation.Presentation, error) {
	if t.Format != FormatPPTX {
		return nil, fmt.Errorf("template %q is format %s, not PPTX", t.Name, t.Format)
	}

	pres := presentation.NewWithTheme(t.Theme)

	switch t.Category {
	case CategoryPitchDeck:
		generatePitchDeckPptx(pres, t)
	default:
		generatePitchDeckPptx(pres, t)
	}

	return pres, nil
}

// GeneratePdf creates a PDF document from this template
func (t *Template) GeneratePdf() (*pdf.Document, error) {
	if t.Format != FormatPDF {
		return nil, fmt.Errorf("template %q is format %s, not PDF", t.Name, t.Format)
	}

	doc := pdf.NewWithTheme(t.Theme)

	switch t.Category {
	case CategoryReport:
		generateReportPdf(doc, t)
	case CategoryInvoice:
		generateInvoicePdf(doc, t)
	case CategoryCertificate:
		generateCertificatePdf(doc, t)
	default:
		generateReportPdf(doc, t)
	}

	return doc, nil
}

// --- DOCX Generators ---

func generateReportDocx(doc *document.Document, t *Template) {
	doc.AddHeading(t.Name, 1)
	doc.AddText(t.Description)
	doc.AddHeading("Executive Summary", 2)

	p := doc.AddParagraph()
	p.SetAlignment(common.TextAlignJustify)
	r := p.AddRun()
	r.SetFont(t.Theme.Typography.BodyFont)
	r.SetText("This report provides a comprehensive overview of key metrics and performance indicators. The data presented covers the most recent reporting period and highlights areas of growth and improvement.")

	doc.AddHeading("Key Metrics", 2)
	tbl := doc.AddTable(4, 3)
	tbl.Cell(0, 0).SetText("Metric")
	tbl.Cell(0, 1).SetText("Current")
	tbl.Cell(0, 2).SetText("Target")
	tbl.Cell(1, 0).SetText("Revenue")
	tbl.Cell(1, 1).SetText("$1,250,000")
	tbl.Cell(1, 2).SetText("$1,500,000")
	tbl.Cell(2, 0).SetText("Active Users")
	tbl.Cell(2, 1).SetText("45,000")
	tbl.Cell(2, 2).SetText("60,000")
	tbl.Cell(3, 0).SetText("Satisfaction")
	tbl.Cell(3, 1).SetText("4.2/5.0")
	tbl.Cell(3, 2).SetText("4.5/5.0")

	for col := 0; col < 3; col++ {
		tbl.Cell(0, col).SetShading(t.Theme.Palette.Primary)
	}

	doc.AddPageBreak()
	doc.AddHeading("Recommendations", 2)
	doc.AddText("Based on the analysis above, we recommend focusing on user acquisition and engagement strategies to meet the established targets.")
}

func generateInvoiceDocx(doc *document.Document, t *Template) {
	doc.AddHeading("INVOICE", 1)

	info := doc.AddParagraph()
	info.AddText("Invoice #: INV-2026-001").SetFont(t.Theme.Typography.BodyFont)

	doc.AddParagraph().AddText("Date: March 13, 2026")
	doc.AddParagraph().AddText("Due Date: April 13, 2026")
	doc.AddParagraph().AddText("")

	doc.AddHeading("Bill To", 2)
	doc.AddText("Client Name")
	doc.AddText("123 Business Street")
	doc.AddText("City, State 12345")

	tbl := doc.AddTable(5, 4)
	tbl.Cell(0, 0).SetText("Item")
	tbl.Cell(0, 1).SetText("Description")
	tbl.Cell(0, 2).SetText("Qty")
	tbl.Cell(0, 3).SetText("Amount")
	tbl.Cell(1, 0).SetText("Service A")
	tbl.Cell(1, 1).SetText("Professional consulting")
	tbl.Cell(1, 2).SetText("10 hrs")
	tbl.Cell(1, 3).SetText("$1,500.00")
	tbl.Cell(2, 0).SetText("Service B")
	tbl.Cell(2, 1).SetText("Development work")
	tbl.Cell(2, 2).SetText("20 hrs")
	tbl.Cell(2, 3).SetText("$4,000.00")
	tbl.Cell(3, 0).SetText("")
	tbl.Cell(3, 1).SetText("")
	tbl.Cell(3, 2).SetText("Subtotal")
	tbl.Cell(3, 3).SetText("$5,500.00")
	tbl.Cell(4, 0).SetText("")
	tbl.Cell(4, 1).SetText("")
	tbl.Cell(4, 2).SetText("TOTAL")
	tbl.Cell(4, 3).SetText("$5,500.00")

	for col := 0; col < 4; col++ {
		tbl.Cell(0, col).SetShading(t.Theme.Palette.Primary)
	}
}

func generateResumeDocx(doc *document.Document, t *Template) {
	doc.AddHeading("John Doe", 1)
	doc.AddText("Senior Software Engineer | john.doe@email.com | (555) 123-4567")

	doc.AddHeading("Professional Summary", 2)
	doc.AddText("Results-driven software engineer with 8+ years of experience building scalable applications. Expertise in Go, TypeScript, and cloud infrastructure.")

	doc.AddHeading("Experience", 2)
	doc.AddHeading("Senior Software Engineer — TechCorp", 3)
	doc.AddText("January 2022 — Present")
	doc.AddText("Led development of microservices architecture serving 1M+ users.")

	doc.AddHeading("Software Engineer — StartupCo", 3)
	doc.AddText("June 2018 — December 2021")
	doc.AddText("Built core platform features and mentored junior developers.")

	doc.AddHeading("Education", 2)
	doc.AddText("B.S. Computer Science — State University, 2018")

	doc.AddHeading("Skills", 2)
	doc.AddText("Go, TypeScript, Python, PostgreSQL, Redis, Docker, Kubernetes, AWS")
}

func generateLetterDocx(doc *document.Document, t *Template) {
	doc.AddText("March 13, 2026")
	doc.AddText("")
	doc.AddText("Recipient Name")
	doc.AddText("Company Name")
	doc.AddText("123 Address Street")
	doc.AddText("City, State 12345")
	doc.AddText("")
	doc.AddText("Dear Recipient,")
	doc.AddText("")

	p := doc.AddParagraph()
	p.SetAlignment(common.TextAlignJustify)
	r := p.AddRun()
	r.SetFont(t.Theme.Typography.BodyFont)
	r.SetText("I am writing to express my interest in the opportunity at your organization. With my background in software development and project management, I believe I would be a valuable addition to your team.")

	doc.AddText("")
	doc.AddText("Thank you for your time and consideration. I look forward to hearing from you.")
	doc.AddText("")
	doc.AddText("Sincerely,")
	doc.AddText("Your Name")
}

func generateNewsletterDocx(doc *document.Document, t *Template) {
	doc.AddHeading("Monthly Newsletter", 1)
	doc.AddText("Volume 1, Issue 3 — March 2026")

	doc.AddHeading("Featured Story", 2)
	doc.AddText("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.")

	doc.AddHeading("Team Updates", 2)
	doc.AddText("Our team has been working hard on several exciting initiatives this quarter.")

	doc.AddHeading("Upcoming Events", 2)
	tbl := doc.AddTable(3, 3)
	tbl.Cell(0, 0).SetText("Date")
	tbl.Cell(0, 1).SetText("Event")
	tbl.Cell(0, 2).SetText("Location")
	tbl.Cell(1, 0).SetText("Mar 20")
	tbl.Cell(1, 1).SetText("Team Workshop")
	tbl.Cell(1, 2).SetText("Main Office")
	tbl.Cell(2, 0).SetText("Apr 5")
	tbl.Cell(2, 1).SetText("Company All-Hands")
	tbl.Cell(2, 2).SetText("Virtual")
}

// --- XLSX Generators ---

func generateDashboardXlsx(wb *spreadsheet.Workbook, t *Template) {
	s := wb.AddSheet("Dashboard")

	headers := []string{"Month", "Revenue", "Costs", "Profit", "Users", "Satisfaction"}
	for i, h := range headers {
		s.Cell(1, i+1).SetString(h)
		s.Cell(1, i+1).SetFont(t.Theme.Typography.HeadingFont.WithSize(11))
		s.Cell(1, i+1).SetBackgroundColor(t.Theme.Palette.Primary)
	}

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
	revenues := []float64{85000, 92000, 88000, 105000, 110000, 120000}
	costs := []float64{60000, 62000, 65000, 68000, 70000, 72000}
	users := []float64{12000, 13500, 14200, 16000, 17500, 19000}
	satisfaction := []float64{4.1, 4.2, 4.0, 4.3, 4.4, 4.5}

	for i, m := range months {
		row := i + 2
		s.SetValue(row, 1, m)
		s.SetValue(row, 2, revenues[i])
		s.SetValue(row, 3, costs[i])
		s.SetValue(row, 4, revenues[i]-costs[i])
		s.SetValue(row, 5, users[i])
		s.SetValue(row, 6, satisfaction[i])
	}

	// Summary sheet
	sum := wb.AddSheet("Summary")
	sum.SetValue(1, 1, "Total Revenue")
	sum.Cell(1, 2).SetFormula("SUM(Dashboard!B2:B7)")
	sum.SetValue(2, 1, "Total Costs")
	sum.Cell(2, 2).SetFormula("SUM(Dashboard!C2:C7)")
	sum.SetValue(3, 1, "Total Profit")
	sum.Cell(3, 2).SetFormula("SUM(Dashboard!D2:D7)")
	sum.SetValue(4, 1, "Avg Satisfaction")
	sum.Cell(4, 2).SetFormula("AVERAGE(Dashboard!F2:F7)")
}

// --- PPTX Generators ---

func generatePitchDeckPptx(pres *presentation.Presentation, t *Template) {
	// Slide 1: Title
	s1 := pres.AddSlide()
	s1.SetBackground(t.Theme.Palette.Primary)
	title := s1.AddTextBox(common.In(1.5), common.In(2), common.In(10), common.In(2.5))
	tp := title.AddParagraph()
	tp.AddRun("Company Name", t.Theme.Typography.HeadingFont.WithSize(44).WithColor(common.White))
	tp.SetAlignment(common.TextAlignCenter)

	sub := s1.AddTextBox(common.In(3), common.In(5), common.In(7), common.In(1))
	sp := sub.AddParagraph()
	sp.AddRun("Investor Presentation — March 2026", t.Theme.Typography.BodyFont.WithColor(common.LightGray))
	sp.SetAlignment(common.TextAlignCenter)

	// Slide 2: Problem
	s2 := pres.AddSlide()
	tb2 := s2.AddTextBox(common.In(1), common.In(0.5), common.In(10), common.In(1.5))
	p2 := tb2.AddParagraph()
	p2.AddRun("The Problem", t.Theme.Typography.HeadingFont.WithSize(32))

	content2 := s2.AddTextBox(common.In(1), common.In(2.5), common.In(10), common.In(4))
	for _, point := range []string{
		"Current solutions are slow and expensive",
		"Teams waste 40% of time on manual processes",
		"No unified platform exists for this workflow",
	} {
		bp := content2.AddParagraph()
		bp.AddRun("• "+point, t.Theme.Typography.BodyFont.WithSize(18))
	}

	// Slide 3: Solution
	s3 := pres.AddSlide()
	tb3 := s3.AddTextBox(common.In(1), common.In(0.5), common.In(10), common.In(1.5))
	p3 := tb3.AddParagraph()
	p3.AddRun("Our Solution", t.Theme.Typography.HeadingFont.WithSize(32))

	s3.AddShape(presentation.ShapeRoundedRectangle,
		common.In(1), common.In(2.5), common.In(5), common.In(4))

	desc := s3.AddTextBox(common.In(7), common.In(2.5), common.In(5), common.In(4))
	dp := desc.AddParagraph()
	dp.AddRun("An all-in-one platform that automates workflows, reduces costs by 60%, and improves team productivity.", t.Theme.Typography.BodyFont.WithSize(16))

	// Slide 4: Traction
	s4 := pres.AddSlide()
	tb4 := s4.AddTextBox(common.In(1), common.In(0.5), common.In(10), common.In(1.5))
	p4 := tb4.AddParagraph()
	p4.AddRun("Traction", t.Theme.Typography.HeadingFont.WithSize(32))

	metrics := []struct{ label, value string }{
		{"Revenue", "$1.2M ARR"},
		{"Users", "45,000+"},
		{"Growth", "+23% MoM"},
	}
	for i, m := range metrics {
		x := float64(1) + float64(i)*4
		box := s4.AddShape(presentation.ShapeRoundedRectangle,
			common.In(x), common.In(2.5), common.In(3.5), common.In(3))
		box.SetFill(t.Theme.Palette.Surface)
		box.SetText(m.value+"\n"+m.label, t.Theme.Typography.HeadingFont.WithSize(20))
	}

	// Slide 5: Thank You
	s5 := pres.AddSlide()
	s5.SetBackground(t.Theme.Palette.Primary)
	thanks := s5.AddTextBox(common.In(2), common.In(2.5), common.In(9), common.In(2))
	thp := thanks.AddParagraph()
	thp.AddRun("Thank You", t.Theme.Typography.HeadingFont.WithSize(48).WithColor(common.White))
	thp.SetAlignment(common.TextAlignCenter)

	contact := s5.AddTextBox(common.In(3), common.In(5), common.In(7), common.In(1))
	cp := contact.AddParagraph()
	cp.AddRun("contact@company.com | www.company.com", t.Theme.Typography.BodyFont.WithColor(common.LightGray))
	cp.SetAlignment(common.TextAlignCenter)
}

// --- PDF Generators ---

func generateReportPdf(doc *pdf.Document, t *Template) {
	p := doc.AddPage()

	// Header bar
	p.AddRectangle(0, 0, 595, 70, t.Theme.Palette.Primary, nil)
	p.AddText(t.Name, 72, 25, t.Theme.Typography.HeadingFont.WithSize(24).WithColor(common.White))

	// Body
	p.AddText(t.Description, 72, 100, t.Theme.Typography.BodyFont)
	p.AddLine(72, 120, 523, 120, t.Theme.Palette.Accent, 1)

	p.AddText("Key Findings", 72, 150, t.Theme.Typography.HeadingFont.WithSize(16))

	tbl := p.AddTable(72, 180, 4, 3)
	tbl.SetCellSize(150, 25)
	tbl.SetHeaderBackground(t.Theme.Palette.Primary)
	tbl.SetFont(t.Theme.Typography.BodyFont.WithSize(10))
	tbl.SetCell(0, 0, "Metric")
	tbl.SetCell(0, 1, "Value")
	tbl.SetCell(0, 2, "Status")
	tbl.SetCell(1, 0, "Performance")
	tbl.SetCell(1, 1, "92%")
	tbl.SetCell(1, 2, "On Track")
	tbl.SetCell(2, 0, "Budget")
	tbl.SetCell(2, 1, "$1.2M")
	tbl.SetCell(2, 2, "Under")
	tbl.SetCell(3, 0, "Timeline")
	tbl.SetCell(3, 1, "Q2 2026")
	tbl.SetCell(3, 2, "On Time")
}

func generateInvoicePdf(doc *pdf.Document, t *Template) {
	p := doc.AddPage()

	p.AddText("INVOICE", 72, 50, t.Theme.Typography.HeadingFont.WithSize(32))
	p.AddText("Invoice #: INV-2026-001", 350, 50, t.Theme.Typography.BodyFont)
	p.AddText("Date: March 13, 2026", 350, 70, t.Theme.Typography.BodyFont)

	p.AddLine(72, 95, 523, 95, t.Theme.Palette.Primary, 2)

	p.AddText("Bill To:", 72, 120, t.Theme.Typography.BodyFont.Bold())
	p.AddText("Client Name", 72, 140, t.Theme.Typography.BodyFont)
	p.AddText("123 Business Street", 72, 155, t.Theme.Typography.BodyFont)

	tbl := p.AddTable(72, 200, 5, 4)
	tbl.SetCellSize(115, 25)
	tbl.SetHeaderBackground(t.Theme.Palette.Primary)
	tbl.SetFont(t.Theme.Typography.BodyFont.WithSize(10))
	tbl.SetCell(0, 0, "Item")
	tbl.SetCell(0, 1, "Description")
	tbl.SetCell(0, 2, "Qty")
	tbl.SetCell(0, 3, "Amount")
	tbl.SetCell(1, 0, "Consulting")
	tbl.SetCell(1, 1, "Strategy session")
	tbl.SetCell(1, 2, "10 hrs")
	tbl.SetCell(1, 3, "$2,500")
	tbl.SetCell(2, 0, "Development")
	tbl.SetCell(2, 1, "Feature build")
	tbl.SetCell(2, 2, "40 hrs")
	tbl.SetCell(2, 3, "$8,000")
	tbl.SetCell(3, 0, "")
	tbl.SetCell(3, 2, "Subtotal")
	tbl.SetCell(3, 3, "$10,500")
	tbl.SetCell(4, 0, "")
	tbl.SetCell(4, 2, "TOTAL")
	tbl.SetCell(4, 3, "$10,500")
}

func generateCertificatePdf(doc *pdf.Document, t *Template) {
	landscapeSize := common.PageSize{Width: common.In(11), Height: common.In(8.5)}
	p := doc.AddPageWithSize(landscapeSize, common.UniformMargins(common.In(0.5)))

	// Border
	stroke := t.Theme.Palette.Primary
	p.AddRectangle(20, 20, 752, 572, common.White, &stroke)

	p.AddText("CERTIFICATE OF ACHIEVEMENT", 180, 80,
		t.Theme.Typography.HeadingFont.WithSize(28).WithColor(t.Theme.Palette.Primary))

	p.AddText("This certifies that", 300, 180,
		t.Theme.Typography.BodyFont.WithSize(14))

	p.AddText("Recipient Name", 250, 230,
		t.Theme.Typography.HeadingFont.WithSize(32))

	p.AddLine(200, 270, 592, 270, t.Theme.Palette.Accent, 1)

	p.AddText("has successfully completed the requirements for", 220, 300,
		t.Theme.Typography.BodyFont.WithSize(14))

	p.AddText("Program Name", 290, 340,
		t.Theme.Typography.HeadingFont.WithSize(22).WithColor(t.Theme.Palette.Primary))

	p.AddText("Date: March 13, 2026", 310, 420,
		t.Theme.Typography.CaptionFont)
}
