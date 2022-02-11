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

// 用户的数据权限
type UserDataItemController struct {
	controllers.Controller
}

// 编辑用户的数据权限
func (this *UserDataItemController) Edit(c *gin.Context) {
	userId := utils.ToInt(c.Query("userId"))
	userServiceId := utils.ToInt(c.Query("userServiceId"))
	if userServiceId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-service/list?userId=%d", userId))

	userService := &models.UserService{}
	re := srbac.Db.First(userService, userServiceId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	srbac.CheckError(re.Error)

	models.UserServicesLoadServices([]*models.UserService{userService})

	user := &models.User{}
	re = srbac.Db.First(user, userService.UserId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	srbac.CheckError(re.Error)

	dataItems := []*models.DataItem{}
	re = srbac.Db.Where("service_id = ?", userService.ServiceId).Order("id asc").Limit(1000).Find(&dataItems)
	srbac.CheckError(re.Error)

	userDataItems := []*models.UserDataItem{}
	re = srbac.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Limit(1000).Find(&userDataItems)
	srbac.CheckError(re.Error)

	dataItemIds := []int64{}
	for _, userDataItem := range userDataItems {
		dataItemIds = append(dataItemIds, userDataItem.DataItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		srbac.CheckError(err)
		newDataItemIds := utils.ToSliceInt64(c.Request.PostForm["data_item_id[]"])

		// 删除
		for _, userDataItem := range userDataItems {
			if !utils.InSlice(userDataItem.DataItemId, newDataItemIds) {
				srbac.Db.Delete(userDataItem)
			}
		}

		// 新增
		hasErr := false
		for _, dataItemId := range newDataItemIds {
			if !utils.InSlice(dataItemId, dataItemIds) {
				userDataItem := models.NewUserDataItem(map[string]interface{}{
					"user_id": userService.UserId,
					"service_id": userService.ServiceId,
					"data_item_id": dataItemId,
				})
				if !(userDataItem.Validate() && userDataItem.Create()) {
					hasErr = true
					this.SetFailed(c, userDataItem.GetError())
					break
				}
			}
		}
		cache.SetUserDataItemIds(userService.UserId, userService.ServiceId, newDataItemIds)
		if !hasErr {
			this.Redirect(c, referer)
		}
	}

	this.HTML(c, "./views/admin/user-data-item/edit.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name + " > " + userService.GetService().Name,
		"user": user,
		"dataItems": dataItems,
		"dataItemIds": dataItemIds,
	})
}