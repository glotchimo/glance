package main

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

type Product struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Price     int    `json:"price"`
	InvoiceID int    `json:"invoice_id"`
}

func (p Product) Map() map[string]any {
	return map[string]any{
		"name":       p.Name,
		"category":   p.Category,
		"price":      p.Price,
		"invoice_id": p.InvoiceID,
	}
}

func (p *Product) Create() error {
	stmt := goqu.Dialect("products").
		Insert("products").
		Rows(goqu.Record{"name": p.Name, "category": p.Category, "price": p.Price}).
		Returning("id")
	q, args, err := stmt.ToSQL()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	if err := DB.QueryRow(q, args...).Scan(&p.ID); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func (p *Product) Read(id int) error {
	stmt := goqu.Dialect("postgres").
		Select("id", "name", "category", "price", "invoice_id").
		From("products").
		Where(goqu.I("id").Eq(id))
	q, args, err := stmt.ToSQL()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	if err := DB.QueryRow(q, args...).Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.InvoiceID); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func (p *Product) ReadOne(name string) error {
	stmt := goqu.Dialect("postgres").
		Select("id", "name", "category", "price").
		From("products").
		Where(goqu.And(
			goqu.I("name").Eq(name)),
			goqu.I("invoice_id").IsNull()).
		Order(goqu.I("id").Asc()).
		Limit(1)
	q, args, err := stmt.ToSQL()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	if err := DB.QueryRow(q, args...).Scan(&p.ID, &p.Name, &p.Category, &p.Price); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func (p *Product) Update() error {
	stmt := goqu.Dialect("products").
		Update("products").
		Set(goqu.Record{"invoice_id": p.InvoiceID}).
		Where(goqu.I("id").Eq(p.ID))
	q, args, err := stmt.ToSQL()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	if _, err := DB.Exec(q, args...); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}