package models

import (
	"reflect"
	"srbac/app"
	"srbac/app/utils"
	"time"
)

// 用户角色分配
type UserRole struct {
	Model
	role *Role
	Id int64 `label:"ID"`
	UserId int64 `label:"用户" validate:"required"`
	RoleId int64 `label:"角色" validate:"required"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewUserRole(data map[string]interface{}) *UserRole {
	userRole := &UserRole{
		UserId: utils.ToInt64(data["user_id"]),
		RoleId: utils.ToInt64(data["role_id"]),
		CreatedAt: time.Now().Unix(),
	}
	userRole.SetRefValue()
	return userRole
}

// 表名
func (this *UserRole) TableName() string {
	return "user_role"
}

// 设置模型反射
func (this *UserRole) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *UserRole) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *UserRole) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *UserRole) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *UserRole) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(app.TimeYmdhis)
}

// 实现 RoleRelation
func (this *UserRole) GetRoleId() int64 {
	return this.RoleId
}

func (this *UserRole) SetRole(role *Role) {
	this.role = role
}

func (this *UserRole) GetRole() *Role {
	if this.role == nil {
		this.role = &Role{}
	}
	return this.role
}