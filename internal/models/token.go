package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Key       string    `gorm:"uniqueIndex"`
	UserID    uuid.UUID
	User      User `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	Expiry    time.Time
}
