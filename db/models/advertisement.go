package models

import "time"

type Advertisement struct {
	ID                 uint   `gorm:"primaryKey"`
	UserID             uint   `gorm:"index"`
	User               User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Text               string `gorm:"type:text"`
	ImageID            string `gorm:"size:255;default:''"`
	DeletedFromChannel bool
	MassgeID           int
	Status             uint8
	CostUser           uint
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (Advertisement) TableName() string {
	return "Advertisements"
}
