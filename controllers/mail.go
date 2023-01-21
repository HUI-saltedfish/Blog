package controllers

import (
	"Blog/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SendMail(c *gin.Context) {
	var (
		err        error
		res        = gin.H{}
		uid        uint64
		subscriber *models.Subscribe
	)
	defer writeJSON(c, res)
	subject := c.PostForm("subject")
	content := c.PostForm("content")
	userId := c.Query("userId")

	if len(subject) == 0 || len(content) == 0 || len(userId) == 0 {
		res["message"] = "error parameter"
		return
	}

	uid, err = strconv.ParseUint(userId, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	subscriber, err = models.GetSubscribeById(uint(uid))
	if err != nil {
		res["message"] = err.Error()
		return
	}
	err = sendMail(subscriber.Email, subject, content)
	if err != nil {
		res["message"] = err.Error()
		return
	}

	res["succeed"] = true
}

func SendBatchMail(c *gin.Context) {
	var (
		err         error
		res         = gin.H{}
		subscribers []*models.Subscribe
		emails      = make([]string, 0)
	)
	defer writeJSON(c, res)
	subject := c.PostForm("subject")
	content := c.PostForm("content")
	if len(subject) == 0 || len(content) == 0 {
		res["message"] = "error parameter"
		return
	}
	subscribers, err = models.ListSubscribe(true)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	for _, subscriber := range subscribers {
		emails = append(emails, subscriber.Email)
	}

	err = sendMail(strings.Join(emails, ";"), subject, content)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
