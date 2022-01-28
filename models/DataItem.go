package models

import (
	"reflect"
	"srbac/libraries/utils"
	"srbac/srbac"
	"time"
)

// 数据节点
type DataItem struct {
	Model
	service *Service
	Id int64 `label:"ID"`
	ServiceId int64 `label:"所属服务" validate:"required"`
	Key string `label:"权限标识" validate:"required"`
	Name string `label:"权限名称" validate:"required"`
	UpdatedAt int64 `label:"更新时间"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewDataItem(data map[string]interface{}) *DataItem {
	dataItem := &DataItem{
		ServiceId: utils.ToInt64(data["service_id"]),
		Key: utils.ToString(data["key"]),
		Name: utils.ToString(data["name"]),
		UpdatedAt: 0,
		CreatedAt: time.Now().Unix(),
	}
	dataItem.SetRefValue()
	return dataItem
}

// 表名
func (this *DataItem) TableName() string {
	return "data_item"
}

// 设置模型反射
func (this *DataItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *DataItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "updated_at")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *DataItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *DataItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *DataItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(srbac.TimeYmdhis)
}

// 所属服务名称
func (this *DataItem) GetServiceName() string {
	if this.service == nil {
		return ""
	} else {
		return this.service.Name
	}
}

// 实现 ServiceRelation
func (this *DataItem) GetServiceId() int64 {
	return this.ServiceId
}

func (this *DataItem) SetService(service *Service) {
	this.service = service
}