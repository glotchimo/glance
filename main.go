package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var (
	//go:embed web/dist
	dist embed.FS
	//go:embed web/dist/index.html
	index   embed.FS
	distFS  fs.FS
	indexFS fs.FS

	store *Store
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

	store, err = newStore(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
}

func handleListProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	products, err := store.listProducts()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func handleCreateInvoice(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	var in createInvoiceIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	invoice := in.Invoice
	products := []Product{}
	for _, p := range in.Products {
		product, err := store.getProduct(p.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for range p.Quantity {
			products = append(products, product)
		}
	}

	if err := invoice.Write(w, products); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/products", handleListProducts)
	http.HandleFunc("/invoices", handleCreateInvoice)
	http.Handle("/", http.FileServerFS(distFS))
	http.Handle("/index.html", http.FileServerFS(indexFS))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
