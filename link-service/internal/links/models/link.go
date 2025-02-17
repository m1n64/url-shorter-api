package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Link struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Slug   string    `gorm:"type:varchar(255);not null"`
	URL    string    `gorm:"type:varchar(255);not null"`
	gorm.Model
}

func (u *Link) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
