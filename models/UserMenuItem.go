package models

import (
	"reflect"
	"srbac/libraries/utils"
	"srbac/srbac"
	"time"
)

// 用户菜单节点分配
type UserMenuItem struct {
	Model
	Id int64 `label:"ID"`
	UserId int64 `label:"用户" validate:"required"`
	ServiceId int64 `label:"服务" validate:"required"`
	MenuItemId int64 `label:"菜单节点" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewUserMenuItem(data map[string]interface{}) *UserMenuItem {
	userMenuItem := &UserMenuItem{
		UserId: utils.ToInt64(data["user_id"]),
		ServiceId: utils.ToInt64(data["service_id"]),
		MenuItemId: utils.ToInt64(data["menu_item_id"]),
		CreatedAt: time.Now().Unix(),
	}
	userMenuItem.SetRefValue()
	return userMenuItem
}

// 表名
func (this *UserMenuItem) TableName() string {
	return "user_menu_item"
}

// 设置模型反射
func (this *UserMenuItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *UserMenuItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *UserMenuItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *UserMenuItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *UserMenuItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(srbac.TimeYmdhis)
}