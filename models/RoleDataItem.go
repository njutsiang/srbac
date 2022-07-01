package models

import (
	"reflect"
	"srbac/app"
	"srbac/utils"
	"time"
)

// 角色数据节点分配
type RoleDataItem struct {
	Model
	Id int64 `label:"ID"`
	RoleId int64 `label:"角色" validate:"required"`
	ServiceId int64 `label:"服务" validate:"required"`
	DataItemId int64 `label:"数据节点" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewRoleDataItem(data map[string]interface{}) *RoleDataItem {
	roleDataItem := &RoleDataItem{
		RoleId: utils.ToInt64(data["role_id"]),
		ServiceId: utils.ToInt64(data["service_id"]),
		DataItemId: utils.ToInt64(data["data_item_id"]),
		CreatedAt: time.Now().Unix(),
	}
	roleDataItem.SetRefValue()
	return roleDataItem
}

// 表名
func (this *RoleDataItem) TableName() string {
	return "role_data_item"
}

// 设置模型反射
func (this *RoleDataItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *RoleDataItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *RoleDataItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *RoleDataItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *RoleDataItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(app.TimeYmdhis)
}