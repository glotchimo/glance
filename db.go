package main

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

const (
	PRODUCTS_TBL = `CREATE TABLE IF NOT EXISTS products (
		name TEXT PRIMARY KEY,
		category TEXT NOT NULL,
		package TEXT NOT NULL,
		price INTEGER NOT NULL,
		retail INTEGER NOT NULL
	);`
)

type Store struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func newStore(dsn string) (*Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if _, err := db.Exec(PRODUCTS_TBL); err != nil {
		return nil, fmt.Errorf("error creating products table: %w", err)
	}

	cache := sq.NewStmtCache(db)
	store := Store{db: db, builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(cache)}

	return &store, nil
}

func (s Store) createProduct(product Product) error {
	q := s.builder.Insert("products").SetMap(product.Map())
	if _, err := q.Exec(); err != nil {
		return fmt.Errorf("error inserting product: %w", err)
	}

	return nil
}

func (s Store) getProduct(name string) (Product, error) {
	var product Product

	q := s.builder.Select("name", "category", "package", "price", "retail").From("products").Where(sq.Eq{"name": name})
	if err := q.QueryRow().Scan(&product.Name, &product.Category, &product.Package, &product.Price, &product.Retail); err != nil {
		return product, fmt.Errorf("error scanning product: %w", err)
	}

	return product, nil
}

func (s Store) listProducts() ([]Product, error) {
	q := s.builder.Select("name", "category", "package", "price", "retail").From("products").OrderBy("name ASC")
	rows, err := q.Query()
	if err != nil {
		return nil, fmt.Errorf("error getting products: %w", err)
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

func (s Store) updateProduct(name string, updates map[string]any) error {
	q := s.builder.Update("products").SetMap(updates).Where(sq.Eq{"name": name})
	if _, err := q.Exec(); err != nil {
		return fmt.Errorf("error updating product: %w", err)
	}

	return nil
}

func (s Store) deleteProduct(name string) error {
	q := s.builder.Delete("products").Where(sq.Eq{"name": name})
	if _, err := q.Exec(); err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	return nil
}
