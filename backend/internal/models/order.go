package models

import (
	"time"
)

// Order represents an order in the database.
type Order struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"user_id"`
	StripePaymentIntentID string    `json:"stripe_payment_intent_id"`
	Status                string    `json:"status"` // pending, succeeded, failed
	TotalAmount           float64   `json:"total_amount"`
	Currency              string    `json:"currency"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// OrderItem represents a line item in an order.
type OrderItem struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderWithItems is an order with its line items.
type OrderWithItems struct {
	Order      *Order      `json:"order"`
	OrderItems []OrderItem `json:"items"`
}

// CartItem represents an item in the shopping cart.
type CartItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Name      string  `json:"name"`
}

// CheckoutRequest is the payload for the checkout endpoint.
type CheckoutRequest struct {
	Items      []CartItem `json:"items"`
	SuccessURL string     `json:"success_url"`
	CancelURL  string     `json:"cancel_url"`
}
