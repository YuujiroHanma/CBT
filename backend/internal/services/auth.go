package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"cbt/backend/internal/models"
	"cbt/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic.
type AuthService struct {
	UserRepo  *repository.UserRepository
	jwtSecret string
}

// NewAuthService creates a new auth service.
func NewAuthService(UserRepo *repository.UserRepository) *AuthService {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Fallback for testing (NEVER do this in production)
		jwtSecret = "your_super_secret_jwt_key_change_in_production"
	}
	return &AuthService{
		UserRepo:  UserRepo,
		jwtSecret: jwtSecret,
	}
}

// HashPassword hashes a password using bcrypt.
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// ComparePassword compares a plain password with a hash.
func (s *AuthService) ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Claims is a custom JWT claims structure.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a user (expires in 24 hours).
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken verifies a JWT token and returns the claims.
func (s *AuthService) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// RegisterResponse is the response returned after successful registration/login.
type RegisterResponse struct {
	User  *models.UserResponse `json:"user"`
	Token string               `json:"token"`
}

// MarshalJSON customizes JSON marshaling to exclude the password.
func (r *RegisterResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"user": map[string]string{
			"id":    r.User.ID,
			"email": r.User.Email,
		},
		"token": r.Token,
	})
}
