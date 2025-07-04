package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Flag struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"unique" json:"name"`
	IsActive  bool      `gorm:"default:false" json:"is_active"`
	IsRollout bool      `gorm:"default:false" json:"is_rollout"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Flag) TableName() string {
	return "flags"
}
