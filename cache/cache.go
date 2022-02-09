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
		delApiItem(apiItem.GetService().Key, old.Method, old.Uri)
	}
	key := fmt.Sprintf("auth:service:%s:apis", apiItem.GetService().Key)
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
	delApiItem(apiItem.GetService().Key, apiItem.Method, apiItem.Uri)
}

// 将接口节点从缓存中删除
func delApiItem(service string, method string, uri string) {
	key := fmt.Sprintf("auth:service:%s:apis", service)
	field := fmt.Sprintf("%s%s", method, uri)
	_, err := srbac.Rdb.HDel(ctx, key, field).Result()
	srbac.CheckError(err)
}

// 将服务下的所有接口节点保存到缓存
func SetService(service *models.Service) {
	old := service.GetOld()
	if old.Key != "" && old.Key != service.Key {
		delService(old.Key)

		apiItems := []*models.ApiItem{}
		re := srbac.Db.Where("service_id = ?", service.Id).Limit(1000).Find(&apiItems)
		srbac.CheckError(re.Error)

		key := fmt.Sprintf("auth:service:%s:apis", service.Key)
		values := map[string]string{}
		for _, apiItem := range apiItems {
			field := fmt.Sprintf("%s%s", apiItem.Method, apiItem.Uri)
			value := "1"
			if apiItem.IsAnonymousAccess == 1 {
				value = "0"
			}
			values[field] = value
		}

		_, err := srbac.Rdb.HMSet(ctx, key, values).Result()
		srbac.CheckError(err)
	}
}

// 将服务下的所有接口节点从缓存中删除
func DelService(service *models.Service) {
	delService(service.Key)
}

// 将服务下的所有接口节点从缓存中删除
func delService(service string) {
	key := fmt.Sprintf("auth:service:%s:apis", service)
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