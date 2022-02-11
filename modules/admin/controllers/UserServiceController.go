package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"srbac/cache"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/libraries/utils"
	"srbac/logics"
	"srbac/models"
	"srbac/srbac"
)

// 用户服务关系
type UserServiceController struct {
	controllers.Controller
}

// 用户服务关系列表
func (this *UserServiceController) List(c *gin.Context) {
	userId := utils.ToInt64(c.Query("userId"))
	if userId <= 0 {
		exception.NewException(code.ParamsError)
	}

	params := c.Request.URL.Query()
	page, perPage := utils.GetPageInfo(params)

	user := &models.User{}
	re := srbac.Db.First(user, userId)
	srbac.CheckError(re.Error)

	count := int64(0)
	query := srbac.Db.Model(&models.UserService{}).Where("user_id = ?", userId).Count(&count)
	srbac.CheckError(query.Error)

	userServices := []*models.UserService{}
	re = query.Order("id asc").Offset((page - 1) * perPage).Limit(perPage).Find(&userServices)
	srbac.CheckError(re.Error)

	userServicesMap := map[int64]bool{}
	for _, userService := range userServices {
		userServicesMap[userService.ServiceId] = true
	}

	hasNew := false
	roleServiceIds := logics.FindRoleServiceIdsByUserId(userId)
	for _, serviceId := range roleServiceIds {
		if _, ok := userServicesMap[serviceId]; !ok {
			hasNew = true
			userService := models.NewUserService(map[string]interface{}{
				"user_id": userId,
				"service_id": serviceId,
			})
			if !(userService.Validate() && userService.Create()) {
				this.SetFailed(c, userService.GetError())
				break
			}
		}
	}
	if hasNew {
		userServices = []*models.UserService{}
		re = query.Order("id asc").Offset((page - 1) * perPage).Limit(perPage).Find(&userServices)
		srbac.CheckError(re.Error)
	}

	models.UserServicesLoadServices(userServices)

	this.HTML(c, "./views/admin/user-service/list.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name,
		"pager": utils.GetPageHtml(count, page, perPage, params, "/admin/user-service/list"),
		"user": user,
		"userServices": userServices,
		"roleServiceIds": roleServiceIds,
	})
}

// 编辑用户服务关系
func (this *UserServiceController) Edit(c *gin.Context) {
	userId := utils.ToInt64(c.Query("userId"))
	if userId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-service/list?userId=%d", userId))

	user := &models.User{}
	re := srbac.Db.First(user, userId)
	srbac.CheckError(re.Error)

	services := []*models.Service{}
	re = srbac.Db.Order("id asc").Limit(1000).Find(&services)
	srbac.CheckError(re.Error)

	userServices := []*models.UserService{}
	re = srbac.Db.Where("user_id = ?", userId).Find(&userServices)
	srbac.CheckError(re.Error)

	// 用户的角色拥有的服务 ids
	roleServiceIds := logics.FindRoleServiceIdsByUserId(userId)

	serviceIds := []int64{}
	for _, userService := range userServices {
		serviceIds = append(serviceIds, userService.ServiceId)
	}
	for _, serviceId := range roleServiceIds {
		if !utils.InSlice(serviceId, serviceIds) {
			serviceIds = append(serviceIds, serviceId)
		}
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		srbac.CheckError(err)
		newServiceIds := utils.ToSliceInt64(c.Request.PostForm["service_id[]"])

		// 删除
		for _, userService := range userServices {
			if !utils.InSlice(userService.ServiceId, newServiceIds) && !utils.InSlice(userService.ServiceId, roleServiceIds) {
				srbac.Db.Delete(userService)
				cache.DelUserApiItemsByUserService(userService)
				cache.DelUserDataItemsByUserService(userService)
				cache.DelUserMenuItemsByUserService(userService)
			}
		}

		// 新增
		hasErr := false
		for _, serviceId := range newServiceIds {
			if !utils.InSlice(serviceId, serviceIds) {
				userService := models.NewUserService(map[string]interface{}{
					"user_id": userId,
					"service_id": serviceId,
				})
				if !(userService.Validate() && userService.Create()) {
					hasErr = true
					this.SetFailed(c, userService.GetError())
					break
				}
			}
		}
		if !hasErr {
			this.Redirect(c, referer)
		}
	}
	this.HTML(c, "./views/admin/user-service/edit.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name,
		"user": user,
		"services": services,
		"serviceIds": serviceIds,
		"roleServiceIds": roleServiceIds,
	})
}

// 删除用户服务关系
func (this *UserServiceController) Delete(c *gin.Context) {
	id := utils.ToInt(c.Query("id"))
	userId := utils.ToInt(c.Query("userId"))
	if id <= 0 || userId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-service/list?userId=%d", userId))

	userService := &models.UserService{}
	re := srbac.Db.First(userService, id)
	srbac.CheckError(re.Error)

	re = srbac.Db.Delete(userService)
	srbac.CheckError(re.Error)

	cache.DelUserApiItemsByUserService(userService)
	cache.DelUserDataItemsByUserService(userService)
	cache.DelUserMenuItemsByUserService(userService)
	this.Redirect(c, referer)
}