package controllers

import (
	"Blog/helpers"
	"Blog/models"
	"Blog/system"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/cihub/seelog"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type GithubUserInfo struct {
	AvatarURL         string      `json:"avatar_url"`
	Bio               interface{} `json:"bio"`
	Blog              string      `json:"blog"`
	Company           interface{} `json:"company"`
	CreatedAt         string      `json:"created_at"`
	Email             interface{} `json:"email"`
	EventsURL         string      `json:"events_url"`
	Followers         int         `json:"followers"`
	FollowersURL      string      `json:"followers_url"`
	Following         int         `json:"following"`
	FollowingURL      string      `json:"following_url"`
	GistsURL          string      `json:"gists_url"`
	GravatarID        string      `json:"gravatar_id"`
	Hireable          interface{} `json:"hireable"`
	HtmlURL           string      `json:"html_url"`
	ID                int         `json:"id"`
	Location          interface{} `json:"location"`
	Login             string      `json:"login"`
	Name              string      `json:"name"`
	NodeID            string      `json:"node_id"`
	OrganizationsURL  string      `json:"organizations_url"`
	PublicGists       int         `json:"public_gists"`
	PublicRepos       int         `json:"public_repos"`
	ReceivedEventsURL string      `json:"received_events_url"`
	ReposURL          string      `json:"repos_url"`
	SiteAdmin         bool        `json:"site_admin"`
	StarredURL        string      `json:"starred_url"`
	SubscriptionsURL  string      `json:"subscriptions_url"`
	TwitterUsername   interface{} `json:"twitter_username"`
	Type              string      `json:"type"`
	UpdatedAt         string      `json:"updated_at"`
	URL               string      `json:"url"`
}

func SigninGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signin.html", nil)
}

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/signup.html", nil)
}

func SignupPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	email := c.PostForm("email")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")
	user := &models.User{
		Email:     email,
		Telephone: telephone,
		Password:  password,
		IsAdmin:   false,
	}

	if len(user.Email) == 0 || len(user.Password) == 0 {
		res["message"] = "Email and password can not be empty"
		return
	}

	user.Password = helpers.Md5(user.Email + user.Password)
	err = user.Insert()
	if err != nil {
		res["message"] = "Email already exists"
		return
	}
	res["succeed"] = true
	res["message"] = "请返回主页重新登录"
}

func SigninPost(c *gin.Context) {
	var (
		err  error
		user *models.User
	)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if len(username) == 0 || len(password) == 0 {
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": "Email and password can not be empty",
		})
		return
	}
	user, err = models.GetUserByUserName(username)
	if err != nil {
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": "User does not exist",
		})
		return
	}
	if user.Password != helpers.Md5(user.Email+password) {
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": "Password is incorrect",
		})
		return
	}
	if user.LockState {
		c.HTML(http.StatusOK, "auth/signin.html", gin.H{
			"message": "User is locked",
		})
		return
	}
	s := sessions.Default(c)
	s.Clear()
	s.Set(SESSION_KEY, user.ID)
	s.Save()

	if user.IsAdmin {
		c.Redirect(http.StatusFound, "/admin/index")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

func LogoutGet(c *gin.Context) {
	s := sessions.Default(c)
	s.Clear()
	s.Save()
	c.Redirect(http.StatusFound, "/")
}

func Oauth2Callback(c *gin.Context) {
	var (
		userInfo *GithubUserInfo
		user     *models.User
	)
	code := c.Query("code")
	state := c.Query("state")

	// validate state
	session := sessions.Default(c)
	if len(state) == 0 || state != session.Get(SESSION_GITHUB_STATE) {
		c.Abort()
		return
	}
	// remove state from session
	session.Delete(SESSION_GITHUB_STATE)
	session.Save()

	// exchange accesstoken by code
	token, err := exchangeTokenByCode(code)
	if err != nil {
		seelog.Errorf("exchangeTokenByCode error: %v", err)
		c.Redirect(http.StatusFound, "/signin")
		return
	}
	// fmt.Println("token: ", token)

	//get github userinfo by accesstoken
	userInfo, err = getGithubUserInfoByAccessToken(token)
	if err != nil {
		seelog.Errorf("getGithubUserInfoByAccessToken error: %v", err)
		c.Redirect(http.StatusFound, "/signin")
		return
	}
	// fmt.Printf("userInfo: %+v\n", userInfo)

	sessionUser, exists := c.Get(CONTEXT_USER_KEY)
	if exists {
		// 已登录
		user = sessionUser.(*models.User)
		_, err := models.IsGithubIdExists(userInfo.Login, int(user.ID))
		if err != nil {
			if user.IsAdmin {
				user.GithubLoginId = userInfo.Login
			}
			user.AvatarUrl = userInfo.AvatarURL
			user.GithubUrl = userInfo.HtmlURL
			user.Nickname = userInfo.Login
			err = user.UpdateGithubUserInfo()
			if err != nil {
				seelog.Errorf("update github user info error: %v", err)
				c.Redirect(http.StatusFound, "/signin")
				return
			}
		} else {
			c.HTML(http.StatusOK, "auth/signin.html", gin.H{
				"message": "this github loginId has bound another account.",
			})
			return
		}
	} else {
		user = &models.User{
			Nickname:      userInfo.Login,
			GithubLoginId: userInfo.Login,
			AvatarUrl:     userInfo.AvatarURL,
			GithubUrl:     userInfo.HtmlURL,
		}
		err = user.FirstOrCreate()
		if err != nil {
			seelog.Errorf("create github user error: %v", err)
			c.Redirect(http.StatusFound, "/signin")
			return
		}
		if user.LockState {
			c.HTML(http.StatusOK, "auth/signin.html", gin.H{
				"message": "User is locked",
			})
			return
		}
	}
	s := sessions.Default(c)
	s.Clear()
	s.Set(SESSION_KEY, user.ID)
	s.Save()
	if user.IsAdmin {
		c.Redirect(http.StatusFound, "/admin/index")
	} else {
		c.Redirect(http.StatusFound, "/")
	}

}

func exchangeTokenByCode(code string) (accessToken string, err error) {
	type Token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"` // 这个字段下面没用到
		Scope       string `json:"scope"`      // 这个字段下面也没用到
	}

	url := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		system.GetConfiguration().GithubClientID, system.GetConfiguration().GithubClientSecret, code)

	// 形成请求
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, url, nil); err != nil {
		seelog.Errorf("new request error: %v", err)
		return "", err
	}
	req.Header.Set("accept", "application/json")

	// 发送请求并获得响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		seelog.Errorf("do request error: %v", err)
		return "", err
	}

	// 将响应体解析为 token，并返回
	var token Token
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		seelog.Errorf("decode token error: %v", err)
		return "", err
	}
	return token.AccessToken, nil
}

