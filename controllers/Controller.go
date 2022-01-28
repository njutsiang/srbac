package controllers

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"srbac/code"
	"srbac/libraries/log"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
)

// 控制器基类
type Controller struct {
}

// 输出 HTML
func (this *Controller) HTML(ctx *gin.Context, filename string, params ...map[string]interface{}) {
	this.HtmlStatus(ctx, http.StatusOK, filename, params...)
}

// 输出 HTML，并指定 Http 状态码
func (this *Controller) HtmlStatus(ctx *gin.Context, code int, filename string, params ...map[string]interface{}) {
	// 响应状态和响应头
	ctx.Status(code)
	ctx.Header("Content-Type", "text/html; charset=utf-8")

	// 载入模板文件
	filenames := []string{"",
		"./views/admin/layout/head.gohtml",
		"./views/admin/layout/header.gohtml",
		"./views/admin/layout/menu.gohtml",
		"./views/admin/layout/footer.gohtml",
	}
	copy(filenames, []string{
		filename,
	})
	tmpl, err := template.ParseFiles(filenames...)
	if err != nil {
		log.Error(err)
		ctx.Status(http.StatusNotFound)
		if _, err := ctx.Writer.Write([]byte(err.Error())); err != nil {
			log.Error(err)
		}
		return
	}

	// 处理数据和页面标题
	data := map[string]interface{}{}
	if len(params) >= 1 {
		data = params[0]
	}
	title := utils.ToString(data["title"])
	data["title"] = title
	data["failed"] = template.HTML(this.GetFailed(ctx))
	data["success"] = template.HTML(this.GetSuccess(ctx))
	data["uri"] = ctx.Request.URL.RequestURI()
	if len(title) >= 1 {
		title += " - "
	}
	title += "SRBAC 基于服务和角色的访问控制"
	data["headTitle"] = title
	data["path"] = ctx.Request.URL.Path

	// 载入数据，并执行模板文件
	err = tmpl.Execute(ctx.Writer, data)
	srbac.CheckError(err)
}

// 获取 POST 表单数据
func (this *Controller) GetPostForm(ctx *gin.Context) map[string]interface{} {
	if ctx.Request.Method == "POST" && ctx.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		if err := ctx.Request.ParseForm(); err != nil {
			log.Error(err)
			return map[string]interface{}{}
		}
		return utils.ToMapInterfaces(ctx.Request.PostForm)
	}
	return map[string]interface{}{}
}

// 获取 POST 多段表单数据
func (this *Controller) GetPostMultipartForm(ctx *gin.Context) *multipart.Form {
	contentType := ctx.Request.Header.Get("Content-Type")
	if ctx.Request.Method == "POST" && len(contentType) >= 19 && contentType[0:19] == "multipart/form-data" {
		if err := ctx.Request.ParseMultipartForm(srbac.Engine.MaxMultipartMemory); err != nil {
			log.Error(err)
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
			log.Error(data)
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
			log.Error(data)
		}
		return data
	}
	return []interface{}{}
}

// 获取 POST 原始数据
func (this *Controller) GetPostRaw(ctx *gin.Context) string {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error(err)
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
	panic(srbac.Redirect(location))
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
	panic(srbac.Response(http.StatusBadRequest))
}

// 获取当前登录用户
func (this *Controller) GetUser(ctx *gin.Context) *models.User {
	value, exists := ctx.Get("user")
	if !exists {
		panic(srbac.NewJsonError(code.UserNotLogin))
	}
	user, ok := value.(*models.User)
	if !ok {
		panic(srbac.NewJsonError(code.UserNotLogin))
	}
	return user
}
