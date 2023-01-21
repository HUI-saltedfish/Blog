package models

import (
	"Blog/system"
	"html/template"
	"sort"
	"strconv"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gorm.io/gorm"
)

// table posts
type Post struct {
	gorm.Model
	UserID       uint       `gorm:"type:int(10) unsigned;not null;index"`
	Title        string     `gorm:"type:varchar(100);not null;uniqueIndex"`
	Body         string     `gorm:"type:longtext;not null"`
	View         int64      `gorm:"type:bigint;not null;default:0"`
	IsPublished  bool       `gorm:"type:tinyint(1);not null;default:0"`
	Tags         []*Tag     `gorm:"many2many:posts_tags;"`
	Comments     []*Comment `gorm:"foreignKey:PostID"`
	CommentTotal int64      `gorm:"type:bigint;not null;default:0"`
	ArchiveID    uint       `gorm:"type:int(10) unsigned;not null;index"`
}

func (post *Post) Insert() error {
	return DB.Create(post).Error
}

func (post *Post) Update() error {
	return DB.Save(post).Error
}

func (post *Post) UpdateView() error {
	res := DB.Model(post).Update("View", post.View)
	return res.Error
}

func (post *Post) Delete() error {
	return DB.Delete(post).Error
}

func (post *Post) Excerpt() template.HTML {
	policy := bluemonday.StrictPolicy()                                                 // 严格的策略
	sanitized := policy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Body)))) // 清理
	runes := []rune(sanitized)
	if len(runes) > 300 {
		return template.HTML(string(runes[:300]) + "...")
	}
	excerpt := template.HTML(sanitized + "...")
	return excerpt
}

func ListPublishedPost(tag string, pageIndex, pageSize int) ([]*Post, error) {
	return _listPost(tag, true, pageIndex, pageSize)
}

func ListAllPost(tag string) ([]*Post, error) {
	return _listPost(tag, false, 0, 0)
}

func _listPost(tag string, published bool, pageIndex, pageSize int) ([]*Post, error) {
	var posts []*Post
	var err error
	if len(tag) > 0 {
		tagId, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			return nil, err
		}
		// var rows *sql.Rows
		tagModel := Tag{}
		tagModel.ID = uint(tagId)
		if published {
			if pageIndex > 0 {
				// sqlQuery := fmt.Sprintf("select posts.* from posts inner join posts_tags on posts.id = posts_tags.post_id where posts_tags.tag_id = %d and posts.is_published = 1 order by posts.id desc limit %d offset %d", tagId, pageSize, (pageIndex-1)*pageSize)
				// rows, err = DB.Raw(sqlQuery).Rows()
				err = DB.Preload("Posts", func(db *gorm.DB) *gorm.DB {
					return db.Where("is_published = 1").Order("id desc").Limit(pageSize).Offset((pageIndex - 1) * pageSize)
				}).First(&tagModel).Error
				posts = tagModel.Posts
			} else {
				// sqlQuery := fmt.Sprintf("select posts.* from posts inner join posts_tags on posts.id = posts_tags.post_id where posts_tags.tag_id = %d and posts.is_published = 1 order by posts.id desc", tagId)
				// rows, err = DB.Raw(sqlQuery).Rows()
				err = DB.Preload("Posts", func(db *gorm.DB) *gorm.DB {
					return db.Where("is_published = 1").Order("id desc")
				}).First(&tagModel).Error
				posts = tagModel.Posts
			}
		} else {
			// rows, err = DB.Table("posts").Select("posts.*").Joins("inner join posts_tags on posts.id = posts_tags.post_id").Where("posts_tags.tag_id = ?", tagId).Order("posts.id desc").Scan(&posts).Rows()
			err = DB.Preload("Posts", func(db *gorm.DB) *gorm.DB {
				return db.Order("id desc")
			}).First(&tagModel).Error
			posts = tagModel.Posts
		}
		if err != nil {
			return nil, err
		}
		// defer rows.Close()
		// for rows.Next() {
		// 	var post Post
		// 	DB.ScanRows(rows, &post)
		// 	posts = append(posts, &post)
		// }
	} else {
		if published {
			if pageIndex > 0 {
				err = DB.Where("is_published = ?", true).Order("id desc").Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&posts).Error
			} else {
				err = DB.Where("is_published = ?", true).Order("id desc").Find(&posts).Error
			}
		} else {
			err = DB.Order("id desc").Find(&posts).Error
		}
	}
	return posts, err
}

func MustListMaxReadPost() (posts []*Post) {
	posts, err := ListMaxReadPost()
	if err != nil {
		panic(err)
	}

	maxReadPost := system.GetConfiguration().MaxShowRead
	if len(posts) > maxReadPost {
		posts = posts[:maxReadPost]
	}
	return
}

func ListMaxReadPost() (posts []*Post, err error) {
	err = DB.Where("is_published = ?", true).Order("view desc").Find(&posts).Error
	return
}

