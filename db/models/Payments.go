package models

import "time"

type Payments struct {
	ID        uint   `gorm:"primaryKey"`
	Metod     string `gorm:"size:300"`
	Amount    uint   `gorm:"index;not null"`
	UserID    uint   `gorm:"index"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Status    uint8  `gorm:"index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Payments) TableName() string {
	return "Payments"
}
