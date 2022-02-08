package logics

import (
	"srbac/models"
	"srbac/srbac"
)

// 查询用户的角色所拥有的服务 ids
func FindRoleServiceIdsByUserId(userId int64) []int64 {
	roleServices := []*models.RoleService{}
	re := srbac.Db.
		Distinct("role_service.service_id").
		Joins("left join user_role on user_role.role_id = role_service.role_id").
		Where("user_role.user_id = ?", userId).
		Find(&roleServices)
	srbac.CheckError(re.Error)
	serviceIds := []int64{}
	for _, roleService := range roleServices {
		serviceIds = append(serviceIds, roleService.ServiceId)
	}
	return serviceIds
}