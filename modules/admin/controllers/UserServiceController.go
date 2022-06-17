package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/cache"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/libraries/utils"
	"srbac/logics"
	"srbac/models"
)

// 用户服务关系
type UserServiceController struct {
	controllers.Controller
}

// 用户服务关系列表
func (this *UserServiceController) List(c *gin.Context) {
	referer := "/admin/user/list"
	userId := utils.ToInt64(c.Query("userId"))
	if userId <= 1 {
		this.Redirect(c, referer)
	}

	params := c.Request.URL.Query()
	page, perPage := utils.GetPageInfo(params)

	user := &models.User{}
	re := app.Db.First(user, userId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	count := int64(0)
	query := app.Db.Model(&models.UserService{}).Where("user_id = ?", userId).Count(&count)
	app.CheckError(query.Error)

	userServices := []*models.UserService{}
	re = query.Order("id asc").Offset((page - 1) * perPage).Limit(perPage).Find(&userServices)
	app.CheckError(re.Error)

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
		app.CheckError(re.Error)
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
	re := app.Db.First(user, userId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	services := []*models.Service{}
	re = app.Db.Order("id asc").Limit(1000).Find(&services)
	app.CheckError(re.Error)

	userServices := []*models.UserService{}
	re = app.Db.Where("user_id = ?", userId).Find(&userServices)
	app.CheckError(re.Error)

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
		app.CheckError(err)
		newServiceIds := utils.ToSliceInt64(c.Request.PostForm["service_id[]"])
		deleteUserServices := []*models.UserService{}
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, userService := range userServices {
				if !utils.InSlice(userService.ServiceId, newServiceIds) && !utils.InSlice(userService.ServiceId, roleServiceIds) {
					if err := db.Delete(userService).Error; err != nil {
						return err
					}
					deleteUserServices = append(deleteUserServices, userService)
				}
			}
			// 新增
			for _, serviceId := range newServiceIds {
				if !utils.InSlice(serviceId, serviceIds) {
					userService := models.NewUserService(map[string]interface{}{
						"user_id": userId,
						"service_id": serviceId,
					})
					userService.SetDb(db)
					if !(userService.Validate() && userService.Create()) {
						return errors.New(userService.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			for _, deleteUserService := range deleteUserServices {
				cache.DelUserApiItemsByUserService(deleteUserService)
				cache.DelUserDataItemsByUserService(deleteUserService)
				cache.DelUserMenuItemsByUserService(deleteUserService)
			}
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
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
	id := utils.ToInt(this.GetPostForm(c)["id"])
	userId := utils.ToInt(this.GetPostForm(c)["userId"])
	if id <= 0 || userId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-service/list?userId=%d", userId), false)

	userService := &models.UserService{}
	re := app.Db.First(userService, id)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	err := app.Db.Transaction(func(db *gorm.DB) error {
		if err := db.Delete(userService).Error; err != nil {
			return err
		}
		if err := db.Where("user_id = ?", userService.UserId).
			Where("service_id = ?", userService.ServiceId).
			Delete(&models.UserApiItem{}).Error; err != nil {
				return err
		}
		if err := db.Where("user_id = ?", userService.UserId).
			Where("service_id = ?", userService.ServiceId).
			Delete(&models.UserDataItem{}).Error; err != nil {
				return err
		}
		if err := db.Where("user_id = ?", userService.UserId).
			Where("service_id = ?", userService.ServiceId).
			Delete(&models.UserMenuItem{}).Error; err != nil {
				return err
		}
		return nil
	})
	app.CheckError(err)

	cache.DelUserApiItemsByUserService(userService)
	cache.DelUserDataItemsByUserService(userService)
	cache.DelUserMenuItemsByUserService(userService)
	this.Redirect(c, referer)
}