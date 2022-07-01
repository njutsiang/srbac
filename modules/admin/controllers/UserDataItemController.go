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
	"srbac/models"
	"srbac/utils"
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

	dataItems := []*models.DataItem{}
	re = app.Db.Where("service_id = ?", userService.ServiceId).Order("id asc").Limit(1000).Find(&dataItems)
	app.CheckError(re.Error)

	userDataItems := []*models.UserDataItem{}
	re = app.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Limit(1000).Find(&userDataItems)
	app.CheckError(re.Error)

	dataItemIds := []int64{}
	for _, userDataItem := range userDataItems {
		dataItemIds = append(dataItemIds, userDataItem.DataItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		app.CheckError(err)
		newDataItemIds := utils.ToSliceInt64(c.Request.PostForm["data_item_id[]"])
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, userDataItem := range userDataItems {
				if !utils.InSlice(userDataItem.DataItemId, newDataItemIds) {
					if err := db.Delete(userDataItem).Error; err != nil {
						return err
					}
				}
			}
			// 新增
			for _, dataItemId := range newDataItemIds {
				if !utils.InSlice(dataItemId, dataItemIds) {
					userDataItem := models.NewUserDataItem(map[string]interface{}{
						"user_id": userService.UserId,
						"service_id": userService.ServiceId,
						"data_item_id": dataItemId,
					})
					userDataItem.SetDb(db)
					if !(userDataItem.Validate() && userDataItem.Create()) {
						return errors.New(userDataItem.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			cache.SetUserDataItemIds(userService.UserId, userService.ServiceId, newDataItemIds)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
		}
	}

	this.HTML(c, "./views/admin/user-data-item/edit.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name,
		"subTitle": userService.GetService().Name,
		"user": user,
		"dataItems": dataItems,
		"dataItemIds": dataItemIds,
	})
}