package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/app/cache"
	"srbac/controllers"
	"srbac/logics"
	"srbac/models"
	"srbac/utils"
	"time"
)

// 接口节点
type ApiItemController struct {
	controllers.Controller
}

// 接口节点列表
func (this *ApiItemController) List(c *gin.Context) {
	query := c.Request.URL.Query()
	page, per_page := utils.GetPageInfo(query)
	serviceId := utils.ToInt(c.Query("serviceId"))

	count := int64(0)
	find := app.Db.Model(&models.ApiItem{})
	if serviceId > 0 {
		find = find.Where("service_id = ?", serviceId)
	}
	re := find.Count(&count)
	app.CheckError(re.Error)

	apiItems := []*models.ApiItem{}
	re = logics.WithApiItemsOrder(find).Offset((page - 1) * per_page).Limit(per_page).Find(&apiItems)
	app.CheckError(re.Error)

	models.ApiItemsLoadServices(apiItems)
	serviceIds := logics.ServiceIds()

	this.HTML(c, "./views/admin/api-item/list.gohtml", map[string]interface{}{
		"menu": "api-item",
		"title": "接口节点列表",
		"pager": utils.GetPageHtml(count, page, per_page, query, "/admin/api-item/list"),
		"apiItems": apiItems,
		"serviceId": serviceId,
		"serviceIds": serviceIds,
	})
}

// 添加接口节点
func (this *ApiItemController) Add(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/api-item/list")

	apiItem := &models.ApiItem{}
	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		apiItem = models.NewApiItem(params)
		if apiItem.Validate() && apiItem.Create() {
			cache.SetApiItem(apiItem)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, apiItem.GetError())
		}
	}

	serviceIds := logics.ServiceIds()
	methods := models.ApiItemMethods()

	this.HTML(c, "./views/admin/api-item/add.gohtml", map[string]interface{}{
		"menu": "api-item",
		"title": "添加接口节点",
		"apiItem": apiItem,
		"serviceIds": serviceIds,
		"methods": methods,
	})
}

// 编辑接口节点
func (this *ApiItemController) Edit(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/api-item/list")
	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	apiItem := &models.ApiItem{}
	re := app.Db.First(apiItem, id)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		apiItem.IsAnonymousAccess = 0
		apiItem.SetAttributes(params)
		apiItem.UpdatedAt = time.Now().Unix()
		if apiItem.Validate() && apiItem.Update() {
			cache.SetApiItem(apiItem)
			cache.SetRoleApiItemsByApiItem(apiItem)
			cache.SetUserApiItemsByApiItem(apiItem)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, apiItem.GetError())
		}
	}

	serviceIds := logics.ServiceIds()
	methods := models.ApiItemMethods()

	this.HTML(c, "./views/admin/api-item/add.gohtml", map[string]interface{}{
		"menu": "api-item",
		"title": "编辑接口节点",
		"apiItem": apiItem,
		"serviceIds": serviceIds,
		"methods": methods,
	})
}

// 删除接口节点
func (this *ApiItemController) Delete(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/api-item/list", false)
	id := utils.ToInt64(this.GetPostForm(c)["id"])
	if id <= 0 {
		this.Redirect(c, referer)
	}

	apiItem := &models.ApiItem{}
	re := app.Db.First(apiItem, id)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	roleApiItems := []*models.RoleApiItem{}
	re = app.Db.Distinct("role_id", "service_id").Where("api_item_id = ?", apiItem.Id).Find(&roleApiItems)
	app.CheckError(re.Error)

	userApiItems := []*models.UserApiItem{}
	re = app.Db.Distinct("user_id", "service_id").Where("api_item_id = ?", apiItem.Id).Find(&userApiItems)
	app.CheckError(re.Error)

	re = app.Db.Delete(apiItem)
	app.CheckError(re.Error)

	cache.DelApiItem(apiItem)
	cache.SetRoleApiItemsByRoleApiItems(roleApiItems)
	cache.SetUserApiItemsByUserApiItems(userApiItems)
	this.Redirect(c, referer)
}