package logics

import (
	"srbac/models"
	"srbac/srbac"
)

// 所有服务 ids 枚举值
func ServiceIds() []srbac.IntValue {
	services := []*models.Service{}
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