package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"cbt/backend/internal/repository"

	"github.com/stripe/stripe-go/v78/webhook"
)

// WebhookHandler handles Stripe webhook events.
type WebhookHandler struct {
	orderRepo *repository.OrderRepository
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(orderRepo *repository.OrderRepository) *WebhookHandler {
	return &WebhookHandler{orderRepo: orderRepo}
}

// StripeEvent represents a Stripe webhook event (simplified).
type StripeEvent struct {
	Type string `json:"type"`
	Data struct {
		Object struct {
			ID                string `json:"id"`
			PaymentStatus     string `json:"payment_status"`
			ClientReferenceID string `json:"client_reference_id"`
		} `json:"object"`
	} `json:"data"`
}

// HandleStripeWebhook handles POST /api/webhook/stripe requests.
// This endpoint is public and verifies the Stripe webhook signature.
func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Read raw body (MUST be done before verifying signature)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading webhook body: %v", err)
		http.Error(w, `{"error": "invalid payload"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Get signature from header
	signature := r.Header.Get("Stripe-Signature")
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Println("warning: STRIPE_WEBHOOK_SECRET not set")
	}

	// Verify signature
	event, err := webhook.ConstructEvent(body, signature, webhookSecret)
	if err != nil {
		log.Printf("webhook signature verification failed: %v", err)
		http.Error(w, `{"error": "invalid signature"}`, http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Handle different event types
	switch event.Type {
	case "checkout.session.completed":
		h.handleCheckoutSessionCompleted(w, r, ctx, body)
	default:
		// Acknowledge receipt of event even if we don't handle it
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"received": true}`))
	}
}

// handleCheckoutSessionCompleted processes a completed checkout session event.
func (h *WebhookHandler) handleCheckoutSessionCompleted(
	w http.ResponseWriter,
	r *http.Request,
	ctx interface{},
	bodyBytes []byte,
) {
	// Parse the event
	var event StripeEvent
	if err := json.Unmarshal(bodyBytes, &event); err != nil {
		log.Printf("error parsing webhook event: %v", err)
		http.Error(w, `{"error": "invalid event format"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Extract session ID and payment status
	sessionID := event.Data.Object.ID
	paymentStatus := event.Data.Object.PaymentStatus

	// Update order status if payment was successful
	if paymentStatus == "paid" {
		ctxVal := r.Context()
		if err := h.orderRepo.UpdateOrderStatus(ctxVal, sessionID, "succeeded"); err != nil {
			log.Printf("error updating order status: %v", err)
			http.Error(w, `{"error": "failed to update order"}`, http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			return
		}

		log.Printf("✓ Order %s marked as succeeded", sessionID)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"received": true}`))
}
