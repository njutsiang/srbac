package models

import (
	"reflect"
	"srbac/libraries/utils"
	"srbac/srbac"
	"time"
)

// 菜单节点
type MenuItem struct {
	Model
	service *Service
	old struct {
		Key string
	}
	Id int64 `label:"ID"`
	ServiceId int64 `label:"所属服务" validate:"required"`
	Key string `label:"菜单标识" validate:"required"`
	Name string `label:"菜单名称" validate:"required"`
	UpdatedAt int64 `label:"更新时间"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewMenuItem(data map[string]interface{}) *MenuItem {
	menuItem := &MenuItem{
		ServiceId: utils.ToInt64(data["service_id"]),
		Key: utils.ToString(data["key"]),
		Name: utils.ToString(data["name"]),
		UpdatedAt: 0,
		CreatedAt: time.Now().Unix(),
	}
	menuItem.SetRefValue()
	return menuItem
}

// 表名
func (this *MenuItem) TableName() string {
	return "menu_item"
}

// 设置模型反射
func (this *MenuItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *MenuItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "service_id")
	delete(data, "updated_at")
	delete(data, "created_at")
	this.old.Key = this.Key
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *MenuItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *MenuItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *MenuItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(srbac.TimeYmdhis)
}

// 所属服务名称
func (this *MenuItem) GetServiceName() string {
	if this.service == nil {
		return ""
	} else {
		return this.service.Name
	}
}

// 实现 ServiceRelation
func (this *MenuItem) GetServiceId() int64 {
	return this.ServiceId
}

func (this *MenuItem) SetService(service *Service) {
	this.service = service
}

func (this *MenuItem) GetOld() struct{Key string} {
	return this.old
}