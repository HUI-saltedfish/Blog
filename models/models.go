package models

import (
	"Blog/system"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// // table posts_tags
// type PostTag struct {
// 	gorm.Model
// 	PostID uint `gorm:"type:int(11);not null;index"`
// 	TagID  uint `gorm:"type:int(11);not null;index"`
// }

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	var err error
	config := system.GetConfiguration()
	mysqlConfig := mysql.Config{
		DSN:                       config.DSN,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
	DB, err = gorm.Open(mysql.New(mysqlConfig), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// delete table
	// DB.Migrator().DropTable(&Page{}, &Post{}, &Tag{}, &User{}, &Comment{}, &Subscribe{}, &Link{}, &QrArchive{})
	// create table
	DB.AutoMigrate(&Page{}, &Post{}, &Tag{}, &User{}, &Comment{}, &Subscribe{}, &Link{}, &QrArchive{})
	return DB, nil

}
