package models

type AdvertisementType struct {
	ID     uint   `gorm:"primaryKey"`
	IsFree bool   `gorm:"not null"`
	Cost   uint   `gorm:"index;not null; default:0"`
	Name   string `gorm:"size:255"`
}

func (AdvertisementType) TableName() string {
	return "AdvertisementType"
}
