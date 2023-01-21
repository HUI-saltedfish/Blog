package blogtest

import (
	"Blog/models"
	"testing"

	"gorm.io/gorm"
)

func TestCreateUsers(t *testing.T) {
	t.Log("test create users")
	// create user
	user := models.User{
		Email:    "443487999@qq.com",
		Password: "123456",
	}

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

	// create user
	result := db.Create(&user)
	t.Logf("result: %+v", user)
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	t.Log("create user success")
}

func TestGetUserByID(t *testing.T) {
	t.Log("test get user by id")
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

	// get user by id
	user, err := models.GetUserByID(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user: %+v", user)
	t.Log("get user success")

}
