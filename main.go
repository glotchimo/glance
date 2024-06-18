package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

var (
	//go:embed web/dist
	dist embed.FS
	//go:embed web/dist/index.html
	index   embed.FS
	distFS  fs.FS
	indexFS fs.FS

	store Store
)

func init() {
	var err error

	distFS, err = fs.Sub(dist, "web/dist")
	if err != nil {
		panic(err)
	}

	indexFS, err = fs.Sub(dist, "web/dist/index.html")
	if err != nil {
		panic(err)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := store.listProducts()
	if err != nil {

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func invoicesHandler(w http.ResponseWriter, r *http.Request) {
	var invoice Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, "invoice.pdf")
}

func main() {
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/invoices", invoicesHandler)
	http.Handle("/", http.FileServerFS(distFS))
	http.Handle("/index.html", http.FileServerFS(indexFS))

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
