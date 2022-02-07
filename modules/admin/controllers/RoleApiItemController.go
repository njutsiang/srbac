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

// 角色的接口权限
type RoleApiItemController struct {
	controllers.Controller
}

// 编辑角色的接口权限
func (this *RoleApiItemController) Edit(c *gin.Context) {
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

	apiItems := []*models.ApiItem{}
	re = srbac.Db.Where("service_id = ?", roleService.ServiceId).Order("uri asc").Limit(1000).Find(&apiItems)
	srbac.CheckError(re.Error)

	// 角色和接口节点的关联
	roleApiItems := []*models.RoleApiItem{}
	re = srbac.Db.Where("role_id = ? AND service_id = ?", roleService.RoleId, roleService.ServiceId).Limit(1000).Find(&roleApiItems)
	srbac.CheckError(re.Error)

	// 角色关联的接口节点 ids
	apiItemIds := []int64{}
	for _, roleApiItem := range roleApiItems {
		apiItemIds = append(apiItemIds, roleApiItem.ApiItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		srbac.CheckError(err)
		NewApiItemIds := utils.ToSliceInt64(c.Request.PostForm["api_item_id[]"])

		// 删除
		for _, roleApiItem := range roleApiItems {
			if !utils.InSlice(roleApiItem.ApiItemId, NewApiItemIds) {
				srbac.Db.Delete(roleApiItem)
			}
		}

		// 新增
		for _, apiItemId := range NewApiItemIds {
			if !utils.InSlice(apiItemId, apiItemIds) {
				roleApiItem := models.NewRoleApiItem(map[string]interface{}{
					"role_id": roleService.RoleId,
					"service_id": roleService.ServiceId,
					"api_item_id": apiItemId,
				})
				if !(roleApiItem.Validate() && roleApiItem.Create()) {
					this.SetFailed(c, roleApiItem.GetError())
					break
				}
			}
		}
		this.Redirect(c, referer)
	}

	this.HTML(c, "./views/admin/role-api-item/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": role.Name + " > " + roleService.GetServiceName(),
		"apiItems": apiItems,
		"apiItemIds": apiItemIds,
	})
}