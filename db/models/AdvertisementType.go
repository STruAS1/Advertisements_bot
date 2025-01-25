package models

type AdvertisementType struct {
	ID       uint   `gorm:"primaryKey"`
	Priority uint   `gorm:"index;not null;default:1"`
	IsFree   bool   `gorm:"not null"`
	Cost     uint   `gorm:"index;not null; default:0"`
	Name     string `gorm:"size:255"`
}

func (AdvertisementType) TableName() string {
	return "AdvertisementType"
}
