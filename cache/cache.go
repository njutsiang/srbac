package cache

import (
	"context"
	"fmt"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
)

var ctx = context.Background()

// 将接口节点保存到缓存
func SetApiItem(apiItem *models.ApiItem) {
	if apiItem.GetService() == nil {
		models.ApiItemsLoadServices([]*models.ApiItem{apiItem})
		if apiItem.GetService() == nil {
			return
		}
	}
	old := apiItem.GetOld()
	if old.Method != "" && old.Uri != "" && (old.Method != apiItem.Method || old.Uri != apiItem.Uri) {
		delApiItem(apiItem.GetService().Id, old.Method, old.Uri)
	}
	key := fmt.Sprintf("auth:service:%d:apis", apiItem.GetService().Id)
	field := fmt.Sprintf("%s%s", apiItem.Method, apiItem.Uri)
	value := "1"
	if apiItem.IsAnonymousAccess == 1 {
		value = "0"
	}
	_, err := srbac.Rdb.HSet(ctx, key, field, value).Result()
	srbac.CheckError(err)
}

// 将接口节点从缓存中删除
func DelApiItem(apiItem *models.ApiItem) {
	if apiItem.GetService() == nil {
		models.ApiItemsLoadServices([]*models.ApiItem{apiItem})
		if apiItem.GetService() == nil {
			return
		}
	}
	delApiItem(apiItem.GetService().Id, apiItem.Method, apiItem.Uri)
}

// 将接口节点从缓存中删除
func delApiItem(serviceId int64, method string, uri string) {
	key := fmt.Sprintf("auth:service:%d:apis", serviceId)
	field := fmt.Sprintf("%s%s", method, uri)
	_, err := srbac.Rdb.HDel(ctx, key, field).Result()
	srbac.CheckError(err)
}

