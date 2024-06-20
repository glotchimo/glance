package main

import (
	"fmt"
	"io"
	"time"

	"github.com/go-pdf/fpdf"
)

type Invoice struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func (i Invoice) Map() map[string]any {
	return map[string]any{
		"id":    i.ID,
		"name":  i.Name,
		"email": i.Address,
		"phone": i.Phone,
	}
}

func (i Invoice) Write(w io.Writer, products []Product) error {
	font := "Arial"
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Calculate the center position
	pageWidth, _ := pdf.GetPageSize()
	centerX := pageWidth / 2

	// Add logo at the top center
	logoWidth := 35.0
	logoHeight := 33.0
	pdf.Image("logo.png", centerX-(logoWidth/2), 10, logoWidth, logoHeight, false, "", 0, "")

	// Add name centered below the logo
	pdf.SetXY(0, logoHeight+11)
	pdf.SetFont(font, "B", 16)
	pdf.CellFormat(pageWidth-1, 10, "Damien Tomeo", "", 1, "C", false, 0, "")

	// Add phone number centered below the name
	pdf.SetFont(font, "", 12)
	pdf.CellFormat(pageWidth-21, 5, "(509) 209-6079", "", 1, "C", false, 0, "")
	pdf.Ln(15)

	// Add table header
	pdf.SetFont(font, "B", 13)
	pdf.Cell(90, 10, "Item Name")
	pdf.Cell(40, 10, "Quantity")
	pdf.Cell(40, 10, "Price")
	pdf.Cell(40, 10, "Total")
	pdf.Ln(10)

	// Group products together
	mapped := make(map[string]Product)
	quantities := make(map[string]int)
	for _, product := range products {
		mapped[product.Name] = product
		quantities[product.Name]++
	}

	// Add item rows
	var total int
	pdf.SetFont(font, "", 12)
	for name, product := range mapped {
		pdf.CellFormat(90, 10, product.Name, "", 0, "L", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprint(quantities[name]), "", 0, "L", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprint(product.Price), "", 0, "L", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprint(quantities[name]*product.Price), "", 0, "L", false, 0, "")
		pdf.Ln(8)
		total += quantities[name] * product.Price
	}
	pdf.Ln(2)

	// Add total row
	pdf.SetFont(font, "B", 12)
	pdf.CellFormat(170, 10, "Total: ", "", 0, "R", false, 0, "")
	pdf.SetFont(font, "", 12)
	pdf.CellFormat(40, 10, fmt.Sprint(total), "", 0, "L", false, 0, "")
	pdf.Ln(20)

	// Add client's personal info section
	pdf.SetFont(font, "", 12)
	pdf.CellFormat(pageWidth, 5, i.Name, "", 1, "L", false, 0, "")
	pdf.CellFormat(pageWidth, 5, i.Address, "", 1, "L", false, 0, "")
	pdf.CellFormat(pageWidth, 5, i.Phone, "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Add timestamp
	now := time.Now().In(time.FixedZone("PST", -7*60*60))
	pdf.SetFont(font, "I", 10)
	pdf.Cell(97, 10, "Invoice generated at "+now.Format("3:04 PM on January 2, 2006 PST"))
	pdf.SetFont(font, "", 12)

	// Add signature
	pdf.Cell(21, 10, "Signature: ")
	pdf.Line(pdf.GetX()+3, pdf.GetY()+7, pdf.GetX()+67, pdf.GetY()+7)
	pdf.Ln(20)

	return pdf.Output(w)
}
