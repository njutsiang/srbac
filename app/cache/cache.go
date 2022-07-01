package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"srbac/app"
	"srbac/models"
	"srbac/utils"
	"time"
)

var ctx = context.Background()

// 将服务保存到缓存
func SetService(service *models.Service) {
	old := service.GetOld()
	if old.Key != "" && old.Key != service.Key {
		key1 := fmt.Sprintf("auth:service:%s", old.Key)
		_, err := app.Rdb.Del(ctx, key1).Result()
		app.CheckError(err)
	}
	key2 := fmt.Sprintf("auth:service:%s", service.Key)
	value := fmt.Sprintf("%d", service.Id)
	_, err := app.Rdb.Set(ctx, key2, value, 0).Result()
	app.CheckError(err)
}

// 将服务保存到缓存
// 将服务下的所有接口节点从缓存中删除
func DelService(service *models.Service) {
	key1 := fmt.Sprintf("auth:service:%s", service.Key)
	key2 := fmt.Sprintf("auth:service:%d:apis", service.Id)
	_, err := app.Rdb.Del(ctx, key1, key2).Result()
	app.CheckError(err)
}

// 将接口节点保存到缓存
func SetApiItem(apiItem *models.ApiItem) {
	old := apiItem.GetOld()
	if old.Method != "" && old.Uri != "" && (old.Method != apiItem.Method || old.Uri != apiItem.Uri) {
		delApiItem(apiItem.ServiceId, old.Method, old.Uri)
	}
	key := fmt.Sprintf("auth:service:%d:apis", apiItem.ServiceId)
	field := fmt.Sprintf("%s%s", apiItem.Method, apiItem.Uri)
	value := "1"
	if apiItem.IsAnonymousAccess == 1 {
		value = "0"
	}
	_, err := app.Rdb.HSet(ctx, key, field, value).Result()
	app.CheckError(err)
}

// 将接口节点从缓存中删除
func DelApiItem(apiItem *models.ApiItem) {
	delApiItem(apiItem.ServiceId, apiItem.Method, apiItem.Uri)
}

// 将接口节点从缓存中删除
func delApiItem(serviceId int64, method string, uri string) {
	key := fmt.Sprintf("auth:service:%d:apis", serviceId)
	field := fmt.Sprintf("%s%s", method, uri)
	_, err := app.Rdb.HDel(ctx, key, field).Result()
	app.CheckError(err)
}

// 将用户拥有的角色保存到缓存
func SetUserRoleIds(userId int64, roleIds []int64) {
	values := []string{}
	for _, roleId := range roleIds {
		values = append(values, utils.ToString(roleId))
	}
	setUserRoleIds(userId, values)
}

// 将用户拥有的角色保存到缓存
func SetUserRoles(userId int64) {
	userRoles := []*models.UserRole{}
	re := app.Db.Where("user_id = ?", userId).Limit(1000).Find(&userRoles)
	app.CheckError(re.Error)
	values := []string{}
	for _, userRole := range userRoles {
		values = append(values, utils.ToString(userRole.RoleId))
	}
	setUserRoleIds(userId, values)
}

