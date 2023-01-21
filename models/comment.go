package models

import (
	"strconv"

	"gorm.io/gorm"
)

// table comments
type Comment struct {
	gorm.Model
	UserID    uint   `gorm:"type:int(11);not null;index"`
	Content   string `gorm:"type:longtext;not null"`
	PostID    uint   `gorm:"type:int(11);not null;index"`
	ReadState bool   `gorm:"type:tinyint(1);not null;default:0"`
	NickName  string `gorm:"type:varchar(100);default:游客"`
	AvatarUrl string `gorm:"type:varchar(255);default:/static/libs/AdminLTE/img/user2-160x160.jpg"`
	GithubUrl string `gorm:"type:varchar(255);default:null"`
}

func (comment *Comment) Insert() error {
	return DB.Create(comment).Error
}

func (comment *Comment) Update() error {
	return DB.Model(comment).UpdateColumn("read_state", true).Error
}

func ListUnreadComment() ([]*Comment, error) {
	var comments []*Comment
	err := DB.Where("read_state = ?", false).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func MustListUnreadComment() []*Comment {
	comments, err := ListUnreadComment()
	if err != nil {
		panic(err)
	}
	return comments
}

func (comment *Comment) Delete() error {
	return DB.Delete(comment, "id = ?", comment.ID).Error
}

func ListCommentByPostID(postID string) ([]*Comment, error) {
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		return nil, err
	}
	var comments []*Comment
	err = DB.Where("post_id = ?", pid).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func CountComment() int {
	var count int64
	DB.Model(&Comment{}).Count(&count)
	return int(count)
}

func SetAllCommentRead() error {
	return DB.Model(&Comment{}).Where("read_state = ?", false).Update("read_state", true).Error
}
