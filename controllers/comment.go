package controllers

import (
	"Blog/models"
	"Blog/system"
	"fmt"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CommentPost(c *gin.Context) {
	var (
		err  error
		res  = gin.H{}
		post *models.Post
	)
	defer writeJSON(c, res)
	s := sessions.Default(c)
	sessionUserID := s.Get(SESSION_KEY)
	if sessionUserID == nil {
		res["message"] = "请先登录"
		return
	}
	userID := sessionUserID.(uint)
	user, err := models.GetUserByID(userID)
	if err != nil {
		res["message"] = "用户不存在"
		return
	}

	verifyCode := c.PostForm("verifyCode")
	captchaId := s.Get(SESSION_CAPTCHA)
	s.Delete(SESSION_CAPTCHA)
	_captchaId, _ := captchaId.(string)
	if !captcha.VerifyString(_captchaId, verifyCode) {
		res["message"] = "验证码错误"
		return
	}

	postID := c.PostForm("postId")
	content := c.PostForm("content")
	if len(content) == 0 {
		res["message"] = "评论内容不能为空"
		return
	}

	post, err = models.GetPostById(postID)
	if err != nil {
		res["message"] = "文章不存在"
		return
	}

	pid, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		res["message"] = "postID错误"
		return
	}

	comment := models.Comment{
		UserID:    userID,
		Content:   content,
		PostID:    uint(pid),
		NickName:  user.Email,
		AvatarUrl: user.AvatarUrl,
		GithubUrl: user.GithubUrl,
	}
	err = comment.Insert()
	if err != nil {
		res["message"] = "评论失败"
		return
	}
	NotifyEmail("[Blog]您有一条新的评论", fmt.Sprintf("<a href='%s/post/%d' target='_blank'>%s</a>: %s", system.GetConfiguration().Domain, post.ID, post.Title, content))
	res["succeed"] = true
}

func CommentDelete(c *gin.Context) {
	var (
		err error
		res = gin.H{}
		cid uint64
	)
	defer writeJSON(c, res)
	s := sessions.Default(c)
	sessionUserID := s.Get(SESSION_KEY)
	if sessionUserID == nil {
		res["message"] = "请先登录"
		return
	}
	userID := sessionUserID.(uint)

	commentID := c.Param("id")
	cid, err = strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		res["message"] = "commentID错误"
		return
	}

	comment := models.Comment{
		UserID: userID,
	}
	comment.ID = uint(cid)
	err = comment.Delete()
	if err != nil {
		res["message"] = "删除失败"
		return
	}
	res["succeed"] = true

}

func CommentReadPost(c *gin.Context) {
	var (
		id  string
		_id uint64
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	id = c.Param("id")
	_id, err = strconv.ParseUint(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	comment := &models.Comment{
		Model: gorm.Model{ID: uint(_id)},
	}
	err = comment.Update()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func CommentReadAllPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	err = models.SetAllCommentRead()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
