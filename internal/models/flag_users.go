package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlagUser struct {
	ID     uuid.UUID `gorm:"type:char(36)" json:"id"`
	FlagID uuid.UUID `gorm:"type:char(36)" json:"flag_id"`
	UserID uuid.UUID `gorm:"type:char(35)" json:"user_id"`

	CreatedAt time.Time      `gorm:"type:timestamp"`
	UpdatedAt time.Time      `gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User User `gorm:"foreignKey:UserID;references:ID;contraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
	Flag Flag `gorm:"foreignKey:FlagID;references:ID;contraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
}

func (FlagUser) TableName() string {
	return "flag_users"
}
