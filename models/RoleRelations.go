package models

import "srbac/srbac"

// 角色的关联数据模型
type RoleRelation interface {
	GetRoleId() int64
	SetRole(role *Role)
}

// 角色的关联数据模型列表
type RoleRelations []RoleRelation

// 为数据模型列表载入所属服务
func (models RoleRelations) LoadRoles() {
	if len(models) == 0 {
		return
	}

	// 查询服务列表
	roleIds := []int64{}
	for _, model := range models {
		roleIds = append(roleIds, model.GetRoleId())
	}
	roles := []*Role{}
	re := srbac.Db.Find(&roles, roleIds)
	srbac.CheckError(re.Error)

	// 处理成以 Id 为键的 Map
	roleMap := map[int64]*Role{}
	for _, role := range roles {
		roleMap[role.Id] = role
	}

	// 载入到数据模型
	for _, model := range models {
		role, ok := roleMap[model.GetRoleId()]
		if ok {
			model.SetRole(role)
		}
	}
}

// 用户角色载入角色
func UserRolesLoadRoles(userRoles []*UserRole) {
	models := RoleRelations{}
	for _, userRole := range userRoles {
		models = append(models, userRole)
	}
	models.LoadRoles()
}