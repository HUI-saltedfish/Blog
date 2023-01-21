package main

import (
	"Blog/controllers"
	"Blog/helpers"
	"Blog/models"
	"Blog/system"
	"flag"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/cihub/seelog"
	"github.com/claudiu/gocron"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	// DB is the global database connection
	DB *gorm.DB
)

func init() {
	configFilePath := flag.String("config", "conf/conf.yaml", "config file path")
	logConfigPath := flag.String("log", "conf/seelog.xml", "log config file path")
	flag.Parse()

	logger, err := seelog.LoggerFromConfigAsFile(*logConfigPath)
	if err != nil {
		panic(err)
	}
	seelog.ReplaceLogger(logger)

	if err := system.LoadConfiguration(*configFilePath); err != nil {
		seelog.Critical(err)
		panic(err)
	}

	db, err := models.InitDB()
	if err != nil {
		seelog.Critical(err)
		panic(err)
	}
	DB = db
	seelog.Info("database connected")

}

func main() {
	// 刷新缓冲区
	defer seelog.Flush()
	// 关闭数据库连接
	sqlDB, _ := DB.DB()
	defer sqlDB.Close()

	// 初始化路由
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// 加载模板函数
	setTemplateFunc(router)
	// 设置Session
	setSession(router)
	router.Use(ShareData())

	// Periodic tasks
	gocron.Every(1).Day().At("00:00").Do(controllers.CreateXMLSitemap)
	// mysql backup
	gocron.Every(7).Days().At("00:00").Do(controllers.BackupMysql)
	// start all the pending jobs
	gocron.Start()

	// 加载静态文件
	router.Static("/static", filepath.Join("", "./static"))

	// 加载路由
	router.NoRoute(controllers.Handle404)
	router.GET("/", controllers.IndexGet)
	router.GET("/index", controllers.IndexGet)
	router.GET("/rss", controllers.RssGet)

	if system.GetConfiguration().SignupEnabled {
		router.GET("/signup", controllers.SignupGet)
		router.POST("/signup", controllers.SignupPost)
	}

	// user signin and logout
	router.GET("/signin", controllers.SigninGet)
	router.POST("/signin", controllers.SigninPost)
	router.GET("/logout", controllers.LogoutGet)
	router.GET("/oauth2callback", controllers.Oauth2Callback)
	router.GET("/auth/:authType", controllers.AuthGet)

	// captcha
	router.GET("/captcha", controllers.CaptchaGet)

	// visitor
	visitor := router.Group("/visitor")
	visitor.Use(AuthRequired())
	{
		visitor.POST("/new_comment", controllers.CommentPost)
		visitor.POST("/comment/:id/delete", controllers.CommentDelete)
	}

	// subscriber
	router.GET("/subscribe", controllers.SubscribeGet)
	router.POST("/subscribe", controllers.Subscribe)
	router.GET("/active", controllers.ActiveSubscribe)
	router.GET("/unsubscribe", controllers.UnSubscribe)

	router.GET("/page/:id", controllers.PageGet)
	router.GET("/post/:id", controllers.PostGet)
	router.GET("/tag/:id", controllers.TagGet)
	router.GET("/archives/:year/:month", controllers.ArchivesGet)

	router.GET("/link/:id", controllers.LinkGet)

	// admin
	authorized := router.Group("/admin")
	authorized.Use(AdminScopeRequired())
	{
		// index
		authorized.GET("/index", controllers.AdminIndexGet)

		// image upload
		authorized.POST("/upload", controllers.UploadGet)

		// page
		authorized.GET("/page", controllers.PageIndexGet)
		authorized.GET("/new_page", controllers.PageNewGet)
		authorized.POST("/new_page", controllers.PageNewPost)
		authorized.GET("/page/:id/edit", controllers.PageEditGet)
		authorized.POST("/page/:id/edit", controllers.PageEditPost)
		authorized.POST("/page/:id/publish", controllers.PagePublish)
		authorized.POST("/page/:id/delete", controllers.PageDelete)

		// post
		authorized.GET("/post", controllers.PostIndexGet)
		authorized.GET("/new_post", controllers.PostNewGet)
		authorized.POST("/new_post", controllers.PostNewPost)
		authorized.GET("/post/:id/edit", controllers.PostEditGet)
		authorized.POST("/post/:id/edit", controllers.PostEditPost)
		authorized.POST("/post/:id/publish", controllers.PostPublish)
		authorized.POST("/post/:id/delete", controllers.PostDelete)

		// tag
		authorized.POST("/new_tag", controllers.TagNewPost)

		// user page
		authorized.GET("/user", controllers.UserIndexGet)
		authorized.POST("/user/:id/lock", controllers.UserLockPost)

		// profile
		authorized.GET("/profile", controllers.ProfileGet)
		authorized.POST("/profile", controllers.ProfilePost)
		authorized.POST("/profile/email/bind", controllers.ProfileEmailBind)
		authorized.POST("/profile/email/unbind", controllers.ProfileEmailUnbind)   // 解绑邮箱 未用到
		authorized.POST("/profile/github/unbind", controllers.ProfileGithubUnbind) // 解绑github 未用到

		// subscriber
		authorized.GET("/subscriber", controllers.SubscriberIndexGet)
		authorized.POST("/subscriber", controllers.SubscriberPost) // 添加订阅者 未用到

		//link
		authorized.GET("/link", controllers.LinkIndexGet)
		authorized.POST("/new_link", controllers.LinkNewPost)
		authorized.POST("/link/:id/edit", controllers.LinkEditPost)
		authorized.POST("/link/:id/delete", controllers.LinkDelete)

		// comment
		authorized.POST("comment/:id", controllers.CommentReadPost)
		authorized.POST("/read_all", controllers.CommentReadAllPost)

		// backup
		authorized.POST("/backup", controllers.BackupPost)   // 未用到
		authorized.POST("/restore", controllers.RestorePost) // 未用到

		// mail
		authorized.POST("/new_mail", controllers.SendMail)
		authorized.POST("/new_batchmail", controllers.SendBatchMail)
	}

	// 启动路由
	router.Run(system.GetConfiguration().Addr)

}

