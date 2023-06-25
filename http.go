package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/doug-martin/goqu/v9"
)

func serve(ch chan error) {
	log.Printf("starting HTTP server on port %s\n", PORT)

	http.HandleFunc("/", SendIndex)
	http.HandleFunc("/products", SendProducts)
	http.HandleFunc("/products/create", TakeProducts)
	http.HandleFunc("/invoices", TakeInvoice)

	ch <- http.ListenAndServe(":"+PORT, nil)
}

func SendIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "GET" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := INDEX_TMPL.Execute(w, nil); err != nil {
		log.Println(err.Error())
		http.Error(w, "error executing template", http.StatusInternalServerError)
	}
}

func SendProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "GET" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Load products from the database
	stmt := goqu.Dialect("postgres").
		Select("name").
		Distinct("name").
		From("products").
		Where(goqu.I("invoice_id").IsNull()).
		Order(goqu.I("name").Asc())
	q, args, err := stmt.ToSQL()
	if err != nil {
		log.Println("error building query:", err.Error())
		http.Error(w, "error building query", http.StatusInternalServerError)
		return
	}

	rows, err := DB.Query(q, args...)
	if err != nil {
		log.Println("error executing query:", err.Error())
		http.Error(w, "error executing query", http.StatusInternalServerError)
		return
	}

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Name); err != nil {
			log.Println("error scanning row:", err.Error())
			http.Error(w, "error scanning row", http.StatusInternalServerError)
			return
		}

		products = append(products, p)
	}

	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Println("error encoding products:", err.Error())
		http.Error(w, "error encoding products", http.StatusInternalServerError)
		return
	}
}

func TakeProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "POST" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var products []Product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		log.Println("error decoding request body:", err.Error())
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	for _, product := range products {
		if err := product.Create(); err != nil {
			log.Println("error storing product:", err.Error())
			http.Error(w, "error storing product", http.StatusInternalServerError)
			return
		}
	}
}

func TakeInvoice(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "POST" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data InvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println("error decoding request body:", err.Error())
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	// Create invoice record
	stmt := goqu.Dialect("postgres").
		Insert("invoices").
		Rows(data.Invoice.Map()).
		Returning("id")
	q, args, err := stmt.ToSQL()
	if err != nil {
		log.Println("error building invoice query:", err.Error())
		http.Error(w, "error building invoice query", http.StatusInternalServerError)
		return
	}

	invoice := data.Invoice
	if err := DB.QueryRow(q, args...).Scan(&invoice.ID); err != nil {
		log.Println("error executing invoice query:", err.Error())
		http.Error(w, "error executing invoice query", http.StatusInternalServerError)
		return
	}

	// Assign products to the invoice
	tx, err := DB.Begin()
	if err != nil {
		tx.Rollback()
		log.Println("error starting products transaction:", err.Error())
		http.Error(w, "error starting products transaction", http.StatusInternalServerError)
		return
	}

	var products []Product
	for _, order := range data.Orders {
		for i := 0; i < order.Quantity; i++ {
			var product Product
			if err := product.ReadOne(order.Name); err != nil {
				tx.Rollback()
				log.Println("error getting product:", err.Error())
				http.Error(w, "error getting product", http.StatusBadRequest)
				return
			}

			product.InvoiceID = invoice.ID
			if err := product.Update(); err != nil {
				tx.Rollback()
				log.Println("error updating product:", err.Error())
				http.Error(w, "error updating product", http.StatusBadRequest)
				return
			}

			products = append(products, product)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Println("error committing products transaction:", err.Error())
		http.Error(w, "error committing products transaction", http.StatusInternalServerError)
		return
	}

	// Generate invoice document from finalized numbers
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="invoice.%d.pdf"`, invoice.ID))
	if err := writeInvoicePDF(w, invoice, products); err != nil {
		log.Println("error generating invoice document:", err.Error())
		http.Error(w, "error generating invoice document", http.StatusInternalServerError)
		return
	}
}
