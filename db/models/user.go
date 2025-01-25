package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	TelegramID   int64  `gorm:"uniqueIndex"`
	Balance      uint   `gorm:"index;not null;default:0"`
	Username     string `gorm:"size:100"`
	FirstName    string `gorm:"size:100"`
	LastName     string `gorm:"size:100"`
	Phone        string `gorm:"size:100"`
	City         string `gorm:"size:100"`
	Verification bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (User) TableName() string {
	return "Users"
}
