package models

import (
	"time"
)

// User represents a user in the database.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never marshal the password hash
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserResponse is the user info returned in API responses (no password).
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
