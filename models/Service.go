package models

import (
	"reflect"
	"srbac/libraries/utils"
	"srbac/srbac"
	"time"
)

// 服务
type Service struct {
	Model
	old struct {
		Key string
	}
	Id int64 `label:"ID"`
	Key string `label:"服务标识" validate:"required"`
	Name string `label:"服务名称"`
	UpdatedAt int64 `label:"更新时间"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewService(data map[string]interface{}) *Service {
	service := &Service{
		Key: utils.ToString(data["key"]),
		Name: utils.ToString(data["name"]),
		UpdatedAt: 0,
		CreatedAt: time.Now().Unix(),
	}
	service.SetRefValue()
	return service
}

// 所有服务 ids 枚举值
func ServiceIds() []srbac.IntValue {
	services := []*Service{}
	re := srbac.Db.Order("id asc").Limit(1000).Find(&services)
	srbac.CheckError(re.Error)
	serviceIds := []srbac.IntValue{}
	for _, service := range services {
		serviceIds = append(serviceIds, srbac.IntValue{
			Key: service.Id,
			Value: service.Name,
		})
	}
	return serviceIds
}

// 表名
func (this *Service) TableName() string {
	return "service"
}

// 设置模型反射
func (this *Service) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *Service) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "updated_at")
	delete(data, "created_at")
	this.old.Key = this.Key
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *Service) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *Service) ErrorMessages() map[string]string {
	return map[string]string{
		"Key.required": "服务标识不能为空",
		"Name.required": "服务名称不能为空",
	}
}

// 格式化创建时间
func (this *Service) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(srbac.TimeYmdhis)
}

func (this *Service) GetOld() struct{Key string} {
	return this.old
}