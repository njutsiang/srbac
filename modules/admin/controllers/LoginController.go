package admin

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"srbac/cache"
	"srbac/controllers"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
	"time"
)

type LoginController struct {
	controllers.Controller
}

// 登录
func (this *LoginController) Login(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user.id") != nil {
		this.Redirect(c, "/admin")
	}
	form := struct {
		Username string
		Password string
		RememberMe int
	}{}
	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		form.Username = utils.ToString(params["username"])
		form.Password = utils.ToString(params["password"])
		form.RememberMe = utils.ToInt(params["remember_me"])
		maxAge := 24 * 3600
		if form.RememberMe == 1 {
			maxAge *= 30
		}
		session.Options(sessions.Options{
			MaxAge: maxAge,
		})
		hasErr := false
		if form.Username == "" {
			hasErr = true
			this.SetFailed(c, "账号不能为空")
		}
		if !hasErr && form.Password == "" {
			hasErr = true
			this.SetFailed(c, "密码不能为空")
		}
		if !hasErr {
			user := &models.User{}
			re := srbac.Db.Where("username = ?", form.Username).First(user)
			if errors.Is(re.Error, gorm.ErrRecordNotFound) {
				hasErr = true
				this.SetFailed(c, "用户不存在")
			}
			if !hasErr && re.Error != nil {
				hasErr = true
				this.SetFailed(c, re.Error.Error())
			}
			if !hasErr && user.Status != 1 {
				hasErr = true
				this.SetFailed(c, "用户已被禁用")
			}
			if !hasErr && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)) != nil  {
				hasErr = true
				this.SetFailed(c, "密码错误")
			}
			if !hasErr {
				token := uuid.NewString()
				token = fmt.Sprintf("%x", md5.Sum([]byte(token)))
				cache.SetUserToken(token, user, maxAge)
				session.Set("user.id", user.Id)
				session.Set("user.name", user.Name)
				session.Set("user.username", user.Username)
				session.Set("user.status", user.Status)
				session.Set("user.token", token)
				if err := session.Save(); err == nil {
					http.SetCookie(c.Writer, &http.Cookie{
						Name: "user_token",
						Value: token,
						Expires: time.Unix(time.Now().Unix() + int64(maxAge), 0),
						MaxAge: maxAge,
						Path: "/",
					})
					this.Redirect(c, "/admin")
				} else {
					hasErr = true
					this.SetFailed(c, err.Error())
				}
			}
		}
	}
	this.HTML(c, "./views/admin/login/login.gohtml", map[string]interface{}{
		"title": "登录",
		"form": form,
	})
}

// 退出登录
func (this *LoginController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user.id")
	session.Delete("user.name")
	session.Delete("user.username")
	session.Delete("user.status")
	if token := session.Get("user.token"); token != nil {
		cache.DelUserToken(utils.ToString(token))
	}
	session.Delete("user.token")
	err := session.Save()
	srbac.CheckError(err)
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "user_token",
		MaxAge: -1,
	})
	this.Redirect(c, "/admin/login")
}