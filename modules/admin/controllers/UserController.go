package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"srbac/cache"
	"srbac/controllers"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
	"time"
)

// 用户
type UserController struct {
	controllers.Controller
}

// 用户列表
func (this *UserController) List(c *gin.Context) {
	query := c.Request.URL.Query()
	page, per_page := utils.GetPageInfo(query)

	count := int64(0)
	re := srbac.Db.Model(&models.User{}).Count(&count)
	srbac.CheckError(re.Error)

	users := []*models.User{}
	re = srbac.Db.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&users)
	srbac.CheckError(re.Error)

	this.HTML(c, "./views/admin/user/list.gohtml", map[string]interface{}{
		"menu": "user",
		"title": "用户列表",
		"pager": utils.GetPageHtml(count, page, per_page, query, "/admin/user/list"),
		"users": users,
	})
}

// 添加用户
func (this *UserController) Add(c *gin.Context) {
	user := &models.User{}
	user.Status = 1
	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		user = models.NewUser(params)
		if user.Validate() && user.Create() {
			c.Redirect(http.StatusFound, "/admin/user/list")
			return
		} else {
			this.SetFailed(c, user.GetError())
		}
	}
	this.HTML(c, "./views/admin/user/add.gohtml", map[string]interface{}{
		"menu": "user",
		"title": "添加用户",
		"user": user,
	})
}

// 编辑用户
func (this *UserController) Edit(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/user/list")
	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	user := &models.User{}
	re := srbac.Db.First(user, id)
	srbac.CheckError(re.Error)

	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		user.Status = 0
		user.SetAttributes(params)
		user.UpdatedAt = time.Now().Unix()
		if user.Validate() && user.Update() {
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, user.GetError())
		}
	}

	this.HTML(c, "./views/admin/user/add.gohtml", map[string]interface{}{
		"menu": "user",
		"title": "编辑用户",
		"user": user,
	})
}

// 删除用户
func (this *UserController) Delete(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/user/list", false)
	id := utils.ToInt64(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	re := srbac.Db.Delete(&models.User{}, id)
	srbac.CheckError(re.Error)
	cache.DelUserRoles(id)
	this.Redirect(c, referer)
}