func MustListMaxCommentPost() (posts []*Post) {
	posts, err := ListMaxCommentPost()
	if err != nil {
		panic(err)
	}
	maxCommentPost := system.GetConfiguration().MaxShowComments
	if len(posts) > maxCommentPost {
		posts = posts[:maxCommentPost]
	}
	return
}

func (post *Post) UpdateCommentTotal() error {
	return DB.Model(post).Update("comment_total", post.CommentTotal).Error
}

func ListMaxCommentPost() (posts []*Post, err error) {
	// sqlQuery := "select p.*,c.total comment_total from posts p inner join (select post_id,count(*) total from comments  group by post_id) c on p.id = c.post_id where p.is_published = ? order by c.total"
	// err = DB.Raw(sqlQuery, true).Scan(&posts).Error
	err = DB.Preload("Comments").Where("is_published = ?", true).Find(&posts).Error
	if err != nil {
		return nil, err
	}
	sort.Slice(posts, func(i, j int) bool {
		return len(posts[i].Comments) > len(posts[j].Comments)
	})
	for _, post := range posts {
		post.CommentTotal = int64(len(post.Comments))
		err = post.UpdateCommentTotal()
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func CountPostByTag(tag string) (count int64, err error) {
	var tagId int64
	if len(tag) > 0 {
		tagId, err = strconv.ParseInt(tag, 10, 64)
		if err != nil {
			return 0, err
		}
		tagModel := Tag{Model: gorm.Model{ID: uint(tagId)}}
		// err = DB.Table("posts").Joins("inner join posts_tags on posts.id = posts_tags.post_id").Where("posts_tags.tag_id = ? and posts.is_published = ?", tagId, true).Count(&count).Error
		err = DB.Preload("Posts", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_published = ?", true)
		}).First(&tagModel).Error
		count = int64(len(tagModel.Posts))
	} else {
		err = DB.Model(&Post{}).Where("is_published = ?", true).Count(&count).Error
	}
	return
}

func CountPost() (count int64) {
	DB.Model(&Post{}).Count(&count)
	return
}

func GetPostById(id string) (*Post, error) {
	var post Post
	err := DB.Where("id = ?", id).First(&post).Error
	return &post, err
}

func MustListPostArchives() (archives []*QrArchive) {
	archives, err := ListPostArchives()
	if err != nil {
		panic(err)
	}

	maxArchives := system.GetConfiguration().MaxShowArchives
	if len(archives) > maxArchives {
		archives = archives[:maxArchives]
	}

	return
}

func ListPostArchives() ([]*QrArchive, error) {
	var archives []*QrArchive
	// err := DB.Raw("select date_format(created_at, '%Y-%m') as date, count(*) as total from posts where is_published = ? group by date_format(created_at, '%Y-%m') order by date desc", true).Scan(&archives).Error
	// for _, archive := range archives {
	// 	archive.ArchiveDate, err = time.Parse("2006-01", string(archive.Date))
	// 	archive.Year = int64(archive.ArchiveDate.Year())
	// 	archive.Month = int64(archive.ArchiveDate.Month())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// err := DB.Preload("Posts", func(db *gorm.DB) *gorm.DB {
	// 	return db.Where("is_published = ?", true)
	// }).Order("year desc, month desc").Find(&archives).Error
	err := DB.Preload("Posts", "is_published = ?", true).Order("year desc, month desc").Find(&archives).Error
	if err != nil {
		return nil, err
	}
	for _, archive := range archives {
		total := len(archive.Posts)
		err = archive.UpdateTotal(total)
		if err != nil {
			return nil, err
		}
	}
	return archives, err
}

func ListPostByArchive(year, month string, pageIndex, pageSize int) (posts []*Post, err error) {
	if len(month) == 1 {
		month = "0" + month
	}
	year_int, err := strconv.Atoi(year)
	if err != nil {
		return nil, err
	}
	month_int, err := strconv.Atoi(month)
	if err != nil {
		return nil, err
	}
	var QrArchive QrArchive
	QrArchive.Year = year_int
	QrArchive.Month = month_int
	err = (&QrArchive).FirstOrCreate()
	if err != nil {
		return nil, err
	}
	DB.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_published = ?", true)
	}).Where("id = ?", QrArchive.ID).First(&QrArchive)
	posts = QrArchive.Posts
	return
}

func CountPostByArchive(year, month string) (count int64, err error) {
	if len(month) == 1 {
		month = "0" + month
	}
	year_int, err := strconv.Atoi(year)
	if err != nil {
		return 0, err
	}
	month_int, err := strconv.Atoi(month)
	if err != nil {
		return 0, err
	}
	var QrArchive QrArchive
	QrArchive.Year = year_int
	QrArchive.Month = month_int
	err = (&QrArchive).FirstOrCreate()
	if err != nil {
		return 0, err
	}
	DB.Preload("Posts", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_published = ?", true)
	}).Where("id = ?", QrArchive.ID).First(&QrArchive)
	count = int64(len(QrArchive.Posts))
	return
}

func DeleteAllTagsByPostId(postId int64) error {
	post := Post{}
	DB.Preload("Tags").Where("id = ?", postId).First(&post)
	return DB.Model(&post).Association("Tags").Clear()
}
