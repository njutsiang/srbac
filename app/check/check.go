package check

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/app/cache"
	"srbac/code"
	"srbac/exception"
	"srbac/models"
	"srbac/utils"
	"time"
)

// 初始化 SRBAC 数据
func InitSrbacData() {
	service := initSrbacService()
	initSrbacRouters(service)
	initSrbacSuperUser()
}

// 初始化 SRBAC 服务
func initSrbacService() *models.Service {
	serviceKey := "srbac-service"
	service := &models.Service{}
	re := app.Db.Where("`key`  = ?", serviceKey).First(service)
	if re.Error == nil {
		return service
	}
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		service = &models.Service{
			Key: serviceKey,
			Name: "SRBAC",
			UpdatedAt: 0,
			CreatedAt: time.Now().Unix(),
		}
		service.SetRefValue()
		if !(service.Validate() && service.Create()) {
			exception.NewException(code.InternalError, service.GetError())
		}
		cache.SetService(service)
		return service
	}
	app.CheckError(re.Error)
	return service
}

// 初始化 SRBAC 接口节点
func initSrbacRouters(service *models.Service) {
	for sort, route := range app.Routes {
		initSrbacRouter(service, route, sort)
	}
}

// 初始化 SRBAC 接口节点
func initSrbacRouter(service *models.Service, route app.Route, sort int) {
	anonymousUri := []string{
		"/",
		"/admin",
		"/admin/login",
		"/admin/logout",
	}
	apiItem := &models.ApiItem{}
	re := app.Db.Where("service_id = ?", service.Id).
		Where("method = ?", route.Method).
		Where("uri = ?", route.Uri).
		First(apiItem)
	if re.Error == nil {
		if apiItem.Sort != int64(sort) {
			apiItem.Sort = int64(sort)
			apiItem.SetRefValue()
			if !(apiItem.Validate() && apiItem.Update()) {
				exception.NewException(code.InternalError, apiItem.GetError())
			}
		}
		return
	}
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		isAnonymousAccess := int64(0)
		if utils.InSlice(route.Uri, anonymousUri) {
			isAnonymousAccess = 1
		}
		apiItem = &models.ApiItem{
			ServiceId: service.Id,
			Method: route.Method,
			Uri: route.Uri,
			Name: route.Name,
			IsAnonymousAccess: isAnonymousAccess,
			Sort: int64(sort),
			UpdatedAt: 0,
			CreatedAt: time.Now().Unix(),
		}
		apiItem.SetRefValue()
		if !(apiItem.Validate() && apiItem.Create()) {
			exception.NewException(code.InternalError, apiItem.GetError())
		}
		cache.SetApiItem(apiItem)
		return
	}
	app.CheckError(re.Error)
}

// 初始化 SRBAC 超级用户
func initSrbacSuperUser() {
	userId := int64(1)
	username := "admin"
	password := "123456"
	user := &models.User{}
	re := app.Db.First(user, userId)
	if re.Error == nil {
		return
	}
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		app.CheckError(err)
		user = &models.User{
			Id: userId,
			Name: "超级管理员",
			Username: username,
			Password: string(passwordHash),
			Status: 1,
			UpdatedAt: 0,
			CreatedAt: time.Now().Unix(),
		}
		user.SetRefValue()
		if !(user.Validate() && user.Create()) {
			exception.NewException(code.InternalError, user.GetError())
		}
		return
	}
	app.CheckError(re.Error)
}