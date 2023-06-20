package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func createProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s: %s %s (%dB)", r.RemoteAddr, r.Method, r.Host, r.ContentLength)

	if r.Method != "POST" {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var products []Product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		log.Println("error decoding request body:", err.Error())
		http.Error(w, "error making feed", http.StatusInternalServerError)
		return
	}

	for _, product := range products {
		log.Println(product)
	}
}
