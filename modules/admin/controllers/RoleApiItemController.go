package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/app/cache"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/logics"
	"srbac/models"
	"srbac/utils"
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

	apiItems := []*models.ApiItem{}
	re = logics.WithApiItemsOrder(app.Db.Where("service_id = ?", roleService.ServiceId)).Limit(1000).Find(&apiItems)
	app.CheckError(re.Error)

	// 角色和接口节点的关联
	roleApiItems := []*models.RoleApiItem{}
	re = app.Db.Where("role_id = ? AND service_id = ?", roleService.RoleId, roleService.ServiceId).Limit(1000).Find(&roleApiItems)
	app.CheckError(re.Error)

	// 角色关联的接口节点 ids
	apiItemIds := []int64{}
	for _, roleApiItem := range roleApiItems {
		apiItemIds = append(apiItemIds, roleApiItem.ApiItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		app.CheckError(err)
		newApiItemIds := utils.ToSliceInt64(c.Request.PostForm["api_item_id[]"])
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, roleApiItem := range roleApiItems {
				if !utils.InSlice(roleApiItem.ApiItemId, newApiItemIds) {
					if err := db.Delete(roleApiItem).Error; err != nil {
						return err
					}
				}
			}
			// 新增
			for _, apiItemId := range newApiItemIds {
				if !utils.InSlice(apiItemId, apiItemIds) {
					roleApiItem := models.NewRoleApiItem(map[string]interface{}{
						"role_id": roleService.RoleId,
						"service_id": roleService.ServiceId,
						"api_item_id": apiItemId,
					})
					roleApiItem.SetDb(db)
					if !(roleApiItem.Validate() && roleApiItem.Create()) {
						return errors.New(roleApiItem.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			cache.SetRoleApiItemIds(roleService.RoleId, roleService.GetServiceId(), newApiItemIds)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
		}
	}

	this.HTML(c, "./views/admin/role-api-item/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": role.Name,
		"subTitle": roleService.GetServiceName(),
		"apiItems": apiItems,
		"apiItemIds": apiItemIds,
	})
}