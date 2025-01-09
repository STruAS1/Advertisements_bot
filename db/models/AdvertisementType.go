package models

type AdvertisementType struct {
	ID        uint   `gorm:"primaryKey"`
	IsFree    bool   `gorm:"default:true"`
	Name      string `gorm:"size:255"`
	ShortName string `gorm:"size:50"`
}

func (AdvertisementType) TableName() string {
	return "AdvertisementType"
}
