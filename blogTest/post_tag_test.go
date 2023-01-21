package blogtest

import (
	"Blog/models"
	"testing"

	"gorm.io/gorm"
)

func TestCreatePostTag(t *testing.T) {
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

	// create post tag
	post := models.Post{
		Title: "test post",
		Body:  "test content",

		Tags: []*models.Tag{
			{
				Name: "test tag",
			},
		},

		UserID: 1,

		Comments: []*models.Comment{
			{
				Content: "test comment",
				UserID:  1,
			},
		},
	}

	result := db.Create(&post)
	t.Logf("result: %+v", post)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Log("create post tag success")

}

func TestListPost(t *testing.T) {
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

	// list post
	var posts []*models.Post
	err = db.Debug().Where("is_published = ?", true).Order("id desc").Find(&posts).Error
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("posts: %+v", posts)
	t.Log("list post success")

}


func TestDeleteAllTagsByPostId(t *testing.T) {
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

	// delete all tags by post id
	err = models.DeleteAllTagsByPostId(1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("delete all tags by post id success")

}