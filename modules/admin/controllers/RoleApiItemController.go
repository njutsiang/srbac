package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/cache"
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
	roleId := utils.ToInt(c.Query("roleId"))
	roleServiceId := utils.ToInt(c.Query("roleServiceId"))
	if roleId <= 0 || roleServiceId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/role-service/list?roleId=%d", roleId))

	roleService := &models.RoleService{}
	re := srbac.Db.First(roleService, roleServiceId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	srbac.CheckError(re.Error)

	models.RoleServicesLoadServices([]*models.RoleService{roleService})

	role := &models.Role{}
	re = srbac.Db.First(role, roleService.RoleId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
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
		newApiItemIds := utils.ToSliceInt64(c.Request.PostForm["api_item_id[]"])

		// 删除
		for _, roleApiItem := range roleApiItems {
			if !utils.InSlice(roleApiItem.ApiItemId, newApiItemIds) {
				srbac.Db.Delete(roleApiItem)
			}
		}

		// 新增
		hasErr := false
		for _, apiItemId := range newApiItemIds {
			if !utils.InSlice(apiItemId, apiItemIds) {
				roleApiItem := models.NewRoleApiItem(map[string]interface{}{
					"role_id": roleService.RoleId,
					"service_id": roleService.ServiceId,
					"api_item_id": apiItemId,
				})
				if !(roleApiItem.Validate() && roleApiItem.Create()) {
					hasErr = true
					this.SetFailed(c, roleApiItem.GetError())
					break
				}
			}
		}
		cache.SetRoleApiItemIds(roleService.RoleId, roleService.GetServiceId(), newApiItemIds)
		if !hasErr {
			this.Redirect(c, referer)
		}
	}

	this.HTML(c, "./views/admin/role-api-item/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": role.Name + " > " + roleService.GetServiceName(),
		"apiItems": apiItems,
		"apiItemIds": apiItemIds,
	})
}