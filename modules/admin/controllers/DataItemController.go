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

// 数据节点
type DataItemController struct {
	controllers.Controller
}

// 数据节点列表
func (this *DataItemController) List(c *gin.Context) {
	query := c.Request.URL.Query()
	page, per_page := utils.GetPageInfo(query)
	serviceId := utils.ToInt(c.Query("serviceId"))

	count := int64(0)
	find := app.Db.Model(&models.DataItem{})
	if serviceId > 0 {
		find = find.Where("service_id = ?", serviceId)
	}
	re := find.Count(&count)
	app.CheckError(re.Error)

	dataItems := []*models.DataItem{}
	re = find.Order("id asc").Offset((page - 1) * per_page).Limit(per_page).Find(&dataItems)
	app.CheckError(re.Error)

	models.DataItemsLoadServices(dataItems)
	serviceIds := logics.ServiceIds()

	this.HTML(c, "./views/admin/data-item/list.gohtml", map[string]interface{}{
		"menu": "data-item",
		"title": "数据节点列表",
		"pager": utils.GetPageHtml(count, page, per_page, query, "/admin/data-item/list"),
		"dataItems": dataItems,
		"serviceId": serviceId,
		"serviceIds": serviceIds,
	})
}

// 添加数据节点
func (this *DataItemController) Add(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/data-item/list")

	dataItem := &models.DataItem{}
	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		dataItem = models.NewDataItem(params)
		if dataItem.Validate() && dataItem.Create() {
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, dataItem.GetError())
		}
	}

	serviceIds := logics.ServiceIds()

	this.HTML(c, "./views/admin/data-item/add.gohtml", map[string]interface{}{
		"menu": "data-item",
		"title": "添加数据节点",
		"dataItem": dataItem,
		"serviceIds": serviceIds,
	})
}

// 编辑数据节点
func (this *DataItemController) Edit(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/data-item/list")
	id := utils.ToInt(c.Query("id"))
	if id <= 0 {
		this.Redirect(c, referer)
	}

	dataItem := &models.DataItem{}
	re := app.Db.First(dataItem, id)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	if c.Request.Method == "POST" {
		params := this.GetPostForm(c)
		dataItem.SetAttributes(params)
		dataItem.UpdatedAt = time.Now().Unix()
		if dataItem.Validate() && dataItem.Update() {
			cache.SetRoleDataItemsByDataItem(dataItem)
			cache.SetUserDataItemsByDataItem(dataItem)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, dataItem.GetError())
		}
	}

	serviceIds := logics.ServiceIds()

	this.HTML(c, "./views/admin/data-item/add.gohtml", map[string]interface{}{
		"menu": "data-item",
		"title": "编辑数据节点",
		"dataItem": dataItem,
		"serviceIds": serviceIds,
	})
}

// 删除数据节点
func (this *DataItemController) Delete(c *gin.Context) {
	referer := this.GetReferer(c, "/admin/data-item/list", false)
	id := utils.ToInt(this.GetPostForm(c)["id"])
	if id <= 0 {
		this.Redirect(c, referer)
	}

	roleDataItems := []*models.RoleDataItem{}
	re := app.Db.Distinct("role_id", "service_id").Where("data_item_id = ?", id).Find(&roleDataItems)
	app.CheckError(re.Error)

	userDataItems := []*models.UserDataItem{}
	re = app.Db.Distinct("user_id", "service_id").Where("data_item_id = ?", id).Find(&userDataItems)
	app.CheckError(re.Error)

	re = app.Db.Delete(&models.DataItem{}, id)
	app.CheckError(re.Error)

	cache.SetRoleDataItemsByRoleDataItems(roleDataItems)
	cache.SetUserDataItemsByUserDataItems(userDataItems)
	this.Redirect(c, referer)
}