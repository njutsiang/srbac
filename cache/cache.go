package cache

import (
	"context"
	"fmt"
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