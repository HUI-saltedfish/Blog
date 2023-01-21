package blogtest

import (
	"Blog/models"
	"testing"

	"gorm.io/gorm"
)

func TestInsertLink(t *testing.T) {
	// init DB
	var db *gorm.DB
	db, err := models.InitDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("init db success")

	// close database connection
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	link := []*models.Link{
		{
			Name: "test link",
			Url:  "http://test.com",
			Sort: 1,
			View: 1,
		},

		{
			Name: "test link2",
			Url:  "http://test2.com",
			Sort: 2,
			View: 2,
		},
	}

	result := db.Create(&link)
	t.Logf("result: %+v", link)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Log("create link success")

}
