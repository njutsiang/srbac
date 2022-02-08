package admin

import (
	"github.com/gin-gonic/gin"
	"srbac/controllers"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
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
	find := srbac.Db.Model(&models.MenuItem{})
	if serviceId > 0 {
		find = find.Where("service_id = ?", serviceId)
	}
	re := find.Count(&count)
	srbac.CheckError(re.Error)

	menuItems := []*models.MenuItem{}
	re = find.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&menuItems)
	srbac.CheckError(re.Error)

	models.MenuItemsLoadServices(menuItems)
	serviceIds := models.ServiceIds()

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

	serviceIds := models.ServiceIds()

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
	re := srbac.Db.First(menuItem, id)
	srbac.CheckError(re.Error)

	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		menuItem.SetAttributes(params)
		menuItem.UpdatedAt = time.Now().Unix()
		if menuItem.Validate() && menuItem.Update() {
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, menuItem.GetError())
		}
	}

	serviceIds := models.ServiceIds()

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
	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}
	re :=srbac.Db.Delete(&models.MenuItem{}, id)
	srbac.CheckError(re.Error)
	this.Redirect(c, referer)
}