// 将服务下的所有接口节点从缓存中删除
func DelService(serviceId int64) {
	key := fmt.Sprintf("auth:service:%d:apis", serviceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
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
	re := srbac.Db.Where("user_id = ?", userId).Limit(1000).Find(&userRoles)
	srbac.CheckError(re.Error)
	values := []string{}
	for _, userRole := range userRoles {
		values = append(values, utils.ToString(userRole.RoleId))
	}
	setUserRoleIds(userId, values)
}

// 将用户拥有的角色保存到缓存
func setUserRoleIds(userId int64, values []string) {
	key := fmt.Sprintf("auth:user:%d:roles", userId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
	_, err = srbac.Rdb.SAdd(ctx, key, values).Result()
	srbac.CheckError(err)
}

// 将用户与角色的关系从缓存中删除
func DelUserRoles(userId int64) {
	key := fmt.Sprintf("auth:user:%d:roles", userId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
}

// 将角色拥有的接口节点保存到缓存
func SetRoleApiItemIds(roleId int64, serviceId int64, apiItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:role:%d:service:%d:apis", roleId, serviceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
	if len(apiItemIds) >= 1 {
		apiItems := []*models.ApiItem{}
		re := srbac.Db.Where("id IN ?", apiItemIds).Where("service_id = ?", serviceId).Find(&apiItems)
		srbac.CheckError(re.Error)
		values := []string{}
		for _, apiItem := range apiItems {
			values = append(values, fmt.Sprintf("%s%s", apiItem.Method, apiItem.Uri))
		}
		_, err = srbac.Rdb.SAdd(ctx, key, values).Result()
		srbac.CheckError(err)
	}
}

// 将角色和接口节点的关系从缓存中删除
func DelRoleApiItems(roleId int64) {
	roleServices := []*models.RoleService{}
	re := srbac.Db.Where("role_id = ?", roleId).Find(&roleServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:apis", roleId, roleService.ServiceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将角色和接口节点的关系从缓存中删除
func DelRoleApiItemsByServiceId(serviceId int64) {
	roleServices := []*models.RoleService{}
	re := srbac.Db.Where("service_id = ?", serviceId).Find(&roleServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:apis", roleService.RoleId, serviceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将角色和接口节点的关系从缓存中删除
func DelRoleApiItemsByRoleService(roleService *models.RoleService) {
	key := fmt.Sprintf("auth:role:%d:service:%d:apis", roleService.RoleId, roleService.ServiceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
}

// 将角色拥有的数据节点保存到缓存
func SetRoleDataItemIds(roleId int64, serviceId int64, dataItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:role:%d:service:%d:items", roleId, serviceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
	if len(dataItemIds) >= 1 {
		dataItems := []*models.DataItem{}
		re := srbac.Db.Where("id IN ?", dataItemIds).Where("service_id = ?", serviceId).Find(&dataItems)
		srbac.CheckError(re.Error)
		values := []string{}
		for _, dataItem := range dataItems {
			values = append(values, dataItem.Key)
		}
		_, err = srbac.Rdb.SAdd(ctx, key, values).Result()
		srbac.CheckError(err)
	}
}

// 将角色和数据节点的关系从缓存中删除
func DelRoleDataItems(roleId int64) {
	roleServices := []*models.RoleService{}
	re := srbac.Db.Where("role_id = ?", roleId).Find(&roleServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:items", roleId, roleService.ServiceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将角色和数据节点的关系从缓存中删除
func DelRoleDataItemsByServiceId(serviceId int64) {
	roleServices := []*models.RoleService{}
	re := srbac.Db.Where("service_id = ?", serviceId).Find(&roleServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:items", roleService.RoleId, serviceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将角色和数据节点的关系从缓存中删除
func DelRoleDataItemsByRoleService(roleService *models.RoleService) {
	key := fmt.Sprintf("auth:role:%d:service:%d:items", roleService.RoleId, roleService.ServiceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
}

// 将角色拥有的菜单节点保存到缓存
func SetRoleMenuItemIds(roleId int64, serviceId int64, menuItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:role:%d:service:%d:menus", roleId, serviceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
	if len(menuItemIds) >= 1 {
		menuItems := []*models.MenuItem{}
		re := srbac.Db.Where("id IN ?", menuItemIds).Where("service_id = ?", serviceId).Find(&menuItems)
		srbac.CheckError(re.Error)
		values := []string{}
		for _, menuItem := range menuItems {
			values = append(values, menuItem.Key)
		}
		_, err = srbac.Rdb.SAdd(ctx, key, values).Result()
		srbac.CheckError(err)
	}
}

// 将角色和菜单节点的关系从缓存中删除
func DelRoleMenuItems(roleId int64) {
	roleServices := []*models.RoleService{}
	re := srbac.Db.Where("role_id = ?", roleId).Find(&roleServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:menus", roleId, roleService.ServiceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将角色和菜单节点的关系从缓存中删除
func DelRoleMenuItemsByServiceId(serviceId int64) {
	roleServices := []*models.RoleService{}
	re := srbac.Db.Where("service_id = ?", serviceId).Find(&roleServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, roleService := range roleServices {
		keys = append(keys, fmt.Sprintf("auth:role:%d:service:%d:menus", roleService.RoleId, serviceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将角色和菜单节点的关系从缓存中删除
func DelRoleMenuItemsByRoleService(roleService *models.RoleService) {
	key := fmt.Sprintf("auth:role:%d:service:%d:menus", roleService.RoleId, roleService.ServiceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
}

// 将用户拥有的接口节点保存到缓存
func SetUserApiItemIds(userId int64, serviceId int64, apiItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:user:%d:service:%d:apis", userId, serviceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
	if len(apiItemIds) >= 1 {
		apiItems := []*models.ApiItem{}
		re := srbac.Db.Where("id IN ?", apiItemIds).Where("service_id = ?", serviceId).Find(&apiItems)
		srbac.CheckError(re.Error)
		values := []string{}
		for _, apiItem := range apiItems {
			values = append(values, fmt.Sprintf("%s%s", apiItem.Method, apiItem.Uri))
		}
		_, err = srbac.Rdb.SAdd(ctx, key, values).Result()
		srbac.CheckError(err)
	}
}

// 将用户和接口节点的关系从缓存中删除
func DelUserApiItems(userId int64) {
	userServices := []*models.UserService{}
	re := srbac.Db.Where("user_id = ?", userId).Find(&userServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:apis", userId, userService.ServiceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将用户和接口节点的关系从缓存中删除
func DelUserApiItemsByServiceId(serviceId int64) {
	userServices := []*models.UserService{}
	re := srbac.Db.Where("service_id = ?", serviceId).Find(&userServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:apis", userService.UserId, serviceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将用户和接口节点的关系从缓存中删除
func DelUserApiItemsByUserService(userService *models.UserService) {
	key := fmt.Sprintf("auth:user:%d:service:%d:apis", userService.UserId, userService.ServiceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
}

// 将用户拥有的数据节点保存到缓存
func SetUserDataItemIds(userId int64, serviceId int64, dataItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:user:%d:service:%d:items", userId, serviceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
	if len(dataItemIds) >= 1 {
		dataItems := []*models.DataItem{}
		re := srbac.Db.Where("id IN ?", dataItemIds).Where("service_id = ?", serviceId).Find(&dataItems)
		srbac.CheckError(re.Error)
		values := []string{}
		for _, dataItem := range dataItems {
			values = append(values, dataItem.Key)
		}
		_, err = srbac.Rdb.SAdd(ctx, key, values).Result()
		srbac.CheckError(err)
	}
}

// 将用户和数据节点的关系从缓存中删除
func DelUserDataItems(userId int64) {
	userServices := []*models.UserService{}
	re := srbac.Db.Where("user_id = ?", userId).Find(&userServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:items", userId, userService.ServiceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将用户和数据节点的关系从缓存中删除
func DelUserDataItemsByServiceId(serviceId int64) {
	userServices := []*models.UserService{}
	re := srbac.Db.Where("service_id = ?", serviceId).Find(&userServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:items", userService.UserId, serviceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将用户和数据节点的关系从缓存中删除
func DelUserDataItemsByUserService(userService *models.UserService) {
	key := fmt.Sprintf("auth:user:%d:service:%d:items", userService.UserId, userService.ServiceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
}

// 将用户拥有的菜单节点保存到缓存
func SetUserMenuItemIds(userId int64, serviceId int64, menuItemIds []int64) {
	if serviceId == 0 {
		return
	}
	key := fmt.Sprintf("auth:user:%d:service:%d:menus", userId, serviceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
	if len(menuItemIds) >= 1 {
		menuItems := []*models.MenuItem{}
		re := srbac.Db.Where("id IN ?", menuItemIds).Where("service_id = ?", serviceId).Find(&menuItems)
		srbac.CheckError(re.Error)
		values := []string{}
		for _, menuItem := range menuItems {
			values = append(values, menuItem.Key)
		}
		_, err = srbac.Rdb.SAdd(ctx, key, values).Result()
		srbac.CheckError(err)
	}
}

// 将用户和菜单节点的关系从缓存中删除
func DelUserMenuItems(userId int64) {
	userServices := []*models.UserService{}
	re := srbac.Db.Where("user_id = ?", userId).Find(&userServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:menus", userId, userService.ServiceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将用户和菜单节点的关系从缓存中删除
func DelUserMenuItemsByServiceId(serviceId int64) {
	userServices := []*models.UserService{}
	re := srbac.Db.Where("service_id = ?", serviceId).Find(&userServices)
	srbac.CheckError(re.Error)
	keys := []string{}
	for _, userService := range userServices {
		keys = append(keys, fmt.Sprintf("auth:user:%d:service:%d:menus", userService.UserId, serviceId))
	}
	_, err := srbac.Rdb.Del(ctx, keys...).Result()
	srbac.CheckError(err)
}

// 将用户和菜单节点的关系从缓存中删除
func DelUserMenuItemsByUserService(userService *models.UserService) {
	key := fmt.Sprintf("auth:user:%d:service:%d:menus", userService.UserId, userService.ServiceId)
	_, err := srbac.Rdb.Del(ctx, key).Result()
	srbac.CheckError(err)
}