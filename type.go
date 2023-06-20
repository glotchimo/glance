package main

type Invoice struct {
	ID       int
	Name     string
	Phone    string
	Location string
}

type Product struct {
	ID        int
	Name      string
	Price     int
	InvoiceID int
}
