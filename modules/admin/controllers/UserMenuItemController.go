package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/app/utils"
	"srbac/cache"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/models"
)

// 用户的菜单权限
type UserMenuItemController struct {
	controllers.Controller
}

// 编辑用户的菜单权限
func (this *UserMenuItemController) Edit(c *gin.Context) {
	userId := utils.ToInt(c.Query("userId"))
	userServiceId := utils.ToInt(c.Query("userServiceId"))
	if userServiceId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-service/list?userId=%d", userId))

	userService := &models.UserService{}
	re := app.Db.First(userService, userServiceId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	models.UserServicesLoadServices([]*models.UserService{userService})

	user := &models.User{}
	re = app.Db.First(user, userService.UserId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	menuItems := []*models.MenuItem{}
	re = app.Db.Where("service_id = ?", userService.ServiceId).Order("id asc").Limit(1000).Find(&menuItems)
	app.CheckError(re.Error)

	userMenuItems := []*models.UserMenuItem{}
	re = app.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Limit(1000).Find(&userMenuItems)
	app.CheckError(re.Error)

	menuItemIds := []int64{}
	for _, userMenuItem := range userMenuItems {
		menuItemIds = append(menuItemIds, userMenuItem.MenuItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		app.CheckError(err)
		newMenuItemIds := utils.ToSliceInt64(c.Request.PostForm["menu_item_id[]"])
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, userMenuItem := range userMenuItems {
				if !utils.InSlice(userMenuItem.MenuItemId, newMenuItemIds) {
					if err := db.Delete(userMenuItem).Error; err != nil {
						return err
					}
				}
			}
			// 新增
			for _, menuItemId := range newMenuItemIds {
				if !utils.InSlice(menuItemId, menuItemIds) {
					userMenuItem := models.NewUserMenuItem(map[string]interface{}{
						"user_id": userService.UserId,
						"service_id": userService.ServiceId,
						"menu_item_id": menuItemId,
					})
					userMenuItem.SetDb(db)
					if !(userMenuItem.Validate() && userMenuItem.Create()) {
						return errors.New(userMenuItem.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			cache.SetUserMenuItemIds(userService.UserId, userService.ServiceId, newMenuItemIds)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
		}
	}

	this.HTML(c, "./views/admin/user-menu-item/edit.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name,
		"subTitle": userService.GetService().Name,
		"user": user,
		"menuItems": menuItems,
		"menuItemIds": menuItemIds,
	})
}