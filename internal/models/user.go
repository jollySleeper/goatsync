package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"primarykey"`
	Username string    `gorm:"uniqueIndex"`
	Email    string    `gorm:"uniqueIndex"`
	UserInfo *UserInfo `gorm:"foreignKey:OwnerID"`
}

type UserInfo struct {
	ID               uuid.UUID `gorm:"primarykey"`
	OwnerID          uint      `gorm:"uniqueIndex"`
	Salt             []byte
	Version          int
	LoginPubkey      []byte
	Pubkey           []byte
	EncryptedContent []byte
}
