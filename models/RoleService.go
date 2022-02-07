package models

import (
	"reflect"
	"srbac/libraries/utils"
	"srbac/srbac"
	"time"
)

// 角色服务分配
type RoleService struct {
	Model
	service *Service
	Id int64 `label:"ID"`
	RoleId int64 `label:"角色" validate:"required"`
	ServiceId int64 `label:"服务" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewRoleService(data map[string]interface{}) *RoleService {
	roleService := &RoleService{
		RoleId: utils.ToInt64(data["role_id"]),
		ServiceId: utils.ToInt64(data["service_id"]),
		CreatedAt: time.Now().Unix(),
	}
	roleService.SetRefValue()
	return roleService
}

// 表名
func (this *RoleService) TableName() string {
	return "role_service"
}

// 设置模型反射
func (this *RoleService) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *RoleService) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "updated_at")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *RoleService) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *RoleService) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *RoleService) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(srbac.TimeYmdhis)
}

// 实现 ServiceRelation
func (this *RoleService) GetServiceId() int64 {
	return this.ServiceId
}

func (this *RoleService) SetService(service *Service) {
	this.service = service
}

func (this *RoleService) GetServiceKey() string {
	if this.service == nil {
		return ""
	} else {
		return this.service.Key
	}
}

func (this *RoleService) GetServiceName() string {
	if this.service == nil {
		return ""
	} else {
		return this.service.Name
	}
}