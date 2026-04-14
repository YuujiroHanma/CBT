package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"cbt/backend/internal/services"
)

// ProductHandler handles product-related HTTP requests.
type ProductHandler struct {
	service *services.ProductService
}

// NewProductHandler creates a new product handler.
func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetProducts handles GET /api/products requests.
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Fetch all products from the service
	products, err := h.service.GetAllProducts(ctx)
	if err != nil {
		log.Printf("error fetching products: %v", err)
		http.Error(w, `{"error": "failed to fetch products"}`, http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
