package admin

import (
	"context"
	"github.com/gin-gonic/gin"
	"srbac/app"
	"srbac/cache"
	"srbac/controllers"
	"srbac/libraries/log"
	"srbac/models"
)

var ctx = context.Background()

type SystemController struct {
	controllers.Controller
}

// 重建所有缓存
func (this *SystemController) RebuildCache(c *gin.Context) {
	referer := this.GetReferer(c, "/admin", false)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()

		// 删除所有缓存
		keys := []string{}
		cursor := uint64(0)
		var err error
		for {
			keys, cursor, err = app.Rdb.Scan(ctx, cursor, "auth:*", 100).Result()
			app.CheckError(err)

			if len(keys) >= 1 {
				_, err = app.Rdb.Del(ctx, keys...).Result()
				app.CheckError(err)
			}
			if cursor == 0 {
				break
			}
		}

		services := []*models.Service{}
		re := app.Db.Limit(1000).Find(&services)
		app.CheckError(re.Error)
		for _, service := range services {
			cache.SetService(service)
			apiItems := []*models.ApiItem{}
			re = app.Db.Where("service_id = ?", service.Id).Find(&apiItems)
			app.CheckError(re.Error)
			for _, apiItem := range apiItems {
				cache.SetApiItem(apiItem)
			}
		}

		page := 1
		perPage := 1000
		roleServices := []*models.RoleService{}
		for {
			re = app.Db.Order("id ASC").Offset((page - 1) * perPage).Limit(perPage).Find(&roleServices)
			app.CheckError(re.Error)
			if len(roleServices) == 0 {
				break
			}
			page++
			for _, roleService := range roleServices {
				roleApiItems := []*models.RoleApiItem{}
				re = app.Db.Where("role_id = ? AND service_id = ?", roleService.RoleId, roleService.ServiceId).Find(&roleApiItems)
				app.CheckError(re.Error)
				if len(roleApiItems) >= 1 {
					apiItemIds := []int64{}
					for _, roleApiItem := range roleApiItems {
						apiItemIds = append(apiItemIds, roleApiItem.ApiItemId)
					}
					cache.SetRoleApiItemIds(roleService.RoleId, roleService.ServiceId, apiItemIds)
				}

				roleDataItems := []*models.RoleDataItem{}
				re = app.Db.Where("role_id = ? AND service_id = ?", roleService.RoleId, roleService.ServiceId).Find(&roleDataItems)
				app.CheckError(re.Error)
				if len(roleDataItems) >= 1 {
					dataItemIds := []int64{}
					for _, roleDataItem := range roleDataItems {
						dataItemIds = append(dataItemIds, roleDataItem.DataItemId)
					}
					cache.SetRoleDataItemIds(roleService.RoleId, roleService.ServiceId, dataItemIds)
				}

				roleMenuItems := []*models.RoleMenuItem{}
				re = app.Db.Where("role_id = ? AND service_id = ?", roleService.RoleId, roleService.ServiceId).Find(&roleMenuItems)
				app.CheckError(re.Error)
				if len(roleMenuItems) >= 1 {
					menuItemIds := []int64{}
					for _, roleMenuItem := range roleMenuItems {
						menuItemIds = append(menuItemIds, roleMenuItem.MenuItemId)
					}
					cache.SetRoleMenuItemIds(roleService.RoleId, roleService.ServiceId, menuItemIds)
				}
			}
		}

		//page = 1
		//userRoles := []*models.UserRole{}
		//for {
		//	re = srbac.Db.Order("user_id ASC").Offset((page - 1) * perPage).Limit(perPage).Find(&userRoles)
		//	srbac.CheckError(re.Error)
		//	if len(userRoles) == 0 {
		//
		//	}
		//}
		//page = 1
		//userServices := []*models.UserService{}
		//for {
		//	re = srbac.Db.Order("id ASC").Offset((page - 1) * perPage).Limit(perPage).Find(&userServices)
		//	srbac.CheckError(re.Error)
		//	if len(userServices) == 0 {
		//		break
		//	}
		//	page++
		//	for _, userService := range userServices {
		//		userApiItems := []*models.UserApiItem{}
		//		re = srbac.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Find(&userApiItems)
		//		srbac.CheckError(re.Error)
		//		if len(userApiItems) >= 1 {
		//			apiItemIds := []int64{}
		//			for _, userApiItem := range userApiItems {
		//				apiItemIds = append(apiItemIds, userApiItem.ApiItemId)
		//			}
		//			cache.SetUserApiItemIds(userService.UserId, userService.ServiceId, apiItemIds)
		//		}
		//
		//		userDataItems := []*models.UserDataItem{}
		//		re = srbac.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Find(&userDataItems)
		//		srbac.CheckError(re.Error)
		//		if len(userDataItems) >= 1 {
		//			dataItemIds := []int64{}
		//			for _, userDataItem := range userDataItems {
		//				dataItemIds = append(dataItemIds, userDataItem.DataItemId)
		//			}
		//			cache.SetUserDataItemIds(userService.UserId, userService.ServiceId, dataItemIds)
		//		}
		//
		//		userMenuItems := []*models.UserMenuItem{}
		//		re = srbac.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Find(&userMenuItems)
		//		srbac.CheckError(re.Error)
		//		if len(userMenuItems) >= 1 {
		//			menuItemIds := []int64{}
		//			for _, userMenuItem := range userMenuItems {
		//				menuItemIds = append(menuItemIds, userMenuItem.MenuItemId)
		//			}
		//			cache.SetUserMenuItemIds(userService.UserId, userService.ServiceId, menuItemIds)
		//		}
		//	}
		//}
	}()
	this.Redirect(c, referer)
}