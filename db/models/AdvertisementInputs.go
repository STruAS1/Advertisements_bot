package models

type AdvertisementInputs struct {
	ID       uint              `gorm:"primaryKey"`
	Priority uint              `gorm:"index;not null;default:1"`
	Name     string            `gorm:"size:255"`
	Options  string            `gorm:"type:text"`
	Optional bool              `gorm:"default:false"`
	InputID  uint              `gorm:"index;not null"`
	TypeID   uint              `gorm:"index;not null"`
	Type     AdvertisementType `gorm:"foreignKey:TypeID;constraint:OnDelete:CASCADE"`
}

func (AdvertisementInputs) TableName() string {
	return "AdvertisementInputs"
}