// 将用户拥有的角色保存到缓存
func setUserRoleIds(userId int64, values []string) {
	key := fmt.Sprintf("auth:user:%d:roles", userId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
	if len(values) >= 1 {
		_, err = app.Rdb.SAdd(ctx, key, values).Result()
		app.CheckError(err)
	}
}

// 将用户与角色的关系从缓存中删除
func DelUserRoles(userId int64) {
	key := fmt.Sprintf("auth:user:%d:roles", userId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
}

// 将角色拥有的接口节点保存到缓存
func SetRoleApiItemIds(roleId int64, serviceId int64, apiItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:role:%d:service:%d:apis", roleId, serviceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
	if len(apiItemIds) >= 1 {
		apiItems := []*models.ApiItem{}
		re := app.Db.Where("id IN ?", apiItemIds).Where("service_id = ?", serviceId).Find(&apiItems)
		app.CheckError(re.Error)
		values := []string{}
		for _, apiItem := range apiItems {
			values = append(values, fmt.Sprintf("%s%s", apiItem.Method, apiItem.Uri))
		}
		if len(values) >= 1 {
			_, err = app.Rdb.SAdd(ctx, key, values).Result()
			app.CheckError(err)
		}
	}
}

// 将角色拥有的接口节点保存到缓存
func SetRoleApiItemsByApiItem(apiItem *models.ApiItem) {
	old := apiItem.GetOld()
	if old.Method != "" && old.Uri != "" && (old.Method != apiItem.Method || old.Uri != apiItem.Uri) {
		roleApiItems := []*models.RoleApiItem{}
		re := app.Db.Distinct("role_id", "service_id").Where("api_item_id = ?", apiItem.Id).Find(&roleApiItems)
		app.CheckError(re.Error)
		SetRoleApiItemsByRoleApiItems(roleApiItems)
	}
}

// 将角色拥有的接口节点保存到缓存
func SetRoleApiItemsByRoleApiItems(roleApiItems []*models.RoleApiItem) {
	for _, roleApiItem := range roleApiItems {
		currentRoleApiItems := []*models.RoleApiItem{}
		re := app.Db.
			Where("role_id = ?", roleApiItem.RoleId).
			Where("service_id = ?", roleApiItem.ServiceId).
			Find(&currentRoleApiItems)
		app.CheckError(re.Error)
		apiItemIds := []int64{}
		for _, currentRoleApiItem := range currentRoleApiItems {
			apiItemIds = append(apiItemIds, currentRoleApiItem.ApiItemId)
		}
		SetRoleApiItemIds(roleApiItem.RoleId, roleApiItem.ServiceId, apiItemIds)
	}
}

// 将角色和接口节点的关系从缓存中删除
func DelRoleApiItemsByRoleServices(roleServices []*models.RoleService) {
	if len(roleServices) == 0 {
		return
	}
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:apis", roleService.RoleId, roleService.ServiceId))
	}
	_, err := app.Rdb.Del(ctx, keys...).Result()
	app.CheckError(err)
}

// 将角色和接口节点的关系从缓存中删除
func DelRoleApiItemsByRoleService(roleService *models.RoleService) {
	key := fmt.Sprintf("auth:role:%d:service:%d:apis", roleService.RoleId, roleService.ServiceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
}

// 将角色拥有的数据节点保存到缓存
func SetRoleDataItemIds(roleId int64, serviceId int64, dataItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:role:%d:service:%d:items", roleId, serviceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
	if len(dataItemIds) >= 1 {
		dataItems := []*models.DataItem{}
		re := app.Db.Where("id IN ?", dataItemIds).Where("service_id = ?", serviceId).Find(&dataItems)
		app.CheckError(re.Error)
		values := []string{}
		for _, dataItem := range dataItems {
			values = append(values, dataItem.Key)
		}
		if len(values) >= 1 {
			_, err = app.Rdb.SAdd(ctx, key, values).Result()
			app.CheckError(err)
		}
	}
}

// 将角色拥有的数据节点保存到缓存
func SetRoleDataItemsByDataItem(dataItem *models.DataItem) {
	old := dataItem.GetOld()
	if old.Key != "" && old.Key != dataItem.Key {
		roleDataItems := []*models.RoleDataItem{}
		re := app.Db.Distinct("role_id", "service_id").Where("data_item_id = ?", dataItem.Id).Find(&roleDataItems)
		app.CheckError(re.Error)
		SetRoleDataItemsByRoleDataItems(roleDataItems)
	}
}

// 将角色拥有的数据节点保存到缓存
func SetRoleDataItemsByRoleDataItems(roleDataItems []*models.RoleDataItem) {
	for _, roleDataItem := range roleDataItems {
		currentRoleDataItems := []*models.RoleDataItem{}
		re := app.Db.
			Where("role_id = ?", roleDataItem.RoleId).
			Where("service_id = ?", roleDataItem.ServiceId).
			Find(&currentRoleDataItems)
		app.CheckError(re.Error)
		dataItemIds := []int64{}
		for _, currentRoleDataItem := range currentRoleDataItems {
			dataItemIds = append(dataItemIds, currentRoleDataItem.DataItemId)
		}
		SetRoleDataItemIds(roleDataItem.RoleId, roleDataItem.ServiceId, dataItemIds)
	}
}

// 将角色和数据节点的关系从缓存中删除
func DelRoleDataItemsByRoleServices(roleServices []*models.RoleService) {
	if len(roleServices) == 0 {
		return
	}
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:items", roleService.RoleId, roleService.ServiceId))
	}
	_, err := app.Rdb.Del(ctx, keys...).Result()
	app.CheckError(err)
}

// 将角色和数据节点的关系从缓存中删除
func DelRoleDataItemsByRoleService(roleService *models.RoleService) {
	key := fmt.Sprintf("auth:role:%d:service:%d:items", roleService.RoleId, roleService.ServiceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
}

// 将角色拥有的菜单节点保存到缓存
func SetRoleMenuItemIds(roleId int64, serviceId int64, menuItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:role:%d:service:%d:menus", roleId, serviceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
	if len(menuItemIds) >= 1 {
		menuItems := []*models.MenuItem{}
		re := app.Db.Where("id IN ?", menuItemIds).Where("service_id = ?", serviceId).Find(&menuItems)
		app.CheckError(re.Error)
		values := []string{}
		for _, menuItem := range menuItems {
			values = append(values, menuItem.Key)
		}
		if len(values) >= 1 {
			_, err = app.Rdb.SAdd(ctx, key, values).Result()
			app.CheckError(err)
		}
	}
}

// 将角色拥有的菜单节点保存到缓存
func SetRoleMenuItemsByMenuItem(menuItem *models.MenuItem) {
	old := menuItem.GetOld()
	if old.Key != "" && old.Key != menuItem.Key {
		roleMenuItems := []*models.RoleMenuItem{}
		re := app.Db.Distinct("role_id", "service_id").Where("menu_item_id = ?", menuItem.Id).Find(&roleMenuItems)
		app.CheckError(re.Error)
		SetRoleMenuItemsByRoleMenuItems(roleMenuItems)
	}
}

// 将角色拥有的菜单节点保存到缓存
func SetRoleMenuItemsByRoleMenuItems(roleMenuItems []*models.RoleMenuItem) {
	for _, roleMenuItem := range roleMenuItems {
		currentRoleMenuItems := []*models.RoleMenuItem{}
		re := app.Db.
			Where("role_id = ?", roleMenuItem.RoleId).
			Where("service_id = ?", roleMenuItem.ServiceId).
			Find(&currentRoleMenuItems)
		app.CheckError(re.Error)
		menuItemIds := []int64{}
		for _, currentRoleMenuItem := range currentRoleMenuItems {
			menuItemIds = append(menuItemIds, currentRoleMenuItem.MenuItemId)
		}
		SetRoleMenuItemIds(roleMenuItem.RoleId, roleMenuItem.ServiceId, menuItemIds)
	}
}

// 将角色和菜单节点的关系从缓存中删除
func DelRoleMenuItemsByRoleServices(roleServices []*models.RoleService) {
	if len(roleServices) == 0 {
		return
	}
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:menus", roleService.RoleId, roleService.ServiceId))
	}
	_, err := app.Rdb.Del(ctx, keys...).Result()
	app.CheckError(err)
}

// 将角色和菜单节点的关系从缓存中删除
func DelRoleMenuItemsByRoleService(roleService *models.RoleService) {
	key := fmt.Sprintf("auth:role:%d:service:%d:menus", roleService.RoleId, roleService.ServiceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
}

// 将用户拥有的接口节点保存到缓存
func SetUserApiItemIds(userId int64, serviceId int64, apiItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:user:%d:service:%d:apis", userId, serviceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
	if len(apiItemIds) >= 1 {
		apiItems := []*models.ApiItem{}
		re := app.Db.Where("id IN ?", apiItemIds).Where("service_id = ?", serviceId).Find(&apiItems)
		app.CheckError(re.Error)
		values := []string{}
		for _, apiItem := range apiItems {
			values = append(values, fmt.Sprintf("%s%s", apiItem.Method, apiItem.Uri))
		}
		if len(values) >= 1 {
			_, err = app.Rdb.SAdd(ctx, key, values).Result()
			app.CheckError(err)
		}
	}
}

// 将用户拥有的接口节点保存到缓存
func SetUserApiItemsByApiItem(apiItem *models.ApiItem) {
	old := apiItem.GetOld()
	if old.Method != "" && old.Uri != "" && (old.Method != apiItem.Method || old.Uri != apiItem.Uri) {
		userApiItems := []*models.UserApiItem{}
		re := app.Db.Distinct("user_id", "service_id").Where("api_item_id = ?", apiItem.Id).Find(&userApiItems)
		app.CheckError(re.Error)
		SetUserApiItemsByUserApiItems(userApiItems)
	}
}

// 将用户拥有的接口节点保存到缓存
func SetUserApiItemsByUserApiItems(userApiItems []*models.UserApiItem) {
	for _, userApiItem := range userApiItems {
		currentUserApiItems := []*models.UserApiItem{}
		re := app.Db.
			Where("user_id = ?", userApiItem.UserId).
			Where("service_id = ?", userApiItem.ServiceId).
			Find(&currentUserApiItems)
		app.CheckError(re.Error)
		apiItemIds := []int64{}
		for _, currentUserApiItem := range currentUserApiItems {
			apiItemIds = append(apiItemIds, currentUserApiItem.ApiItemId)
		}
		SetUserApiItemIds(userApiItem.UserId, userApiItem.ServiceId, apiItemIds)
	}
}

// 将用户和接口节点的关系从缓存中删除
func DelUserApiItemsByUserServices(userServices []*models.UserService) {
	if len(userServices) == 0 {
		return
	}
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:apis", userService.UserId, userService.ServiceId))
	}
	_, err := app.Rdb.Del(ctx, keys...).Result()
	app.CheckError(err)
}

// 将用户和接口节点的关系从缓存中删除
func DelUserApiItemsByUserService(userService *models.UserService) {
	key := fmt.Sprintf("auth:user:%d:service:%d:apis", userService.UserId, userService.ServiceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
}

// 将用户拥有的数据节点保存到缓存
func SetUserDataItemIds(userId int64, serviceId int64, dataItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:user:%d:service:%d:items", userId, serviceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
	if len(dataItemIds) >= 1 {
		dataItems := []*models.DataItem{}
		re := app.Db.Where("id IN ?", dataItemIds).Where("service_id = ?", serviceId).Find(&dataItems)
		app.CheckError(re.Error)
		values := []string{}
		for _, dataItem := range dataItems {
			values = append(values, dataItem.Key)
		}
		if len(values) >= 1 {
			_, err = app.Rdb.SAdd(ctx, key, values).Result()
			app.CheckError(err)
		}
	}
}

// 将用户拥有的数据节点保存到缓存
func SetUserDataItemsByDataItem(dataItem *models.DataItem) {
	old := dataItem.GetOld()
	if old.Key != "" && old.Key != dataItem.Key {
		userDataItems := []*models.UserDataItem{}
		re := app.Db.Distinct("user_id", "service_id").Where("data_item_id = ?", dataItem.Id).Find(&userDataItems)
		app.CheckError(re.Error)
		SetUserDataItemsByUserDataItems(userDataItems)
	}
}

// 将用户拥有的数据节点保存到缓存
func SetUserDataItemsByUserDataItems(userDataItems []*models.UserDataItem) {
	for _, userDataItem := range userDataItems {
		currentUserDataItems := []*models.UserDataItem{}
		re := app.Db.
			Where("user_id = ?", userDataItem.UserId).
			Where("service_id = ?", userDataItem.ServiceId).
			Find(&currentUserDataItems)
		app.CheckError(re.Error)
		dataItemIds := []int64{}
		for _, currentUserDataItem := range currentUserDataItems {
			dataItemIds = append(dataItemIds, currentUserDataItem.DataItemId)
		}
		SetUserDataItemIds(userDataItem.UserId, userDataItem.ServiceId, dataItemIds)
	}
}

// 将用户和数据节点的关系从缓存中删除
func DelUserDataItemsByUserServices(userServices []*models.UserService) {
	if len(userServices) == 0 {
		return
	}
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:items", userService.UserId, userService.ServiceId))
	}
	_, err := app.Rdb.Del(ctx, keys...).Result()
	app.CheckError(err)
}

// 将用户和数据节点的关系从缓存中删除
func DelUserDataItemsByUserService(userService *models.UserService) {
	key := fmt.Sprintf("auth:user:%d:service:%d:items", userService.UserId, userService.ServiceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
}

// 将用户拥有的菜单节点保存到缓存
func SetUserMenuItemIds(userId int64, serviceId int64, menuItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:user:%d:service:%d:menus", userId, serviceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
	if len(menuItemIds) >= 1 {
		menuItems := []*models.MenuItem{}
		re := app.Db.Where("id IN ?", menuItemIds).Where("service_id = ?", serviceId).Find(&menuItems)
		app.CheckError(re.Error)
		values := []string{}
		for _, menuItem := range menuItems {
			values = append(values, menuItem.Key)
		}
		if len(values) >= 1 {
			_, err = app.Rdb.SAdd(ctx, key, values).Result()
			app.CheckError(err)
		}
	}
}

// 将用户拥有的菜单节点保存到缓存
func SetUserMenuItemsByMenuItem(menuItem *models.MenuItem) {
	old := menuItem.GetOld()
	if old.Key != "" && old.Key != menuItem.Key {
		userMenuItems := []*models.UserMenuItem{}
		re := app.Db.Distinct("user_id", "service_id").Where("menu_item_id = ?", menuItem.Id).Find(&userMenuItems)
		app.CheckError(re.Error)
		SetUserMenuItemsByUserMenuItems(userMenuItems)
	}
}

// 将用户拥有的菜单节点保存到缓存
func SetUserMenuItemsByUserMenuItems(userMenuItems []*models.UserMenuItem) {
	for _, userMenuItem := range userMenuItems {
		currentUserMenuItems := []*models.UserMenuItem{}
		re := app.Db.
			Where("user_id = ?", userMenuItem.UserId).
			Where("service_id = ?", userMenuItem.ServiceId).
			Find(&currentUserMenuItems)
		app.CheckError(re.Error)
		menuItemIds := []int64{}
		for _, currentUserMenuItem := range currentUserMenuItems {
			menuItemIds = append(menuItemIds, currentUserMenuItem.MenuItemId)
		}
		SetUserMenuItemIds(userMenuItem.UserId, userMenuItem.ServiceId, menuItemIds)
	}
}

// 将用户和菜单节点的关系从缓存中删除
func DelUserMenuItemsByUserServices(userServices []*models.UserService) {
	if len(userServices) == 0 {
		return
	}
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:menus", userService.UserId, userService.ServiceId))
	}
	_, err := app.Rdb.Del(ctx, keys...).Result()
	app.CheckError(err)
}

// 将用户和菜单节点的关系从缓存中删除
func DelUserMenuItemsByUserService(userService *models.UserService) {
	key := fmt.Sprintf("auth:user:%d:service:%d:menus", userService.UserId, userService.ServiceId)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
}

// 保存用户 Token 到缓存
func SetUserToken(token string, user *models.User, expiration int) {
	key := fmt.Sprintf("user:token:%s", token)
	value, _ := json.Marshal(map[string]interface{}{
		"id":       user.Id,
		"name":     app.AsciiString(user.Name),
		"username": user.Username,
		"status":   user.Status,
	})
	_, err := app.Rdb.SetEX(ctx, key, value, time.Second * time.Duration(expiration)).Result()
	app.CheckError(err)
}

// 将用户 Token 从缓存中删除
func DelUserToken(token string) {
	key := fmt.Sprintf("user:token:%s", token)
	_, err := app.Rdb.Del(ctx, key).Result()
	app.CheckError(err)
}