package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlagSegment struct {
	ID     uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	FlagID uuid.UUID `gorm:"type:char(36)" json:"flag_id"`
	SegID  uuid.UUID `gorm:"type:char(36)" json:"seg_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Segment Segment `gorm:"foreignKey:SegID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
	Flag    Flag    `gorm:"foreignKey:FlagID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
}

func (FlagSegment) TableName() string {
	return "flag_segments"
}
