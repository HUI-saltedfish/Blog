package models

import (
	"time"

	"gorm.io/gorm"
)

// table users
type User struct {
	gorm.Model
	Email         string     `gorm:"type:varchar(100);uniqueIndex, default:null"`
	Telephone     string     `gorm:"type:varchar(100);uniqueIndex, default:null"`
	Password      string     `gorm:"type:varchar(100);default:null"`
	VerifyState   bool       `gorm:"type:tinyint(1);not null;default:0"`
	SecretKey     string     `gorm:"type:varchar(100);default:null"`
	OutTime       time.Time  `gorm:"type:datetime;default:null"`
	GithubLoginId string     `gorm:"type:varchar(100);default:null"`
	GithubUrl     string     `gorm:"type:varchar(100);default:null"`
	IsAdmin       bool       `gorm:"type:tinyint(1);not null;default:0"`
	AvatarUrl     string     `gorm:"type:varchar(100);default:/static/libs/AdminLTE/img/user2-160x160.jpg"`
	Nickname      string     `gorm:"type:varchar(100);default:游客"`
	LockState     bool       `gorm:"type:tinyint(1);not null;default:0"`
	Comment       []*Comment `gorm:"foreignKey:UserID"`
	Post          []*Post    `gorm:"foreignKey:UserID"`
}

func (user *User) Insert() error {
	result := DB.Create(user)
	return result.Error
}

func (user *User) Update() error {
	result := DB.Save(user)
	return result.Error
}

func GetUserByUserName(username string) (*User, error) {
	var user User
	result := DB.Where("email = ?", username).First(&user)
	return &user, result.Error
}

func (user *User) FirstOrCreate() error {
	result := DB.FirstOrCreate(user, "github_login_id = ?", user.GithubLoginId)
	return result.Error
}

func IsGithubIdExists(githubLoginId string, id int) (*User, error) {
	var user User
	result := DB.Where("github_login_id = ? and id != ?", githubLoginId, id).First(&user)
	return &user, result.Error
}

func GetUserByID(id uint) (*User, error) {
	var user User
	result := DB.First(&user, id)
	return &user, result.Error
}

func (user *User) UpdateProfile(avatarUrl, Nickname string) error {
	return DB.Model(user).Updates(map[string]interface{}{"avatar_url": avatarUrl, "nickname": Nickname}).Error
}

func (user *User) UpdateEmail(email string) error {
	if len(email) > 0 {
		return DB.Model(user).Updates(map[string]interface{}{"email": email}).Error
	} else {
		return DB.Model(user).Updates(map[string]interface{}{"email": nil}).Error
	}
}

func (user *User) UpdateGithubUserInfo() error {
	var githubLoginId interface{}
	if len(user.GithubLoginId) == 0 {
		githubLoginId = nil
	} else {
		githubLoginId = user.GithubLoginId
	}
	res := DB.Model(user).Updates(map[string]interface{}{
		"github_login_id": githubLoginId,
		"github_url":      user.GithubUrl,
		"avatar_url":      user.AvatarUrl,
	})
	return res.Error
}

func (user *User) Lock() error {
	return DB.Model(user).Updates(map[string]interface{}{
		"lock_state": true,
	}).Error
}

func (user *User) Unlock() error {
	return DB.Model(user).Updates(map[string]interface{}{
		"lock_state": false,
	}).Error
}

func ListUsers() ([]*User, error) {
	var users []*User
	err := DB.Find(&users).Error
	return users, err
}
