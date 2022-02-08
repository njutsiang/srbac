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

// 用户角色关系
type UserRoleController struct {
	controllers.Controller
}

// 用户角色关系列表
func (this *UserRoleController) List(c *gin.Context) {
	userId := utils.ToInt(c.Query("userId"))
	if userId <= 0 {
		exception.NewException(code.ParamsError)
	}

	query := c.Request.URL.Query()
	page, perPage := utils.GetPageInfo(query)

	user := &models.User{}
	re := srbac.Db.First(user, userId)
	srbac.CheckError(re.Error)

	count := int64(0)

	userRoles := []*models.UserRole{}
	re = srbac.Db.Where("user_id = ?", userId).Limit(1000).Find(&userRoles)
	srbac.CheckError(re.Error)

	models.UserRolesLoadRoles(userRoles)

	this.HTML(c, "./views/admin/user-role/list.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name,
		"pager": utils.GetPageHtml(count, page, perPage, query, "/admin/user-role/list"),
		"user": user,
		"userRoles": userRoles,
	})
}

// 编辑用户角色关系
func (this *UserRoleController) Edit(c *gin.Context) {
	userId := utils.ToInt(c.Query("userId"))
	if userId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-role/list?userId=%d", userId))

	user := &models.User{}
	re := srbac.Db.First(user, userId)
	srbac.CheckError(re.Error)

	roles := []*models.Role{}
	re = srbac.Db.Order("id asc").Limit(1000).Find(&roles)
	srbac.CheckError(re.Error)

	userRoles := []*models.UserRole{}
	re = srbac.Db.Where("user_id = ?", userId).Find(&userRoles)
	srbac.CheckError(re.Error)

	roleIds := []int64{}
	for _, userRole := range userRoles {
		roleIds = append(roleIds, userRole.RoleId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		srbac.CheckError(err)
		NewRoleIds := utils.ToSliceInt64(c.Request.PostForm["role_id[]"])

		// 删除
		for _, userRole := range userRoles {
			if !utils.InSlice(userRole.RoleId, NewRoleIds) {
				srbac.Db.Delete(userRole)
			}
		}

		// 新增
		hasErr := false
		for _, roleId := range NewRoleIds {
			if !utils.InSlice(roleId, roleIds) {
				userRole := models.NewUserRole(map[string]interface{}{
					"user_id": userId,
					"role_id": roleId,
				})
				if !(userRole.Validate() && userRole.Create()) {
					hasErr = true
					this.SetFailed(c, userRole.GetError())
					break
				}
			}
		}
		if !hasErr {
			this.Redirect(c, referer)
		}
	}

	this.HTML(c, "./views/admin/user-role/edit.gohtml", map[string]interface{}{
		"menu": "role",
		"title": user.Name,
		"roles": roles,
		"roleIds": roleIds,
	})
}

// 删除用户角色关系
func (this *UserRoleController) Delete(c *gin.Context) {
	id := utils.ToInt(c.Query("id"))
	userId := utils.ToInt(c.Query("userId"))
	if id <= 0 || userId <= 0 {
		exception.NewException(code.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-role/list?userId=%d", userId))
	re := srbac.Db.Delete(&models.UserRole{}, id)
	srbac.CheckError(re.Error)
	this.Redirect(c, referer)
}