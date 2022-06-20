package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/app/utils"
	"srbac/cache"
	"srbac/controllers"
	"srbac/logics"
	"srbac/models"
	"time"
)

// 菜单节点
type MenuItemController struct {
	controllers.Controller
}

// 菜单节点列表
func (this *MenuItemController) List(c *gin.Context) {
	query := c.Request.URL.Query()
	page, per_page := utils.GetPageInfo(query)
	serviceId := utils.ToInt(c.Query("serviceId"))

	count := int64(0)
	find := app.Db.Model(&models.MenuItem{})
	if serviceId > 0 {
		find = find.Where("service_id = ?", serviceId)
	}
	re := find.Count(&count)
	app.CheckError(re.Error)

	menuItems := []*models.MenuItem{}
	re = find.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&menuItems)
	app.CheckError(re.Error)

	models.MenuItemsLoadServices(menuItems)
	serviceIds := logics.ServiceIds()

	this.HTML(c, "./views/admin/menu-item/list.gohtml", map[string]interface{}{
		"menu": "menu-item",
		"title": "菜单节点列表",
		"pager": utils.GetPageHtml(count, page, per_page, query, "/admin/menu-item/list"),
		"menuItems": menuItems,
		"serviceId": serviceId,
		"serviceIds": serviceIds,
	})
}

// 添加菜单节点
func (this *MenuItemController) Add(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/menu-item/list")

	menuItem := &models.MenuItem{}
	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		menuItem = models.NewMenuItem(params)
		if menuItem.Validate() && menuItem.Create() {
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, menuItem.GetError())
		}
	}

	serviceIds := logics.ServiceIds()

	this.HTML(c, "./views/admin/menu-item/add.gohtml", map[string]interface{}{
		"menu": "menu-item",
		"title": "添加菜单节点",
		"menuItem": menuItem,
		"serviceIds": serviceIds,
	})
}

// 编辑菜单节点
func (this *MenuItemController) Edit(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/menu-item/list")
	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	menuItem := &models.MenuItem{}
	re := app.Db.First(menuItem, id)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		menuItem.SetAttributes(params)
		menuItem.UpdatedAt = time.Now().Unix()
		if menuItem.Validate() && menuItem.Update() {
			cache.SetRoleMenuItemsByMenuItem(menuItem)
			cache.SetUserMenuItemsByMenuItem(menuItem)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, menuItem.GetError())
		}
	}

	serviceIds := logics.ServiceIds()

	this.HTML(c, "./views/admin/menu-item/add.gohtml", map[string]interface{}{
		"menu": "menu-item",
		"title": "编辑菜单节点",
		"menuItem": menuItem,
		"serviceIds": serviceIds,
	})
}

// 删除菜单节点
func (this *MenuItemController) Delete(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/menu-item/list", false)
	id := utils.ToInt(this.GetPostForm(c)["id"])
	if id <= 0 {
		this.Redirect(c, referer)
	}

	roleMenuItems := []*models.RoleMenuItem{}
	re := app.Db.Distinct("role_id", "service_id").Where("menu_item_id = ?", id).Find(&roleMenuItems)
	app.CheckError(re.Error)

	userMenuItems := []*models.UserMenuItem{}
	re = app.Db.Distinct("user_id", "service_id").Where("menu_item_id = ?", id).Find(&userMenuItems)
	app.CheckError(re.Error)

	re = app.Db.Delete(&models.MenuItem{}, id)
	app.CheckError(re.Error)

	cache.SetRoleMenuItemsByRoleMenuItems(roleMenuItems)
	cache.SetUserMenuItemsByUserMenuItems(userMenuItems)
	this.Redirect(c, referer)
}