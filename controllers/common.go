package controllers

import (
	"Blog/helpers"
	"Blog/models"
	"Blog/system"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/denisbakhtin/sitemap"
	"github.com/gin-gonic/gin"
)

const (
	SESSION_KEY          = "userID"
	CONTEXT_USER_KEY     = "user"
	SESSION_GITHUB_STATE = "GITHUB_STATE"
	SESSION_CAPTCHA      = "CAPTCHA"
)

func Handle404(c *gin.Context) {
	HandleMessage(c, "Sorry, the page you are looking for could not be found.")
}

func HandleMessage(c *gin.Context, message string) {
	c.HTML(http.StatusNotFound, "errors/error.html", gin.H{
		"message": message,
	})
}

func CreateXMLSitemap() {
	configuration := system.GetConfiguration()
	folder := path.Join(configuration.Pubic, "sitemap")
	os.MkdirAll(folder, os.ModePerm)
	domain := configuration.Domain
	now := helpers.GetCurrentTime()
	items := make([]sitemap.Item, 0)

	items = append(items, sitemap.Item{
		Loc:        domain,
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1.0,
	})

	posts, err := models.ListPublishedPost("", 0, 0)
	if err != nil {
		panic(err)
	}
	for _, post := range posts {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s/post/%d", domain, post.ID),
			LastMod:    post.UpdatedAt,
			Changefreq: "weekly",
			Priority:   0.9,
		})
	}

	pages, err := models.ListPublishedPage()
	if err != nil {
		panic(err)
	}
	for _, page := range pages {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s/page/%d", domain, page.ID),
			LastMod:    page.UpdatedAt,
			Changefreq: "monthly",
			Priority:   0.8,
		})
	}

	err = sitemap.SiteMap(path.Join(folder, "sitemap.xml"), items)
	if err != nil {
		panic(err)
	}

	err = sitemap.SiteMapIndex(folder, "sitemap.xml", domain+"/static/sitemap/")
	if err != nil {
		panic(err)
	}

}

func writeJSON(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] = false
	}
	ctx.JSON(http.StatusOK, h)
}

func NotifyEmail(subject string, body string) error {
	notifyEmailsStr := system.GetConfiguration().NotifyEmails
	if len(notifyEmailsStr) == 0 {
		return nil
	}
	notifyEmails := strings.Split(notifyEmailsStr, ";")
	emails := make([]string, 0)
	for _, email := range notifyEmails {
		if len(email) > 0 {
			emails = append(emails, email)
		}
	}
	if len(emails) == 0 {
		return nil
	}
	return sendMail(strings.Join(emails, ";"), subject, body)
}

func sendMail(to string, subject string, body string) error {
	configuration := system.GetConfiguration()
	return helpers.SendToMail(configuration.SmtpUsername, configuration.SmtpPassword, configuration.SmtpHost, to, subject, body, "html")
}