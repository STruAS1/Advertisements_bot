package models

import (
	"time"

	"gorm.io/gorm"
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
	Bans         []Ban `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "Users"
}

type Ban struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"index"`
	UnbanDate *time.Time `gorm:"default:null"`
	Reason    string     `gorm:"type:text"`
	CreatedAt time.Time
}

func (Ban) TableName() string {
	return "Bans"
}

func (u *User) IsBanned() (bool, *time.Time) {
	for _, ban := range u.Bans {
		if ban.UnbanDate == nil || time.Now().Before(*ban.UnbanDate) {
			return true, ban.UnbanDate
		}
	}
	return false, nil
}

func BanUser(db *gorm.DB, userID uint, duration time.Duration, reason string) error {
	unbanTime := time.Now().Add(duration)
	ban := Ban{UserID: userID, UnbanDate: &unbanTime, Reason: reason}
	return db.Create(&ban).Error
}

func BanUserForever(db *gorm.DB, userID uint, reason string) error {
	ban := Ban{UserID: userID, UnbanDate: nil, Reason: reason}
	return db.Create(&ban).Error
}

func GetUserBans(db *gorm.DB, userID uint) ([]Ban, error) {
	var bans []Ban
	err := db.Where("user_id = ?", userID).Find(&bans).Error
	return bans, err
}

func UnbanUser(db *gorm.DB, userID uint) error {
	return db.Delete(&Ban{}, "user_id = ?", userID).Error
}
