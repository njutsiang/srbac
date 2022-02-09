package admin

import (
	"github.com/gin-gonic/gin"
	"srbac/cache"
	"srbac/controllers"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
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
	find := srbac.Db.Model(&models.ApiItem{})
	if serviceId > 0 {
		find = find.Where("service_id = ?", serviceId)
	}
	re := find.Count(&count)
	srbac.CheckError(re.Error)

	apiItems := []*models.ApiItem{}
	re = find.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&apiItems)
	srbac.CheckError(re.Error)

	models.ApiItemsLoadServices(apiItems)
	serviceIds := models.ServiceIds()

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

	serviceIds := models.ServiceIds()
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
	re := srbac.Db.First(apiItem, id)
	srbac.CheckError(re.Error)

	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		apiItem.IsAnonymousAccess = 0
		apiItem.SetAttributes(params)
		apiItem.UpdatedAt = time.Now().Unix()
		if apiItem.Validate() && apiItem.Update() {
			cache.SetApiItem(apiItem)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, apiItem.GetError())
		}
	}

	serviceIds := models.ServiceIds()
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
	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	apiItem := &models.ApiItem{}
	re := srbac.Db.First(apiItem, id)
	srbac.CheckError(re.Error)

	srbac.Db.Delete(apiItem)
	cache.DelApiItem(apiItem)
	this.Redirect(c, referer)
}