package blogtest

import (
	"Blog/models"
	"testing"

	"gorm.io/gorm"
)

func TestInsertComment(t *testing.T) {
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

	comment := models.Comment{
		Content:   "test comment",
		UserID:    1,
		PostID:    1,
		ReadState: false,
	}

	result := db.Create(&comment)
	t.Logf("result: %+v", comment)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Log("create comment success")
}

func TestListMaxCommentPost(t *testing.T) {
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

	var posts []*models.Post
	posts, err = models.ListMaxCommentPost()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("posts: %+v", posts)
}

func TestDeleteComment(t *testing.T){
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

	var comment = models.Comment{
		UserID: 1,
	}
	comment.ID = 1
	result := db.Delete(&comment)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Log("delete comment success")
}