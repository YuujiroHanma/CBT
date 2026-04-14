package services

import (
	"context"
	"fmt"

	"cbt/backend/internal/models"
	"cbt/backend/internal/repository"
)

// ProductService handles business logic for products.
type ProductService struct {
	repo *repository.ProductRepository
}

// NewProductService creates a new product service.
func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAllProducts retrieves all products.
func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	// If no products found, return empty slice instead of nil (for JSON consistency)
	if products == nil {
		products = []models.Product{}
	}

	return products, nil
}

// GetProductByID retrieves a single product by ID.
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}
