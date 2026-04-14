package services

import (
	"fmt"
	"os"

	"cbt/backend/internal/models"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
)

// StripeService handles Stripe operations.
type StripeService struct {
	secretKey string
}

// NewStripeService creates a new Stripe service.
func NewStripeService() *StripeService {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	if secretKey == "" {
		secretKey = "sk_test_your_key_here" // Fallback only for testing
	}
	stripe.Key = secretKey
	return &StripeService{secretKey: secretKey}
}

// CreateCheckoutSession creates a Stripe Checkout Session for the cart items.
func (s *StripeService) CreateCheckoutSession(
	userID string,
	items []models.CartItem,
	successURL string,
	cancelURL string,
) (*stripe.CheckoutSession, error) {
	// Build line items from cart items
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, len(items))

	for i, item := range items {
		// Convert price to cents (Stripe expects integer amounts)
		priceInCents := int64(item.Price * 100)

		lineItems[i] = &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("usd"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.Name),
				},
				UnitAmount: stripe.Int64(priceInCents),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		}
	}

	// Create checkout session
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(cancelURL),
		ClientReferenceID:  stripe.String(userID), // Link to our user ID
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess, nil
}
