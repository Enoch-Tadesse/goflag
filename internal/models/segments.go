package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Segment struct {
	ID   uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name string    `gorm:"unique" json:"name"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Segment) TableName() string {
	return "segments"
}
