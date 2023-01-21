package controllers

import (
	"Blog/models"
	"Blog/system"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
)

func ArchivesGet(c *gin.Context) {
	var (
		year      string
		month     string
		page      string
		pageIndex int
		pageSize  int = system.GetConfiguration().PageSize
		total     int64
		posts     []*models.Post
		err       error
		policy    *bluemonday.Policy
	)
	year = c.Param("year")
	month = c.Param("month")
	page = c.DefaultQuery("page", "1")
	pageIndex, err = strconv.Atoi(page)
	if err != nil {
		pageIndex = 1
	}
	posts, err = models.ListPostByArchive(year, month, pageIndex, pageSize)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	total, err = models.CountPostByArchive(year, month)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostID(strconv.FormatUint(uint64(post.ID), 10))
		post.Body = policy.Sanitize(string(post.Body))
	}
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"posts":           posts,
		"tags":            models.MustListTag(),
		"archives":        models.MustListPostArchives(),
		"links":           models.MustListLinks(),
		"pageIndex":       pageIndex,
		"totalPage":       int(math.Ceil(float64(total) / float64(pageSize))),
		"path":            c.Request.URL.Path,
		"maxReadPosts":    models.MustListMaxReadPost(),
		"maxCommentPosts": models.MustListMaxCommentPost(),
		"user":            user,
	})
}
