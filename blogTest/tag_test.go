package blogtest

import (
	"Blog/models"
	"testing"

	"gorm.io/gorm"
)

func TestTagList(t *testing.T) {
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

	var tags []*models.Tag
	tags, err = models.ListTag()
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range tags {
		t.Logf("tag: %+v", tag)
	}
}
