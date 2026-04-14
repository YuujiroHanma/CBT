package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"cbt/backend/internal/middleware"
	"cbt/backend/internal/models"
	"cbt/backend/internal/repository"
	"cbt/backend/internal/services"
)

// OrderHandler handles order and checkout HTTP requests.
type OrderHandler struct {
	orderRepo   *repository.OrderRepository
	stripeServ  *services.StripeService
	authService *services.AuthService
}

// NewOrderHandler creates a new order handler.
func NewOrderHandler(
	orderRepo *repository.OrderRepository,
	stripeServ *services.StripeService,
	authService *services.AuthService,
) *OrderHandler {
	return &OrderHandler{
		orderRepo:   orderRepo,
		stripeServ:  stripeServ,
		authService: authService,
	}
}

// Checkout handles POST /api/checkout requests.
// This endpoint is protected and requires JWT authentication.
func (h *OrderHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract user ID from context (added by RequireAuth middleware)
	userID, ok := ctx.Value(middleware.UserContextKey).(string)
	if !ok {
		log.Println("error: user ID not found in context")
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Parse request body
	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Validate cart items
	if len(req.Items) == 0 {
		http.Error(w, `{"error": "cart is empty"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range req.Items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// Create Stripe Checkout Session
	stripeSession, err := h.stripeServ.CreateCheckoutSession(userID, req.Items, req.SuccessURL, req.CancelURL)
	if err != nil {
		log.Printf("error creating stripe session: %v", err)
		http.Error(w, `{"error": "failed to create checkout session"}`, http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Create pending order in database
	_, err = h.orderRepo.CreateOrderWithItems(ctx, userID, stripeSession.ID, totalAmount, req.Items)
	if err != nil {
		log.Printf("error creating order: %v", err)
		http.Error(w, `{"error": "failed to create order"}`, http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Return Stripe session URL to frontend
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url": stripeSession.URL,
	})
}
