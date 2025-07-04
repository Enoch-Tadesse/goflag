package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SegmentRule struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	SegID     uuid.UUID `gorm:"type:char(36)" json:"seg_id"`
	Attribute string    `gorm:"not null" json:"attribute"`
	Operator  string    `gorm:"not null" json:"operator"`
	Value     string    `gorm:"not null" json:"value"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Segment Segment `gorm:"foreignKey:SegID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
}

func (SegmentRule) TableName() string {
	return "segment_rules"
}
