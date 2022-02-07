package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
)

// 角色的数据权限
type RoleDataItemController struct {
	controllers.Controller
}

// 编辑角色的数据权限
func (this *RoleDataItemController) Edit(c *gin.Context) {
	roleServiceId := utils.ToInt(c.Query("roleServiceId"))
	if roleServiceId <= 0 {
		exception.NewException(code.ParamsError)
	}

	roleService := &models.RoleService{}
	re := srbac.Db.First(roleService, roleServiceId)
	srbac.CheckError(re.Error)

	models.RoleServiceLoadServices([]*models.RoleService{roleService})

	referer := this.GetReferer(c, fmt.Sprintf("/admin/role-service/list?roleId=%d", roleService.RoleId))

	role := &models.Role{}
	re = srbac.Db.First(role, roleService.RoleId)
	srbac.CheckError(re.Error)

	dataItems := []*models.DataItem{}
	re = srbac.Db.Where("service_id = ?", roleService.ServiceId).Order("`key` ASC").Limit(1000).Find(&dataItems)
	srbac.CheckError(re.Error)

	// 角色和数据节点的关联
	roleDataItems := []*models.RoleDataItem{}
	re = srbac.Db.Where("role_id = ? AND service_id = ?", roleService.RoleId, roleService.ServiceId).Limit(1000).Find(&roleDataItems)
	srbac.CheckError(re.Error)

	// 角色关联的数据节点 ids
	dataItemIds := []int64{}
	for _, roleDataItem := range roleDataItems {
		dataItemIds = append(dataItemIds, roleDataItem.DataItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		srbac.CheckError(err)
		NewDataItemIds := utils.ToSliceInt64(c.Request.PostForm["data_item_id[]"])

		// 删除
		for _, roleDataItem := range roleDataItems {
			if !utils.InSlice(roleDataItem.DataItemId, NewDataItemIds) {
				srbac.Db.Delete(roleDataItem)
			}
		}

		// 新增
		for _, dataItemId := range NewDataItemIds {
			if !utils.InSlice(dataItemId, dataItemIds) {
				roleDataItem := models.NewRoleDataItem(map[string]interface{}{
					"role_id": roleService.RoleId,
					"service_id": roleService.ServiceId,
					"data_item_id": dataItemId,
				})
				if !(roleDataItem.Validate() && roleDataItem.Create()) {
					this.SetFailed(c, roleDataItem.GetError())
					break
				}
			}
		}
		this.Redirect(c, referer)
	}

	this.HTML(c, "./views/admin/role-data-item/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": role.Name + " > " + roleService.GetServiceName(),
		"dataItems": dataItems,
		"dataItemIds": dataItemIds,
	})
}