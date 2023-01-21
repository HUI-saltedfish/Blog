package models

import "gorm.io/gorm"

// table smmsfile
type SmmsFile struct {
	gorm.Model
	FileName  string `json:"filename"`
	StoreName string `json:"storename"`
	Size      int64  `json:"size"`
	Width     int64  `json:"width"`
	Height    int64  `json:"height"`
	Hash      string `json:"hash"`
	Delete    string `json:"delete"`
	Url       string `json:"url"`
	Path      string `json:"path"`
}

func (sf SmmsFile) Insert() error {
	err := DB.Create(sf).Error
	return err
}
