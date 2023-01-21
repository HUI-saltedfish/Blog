package models

import (
	"time"

	"gorm.io/gorm"
)

// table query result
type QrArchive struct {
	gorm.Model
	Year        int       `gorm:"type:int(11);not null;index"`
	Month       int       `gorm:"type:int(11);not null;index"`
	Total       int       `gorm:"type:int(11);not null;default:0"`
	ArchiveDate time.Time `gorm:"type:datetime;not null;index"`
	Posts       []*Post   `gorm:"foreignKey:ArchiveID"`
}

func (qrArchive *QrArchive) FirstOrCreate() error {
	result := DB.FirstOrCreate(qrArchive, "year = ? and month = ?", qrArchive.Year, qrArchive.Month)
	return result.Error
}

func (qrArchive *QrArchive) UpdateTotal(total int) error {
	result := DB.Model(&QrArchive{}).Where("id = ?", qrArchive.ID).Update("total", total)
	return result.Error
}
