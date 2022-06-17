package models

import (
	"reflect"
	"srbac/app"
	"srbac/libraries/utils"
	"time"
)

// 角色
type Role struct {
	Model
	Id int64 `label:"ID"`
	Key string `label:"角色标识" validate:"required,max=32"`
	Name string `label:"角色名称" validate:"required,max=32"`
	UpdatedAt int64 `label:"更新时间"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewRole(data map[string]interface{}) *Role {
	role := &Role{
		Key: utils.ToString(data["key"]),
		Name: utils.ToString(data["name"]),
		UpdatedAt: 0,
		CreatedAt: time.Now().Unix(),
	}
	role.SetRefValue()
	return role
}

// 表名
func (this *Role) TableName() string {
	return "role"
}

// 设置模型反射
func (this *Role) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *Role) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "updated_at")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *Role) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *Role) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *Role) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(app.TimeYmdhis)
}