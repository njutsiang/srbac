package models

import (
	"reflect"
	"srbac/app"
	"srbac/app/utils"
	"time"
)

// 接口节点
type ApiItem struct {
	Model
	service *Service
	old struct {
		Method string
		Uri string
	}
	Id int64 `label:"ID"`
	ServiceId int64 `label:"所属服务" validate:"required"`
	Method string `label:"请求方式" validate:"required,max=8"`
	Uri string `label:"接口路由" validate:"required,max=128"`
	Name string `label:"接口名称" validate:"max=32"`
	IsAnonymousAccess int64 `label:"是否允许匿名文档"`
	Sort int64 `label:"排序值"`
	UpdatedAt int64 `label:"更新时间"`
	CreatedAt int64 `label:"创建时间" validate:"required"`
}

// 实例化
func NewApiItem(data map[string]interface{}) *ApiItem {
	apiItem := &ApiItem{
		ServiceId: utils.ToInt64(data["service_id"]),
		Method: utils.ToString(data["method"]),
		Uri: utils.ToString(data["uri"]),
		Name: utils.ToString(data["name"]),
		IsAnonymousAccess: utils.ToInt64(data["is_anonymous_access"]),
		Sort: 0,
		UpdatedAt: 0,
		CreatedAt: time.Now().Unix(),
	}
	apiItem.SetRefValue()
	return apiItem
}

// 请求方式枚举值
func ApiItemMethods() []string {
	return []string{
		"*", "GET", "POST", "PUT", "DELETE",
	}
}

// 表名
func (this *ApiItem) TableName() string {
	return "api_item"
}

// 设置模型反射
func (this *ApiItem) SetRefValue() {
	this.refValue = reflect.ValueOf(this)
}

// 向实例载入属性
func (this *ApiItem) SetAttributes(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "service_id")
	delete(data, "updated_at")
	delete(data, "created_at")
	this.old.Method = this.Method
	this.old.Uri = this.Uri
	this.SetRefValue()
	this.setAttributes(data)
}

// 校验数据
func (this *ApiItem) Validate() bool {
	this.InitValidate()
	this.errorMessages = this.ErrorMessages()
	this.error = this.validate.Struct(this)
	return this.error == nil
}

// 错误提示语
func (this *ApiItem) ErrorMessages() map[string]string {
	return map[string]string{}
}

// 格式化创建时间
func (this *ApiItem) GetCreatedAt() string {
	return time.Unix(this.CreatedAt, 0).Format(app.TimeYmdhis)
}

// 所属服务名称
func (this *ApiItem) GetServiceName() string {
	if this.service == nil {
		return ""
	} else {
		return this.service.Name
	}
}

// 实现 ServiceRelation
func (this *ApiItem) GetServiceId() int64 {
	return this.ServiceId
}

func (this *ApiItem) SetService(service *Service) {
	this.service = service
}

func (this *ApiItem) GetService() *Service {
	return this.service
}

func (this *ApiItem) GetOld() struct{Method string; Uri string} {
	return this.old
}