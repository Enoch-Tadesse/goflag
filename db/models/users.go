package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	Country   string     `json:"country"`
	Paid      bool       `json:"paid"`
	SignUp    time.Time  `json:"sign_up"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // pointer to allow null
}
