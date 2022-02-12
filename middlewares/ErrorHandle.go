package middlewares

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"srbac/libraries/log"
	"srbac/srbac"
)

// 处理 panic 抛出的异常
func ErrorHandle(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
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
				message = err.(error).Error()
			}
			log.Error(err)
			handleHtmlError(ctx, message)
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
	// 执行 Abort() 结束当前请求，否则程序会继续执行 Controller
	ctx.Abort()
}

// 处理错误提示语
func handleHtmlError(ctx *gin.Context, message string) {
	ctx.Status(http.StatusInternalServerError)
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	if _, err := ctx.Writer.Write([]byte(message)); err != nil {
		log.Error(err)
	}
	ctx.Abort()
}
