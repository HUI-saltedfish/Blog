package blogtest

import (
	"Blog/models"
	"testing"

	"gorm.io/gorm"
)

func TestInsertPost(t *testing.T) {
	// init DB
	var err error
	var db *gorm.DB
	db, err = models.InitDB()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("init db success")

	// close database connection
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// insert post
	post := &models.Post{
		UserID:      1,
		Title:       "test title 222",
		Body:        "test content",
		IsPublished: true,
		Tags: []*models.Tag{
			{
				Name: "test tag  111",
			},
		},
		Comments: []*models.Comment{
			{
				UserID:  1,
				PostID:  2,
				Content: "test comment aaa",
			},
		},
	}
	err = db.Create(post).Error
	if err != nil {
		t.Fatal(err)
	}

	t.Log("insert post success")

}

// func TestListPostArchives(t *testing.T) {
// 	// init DB
// 	var err error
// 	var db *gorm.DB
// 	db, err = models.InitDB()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log("init db success")

// 	// close database connection
// 	sqlDB, _ := db.DB()
// 	defer sqlDB.Close()

// 	var archives []*models.QrArchive
// 	db.Raw("select date_format(created_at, '%Y-%m') as date, count(*) as total from posts where is_published = ? group by date_format(created_at, '%Y-%m') order by date desc", true).Debug().Scan(&archives)

// 	for _, archive := range archives {
// 		archive.ArchiveDate, err = time.Parse("2006-01", string(archive.Date))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}

// 	for _, archive := range archives {
// 		t.Logf("archive: %+v", archive)
// 	}

// }
