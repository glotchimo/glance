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

func (i Invoice) Write(products []Product, w io.Writer) error {
	font := "sans"

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add document header
	pdf.SetFont(font, "B", 16)
	pdf.Cell(140, 10, "Invoice")
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
	pdf.SetFont(font, "I", 10)
	pdf.Cell(97, 10, "Invoice generated at "+time.Now().Format("3:04 PM on January 2, 2006"))
	pdf.SetFont(font, "", 12)

	return pdf.Output(w)
}
