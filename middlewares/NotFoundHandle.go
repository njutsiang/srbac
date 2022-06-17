package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"srbac/app"
)

// 处理 404 页面不存在
func NotFoundHandle(c *gin.Context) {
	app.HtmlStatus(c, http.StatusNotFound, "./views/admin/error/error.gohtml", map[string]interface{}{
		"code": 404,
		"message": "页面不存在",
		"referer": getReferer(c, "/admin"),
	})
}