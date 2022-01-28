package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"srbac/controllers"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
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
	re := srbac.Db.Model(&models.Role{}).Count(&count)
	srbac.CheckError(re.Error)

	roles := []*models.Role{}
	re = srbac.Db.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&roles)
	srbac.CheckError(re.Error)

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
	re := srbac.Db.First(role, id)
	srbac.CheckError(re.Error)

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
	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	re := srbac.Db.Delete(&models.Role{}, id)
	srbac.CheckError(re.Error)
	this.Redirect(c, referer)
}