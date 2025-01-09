package models

import "time"

type Cities struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Cities) TableName() string {
	return "Cities"
}
