package controllers

import (
	"Blog/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LinkGet(c *gin.Context) {
	id := c.Param("id")
	_id, _ := strconv.ParseInt(id, 10, 64)
	link, err := models.GetLinkById(uint(_id))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Printf("%+v\n", link)
	link.View += 1
	link.Update()
	c.Redirect(http.StatusFound, link.Url)
}

func LinkIndexGet(c *gin.Context) {
	links, err := models.ListLinks()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/link.html", gin.H{
		"links":    links,
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func LinkNewPost(c *gin.Context) {
	var (
		err   error
		res   = gin.H{}
		_sort int64
	)
	defer writeJSON(c, res)
	name := c.PostForm("name")
	url := c.PostForm("url")
	sort := c.PostForm("sort")
	_sort, err = strconv.ParseInt(sort, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	link := &models.Link{
		Name: name,
		Url:  url,
		Sort: _sort,
	}
	err = link.Insert()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func LinkEditPost(c *gin.Context) {
	var (
		_id   uint64
		err   error
		res   = gin.H{}
		_sort int64
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	name := c.PostForm("name")
	url := c.PostForm("url")
	sort := c.PostForm("sort")
	if len(id) == 0 || len(name) == 0 || len(url) == 0 {
		res["message"] = "参数错误"
		return
	}
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	_sort, err = strconv.ParseInt(sort, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	link := &models.Link{
		Model: gorm.Model{ID: uint(_id)},
		Name:  name,
		Url:   url,
		Sort:  _sort,
	}
	err = link.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func LinkDelete(c *gin.Context) {
	var (
		err error
		_id uint64
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	if len(id) == 0 {
		res["message"] = "参数错误"
		return
	}
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	link := &models.Link{
		Model: gorm.Model{ID: uint(_id)},
	}
	err = link.Delete()
	if err != nil {
		res["message"] = err.Error()
		return
	}

	res["succeed"] = true
}
