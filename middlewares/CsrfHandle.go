package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"srbac/app"
	"srbac/exception"
	"srbac/utils"
)

// 处理 CSRF 问题
func CsrfHandle(c *gin.Context) {
	session := sessions.Default(c)
	csrfToken := ""
	sessionCsrfToken := ""
	if _csrfToken := session.Get("csrf.token"); _csrfToken == nil {
		maxAge := 24 * 3600
		if c.Request.Method == "POST" {
			err := c.Request.ParseForm()
			app.CheckError(err)
			rememberMe := utils.ToInt(c.Request.PostForm.Get("remember_me"))
			if rememberMe == 1 {
				maxAge *= 30
			}
		}
		session.Options(sessions.Options{
			MaxAge: maxAge,
		})
		csrfToken = generateCsrfToken()
		session.Set("csrf.token", csrfToken)
		err := session.Save()
		app.CheckError(err)
	} else {
		csrfToken = utils.ToString(_csrfToken)
		sessionCsrfToken = csrfToken
	}
	if c.Request.Method == "POST" && sessionCsrfToken != getPostCsrfToken(c) {
		exception.Throw(exception.CsrfTokenError)
	}
	c.Header("X-Csrf-Token", csrfToken)
	c.Next()
}

// 生成一个新的 CsrfToken
func generateCsrfToken() string {
	return uuid.NewString()
}

// 获取 POST 请求中的表单令牌
func getPostCsrfToken(c *gin.Context) string {
	err := c.Request.ParseForm()
	app.CheckError(err)
	return c.Request.PostForm.Get("_csrf_token")
}