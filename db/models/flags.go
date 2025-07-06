package models

import (
	"time"

	"github.com/google/uuid"
)

type Flag struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	IsActive  bool       `json:"is_active"`
	IsRollout bool       `json:"is_rollout"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
