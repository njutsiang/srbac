package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/app/cache"
	"srbac/controllers"
	"srbac/exception"
	"srbac/models"
	"srbac/utils"
)

// 角色的菜单权限
type RoleMenuItemController struct {
	controllers.Controller
}

// 编辑角色的菜单权限
func (this *RoleMenuItemController) Edit(c *gin.Context) {
	roleId := utils.ToInt(c.Query("roleId"))
	roleServiceId := utils.ToInt(c.Query("roleServiceId"))
	if roleServiceId <= 0 {
		exception.Throw(exception.ParamsError)
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

	menuItems := []*models.MenuItem{}
	re = app.Db.Where("service_id = ?", roleService.ServiceId).Order("`key` ASC").Limit(1000).Find(&menuItems)
	app.CheckError(re.Error)

	// 角色和菜单节点的关联
	roleMenuItems := []*models.RoleMenuItem{}
	re = app.Db.Where("role_id = ? AND service_id = ?", roleService.RoleId, roleService.ServiceId).Limit(1000).Find(&roleMenuItems)
	app.CheckError(re.Error)

	// 角色关联的菜单节点 ids
	menuItemIds := []int64{}
	for _, roleMenuItem := range roleMenuItems {
		menuItemIds = append(menuItemIds, roleMenuItem.MenuItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		app.CheckError(err)
		newMenuItemIds := utils.ToSliceInt64(c.Request.PostForm["menu_item_id[]"])
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, roleMenuItem := range roleMenuItems {
				if !utils.InSlice(roleMenuItem.MenuItemId, newMenuItemIds) {
					if err := db.Delete(roleMenuItem).Error; err != nil {
						return err
					}
				}
			}
			// 新增
			for _, menuItemId := range newMenuItemIds {
				if !utils.InSlice(menuItemId, menuItemIds) {
					roleMenuItem := models.NewRoleMenuItem(map[string]interface{}{
						"role_id": roleService.RoleId,
						"service_id": roleService.ServiceId,
						"menu_item_id": menuItemId,
					})
					roleMenuItem.SetDb(db)
					if !(roleMenuItem.Validate() && roleMenuItem.Create()) {
						return errors.New(roleMenuItem.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			cache.SetRoleMenuItemIds(roleService.RoleId, roleService.ServiceId, newMenuItemIds)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
		}
	}

	this.HTML(c, "./views/admin/role-menu-item/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": role.Name,
		"subTitle": roleService.GetServiceName(),
		"menuItems": menuItems,
		"menuItemIds": menuItemIds,
	})
}