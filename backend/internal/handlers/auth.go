package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"cbt/backend/internal/models"
	"cbt/backend/internal/services"

	"github.com/jackc/pgx/v5"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRequest is the payload for the register endpoint.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the payload for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles POST /api/auth/register requests.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "email and password are required"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, `{"error": "password must be at least 6 characters"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Hash password
	passwordHash, err := h.authService.HashPassword(req.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		http.Error(w, `{"error": "failed to register user"}`, http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Create user in database
	user, err := h.authService.UserRepo.CreateUser(ctx, req.Email, passwordHash)
	if err != nil {
		log.Printf("error creating user: %v", err)
		// Check if it's a duplicate email error
		if errors.Is(err, pgx.ErrNoRows) || err.Error() == "failed to create user: ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" {
			http.Error(w, `{"error": "email already registered"}`, http.StatusConflict)
		} else {
			http.Error(w, `{"error": "failed to register user"}`, http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		log.Printf("error generating token: %v", err)
		http.Error(w, `{"error": "failed to generate token"}`, http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(services.RegisterResponse{
		User: &models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}

// Login handles POST /api/auth/login requests.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "email and password are required"}`, http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Fetch user from database
	user, err := h.authService.UserRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("error fetching user: %v", err)
		http.Error(w, `{"error": "invalid email or password"}`, http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Compare password
	if !h.authService.ComparePassword(user.Password, req.Password) {
		http.Error(w, `{"error": "invalid email or password"}`, http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		log.Printf("error generating token: %v", err)
		http.Error(w, `{"error": "failed to generate token"}`, http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services.RegisterResponse{
		User: &models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}
