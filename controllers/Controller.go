package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"srbac/app"
	"srbac/exception"
	"srbac/models"
	"srbac/utils"
)

// 控制器基类
type Controller struct {
}

// 输出 HTML
func (this *Controller) HTML(ctx *gin.Context, filename string, params ...map[string]interface{}) {
	app.HtmlStatus(ctx, http.StatusOK, filename, params...)
}

// 获取 POST 表单数据
func (this *Controller) GetPostForm(ctx *gin.Context) map[string]interface{} {
	if ctx.Request.Method == "POST" && ctx.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		if err := ctx.Request.ParseForm(); err != nil {
			app.Error(err)
			return map[string]interface{}{}
		}
		return utils.ToMap(ctx.Request.PostForm)
	}
	return map[string]interface{}{}
}

// 获取 POST 多段表单数据
func (this *Controller) GetPostMultipartForm(ctx *gin.Context) *multipart.Form {
	contentType := ctx.Request.Header.Get("Content-Type")
	if ctx.Request.Method == "POST" && len(contentType) >= 19 && contentType[0:19] == "multipart/form-data" {
		if err := ctx.Request.ParseMultipartForm(app.Engine.MaxMultipartMemory); err != nil {
			app.Error(err)
		}
	}
	return ctx.Request.MultipartForm
}

// 获取 POST JSON Map 数据
func (this *Controller) GetPostJson(ctx *gin.Context) map[string]interface{} {
	contentType := ctx.Request.Header.Get("Content-Type")
	if ctx.Request.Method == "POST" && len(contentType) >= 16 && contentType[0:16] == "application/json" {
		data := map[string]interface{}{}
		if err := ctx.BindJSON(&data); err != nil {
			app.Error(data)
		}
		return data
	}
	return map[string]interface{}{}
}

// 获取 POST JSON Slice 数据
func (this *Controller) GetPostJsonSlice(ctx *gin.Context) []interface{} {
	contentType := ctx.Request.Header.Get("Content-Type")
	if ctx.Request.Method == "POST" && len(contentType) >= 16 && contentType[0:16] == "application/json" {
		data := []interface{}{}
		if err := ctx.BindJSON(&data); err != nil {
			app.Error(data)
		}
		return data
	}
	return []interface{}{}
}

// 获取 POST 原始数据
func (this *Controller) GetPostRaw(ctx *gin.Context) string {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		app.Error(err)
	}
	return string(body)
}

// 获取上一个页面的地址
// defaultUri，当 Query 和 Header 都不存在 referer 时，将该值作为 referer
// redirect，默认 true，当 Query 和 Header 都不存在 referer 时，是否重定向，并且在 Query 中带上 referer
func (this *Controller) GetReferer(ctx *gin.Context, defaultUri string, redirect ...bool) string {
	referer := ctx.Query("referer")
	if referer == "" {
		referer = ctx.Request.Header.Get("Referer")
		if referer == "" || referer == ctx.Request.URL.RequestURI() {
			referer = defaultUri
		}
		if !(len(redirect) >= 1 && !redirect[0]) {
			query := ctx.Request.URL.Query()
			query.Set("referer", referer)
			this.Redirect(ctx, ctx.Request.URL.Path+"?"+query.Encode())
		}
	}
	return referer
}

// 跳转至指定的地址，并且 panic
func (this *Controller) Redirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusFound, location)
	panic(app.Redirect(location))
}

// 设置错误信息
func (this *Controller) SetFailed(ctx *gin.Context, message string) {
	ctx.Set("failed", message)
}

// 获取错误信息
func (this *Controller) GetFailed(ctx *gin.Context) string {
	if failed, exists := ctx.Get("failed"); exists {
		return utils.ToString(failed)
	} else {
		return ""
	}
}

// 设置成功信息
func (this *Controller) SetSuccess(ctx *gin.Context, message string) {
	ctx.Set("success", message)
}

// 获取成功信息
func (this *Controller) GetSuccess(ctx *gin.Context) string {
	if success, exists := ctx.Get("success"); exists {
		return utils.ToString(success)
	} else {
		return ""
	}
}

// 响应错误，并且 panic
func (this *Controller) ResponseErrorJson(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusBadRequest, data)
	panic(app.Response(http.StatusBadRequest))
}

// 获取当前登录用户
func (this *Controller) GetUser(ctx *gin.Context) *models.User {
	value, exists := ctx.Get("user")
	if !exists {
		panic(app.NewJsonError(exception.UserNotLogin))
	}
	user, ok := value.(*models.User)
	if !ok {
		panic(app.NewJsonError(exception.UserNotLogin))
	}
	return user
}

// 当前登录用户 id
func (this *Controller) GetUserId(ctx *gin.Context) int64 {
	session := sessions.Default(ctx)
	return utils.ToInt64(session.Get("user.id"))
}