func setTemplateFunc(router *gin.Engine) {
	funcMap := template.FuncMap{
		"dateFormat": helpers.DateFormat,
		"substring":  helpers.Substring,
		"isOdd":      helpers.IsOdd,
		"isEven":     helpers.IsEven,
		"truncate":   helpers.Truncate,
		"add":        helpers.Add,
		"minus":      helpers.Minus,
		"listTag":    helpers.ListTag,
	}
	router.SetFuncMap(funcMap)
	router.LoadHTMLGlob("views/**/*")
}

func setSession(router *gin.Engine) {
	config := system.GetConfiguration()
	store := cookie.NewStore([]byte(config.SessSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	})
	router.Use(sessions.Sessions("blog", store))
}

func ShareData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(controllers.SESSION_KEY)
		if userID != nil {
			user, err := models.GetUserByID(userID.(uint))
			if err != nil {
				c.AbortWithError(500, err)
			}
			c.Set(controllers.CONTEXT_USER_KEY, user)
		}
		if system.GetConfiguration().SignupEnabled {
			c.Set("signupEnabled", true)
		}
		c.Next()
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get(controllers.CONTEXT_USER_KEY)
		if _, ok := user.(*models.User); ok {
			c.Next()
			return
		}
		seelog.Warnf("You have not logged in yet, So you can't access this page %v", c.Request.URL)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "You have not logged in yet, So you can't access this page",
		})
		c.Abort()
	}
}

// AdminScopeRequired check if the user has admin scope
func AdminScopeRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get(controllers.CONTEXT_USER_KEY)
		if u, ok := user.(*models.User); ok && u.IsAdmin {
			c.Next()
			return
		}
		seelog.Warnf("You have no permission to access this page %v", c.Request.URL)
		c.HTML(http.StatusForbidden, "errors/error.html", gin.H{
			"message": "You have no permission to access this page",
		})
		c.Abort()
	}
}
