package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"srbac/app"
	"srbac/app/cache"
	"srbac/controllers"
	"srbac/models"
	"srbac/utils"
	"time"
)

// 角色
type RoleController struct {
	controllers.Controller
}

// 角色列表
func (this *RoleController) List(c *gin.Context) {
	query := c.Request.URL.Query()
	page, per_page := utils.GetPageInfo(query)

	count := int64(0)
	re := app.Db.Model(&models.Role{}).Count(&count)
	app.CheckError(re.Error)

	roles := []*models.Role{}
	re = app.Db.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&roles)
	app.CheckError(re.Error)

	this.HTML(c, "./views/admin/role/list.gohtml", map[string]interface{}{
		"menu": "role",
		"title": "角色列表",
		"pager": utils.GetPageHtml(count, page, per_page, query, "/admin/role/list"),
		"roles": roles,
	})
}

// 添加角色
func (this *RoleController) Add(c *gin.Context) {
	role := &models.Role{}
	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		role = models.NewRole(params)
		if role.Validate() && role.Create() {
			c.Redirect(http.StatusFound, "/admin/role/list")
			return
		} else {
			this.SetFailed(c, role.GetError())
		}
	}
	this.HTML(c, "./views/admin/role/add.gohtml", map[string]interface{}{
		"menu": "role",
		"title": "添加角色",
		"role": role,
	})
}

// 编辑角色
func (this *RoleController) Edit(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/role/list")
	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	role := &models.Role{}
	re := app.Db.First(role, id)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		role.SetAttributes(params)
		role.UpdatedAt = time.Now().Unix()
		if role.Validate() && role.Update() {
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, role.GetError())
		}
	}

	this.HTML(c, "./views/admin/role/add.gohtml", map[string]interface{}{
		"menu": "role",
		"title": "编辑角色",
		"role": role,
	})
}

// 删除角色
func (this *RoleController) Delete(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/role/list", false)
	id := utils.ToInt64(this.GetPostForm(c)["id"])
	if id <= 0 {
		this.Redirect(c, referer)
	}

	roleServices := []*models.RoleService{}
	re := app.Db.Where("role_id = ?", id).Find(&roleServices)
	app.CheckError(re.Error)

	re = app.Db.Delete(&models.Role{}, id)
	app.CheckError(re.Error)

	cache.DelRoleApiItemsByRoleServices(roleServices)
	cache.DelRoleDataItemsByRoleServices(roleServices)
	cache.DelRoleMenuItemsByRoleServices(roleServices)
	this.Redirect(c, referer)
}