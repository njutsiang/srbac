package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 指定 URI 是否允许匿名访问
// true=允许匿名访问
// false=不允许匿名访问
// 缺省不允许匿名访问
var uris = map[string]bool{
	"/favicon.ico": true,
	"/admin/login": true,
	"/admin/logout": true,
}

// 检查登录状态
func SessionHandle(c *gin.Context) {
	isAnonymous, ok := uris[c.Request.RequestURI]
	if !(len(c.Request.RequestURI) >= 8 && c.Request.RequestURI[0:8] == "/assets/") && (!isAnonymous || !ok) {
		session := sessions.Default(c)
		if session.Get("user.id") == nil {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
		}
	}
	c.Next()
}