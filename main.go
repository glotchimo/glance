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

func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := store.createProduct(product); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
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

func handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	name := r.URL.Query().Get("name")

	var updates map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := store.updateProduct(name, updates); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	name := r.URL.Query().Get("name")

	var updates map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := store.deleteProduct(name); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func handleGenerateInvoice(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	type input struct {
		Invoice
		Products []struct {
			Name     string
			Quantity int
		}
	}

	var in input
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	invoice := in.Invoice
	orderedProducts := []OrderedProduct{}
	for _, p := range in.Products {
		product, err := store.getProduct(p.Name)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orderedProducts = append(orderedProducts, OrderedProduct{
			Product:  product,
			Quantity: p.Quantity,
		})
	}

	if err := invoice.Write(w, orderedProducts); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/api/products/create", handleCreateProduct)
	http.HandleFunc("/api/products/list", handleListProducts)
	http.HandleFunc("/api/products/update/:name", handleUpdateProduct)
	http.HandleFunc("/api/products/delete/:name", handleDeleteProduct)
	http.HandleFunc("/api/invoices", handleGenerateInvoice)
	http.Handle("/", http.FileServerFS(distFS))
	http.Handle("/index.html", http.FileServerFS(indexFS))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
