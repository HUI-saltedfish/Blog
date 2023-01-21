package blogtest

import (
	"Blog/models"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestModels(t *testing.T) {
	t.Log("test models")
	var err error
	mysqlConfig := mysql.Config{
		DSN:                       "root:123456@tcp(localhost:3306)/blog?charset=utf8&parseTime=True&loc=Local",
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
	DB, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// close database connection
	sqlDB, _ := DB.DB()

	// delete table
	// DB.Migrator().DropTable(&models.Page{}, &models.Post{}, &models.Tag{}, &models.User{}, &models.Comment{}, &models.Subscribe{}, &models.Link{})

	// create table
	DB.AutoMigrate(&models.Page{}, &models.Post{}, &models.Tag{}, &models.User{}, &models.Comment{}, &models.Subscribe{}, &models.Link{})

	defer sqlDB.Close()

}
