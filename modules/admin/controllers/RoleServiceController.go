package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/libraries/utils"
	"srbac/models"
	"srbac/srbac"
)

// 角色服务关系
type RoleServiceController struct {
	controllers.Controller
}

// 角色服务关系列表
func (this *RoleServiceController) List(c *gin.Context) {
	query := c.Request.URL.Query()
	page, perPage := utils.GetPageInfo(query)

	roleId := utils.ToInt(c.Query("roleId"))
	if roleId <= 0 {
		exception.NewException(code.ParamsError)
	}

	role := &models.Role{}
	re := srbac.Db.First(role, roleId)
	srbac.CheckError(re.Error)

	count := int64(0)
	re = srbac.Db.Model(&models.RoleService{}).Count(&count)
	srbac.CheckError(re.Error)

	roleServices := []*models.RoleService{}
	re = srbac.Db.Where("role_id = ?", roleId).Order("service_id asc").Limit(perPage).Offset((page - 1) * perPage).Find(&roleServices)
	srbac.CheckError(re.Error)

	models.RoleServiceLoadServices(roleServices)

	this.HTML(c, "./views/admin/role-service/list.gohtml", map[string]interface{}{
		"menu": "role",
		"title": "角色：" + role.Name,
		"pager": utils.GetPageHtml(count, page, perPage, query, "/admin/role-service/list"),
		"role": role,
		"roleServices": roleServices,
	})
}

// 编辑角色服务关系
func (this *RoleServiceController) Edit(c *gin.Context) {
	roleId := utils.ToInt64(c.Query("roleId"))
	if roleId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/role-service/list?roleId=%d", roleId))

	// 当前角色
	role := &models.Role{}
	re := srbac.Db.First(role, roleId)
	srbac.CheckError(re.Error)

	// 所有服务
	services := []*models.Service{}
	re = srbac.Db.Order("id asc").Limit(1000).Find(&services)
	srbac.CheckError(re.Error)

	// 角色服务关系
	roleServices := []*models.RoleService{}
	re = srbac.Db.Where("role_id = ?", roleId).Limit(1000).Find(&roleServices)
	srbac.CheckError(re.Error)

	// 角色关联的服务 ids
	serviceIds := []int64{}
	for _, roleService := range roleServices {
		serviceIds = append(serviceIds, roleService.ServiceId)
	}

	if c.Request.Method == "POST" {
		// 取参数
		err := c.Request.ParseForm()
		srbac.CheckError(err)
		NewServiceIds := utils.ToSliceInt64(c.Request.PostForm["service_id[]"])

		// 删除
		for _, roleService := range roleServices {
			if !utils.InSlice(roleService.ServiceId, NewServiceIds) {
				srbac.Db.Delete(roleService)
			}
		}

		// 新增
		for _, serviceId := range NewServiceIds {
			if !utils.InSlice(serviceId, serviceIds) {
				roleService := models.NewRoleService(map[string]interface{}{
					"role_id": roleId,
					"service_id": serviceId,
				})
				if !(roleService.Validate() && roleService.Create()) {
					this.SetFailed(c, roleService.GetError())
					break
				}
			}
		}

		this.Redirect(c, referer)
	}

	this.HTML(c, "./views/admin/role-service/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": "角色：" + role.Name,
		"services": services,
		"serviceIds": serviceIds,
	})
}

// 删除角色服务关系
func (this *RoleServiceController) Delete(c *gin.Context) {
	id := utils.ToInt(c.Query("id"))
	roleId := utils.ToInt(c.Query("roleId"))
	if id <= 0 || roleId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/role-service/list?roleId=%d", roleId))

	re := srbac.Db.Delete(&models.RoleService{}, id)
	srbac.CheckError(re.Error)

	this.Redirect(c, referer)
}