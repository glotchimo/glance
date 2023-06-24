package main

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

type InvoiceRequest struct {
	Invoice Invoice `json:"invoice"`
	Orders  []struct {
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
	} `json:"orders"`
}

type Invoice struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
}

func (i Invoice) Map() map[string]any {
	return map[string]any{
		"name":     i.Name,
		"phone":    i.Phone,
		"location": i.Location,
	}
}

func (i *Invoice) Read(id int) error {
	stmt := goqu.Dialect("postgres").
		Select("id", "name", "phone", "location").
		From("invoices").
		Where(goqu.I("id").Eq(id))
	q, args, err := stmt.ToSQL()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	if err := DB.QueryRow(q, args...).Scan(&i.ID, &i.Name, &i.Phone, &i.Location); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func (i *Invoice) Create() error {
	stmt := goqu.Dialect("postgres").
		Insert("invoices").
		Rows(i.Map()).
		Returning("id")
	q, args, err := stmt.ToSQL()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	if err := DB.QueryRow(q, args...).Scan(&i.ID); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

type Product struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	InvoiceID int    `json:"invoice_id"`
}

func (p Product) Map() map[string]any {
	return map[string]any{
		"name":       p.Name,
		"price":      p.Price,
		"invoice_id": p.InvoiceID,
	}
}

func (p *Product) Read(id int) error {
	stmt := goqu.Dialect("postgres").
		Select("id", "name", "price", "invoice_id").
		From("products").
		Where(goqu.I("id").Eq(id))
	q, args, err := stmt.ToSQL()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	if err := DB.QueryRow(q, args...).Scan(&p.ID, &p.Name, &p.Price, &p.InvoiceID); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func (p *Product) ReadOne(name string) error {
	stmt := goqu.Dialect("postgres").
		Select("id", "name", "price").
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

	if err := DB.QueryRow(q, args...).Scan(&p.ID, &p.Name, &p.Price); err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func (p *Product) Create() error {
	stmt := goqu.Dialect("products").
		Insert("products").
		Rows(p.Map()).
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

func (p *Product) Update() error {
	stmt := goqu.Dialect("products").
		Update("products").
		Set(p.Map()).
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
