package repository

import (
	"context"
	"fmt"

	"cbt/backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRepository handles order database operations.
type OrderRepository struct {
	db *pgxpool.Pool
}

// NewOrderRepository creates a new order repository.
func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// CreateOrderWithItems creates an order and its associated order items in a transaction.
func (r *OrderRepository) CreateOrderWithItems(
	ctx context.Context,
	userID string,
	stripeSessionID string,
	totalAmount float64,
	items []models.CartItem,
) (*models.Order, error) {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create the order
	var order models.Order
	orderQuery := `
		INSERT INTO orders (user_id, stripe_payment_intent_id, status, total_amount, currency)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, stripe_payment_intent_id, status, total_amount, currency, created_at, updated_at
	`

	if err := tx.QueryRow(ctx, orderQuery, userID, stripeSessionID, "pending", totalAmount, "usd").Scan(
		&order.ID,
		&order.UserID,
		&order.StripePaymentIntentID,
		&order.Status,
		&order.TotalAmount,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items
	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4)
	`

	for _, item := range items {
		if _, err := tx.Exec(ctx, itemQuery, order.ID, item.ProductID, item.Quantity, item.Price); err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &order, nil
}

// UpdateOrderStatus updates the status of an order by Stripe Session ID.
func (r *OrderRepository) UpdateOrderStatus(
	ctx context.Context,
	stripeSessionID string,
	status string,
) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE stripe_payment_intent_id = $2
	`

	result, err := r.db.Exec(ctx, query, status, stripeSessionID)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no order found with session ID: %s", stripeSessionID)
	}

	return nil
}

// GetOrderByStripeSessionID fetches an order by Stripe Session ID.
func (r *OrderRepository) GetOrderByStripeSessionID(ctx context.Context, sessionID string) (*models.Order, error) {
	query := `
		SELECT id, user_id, stripe_payment_intent_id, status, total_amount, currency, created_at, updated_at
		FROM orders
		WHERE stripe_payment_intent_id = $1
	`

	var order models.Order
	if err := r.db.QueryRow(ctx, query, sessionID).Scan(
		&order.ID,
		&order.UserID,
		&order.StripePaymentIntentID,
		&order.Status,
		&order.TotalAmount,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	return &order, nil
}

// GetOrderWithItems fetches an order with all its line items.
func (r *OrderRepository) GetOrderWithItems(ctx context.Context, orderID string) (*models.OrderWithItems, error) {
	// Fetch order
	orderQuery := `
		SELECT id, user_id, stripe_payment_intent_id, status, total_amount, currency, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order
	if err := r.db.QueryRow(ctx, orderQuery, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.StripePaymentIntentID,
		&order.Status,
		&order.TotalAmount,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	// Fetch order items
	itemsQuery := `
		SELECT id, order_id, product_id, quantity, unit_price, created_at
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.Query(ctx, itemsQuery, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order items: %w", err)
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		items = append(items, item)
	}

	return &models.OrderWithItems{
		Order:      &order,
		OrderItems: items,
	}, nil
}
