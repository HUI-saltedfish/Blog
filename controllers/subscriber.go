package controllers

import (
	"Blog/helpers"
	"Blog/models"
	"Blog/system"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

func SubscribeGet(c *gin.Context) {
	count := models.CountSubscribe()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "other/subscribe.html", gin.H{
		"total": count,
		"user":  user,
	})
}

func Subscribe(c *gin.Context) {
	var err error
	mail := c.PostForm("mail")
	if len(mail) > 0 {
		var subscriber *models.Subscribe
		subscriber, err = models.GetSubscribeByEmail(mail)
		if err == nil {
			if !subscriber.VerifyState && helpers.GetCurrentTime().After(subscriber.OutTime) { // 激活链接过期
				err = sendActiveEmail(subscriber)
				if err == nil {
					count := models.CountSubscribe()
					c.HTML(http.StatusOK, "other/subscribe.html", gin.H{
						"total":   count,
						"message": "激活链接已发送，请查收邮件",
					})
					return
				}
			} else if subscriber.VerifyState && !subscriber.SubscribeState { // 已激活，未订阅
				subscriber.SubscribeState = true
				err = subscriber.Update()
				if err == nil {
					err = errors.New("订阅成功")
				}
			} else {
				err = errors.New("邮件已激活或有未激活的邮件在你的邮箱中，请先激活")
			}
		} else {
			subscriber = &models.Subscribe{
				Email: mail,
			}
			err = subscriber.Insert()
			if err == nil {
				err = sendActiveEmail(subscriber)
				if err == nil {
					count := models.CountSubscribe()
					c.HTML(http.StatusOK, "other/subscribe.html", gin.H{
						"total":   count,
						"message": "激活链接已发送，请查收邮件",
					})
					return
				}
			}

		}
	} else {
		err = errors.New("邮箱不能为空")
	}
	count := models.CountSubscribe()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "other/subscribe.html", gin.H{
		"total":   count,
		"message": err.Error(),
		"user":    user,
	})
}

func sendActiveEmail(subscribe *models.Subscribe) (err error) {
	uuid := helpers.UUID()
	duration, _ := time.ParseDuration("30m")
	subscribe.OutTime = helpers.GetCurrentTime().Add(duration)
	subscribe.SecretKey = uuid
	signature := helpers.Md5(subscribe.Email + uuid + subscribe.OutTime.String())
	subscribe.Signature = signature
	err = sendMail(subscribe.Email, "[Blog]激活", fmt.Sprintf("%s/active?sid=%s", system.GetConfiguration().Domain, signature))
	if err != nil {
		seelog.Errorf("send active email error: %s", err.Error())
		return
	}
	err = subscribe.Update()
	if err != nil {
		seelog.Errorf("update subscribe error: %s", err.Error())
		return
	}
	return
}

func ActiveSubscribe(c *gin.Context) {
	var (
		err        error
		subscriber *models.Subscribe
	)
	sid := c.Query("sid")
	if len(sid) == 0 {
		HandleMessage(c, "激活链接无效, 请重新获取")
		return
	}
	subscriber, err = models.GetSubscribeBySignature(sid)
	if err != nil {
		HandleMessage(c, "激活链接无效, 请重新获取")
		return
	}
	if !helpers.GetCurrentTime().Before(subscriber.OutTime) {
		HandleMessage(c, "激活链接已过期, 请重新获取")
		return
	}
	subscriber.VerifyState = true
	err = subscriber.Update()
	if err != nil {
		HandleMessage(c, "激活失败")
		return
	}
	HandleMessage(c, "激活成功")
}

func UnSubscribe(c *gin.Context) {
	sid := c.Query("sid")
	if len(sid) == 0 {
		HandleMessage(c, "sid 未传入")
		return
	}
	subscriber, err := models.GetSubscribeBySignature(sid)
	if err != nil || !subscriber.VerifyState || !subscriber.SubscribeState {
		HandleMessage(c, "取消订阅失败")
		return
	}
	subscriber.SubscribeState = false
	err = subscriber.Update()
	if err != nil {
		HandleMessage(c, "取消订阅失败")
		return
	}
	HandleMessage(c, "取消订阅成功")
}

func SubscriberIndexGet(c *gin.Context) {
	subscribers, _ := models.ListSubscribe(true)
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/subscriber.html", gin.H{
		"subscribers": subscribers,
		"user":        user,
		"comments":    models.MustListUnreadComment(),
	})
}

func sendEmailToSubscribers(subject, body string) (err error) {
	var (
		subscribers []*models.Subscribe
		emails      = make([]string, 0)
	)
	subscribers, err = models.ListSubscribe(true)
	if err != nil {
		return
	}
	for _, subscriber := range subscribers {
		emails = append(emails, subscriber.Email)
	}
	if len(emails) == 0 {
		err = errors.New("没有订阅者")
		return
	}
	err = sendMail(strings.Join(emails, ";"), subject, body)
	return
}

func SubscriberPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	mail := c.PostForm("mail")
	subject := c.PostForm("subject")
	body := c.PostForm("body")
	if len(mail) > 0 {
		err = sendMail(mail, subject, body)
	} else {
		err = sendEmailToSubscribers(subject, body)
	}

	if err != nil {
		res["message"] = err.Error()
		return
	}
}
