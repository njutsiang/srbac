package models

import (
	"srbac/srbac"
)

// 服务的关联数据模型
type ServiceRelation interface {
	GetServiceId() int64
	SetService(service *Service)
}

// 服务的关联数据模型列表
type ServiceRelations []ServiceRelation

// 为数据模型列表载入所属服务
func (models ServiceRelations) LoadServices() {
	if len(models) == 0 {
		return
	}

	// 查询服务列表
	serviceIds := []int64{}
	for _, model := range models {
		serviceIds = append(serviceIds, model.GetServiceId())
	}
	services := []*Service{}
	re := srbac.Db.Find(&services, serviceIds)
	srbac.CheckError(re.Error)

	// 处理成以 Id 为键的 Map
	servicesMap := map[int64]*Service{}
	for _, service := range services {
		servicesMap[service.Id] = service
	}

	// 载入到数据模型
	for _, model := range models {
		service, ok := servicesMap[model.GetServiceId()]
		if ok {
			model.SetService(service)
		}
	}
}

// 接口节点列表载入所属服务
func ApiItemsLoadServices(apiItems []*ApiItem) {
	models := ServiceRelations{}
	for _, apiItem := range apiItems {
		models = append(models, apiItem)
	}
	models.LoadServices()
}

// 数据节点列表载入所属服务
func DataItemsLoadServices(dataItems []*DataItem) {
	models := ServiceRelations{}
	for _, dataItem := range dataItems {
		models = append(models, dataItem)
	}
	models.LoadServices()
}

// 菜单节点列表载入所属服务
func MenuItemsLoadServices(menuItems []*MenuItem) {
	models := ServiceRelations{}
	for _, menuItem := range menuItems {
		models = append(models, menuItem)
	}
	models.LoadServices()
}

// 角色服务分配载入所属服务
func RoleServicesLoadServices(roleServices []*RoleService) {
	models := ServiceRelations{}
	for _, roleService := range roleServices {
		models = append(models, roleService)
	}
	models.LoadServices()
}