package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/app/cache"
	"srbac/code"
	"srbac/controllers"
	"srbac/exception"
	"srbac/models"
	"srbac/utils"
)

// 角色服务关系
type RoleServiceController struct {
	controllers.Controller
}

// 角色服务关系列表
func (this *RoleServiceController) List(c *gin.Context) {
	referer := "/admin/role/list"
	params := c.Request.URL.Query()
	page, perPage := utils.GetPageInfo(params)

	roleId := utils.ToInt(c.Query("roleId"))
	if roleId <= 0 {
		exception.NewException(code.ParamsError)
	}

	role := &models.Role{}
	re := app.Db.First(role, roleId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	count := int64(0)
	query := app.Db.Model(&models.RoleService{}).Where("role_id = ?", roleId).Count(&count)
	app.CheckError(query.Error)

	roleServices := []*models.RoleService{}
	re = query.Order("service_id asc").Limit(perPage).Offset((page - 1) * perPage).Find(&roleServices)
	app.CheckError(re.Error)

	models.RoleServicesLoadServices(roleServices)

	this.HTML(c, "./views/admin/role-service/list.gohtml", map[string]interface{}{
		"menu": "role",
		"title": role.Name,
		"pager": utils.GetPageHtml(count, page, perPage, params, "/admin/role-service/list"),
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
	re := app.Db.First(role, roleId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	// 所有服务
	services := []*models.Service{}
	re = app.Db.Order("id asc").Limit(1000).Find(&services)
	app.CheckError(re.Error)

	// 角色服务关系
	roleServices := []*models.RoleService{}
	re = app.Db.Where("role_id = ?", roleId).Limit(1000).Find(&roleServices)
	app.CheckError(re.Error)

	// 角色关联的服务 ids
	serviceIds := []int64{}
	for _, roleService := range roleServices {
		serviceIds = append(serviceIds, roleService.ServiceId)
	}

	if c.Request.Method == "POST" {
		// 取参数
		err := c.Request.ParseForm()
		app.CheckError(err)
		newServiceIds := utils.ToSliceInt64(c.Request.PostForm["service_id[]"])
		deleteRoleServices := []*models.RoleService{}
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, roleService := range roleServices {
				if !utils.InSlice(roleService.ServiceId, newServiceIds) {
					if err := db.Delete(roleService).Error; err != nil {
						return err
					}
					deleteRoleServices = append(deleteRoleServices, roleService)
				}
			}
			// 新增
			for _, serviceId := range newServiceIds {
				if !utils.InSlice(serviceId, serviceIds) {
					roleService := models.NewRoleService(map[string]interface{}{
						"role_id": roleId,
						"service_id": serviceId,
					})
					roleService.SetDb(db)
					if !(roleService.Validate() && roleService.Create()) {
						return errors.New(roleService.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			for _, deleteRoleService := range deleteRoleServices {
				cache.DelRoleApiItemsByRoleService(deleteRoleService)
				cache.DelRoleDataItemsByRoleService(deleteRoleService)
				cache.DelRoleMenuItemsByRoleService(deleteRoleService)
			}
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
		}
	}

	this.HTML(c, "./views/admin/role-service/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": role.Name,
		"services": services,
		"serviceIds": serviceIds,
	})
}

// 删除角色服务关系
func (this *RoleServiceController) Delete(c *gin.Context) {
	id := utils.ToInt64(this.GetPostForm(c)["id"])
	roleId := utils.ToInt64(this.GetPostForm(c)["roleId"])
	if id <= 0 || roleId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/role-service/list?roleId=%d", roleId), false)

	roleService := &models.RoleService{}
	re := app.Db.First(roleService, id)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	err := app.Db.Transaction(func(db *gorm.DB) error {
		if err := db.Delete(roleService).Error; err != nil {
			return err
		}
		if err := db.Where("role_id = ?", roleService.RoleId).
			Where("service_id = ?", roleService.ServiceId).
			Delete(&models.RoleApiItem{}).Error; err != nil {
				return err
		}
		if err := db.Where("role_id = ?", roleService.RoleId).
			Where("service_id = ?", roleService.ServiceId).
			Delete(&models.RoleDataItem{}).Error; err != nil {
				return err
		}
		if err := db.Where("role_id = ?", roleService.RoleId).
			Where("service_id = ?", roleService.ServiceId).
			Delete(&models.RoleMenuItem{}).Error; err != nil {
				return err
		}
		return nil
	})
	app.CheckError(err)

	cache.DelRoleApiItemsByRoleService(roleService)
	cache.DelRoleDataItemsByRoleService(roleService)
	cache.DelRoleMenuItemsByRoleService(roleService)
	this.Redirect(c, referer)
}