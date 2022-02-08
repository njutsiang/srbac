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

// 用户的接口权限
type UserApiItemController struct {
	controllers.Controller
}

// 编辑用户的接口权限
func (this *UserApiItemController) Edit(c *gin.Context) {
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

	apiItems := []*models.ApiItem{}
	re = srbac.Db.Where("service_id = ?", userService.ServiceId).Order("id asc").Limit(1000).Find(&apiItems)
	srbac.CheckError(re.Error)

	userApiItems := []*models.UserApiItem{}
	re = srbac.Db.Where("user_id = ? AND service_id = ?", userService.UserId, userService.ServiceId).Limit(1000).Find(&userApiItems)
	srbac.CheckError(re.Error)

	apiItemIds := []int64{}
	for _, userApiItem := range userApiItems {
		apiItemIds = append(apiItemIds, userApiItem.ApiItemId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		srbac.CheckError(err)
		newApiItemIds := utils.ToSliceInt64(c.Request.PostForm["api_item_id[]"])

		// 删除
		for _, userApiItem := range userApiItems {
			if !utils.InSlice(userApiItem.ApiItemId, newApiItemIds) {
				srbac.Db.Delete(userApiItem)
			}
		}

		// 新增
		hasErr := false
		for _, apiItemId := range newApiItemIds {
			if !utils.InSlice(apiItemId, apiItemIds) {
				userApiItem := models.NewUserApiItem(map[string]interface{}{
					"user_id": userService.UserId,
					"service_id": userService.ServiceId,
					"api_item_id": apiItemId,
				})
				if !(userApiItem.Validate() && userApiItem.Create()) {
					hasErr = true
					this.SetFailed(c, userApiItem.GetError())
					break
				}
			}
		}
		if !hasErr {
			this.Redirect(c, referer)
		}
	}

	this.HTML(c, "./views/admin/user-api-item/edit.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name + " > " + userService.GetService().Name,
		"user": user,
		"apiItems": apiItems,
		"apiItemIds": apiItemIds,
	})
}