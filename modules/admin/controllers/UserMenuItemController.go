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

// 用户的菜单权限
type UserMenuItemController struct {
	controllers.Controller
}

// 编辑用户的菜单权限
func (this *UserMenuItemController) Edit(c *gin.Context) {
	userServiceId := utils.ToInt(c.Query("userServiceId"))
	if userServiceId <= 0 {
		exception.NewException(code.ParamsError)
	}

	userService := &models.UserService{}
	re := srbac.Db.First(userService, userServiceId)
	srbac.CheckError(re.Error)

	models.UserServicesLoadServices([]*models.UserService{userService})

	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-service/list?userId=%d", userService.UserId))

	user := &models.User{}
	re = srbac.Db.First(user, userService.UserId)
	srbac.CheckError(re.Error)

	menuItems := []*models.MenuItem{}
	re = srbac.Db.Where("service_id = ?", userService.ServiceId).Order("id asc").Limit(1000).Find(&menuItems)
	srbac.CheckError(re.Error)

	userMenuItems := []*models.UserMenuItem{}
	re = srbac.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Limit(1000).Find(&userMenuItems)
	srbac.CheckError(re.Error)

	menuItemIds := []int64{}
	for _, userMenuItem := range userMenuItems {
		menuItemIds = append(menuItemIds, userMenuItem.MenuItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		srbac.CheckError(err)
		newMenuItemIds := utils.ToSliceInt64(c.Request.PostForm["menu_item_id[]"])

		// 删除
		for _, userMenuItem := range userMenuItems {
			if !utils.InSlice(userMenuItem.MenuItemId, newMenuItemIds) {
				srbac.Db.Delete(userMenuItem)
			}
		}

		// 新增
		hasErr := false
		for _, menuItemId := range newMenuItemIds {
			if !utils.InSlice(menuItemId, menuItemIds) {
				userMenuItem := models.NewUserMenuItem(map[string]interface{}{
					"user_id": userService.UserId,
					"service_id": userService.ServiceId,
					"menu_item_id": menuItemId,
				})
				if !(userMenuItem.Validate() && userMenuItem.Create()) {
					hasErr = true
					this.SetFailed(c, userMenuItem.GetError())
					break
				}
			}
		}
		if !hasErr {
			this.Redirect(c, referer)
		}
	}

	this.HTML(c, "./views/admin/user-menu-item/edit.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name + " > " + userService.GetService().Name,
		"user": user,
		"menuItems": menuItems,
		"menuItemIds": menuItemIds,
	})
}