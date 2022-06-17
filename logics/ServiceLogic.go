package logics

import (
	"srbac/app"
	"srbac/models"
)

// 所有服务 ids 枚举值
func ServiceIds() []app.IntValue {
	services := []*models.Service{}
	re := app.Db.Order("id asc").Limit(1000).Find(&services)
	app.CheckError(re.Error)
	serviceIds := []app.IntValue{}
	for _, service := range services {
		serviceIds = append(serviceIds, app.IntValue{
			Key: service.Id,
			Value: service.Name,
		})
	}
	return serviceIds
}