package controllers

import (
	"Blog/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func PostGet(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil || !post.IsPublished {
		Handle404(c)
		return
	}

	post.View = post.View + 1
	post.UpdateView()
	post.Tags, _ = models.ListTagByPostID(id)
	post.Comments, _ = models.ListCommentByPostID(id)
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "post/display.html", gin.H{
		"post": post,
		"user": user,
	})
}

func PostIndexGet(c *gin.Context) {
	posts, err := models.ListAllPost("")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/post.html", gin.H{
		"posts":    posts,
		"Active":   "post",
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func PostNewGet(c *gin.Context) {
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "post/new.html", gin.H{
		"user": user,
	})
}

func PostNewPost(c *gin.Context) {
	var err error
	tags := c.PostForm("tags")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished") == "on"
	user := c.MustGet(CONTEXT_USER_KEY).(*models.User)
	year := time.Now().Year()
	month := int(time.Now().Month())
	archive := &models.QrArchive{}
	archive.Year = year
	archive.Month = month
	// archive.ArchiveDate only save year and month
	archive.ArchiveDate = time.Now()
	if err != nil {
		seelog.Errorf("parse time in new post error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	archive.FirstOrCreate()

	post := &models.Post{
		UserID:      user.ID,
		Title:       title,
		Body:        body,
		IsPublished: isPublished,
		ArchiveID:   archive.ID,
	}

	if len(tags) > 0 {
		tagArr := strings.Split(tags, ",")
		for _, tagId := range tagArr {
			tag, err := models.GetTagByID(tagId)
			if err != nil {
				continue
			}
			post.Tags = append(post.Tags, tag)
		}
	}

	if err := post.Insert(); err != nil {
		c.HTML(http.StatusOK, "post/new.html", gin.H{
			"message": err.Error(),
			"post":    post,
			"user":    user,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/post")
}

func PostEditGet(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil {
		Handle404(c)
		return
	}
	post.Tags, _ = models.ListTagByPostID(id)
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "post/modify.html", gin.H{
		"post": post,
		"user": user,
	})
}

func PostEditPost(c *gin.Context) {
	id := c.Param("id")
	tags := c.PostForm("tags")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished") == "on"
	user := c.MustGet(CONTEXT_USER_KEY).(*models.User)

	post, err := models.GetPostById(id)
	if err != nil {
		Handle404(c)
		return
	}

	post.Title = title
	post.Body = body
	post.IsPublished = isPublished
	post.UserID = user.ID

	pid, _ := strconv.ParseInt(id, 10, 64)
	err = models.DeleteAllTagsByPostId(pid)
	if err != nil {
		seelog.Errorf("DeleteAllTagsByPostId error: %v", err)
	}

	if len(tags) > 0 {
		tagArr := strings.Split(tags, ",")
		for _, tagId := range tagArr {
			tag, err := models.GetTagByID(tagId)
			if err != nil {
				continue
			}
			post.Tags = append(post.Tags, tag)
		}
	}

	if err := post.Update(); err != nil {
		fmt.Println(err)
		user, _ := c.Get(CONTEXT_USER_KEY)
		c.HTML(http.StatusOK, "post/modify.html", gin.H{
			"message": err.Error(),
			"post":    post,
			"user":    user,
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/post")
}

func PostPublish(c *gin.Context) {
	var (
		err  error
		res  = gin.H{}
		post *models.Post
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	post, err = models.GetPostById(id)
	if err != nil {
		res["message"] = err.Error()
		return
	}

	post.IsPublished = !post.IsPublished
	if err = post.Update(); err != nil {
		res["message"] = err.Error()
		return
	}

	res["succeed"] = true
}

func PostDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	post, err := models.GetPostById(id)
	if err != nil {
		res["message"] = err.Error()
		return
	}

	if err = post.Delete(); err != nil {
		res["message"] = err.Error()
		return
	}

	res["succeed"] = true
}
