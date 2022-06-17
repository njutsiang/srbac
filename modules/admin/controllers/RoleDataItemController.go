package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/cache"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/libraries/utils"
	"srbac/models"
)

// 角色的数据权限
type RoleDataItemController struct {
	controllers.Controller
}

// 编辑角色的数据权限
func (this *RoleDataItemController) Edit(c *gin.Context) {
	roleId := utils.ToInt(c.Query("roleId"))
	roleServiceId := utils.ToInt(c.Query("roleServiceId"))
	if roleServiceId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/role-service/list?roleId=%d", roleId))

	roleService := &models.RoleService{}
	re := app.Db.First(roleService, roleServiceId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	models.RoleServicesLoadServices([]*models.RoleService{roleService})

	role := &models.Role{}
	re = app.Db.First(role, roleService.RoleId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	dataItems := []*models.DataItem{}
	re = app.Db.Where("service_id = ?", roleService.ServiceId).Order("`key` ASC").Limit(1000).Find(&dataItems)
	app.CheckError(re.Error)

	// 角色和数据节点的关联
	roleDataItems := []*models.RoleDataItem{}
	re = app.Db.Where("role_id = ? AND service_id = ?", roleService.RoleId, roleService.ServiceId).Limit(1000).Find(&roleDataItems)
	app.CheckError(re.Error)

	// 角色关联的数据节点 ids
	dataItemIds := []int64{}
	for _, roleDataItem := range roleDataItems {
		dataItemIds = append(dataItemIds, roleDataItem.DataItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		app.CheckError(err)
		newDataItemIds := utils.ToSliceInt64(c.Request.PostForm["data_item_id[]"])
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, roleDataItem := range roleDataItems {
				if !utils.InSlice(roleDataItem.DataItemId, newDataItemIds) {
					if err := db.Delete(roleDataItem).Error; err != nil {
						return err
					}
				}
			}
			// 新增
			for _, dataItemId := range newDataItemIds {
				if !utils.InSlice(dataItemId, dataItemIds) {
					roleDataItem := models.NewRoleDataItem(map[string]interface{}{
						"role_id": roleService.RoleId,
						"service_id": roleService.ServiceId,
						"data_item_id": dataItemId,
					})
					roleDataItem.SetDb(db)
					if !(roleDataItem.Validate() && roleDataItem.Create()) {
						return errors.New(roleDataItem.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			cache.SetRoleDataItemIds(roleService.RoleId, roleService.ServiceId, newDataItemIds)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
		}
	}

	this.HTML(c, "./views/admin/role-data-item/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": role.Name,
		"subTitle": roleService.GetServiceName(),
		"dataItems": dataItems,
		"dataItemIds": dataItemIds,
	})
}