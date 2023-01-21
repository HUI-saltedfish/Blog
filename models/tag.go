package models

import (
	"Blog/system"
	"sort"

	"gorm.io/gorm"
)

// table tags
type Tag struct {
	gorm.Model
	Name  string  `gorm:"type:varchar(100);not null;uniqueIndex"`
	Total int64   `gorm:"type:bigint;not null;default:0"`
	Posts []*Post `gorm:"many2many:posts_tags;"`
}

func (tag *Tag) Insert() error {
	return DB.FirstOrCreate(tag, "name = ?", tag.Name).Error
}

func GetTagByID(id string) (*Tag, error) {
	var tag Tag
	err := DB.Where("id = ?", id).First(&tag).Error
	return &tag, err
}

func ListTag() ([]*Tag, error) {
	var tags []*Tag
	err := DB.Preload("Posts", "is_published = ?", true).Find(&tags).Error
	sort.Slice(tags, func(i, j int) bool {
		return len(tags[i].Posts) > len(tags[j].Posts)
	})
	for _, tag := range tags {
		tag.Total = int64(len(tag.Posts))
	}
	return tags, err
}

func MustListTag() []*Tag {
	tags, err := ListTag()
	if err != nil {
		panic(err)
	}
	maxTags := system.GetConfiguration().MaxShowTags
	if len(tags) > maxTags {
		tags = tags[:maxTags]
	}

	return tags
}

func ListTagByPostID(postID string) ([]*Tag, error) {
	var tags []*Tag
	post, err := GetPostById(postID)
	if err != nil {
		return nil, err
	}
	err = DB.Model(post).Association("Tags").Find(&tags)
	return tags, err

}

func CountTag() int64 {
	var count int64
	DB.Model(&Tag{}).Count(&count)
	return count
}

func ListAllTag() ([]*Tag, error) {
	var tags []*Tag
	err := DB.Model(&Tag{}).Find(&tags).Error
	return tags, err
}
