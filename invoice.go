package main

import (
	"fmt"
	"io"
	"time"

	"github.com/go-pdf/fpdf"
)

type Invoice struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (i Invoice) Map() map[string]any {
	return map[string]any{
		"id":    i.ID,
		"name":  i.Name,
		"email": i.Email,
		"phone": i.Phone,
	}
}

func (i Invoice) Write(w io.Writer, products []Product) error {
	font := "Arial"

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add document header
	pdf.SetFont(font, "B", 16)
	pdf.Cell(140, 10, "Sovereign Fireworks / Invoice")
	pdf.SetFont(font, "", 14)
	pdf.Cell(40, 10, i.Name)
	pdf.Ln(20)

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

	// Add timestamp
	pst, _ := time.LoadLocation("America/Los_Angeles")
	now := time.Now().In(pst)
	pdf.SetFont(font, "I", 10)
	pdf.Cell(97, 10, "Invoice generated at "+now.Format("3:04 PM on January 2, 2006 (PST)"))
	pdf.SetFont(font, "", 12)

	// Add signature
	pdf.Cell(21, 10, "Signature: ")
	pdf.Line(pdf.GetX()+3, pdf.GetY()+7, pdf.GetX()+67, pdf.GetY()+7)
	pdf.Ln(20)

	return pdf.Output(w)
}
