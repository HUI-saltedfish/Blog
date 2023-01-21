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

func TagGet(c *gin.Context) {
	var (
		tagID     string
		page      string
		pageIndex int
		pageSize  int = system.GetConfiguration().PageSize
		total     int64
		err       error
		policy    *bluemonday.Policy
		posts     []*models.Post
	)
	tagID = c.Param("id")
	page = c.DefaultQuery("page", "1")
	pageIndex, err = strconv.Atoi(page)
	if err != nil {
		pageIndex = 1
	}
	posts, err = models.ListPublishedPost(tagID, pageIndex, pageSize)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	policy = bluemonday.StrictPolicy()
	for _, post := range posts {
		post.Tags, _ = models.ListTagByPostID(strconv.FormatUint(uint64(post.ID), 10))
		post.Body = policy.Sanitize(string(post.Body))
	}
	total, err = models.CountPostByTag(tagID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
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

func TagNewPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	name := c.PostForm("value")
	tag := &models.Tag{
		Name:  name,
		Total: 0,
	}
	err = tag.Insert()
	if err != nil {
		res["message"] = err.Error()
		return
	}

	res["secceed"] = true
	res["data"] = tag
}
