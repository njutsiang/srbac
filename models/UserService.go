package models

import (
	"reflect"
	"srbac/libraries/utils"
	"srbac/srbac"
	"time"
)

// 用户服务分配
type UserService struct {
	Model
	service *Service
	Id int64 `label:"ID"`
	UserId int64 `label:"用户" validate:"required"`
	ServiceId int64 `label:"服务" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewUserService(data map[string]interface{}) *UserService {
	userService := &UserService{
		UserId: utils.ToInt64(data["user_id"]),
		ServiceId: utils.ToInt64(data["service_id"]),
		CreatedAt: time.Now().Unix(),
	}
	userService.SetRefValue()
	return userService
}

// 表名
func (this *UserService) TableName() string {
	return "user_service"
}

// 设置模型反射
func (this *UserService) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *UserService) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *UserService) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *UserService) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *UserService) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(srbac.TimeYmdhis)
}

// 实现 ServiceRelation
func (this *UserService) GetServiceId() int64 {
	return this.ServiceId
}

func (this *UserService) SetService(service *Service) {
	this.service = service
}

func (this *UserService) GetService() *Service {
	if this.service == nil {
		this.service = &Service{}
	}
	return this.service
}