func getGithubUserInfoByAccessToken(token string) (*GithubUserInfo, error) {
	var (
		req  *http.Request
		body []byte
		err  error
	)
	req, err = http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		seelog.Errorf("new request error: %v", err)
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", "token "+token)

	var client = http.Client{}
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		seelog.Errorf("do request error: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		seelog.Errorf("read body error: %v", err)
		return nil, err
	}

	var userInfo GithubUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		seelog.Errorf("unmarshal error: %v", err)
		return nil, err
	}
	return &userInfo, nil
}

func UserIndexGet(c *gin.Context) {
	users, _ := models.ListUsers()
	user, _ := c.Get(CONTEXT_USER_KEY)
	c.HTML(http.StatusOK, "admin/user.html", gin.H{
		"users":    users,
		"user":     user,
		"comments": models.MustListUnreadComment(),
	})
}

func UserLockPost(c *gin.Context) {
	var (
		err  error
		_id  int64
		res  = gin.H{}
		user *models.User
	)
	defer writeJSON(c, res)
	id := c.Param("id")
	_id, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	user, err = models.GetUserByID(uint(_id))
	if err != nil {
		res["message"] = err.Error()
		return
	}

	if user.LockState {
		err = user.Unlock()
	} else {
		err = user.Lock()
	}
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func ProfileGet(c *gin.Context) {
	sessionUser, exists := c.Get(CONTEXT_USER_KEY)
	if !exists {
		c.Redirect(http.StatusFound, "/admin/index")
		return
	}
	c.HTML(http.StatusOK, "admin/profile.html", gin.H{
		"user":     sessionUser,
		"comments": models.MustListUnreadComment(),
	})
}

func ProfilePost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	avatarUrl := c.PostForm("avatarUrl")
	nickName := c.PostForm("nickName")
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
	user, ok := sessionUser.(*models.User)
	if !ok {
		res["message"] = "session error"
		return
	}
	err = user.UpdateProfile(avatarUrl, nickName)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
	res["user"] = models.User{AvatarUrl: avatarUrl, Nickname: nickName}
}

func ProfileEmailBind(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	email := c.PostForm("email")
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
	user, ok := sessionUser.(*models.User)
	if !ok {
		res["message"] = "session error"
		return
	}
	if len(user.Email) > 0 {
		res["message"] = "email has been bound"
		return
	}
	_, err = models.GetUserByUserName(email)
	if err == nil {
		res["message"] = "email has been bound"
		return
	}
	err = user.UpdateEmail(email)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func ProfileEmailUnbind(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
	user, ok := sessionUser.(*models.User)
	if !ok {
		res["message"] = "session error"
		return
	}
	if len(user.Email) == 0 {
		res["message"] = "email has not been bound"
		return
	}
	err = user.UpdateEmail("")
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func ProfileGithubUnbind(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer writeJSON(c, res)
	sessionUser, _ := c.Get(CONTEXT_USER_KEY)
	user, ok := sessionUser.(*models.User)
	if !ok {
		res["message"] = "session error"
		return
	}
	if len(user.GithubLoginId) == 0 {
		res["message"] = "github has not been bound"
		return
	}
	err = user.UpdateGithubUserInfo()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}
