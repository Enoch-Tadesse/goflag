package models

import (
	"time"

	"github.com/google/uuid"
)

type FlagUser struct {
	ID        uuid.UUID  `json:"id"`
	FlagID    uuid.UUID  `json:"flag_id"`
	UserID    uuid.UUID  `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
