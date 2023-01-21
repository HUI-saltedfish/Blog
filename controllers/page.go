package controllers

import (
	"Blog/models"
	"net/http"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func PageGet(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageByID(id)
	if err != nil || !page.IsPublished {
		Handle404(c)
		return
	}

	page.View = page.View + 1
	page.UpdateView()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "page/display.html", gin.H{
		"page": page,
		"user": user,
	})
}

func PageIndexGet(c *gin.Context) {
	pages, err := models.ListAllPage()
	if err != nil {
		seelog.Errorf("ListAllPage error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/page.html", gin.H{
		"pages":    pages,
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func PageNewGet(c *gin.Context) {
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "page/new.html", gin.H{
		"user": user,
	})
}

func PageNewPost(c *gin.Context) {
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished") == "on"
	page := &models.Page{
		Title:       title,
		Body:        body,
		IsPublished: isPublished,
	}
	user, _ := c.Get(CONTEXT_USER_KEY)
	if err := page.Insert(); err != nil {
		c.HTML(http.StatusOK, "page/new.html", gin.H{
			"message": err.Error(),
			"page":    page,
			"user":    user,
		})
	}
	c.Redirect(http.StatusFound, "/admin/page")
}

func PageEditGet(c *gin.Context) {
	id := c.Param("id")
	page, err := models.GetPageByID(id)
	if err != nil {
		Handle404(c)
		return
	}
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "page/modify.html", gin.H{
		"page": page,
		"user": user,
	})
}

func PageEditPost(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	body := c.PostForm("body")
	isPublished := c.PostForm("isPublished") == "on"
	page, err := models.GetPageByID(id)
	if err != nil {
		Handle404(c)
		return
	}
	page.Title = title
	page.Body = body
	page.IsPublished = isPublished
	user, _ := c.Get(CONTEXT_USER_KEY)
	if err := page.Update(); err != nil {
		c.HTML(http.StatusOK, "page/modify.html", gin.H{
			"message": err.Error(),
			"page":    page,
			"user":    user,
		})
	}
	c.Redirect(http.StatusFound, "/admin/page")
}

func PagePublish(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	page, err := models.GetPageByID(id)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	page.IsPublished = !page.IsPublished
	if err := page.Update(); err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func PageDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	page, err := models.GetPageByID(id)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	if err := page.Delete(); err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
