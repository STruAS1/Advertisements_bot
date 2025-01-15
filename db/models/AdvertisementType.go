package models

type AdvertisementType struct {
	ID     uint   `gorm:"primaryKey"`
	IsFree bool   `gorm:"default:true"`
	Cost   uint   `gorm:"index;not null; default:0"`
	Name   string `gorm:"size:255"`
}

func (AdvertisementType) TableName() string {
	return "AdvertisementType"
}
