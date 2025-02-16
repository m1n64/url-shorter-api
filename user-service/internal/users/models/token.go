package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	gorm.Model
}

func (u *Token) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
