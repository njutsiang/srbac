package logics

import (
	"srbac/app"
	"srbac/models"
)

// 查询用户的角色所拥有的服务 ids
func FindRoleServiceIdsByUserId(userId int64) []int64 {
	roleServices := []*models.RoleService{}
	re := app.Db.
		Distinct("role_service.service_id").
		Joins("left join user_role on user_role.role_id = role_service.role_id").
		Where("user_role.user_id = ?", userId).
		Find(&roleServices)
	app.CheckError(re.Error)
	serviceIds := []int64{}
	for _, roleService := range roleServices {
		serviceIds = append(serviceIds, roleService.ServiceId)
	}
	return serviceIds
}