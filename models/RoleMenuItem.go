package models

import (
	"reflect"
	"srbac/app"
	"srbac/app/utils"
	"time"
)

// 角色菜单节点分配
type RoleMenuItem struct {
	Model
	Id int64 `label:"ID"`
	RoleId int64 `label:"角色" validate:"required"`
	ServiceId int64 `label:"服务" validate:"required"`
	MenuItemId int64 `label:"菜单节点" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewRoleMenuItem(data map[string]interface{}) *RoleMenuItem {
	roleMenuItem := &RoleMenuItem{
		RoleId: utils.ToInt64(data["role_id"]),
		ServiceId: utils.ToInt64(data["service_id"]),
		MenuItemId: utils.ToInt64(data["menu_item_id"]),
		CreatedAt: time.Now().Unix(),
	}
	roleMenuItem.SetRefValue()
	return roleMenuItem
}

// 表名
func (this *RoleMenuItem) TableName() string {
	return "role_menu_item"
}

// 设置模型反射
func (this *RoleMenuItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *RoleMenuItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *RoleMenuItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *RoleMenuItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *RoleMenuItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(app.TimeYmdhis)
}