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
	"srbac/logics"
	"srbac/models"
	"srbac/utils"
)

// 用户的接口权限
type UserApiItemController struct {
	controllers.Controller
}

// 编辑用户的接口权限
func (this *UserApiItemController) Edit(c *gin.Context) {
	userId := utils.ToInt(c.Query("userId"))
	userServiceId := utils.ToInt(c.Query("userServiceId"))
	if userServiceId <= 0 {
		exception.Throw(exception.ParamsError)
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

	apiItems := []*models.ApiItem{}
	re = logics.WithApiItemsOrder(app.Db.Where("service_id = ?", userService.ServiceId)).Limit(1000).Find(&apiItems)
	app.CheckError(re.Error)

	userApiItems := []*models.UserApiItem{}
	re = app.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Limit(1000).Find(&userApiItems)
	app.CheckError(re.Error)

	apiItemIds := []int64{}
	for _, userApiItem := range userApiItems {
		apiItemIds = append(apiItemIds, userApiItem.ApiItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		app.CheckError(err)
		newApiItemIds := utils.ToSliceInt64(c.Request.PostForm["api_item_id[]"])
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, userApiItem := range userApiItems {
				if !utils.InSlice(userApiItem.ApiItemId, newApiItemIds) {
					if err := db.Delete(userApiItem).Error; err != nil {
						return err
					}
				}
			}
			// 新增
			for _, apiItemId := range newApiItemIds {
				if !utils.InSlice(apiItemId, apiItemIds) {
					userApiItem := models.NewUserApiItem(map[string]interface{}{
						"user_id": userService.UserId,
						"service_id": userService.ServiceId,
						"api_item_id": apiItemId,
					})
					userApiItem.SetDb(db)
					if !(userApiItem.Validate() && userApiItem.Create()) {
						return errors.New(userApiItem.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			cache.SetUserApiItemIds(userService.UserId, userService.ServiceId, newApiItemIds)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
		}
	}

	this.HTML(c, "./views/admin/user-api-item/edit.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name,
		"subTitle": userService.GetService().Name,
		"user": user,
		"apiItems": apiItems,
		"apiItemIds": apiItemIds,
	})
}