package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"cbt/backend/internal/database"
	"cbt/backend/internal/handlers"
	"cbt/backend/internal/middleware"
	"cbt/backend/internal/repository"
	"cbt/backend/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Get configuration from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	ctx := context.Background()
	db, err := database.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()
	log.Println("✓ Database connection established")

	// Initialize repositories, services, and handlers
	productRepo := repository.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	orderRepo := repository.NewOrderRepository(db)
	stripeService := services.NewStripeService()
	orderHandler := handlers.NewOrderHandler(orderRepo, stripeService, authService)

	webhookHandler := handlers.NewWebhookHandler(orderRepo)

	// Set up router with chi
	router := chi.NewRouter()

	// Apply CORS middleware globally
	router.Use(middleware.CORSMiddleware)

	// Public routes
	router.Get("/api/products", productHandler.GetProducts)
	router.Post("/api/auth/register", authHandler.Register)
	router.Post("/api/auth/login", authHandler.Login)

	// Stripe webhook (public, must verify signature)
	router.Post("/api/webhook/stripe", webhookHandler.HandleStripeWebhook)

	// Protected routes (require JWT authentication)
	router.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService.VerifyToken))
		r.Post("/api/checkout", orderHandler.Checkout)
	})

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Start server
	log.Printf("🚀 Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
