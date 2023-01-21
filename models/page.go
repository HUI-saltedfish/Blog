package models

import (
	"strconv"

	"gorm.io/gorm"
)

// table pages
type Page struct {
	gorm.Model
	Title       string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Body        string `gorm:"type:longtext;not null"`
	View        int64  `gorm:"type:bigint;not null;default:0"`
	IsPublished bool   `gorm:"type:tinyint(1);not null;default:0"`
}

func (page *Page) Insert() error {
	return DB.Create(page).Error
}

func (page *Page) Update() error {
	res := DB.Model(page).Updates(map[string]interface{}{
		"Title":       page.Title,
		"Body":        page.Body,
		"IsPublished": page.IsPublished,
	})
	return res.Error
}

func (page *Page) UpdateView() error {
	res := DB.Model(page).Update("View", page.View)
	return res.Error
}

func (page *Page) Delete() error {
	return DB.Delete(page).Error
}

func GetPageByID(id string) (*Page, error) {
	pid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	var page Page
	err = DB.First(&page, pid).Error
	return &page, err
}

func ListPublishedPage() ([]*Page, error) {
	return _listPage(true)
}

func ListAllPage() ([]*Page, error) {
	return _listPage(false)
}

func _listPage(published bool) ([]*Page, error) {
	var pages []*Page
	var err error
	if published {
		err = DB.Where("is_published = ?", true).Order("id desc").Find(&pages).Error
	} else {
		err = DB.Order("id desc").Find(&pages).Error
	}
	return pages, err
}

func CountPage() int64 {
	var count int64
	DB.Model(&Page{}).Count(&count)
	return count
}
