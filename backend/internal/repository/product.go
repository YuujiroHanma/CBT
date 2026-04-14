package repository

import (
	"context"
	"fmt"

	"cbt/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ProductRepository handles product database operations.
type ProductRepository struct {
	db *pgxpool.Pool
}

// NewProductRepository creates a new product repository.
func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll fetches all products from the database.
func (r *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT id, name, description, price, image_url, stock_quantity, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageURL,
			&product.StockQuantity,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

// GetByID fetches a single product by ID.
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT id, name, description, price, image_url, stock_quantity, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ImageURL,
		&product.StockQuantity,
		&product.CreatedAt,
		&product.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	return &product, nil
}
