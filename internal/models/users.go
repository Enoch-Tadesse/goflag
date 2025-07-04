package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID      uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Email   string    `gorm:"type:varchar(50);unique" json:"email"`
	Country string    `gorm:"size:50" json:"country"`
	Paid    bool      `gorm:"type:bool" json:"paid"`
	SignUp  time.Time `gorm:"type:timestamp" json:"sign_up"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}
