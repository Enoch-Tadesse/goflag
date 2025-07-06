package models

import (
	"time"

	"github.com/google/uuid"
)

type SegmentRule struct {
	ID        uuid.UUID  `json:"id"`
	SegID     uuid.UUID  `json:"seg_id"`
	Attribute string     `json:"attribute"`
	Operator  string     `json:"operator"`
	Value     string     `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
