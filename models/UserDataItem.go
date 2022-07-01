package models

import (
	"reflect"
	"srbac/app"
	"srbac/utils"
	"time"
)

// 用户数据节点分配
type UserDataItem struct {
	Model
	Id int64 `label:"ID"`
	UserId int64 `label:"用户" validate:"required"`
	ServiceId int64 `label:"服务" validate:"required"`
	DataItemId int64 `label:"数据节点" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewUserDataItem(data map[string]interface{}) *UserDataItem {
	userDataItem := &UserDataItem{
		UserId: utils.ToInt64(data["user_id"]),
		ServiceId: utils.ToInt64(data["service_id"]),
		DataItemId: utils.ToInt64(data["data_item_id"]),
		CreatedAt: time.Now().Unix(),
	}
	userDataItem.SetRefValue()
	return userDataItem
}

// 表名
func (this *UserDataItem) TableName() string {
	return "user_data_item"
}

// 设置模型反射
func (this *UserDataItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *UserDataItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *UserDataItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *UserDataItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *UserDataItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(app.TimeYmdhis)
}