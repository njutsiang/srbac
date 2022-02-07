package models

import (
	"reflect"
	"srbac/libraries/utils"
	"srbac/srbac"
	"time"
)

// 角色接口节点分配
type RoleApiItem struct {
	Model
	Id int64 `label:"ID"`
	RoleId int64 `label:"角色" validate:"required"`
	ServiceId int64 `label:"服务" validate:"required"`
	ApiItemId int64 `label:"接口节点" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewRoleApiItem(data map[string]interface{}) *RoleApiItem {
	roleApiItem := &RoleApiItem{
		RoleId: utils.ToInt64(data["role_id"]),
		ServiceId: utils.ToInt64(data["service_id"]),
		ApiItemId: utils.ToInt64(data["api_item_id"]),
		CreatedAt: time.Now().Unix(),
	}
	roleApiItem.SetRefValue()
	return roleApiItem
}

// 表名
func (this *RoleApiItem) TableName() string {
	return "role_api_item"
}

// 设置模型反射
func (this *RoleApiItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *RoleApiItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *RoleApiItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *RoleApiItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *RoleApiItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(srbac.TimeYmdhis)
}