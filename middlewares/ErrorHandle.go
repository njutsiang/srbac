package middlewares

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"srbac/exception"
	"srbac/libraries/log"
	"srbac/srbac"
)

// 处理 panic 抛出的异常
func ErrorHandle(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			code := 0
			message := "未知错误"
			switch err.(type) {
			case srbac.Redirect, srbac.Response:
				return
			case string:
				message = err.(string)
			case error:
				if jsonError, ok := err.(*srbac.JsonError); ok {
					handleJsonError(ctx, jsonError)
					return
				}
				if ex, ok := err.(*exception.Exception); ok {
					code = int(ex.GetCode())
				}
				message = err.(error).Error()
			}
			log.Error(err)
			handleHtmlError(ctx, code, message)

			// 执行 Abort() 结束当前请求，否则程序会继续执行 Controller
			ctx.Abort()
		}
	}()
	ctx.Next()
}

// 处理 JsonError
func handleJsonError(ctx *gin.Context, jsonError *srbac.JsonError) {
	ctx.Status(jsonError.StatusCode)
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	jsonString, _ := json.Marshal(map[string]interface{}{
		"error": map[string]interface{}{
			"code": jsonError.Code,
			"message": jsonError.Message,
		},
	})
	if _, err := ctx.Writer.Write(jsonString); err != nil {
		log.Error(err)
	}
}

// 处理错误提示语
func handleHtmlError(c *gin.Context, code int, message string) {
	srbac.HtmlStatus(c, http.StatusBadRequest, "./views/admin/error/error.gohtml", map[string]interface{}{
		"code": code,
		"message": message,
		"referer": getReferer(c, "/admin"),
	})
}

// 获取上一个页面的地址
func getReferer(c *gin.Context, defaultUri string) string {
	referer := c.Query("referer")
	if referer == "" {
		referer = c.Request.Header.Get("Referer")
		if referer == "" || referer == c.Request.URL.RequestURI() {
			referer = defaultUri
		}
	}
	return referer
}