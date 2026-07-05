package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Role         UserRole  `json:"role"`
	PasswordHash *string   `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}
