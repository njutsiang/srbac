package models

import (
	"reflect"
	"srbac/app"
	"srbac/libraries/utils"
	"time"
)

// 用户接口节点分配
type UserApiItem struct {
	Model
	Id int64 `label:"ID"`
	UserId int64 `label:"用户" validate:"required"`
	ServiceId int64 `label:"服务" validate:"required"`
	ApiItemId int64 `label:"接口节点" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewUserApiItem(data map[string]interface{}) *UserApiItem {
	userApiItem := &UserApiItem{
		UserId: utils.ToInt64(data["user_id"]),
		ServiceId: utils.ToInt64(data["service_id"]),
		ApiItemId: utils.ToInt64(data["api_item_id"]),
		CreatedAt: time.Now().Unix(),
	}
	userApiItem.SetRefValue()
	return userApiItem
}

// 表名
func (this *UserApiItem) TableName() string {
	return "user_api_item"
}

// 设置模型反射
func (this *UserApiItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *UserApiItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *UserApiItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *UserApiItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *UserApiItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(app.TimeYmdhis)
}