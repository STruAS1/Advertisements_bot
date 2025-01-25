package models

import "time"

type VerificationAplication struct {
	ID             uint   `gorm:"primaryKey"`
	UserID         uint   `gorm:"index"`
	User           User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FirstName      string `gorm:"size:255;default:''"`
	LastName       string `gorm:"size:255;default:''"`
	Patronymic     string `gorm:"size:255;default:''"`
	VisaType       string `gorm:"size:255;default:''"`
	CardIdFileID   string `gorm:"size:255;default:''"`
	DocumentFileID string `gorm:"size:255;default:''"`
	Services       string `gorm:"size:255;default:''"`
	Status         uint8
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (VerificationAplication) TableName() string {
	return "VerificationAplication"
}
