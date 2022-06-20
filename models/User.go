package models

import (
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"srbac/app"
	"srbac/app/utils"
	"time"
)

// 用户
type User struct {
	Model
	rawPassword string
	Id int64 `label:"ID"`
	Name string `label:"姓名" validate:"max=32"`
	Username string `label:"用户名" validate:"required,max=32"`
	Password string `label:"密码" validate:"required,max=128"`
	Status int64 `label:"状态"`
	UpdatedAt int64 `label:"更新时间"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewUser(data map[string]interface{}) *User {
	password := utils.ToString(data["password"])
	password_hash := []byte{}
	if len(password) > 0 {
		var err error
		password_hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		app.CheckError(err)
	}
	user := &User{
		rawPassword: password,
		Name: utils.ToString(data["name"]),
		Username: utils.ToString(data["username"]),
		Password: string(password_hash),
		Status: utils.ToInt64(data["status"]),
		UpdatedAt: 0,
		CreatedAt: time.Now().Unix(),
	}
	user.SetRefValue()
	return user
}

// 表名
func (this *User) TableName() string {
	return "user"
}

// 设置模型反射
func (this *User) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *User) SetAttributes(data map[string]interface{}) {
	password := utils.ToString(data["password"])
	if len(password) > 0 {
		password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		app.CheckError(err)
		this.rawPassword = password
		this.Password = string(password_hash)
	}
	delete(data, "id")
	delete(data, "password")
	delete(data, "updated_at")
	delete(data, "created_at")
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *User) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *User) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *User) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(app.TimeYmdhis)
}

// 原始明文密码
func (this *User) RawPassword() string {
	return this.rawPassword
}