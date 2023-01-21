package models

import (
	"time"

	"gorm.io/gorm"
)

// table Subscribe
type Subscribe struct {
	gorm.Model
	Email          string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	VerifyState    bool      `gorm:"type:tinyint(1);not null;default:0"`
	SubscribeState bool      `gorm:"type:tinyint(1);not null;default:1"`
	OutTime        time.Time `gorm:"type:datetime;default:null"`
	SecretKey      string    `gorm:"type:varchar(100);default:null"`
	Signature      string    `gorm:"type:varchar(100);default:null"`
}

func (s *Subscribe) Insert() error {
	return DB.FirstOrCreate(s, "email = ?", s.Email).Error
}

func (s *Subscribe) Update() error {
	return DB.Model(s).Updates(s).Error
}

func CountSubscribe() int {
	var count int64
	DB.Model(&Subscribe{}).Where("subscribe_state = ?", true).Count(&count)
	return int(count)
}

func ListSubscribe(invalid bool) ([]*Subscribe, error) {
	var subs []*Subscribe
	// err := DB.Where("verify_state = ? and subscribe_state = ?", true, invalid).Find(&subs).Error
	err := DB.Where("subscribe_state = ?", invalid).Find(&subs).Error
	return subs, err
}

func GetSubscribeByEmail(email string) (*Subscribe, error) {
	var sub Subscribe
	err := DB.Where("email = ?", email).First(&sub).Error
	return &sub, err
}

func GetSubscribeBySignature(signature string) (*Subscribe, error) {
	var sub Subscribe
	err := DB.Where("signature = ?", signature).First(&sub).Error
	return &sub, err
}

func GetSubscribeById(id uint) (*Subscribe, error) {
	var sub Subscribe
	err := DB.Where("id = ?", id).First(&sub).Error
	return &sub, err
}

