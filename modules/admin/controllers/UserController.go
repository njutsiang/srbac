package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"srbac/app"
	"srbac/app/utils"
	"srbac/cache"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/models"
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
	re := app.Db.Model(&models.User{}).Count(&count)
	app.CheckError(re.Error)

	users := []*models.User{}
	re = app.Db.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&users)
	app.CheckError(re.Error)

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
	id := utils.ToInt64(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}
	if id == 1 && this.GetUserId(c) != 1 {
		exception.NewException(code.NoPermission)
	}

	user := &models.User{}
	re := app.Db.First(user, id)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

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
	id := utils.ToInt64(this.GetPostForm(c)["id"])
	if id <= 1 {
		this.Redirect(c, referer)
	}

	userServices := []*models.UserService{}
	re := app.Db.Where("user_id = ?", id).Find(&userServices)
	app.CheckError(re.Error)

	re = app.Db.Delete(&models.User{}, id)
	app.CheckError(re.Error)

	cache.DelUserRoles(id)
	cache.DelUserApiItemsByUserServices(userServices)
	cache.DelUserDataItemsByUserServices(userServices)
	cache.DelUserMenuItemsByUserServices(userServices)
	this.Redirect(c, referer)
}