package srbac

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"srbac/libraries/log"
	"srbac/libraries/utils"
	"strings"
)

// 判断是否有错误，有则记录错误日志，并抛出 panic
func CheckError(err interface{}) {
	if err == nil {
		return
	}
	panic(err)
}

// 输出 HTML
func HtmlStatus(c *gin.Context, code int, filename string, params ...map[string]interface{}) {
	session := sessions.Default(c)

	// 响应状态和响应头
	c.Status(code)
	c.Header("Content-Type", "text/html; charset=utf-8")

	// 所有模板文件
	filenames := []string{"",
		"./views/admin/layout/head.gohtml",
		"./views/admin/layout/header.gohtml",
		"./views/admin/layout/menu.gohtml",
		"./views/admin/layout/footer.gohtml",
	}
	copy(filenames, []string{
		filename,
	})

	// 主模板文件
	filenameItems := strings.Split(filename, "/")
	filenameItem := filenameItems[len(filenameItems) - 1]

	// 向模板注册自定义函数
	tmpl := template.New(filenameItem).Funcs(template.FuncMap{
		"InSlice": utils.InSlice,
	})

	// 解析所有模板
	tmpl, err := tmpl.ParseFiles(filenames...)
	if err != nil {
		log.Error(err)
		c.Status(http.StatusNotFound)
		if _, err := c.Writer.Write([]byte(err.Error())); err != nil {
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
	data["failed"] = template.HTML(getFailed(c))
	data["success"] = template.HTML(getSuccess(c))
	data["uri"] = c.Request.URL.RequestURI()
	if len(title) >= 1 {
		title += " - "
	}
	title += "SRBAC 基于服务和角色的访问控制"
	data["headTitle"] = title
	data["path"] = c.Request.URL.Path
	data["sessionUserId"] = utils.ToString(session.Get("user.id"))
	data["sessionUserName"] = utils.ToString(session.Get("user.name"))
	data["csrfTokenKey"] = "_csrf_token"
	data["csrfTokenValue"] = utils.ToString(session.Get("csrf.token"))

	// 载入数据，并执行模板文件
	err = tmpl.Execute(c.Writer, data)
	CheckError(err)
}

// 获取错误信息
func getFailed(c *gin.Context) string {
	if failed, exists := c.Get("failed"); exists {
		return utils.ToString(failed)
	} else {
		return ""
	}
}

// 获取成功信息
func getSuccess(c *gin.Context) string {
	if success, exists := c.Get("success"); exists {
		return utils.ToString(success)
	} else {
		return ""
	}
}