package main

import (
	"log"
)

func main() {
	invoice := Invoice{ID: 1}
	products := []Product{
		{1, "Goliath", 200, 1},
		{2, "Goliath", 200, 1},
		{3, "Excalibur", 500, 1},
		{4, "Excalibur", 500, 1},
		{5, "Excalibur", 500, 1},
		{6, "Hell's Shells", 300, 1},
		{7, "Hell's Shells", 300, 1},
		{8, "Hell's Shells", 300, 1},
	}

	if err := createInvoicePDF(invoice, products); err != nil {
		log.Fatal(err)
	}
}
