package main

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"github.com/rs/xid"
)

const (
	PRODUCTS_TBL = `CREATE TABLE IF NOT EXISTS products (
		name TEXT PRIMARY KEY,
		category TEXT NOT NULL,
		package TEXT NOT NULL,
		price INTEGER NOT NULL,
		retail INTEGER NOT NULL
	);`

	INVOICES_TBL = `CREATE TABLE IF NOT EXISTS invoices (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		phone TEXT NOT NULL
	)`

	INVOICE_PRODUCTS_TBL = `CREATE TABLE IF NOT EXISTS invoice_products (
		invoice_id TEXT REFERENCES invoices(id),
		product_id TEXT REFERENCES products(name)
	)`
)

type Store struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewStore(dsn string) (*Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if _, err := db.Exec(PRODUCTS_TBL); err != nil {
		return nil, fmt.Errorf("error creating products table: %w", err)
	}

	if _, err := db.Exec(INVOICES_TBL); err != nil {
		return nil, fmt.Errorf("error creating invoices table: %w", err)
	}

	if _, err := db.Exec(INVOICE_PRODUCTS_TBL); err != nil {
		return nil, fmt.Errorf("error creating invoice_products table: %w", err)
	}

	cache := sq.NewStmtCache(db)
	store := Store{db: db, builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(cache)}

	return &store, nil
}

func (s Store) listProducts() ([]Product, error) {
	q := s.builder.Select("name", "category", "package", "price", "retail")
	rows, err := q.Query()
	if err != nil {
		return nil, fmt.Errorf("error getting products")
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.Name, &product.Category, &product.Package, &product.Price, &product.Retail); err != nil {
			return nil, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error scanning products: %w", err)
	}

	return products, nil
}

type createInvoiceIn struct {
	Invoice  Invoice
	Products []Product
}

func (s Store) createInvoice(in createInvoiceIn) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	in.Invoice.ID = xid.New().String()
	stmt := sq.Insert("invoices").SetMap(in.Invoice.Map())
	q, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("error building invoice insertion query: %w", err)
	}

	if _, err := tx.Exec(q, args...); err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing invoice insertion query: %w", err)
	}

	for _, p := range in.Products {
		stmt = sq.Insert("invoice_products").SetMap(map[string]any{"invoice_id": in.Invoice.ID, "product_id": p.Name})
		q, args, err = stmt.ToSql()
		if err != nil {
			return fmt.Errorf("error building invoice_product insertion query: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
