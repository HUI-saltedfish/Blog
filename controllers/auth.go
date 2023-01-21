package controllers

import (
	"Blog/helpers"
	"Blog/system"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthGet(c *gin.Context) {
	authType := c.Param("authType")

	session := sessions.Default(c)
	uuid := helpers.UUID()
	session.Delete(SESSION_GITHUB_STATE)
	session.Set(SESSION_GITHUB_STATE, uuid)
	session.Save()

	var authurl string
	switch authType {
	case "github":
		authurl = fmt.Sprintf(system.GetConfiguration().GithubAuthURL, system.GetConfiguration().GithubClientID, uuid)
	case "google":
	case "facebook":
	case "twitter":
	case "weibo":
	case "qq":
	case "wechat":
	default:
		authurl = "/signin"
	}
	// fmt.Println("authurl: ", authurl)
	c.Redirect(http.StatusFound, authurl)

}
