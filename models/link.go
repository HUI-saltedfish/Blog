package models

import (
	"Blog/system"

	"gorm.io/gorm"
)

// table link
type Link struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Url  string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Sort int64  `gorm:"type:bigint;not null;default:0"`
	View int64  `gorm:"type:bigint;not null;default:0"`
}

func (l *Link) Insert() error {
	return DB.FirstOrCreate(l, "url = ?", l.Url).Error
}

func (l *Link) Update() error {
	return DB.Model(l).Updates(l).Error
}

func (l *Link) Delete() error {
	return DB.Delete(l).Error
}

func ListLinks() ([]*Link, error) {
	var links []*Link
	err := DB.Order("sort").Find(&links).Error
	return links, err
}

func MustListLinks() []*Link {
	links, _ := ListLinks()
	maxLinks := system.GetConfiguration().MaxShowLinks
	if len(links) > maxLinks {
		return links[:maxLinks]
	}
	return links
}

func GetLinkById(id uint) (*Link, error) {
	var link Link
	err := DB.Where("id = ?", id).First(&link).Error
	return &link, err
}

func GetLinkByUrl(url string) (*Link, error) {
	var link Link
	err := DB.Where("url = ?", url).First(&link).Error
	return &link, err
}
