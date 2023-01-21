package controllers

import (
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CaptchaGet(ctx *gin.Context) {
	session := sessions.Default(ctx)
	captchaId := captcha.NewLen(4)
	session.Delete(SESSION_CAPTCHA)
	session.Set(SESSION_CAPTCHA, captchaId)
	session.Save()
	captcha.WriteImage(ctx.Writer, captchaId, 100, 40)
}
