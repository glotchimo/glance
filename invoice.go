package main

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

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
