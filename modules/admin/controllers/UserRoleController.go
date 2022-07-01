package admin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"srbac/app"
	"srbac/app/cache"
	"srbac/controllers"
	"srbac/exception"
	"srbac/models"
	"srbac/utils"
)

// 用户角色关系
type UserRoleController struct {
	controllers.Controller
}

// 用户角色关系列表
func (this *UserRoleController) List(c *gin.Context) {
	referer := "/admin/user/list"
	userId := utils.ToInt(c.Query("userId"))
	if userId <= 1 {
		this.Redirect(c, referer)
	}

	params := c.Request.URL.Query()
	page, perPage := utils.GetPageInfo(params)

	user := &models.User{}
	re := app.Db.First(user, userId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	count := int64(0)
	query := app.Db.Model(&models.UserRole{}).Where("user_id = ?", userId).Count(&count)
	app.CheckError(query.Error)

	userRoles := []*models.UserRole{}
	re = query.Order("id asc").Offset((page - 1) * perPage).Limit(perPage).Find(&userRoles)
	app.CheckError(re.Error)

	models.UserRolesLoadRoles(userRoles)

	this.HTML(c, "./views/admin/user-role/list.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name,
		"pager": utils.GetPageHtml(count, page, perPage, params, "/admin/user-role/list"),
		"user": user,
		"userRoles": userRoles,
	})
}

// 编辑用户角色关系
func (this *UserRoleController) Edit(c *gin.Context) {
	userId := utils.ToInt64(c.Query("userId"))
	if userId <= 0 {
		exception.Throw(exception.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-role/list?userId=%d", userId))

	user := &models.User{}
	re := app.Db.First(user, userId)
	if errors.Is(re.Error, gorm.ErrRecordNotFound) {
		this.Redirect(c, referer)
	}
	app.CheckError(re.Error)

	roles := []*models.Role{}
	re = app.Db.Order("id asc").Limit(1000).Find(&roles)
	app.CheckError(re.Error)

	userRoles := []*models.UserRole{}
	re = app.Db.Where("user_id = ?", userId).Find(&userRoles)
	app.CheckError(re.Error)

	roleIds := []int64{}
	for _, userRole := range userRoles {
		roleIds = append(roleIds, userRole.RoleId)
	}

	if c.Request.Method == "POST" {
		err := c.Request.ParseForm()
		app.CheckError(err)
		newRoleIds := utils.ToSliceInt64(c.Request.PostForm["role_id[]"])
		if err := app.Db.Transaction(func(db *gorm.DB) error {
			// 删除
			for _, userRole := range userRoles {
				if !utils.InSlice(userRole.RoleId, newRoleIds) {
					if err := db.Delete(userRole).Error; err != nil {
						return err
					}
				}
			}
			// 新增
			for _, roleId := range newRoleIds {
				if !utils.InSlice(roleId, roleIds) {
					userRole := models.NewUserRole(map[string]interface{}{
						"user_id": userId,
						"role_id": roleId,
					})
					userRole.SetDb(db)
					if !(userRole.Validate() && userRole.Create()) {
						return errors.New(userRole.GetError())
					}
				}
			}
			return nil
		}); err == nil {
			cache.SetUserRoleIds(userId, newRoleIds)
			this.Redirect(c, referer)
		} else {
			this.SetFailed(c, err.Error())
		}
	}

	this.HTML(c, "./views/admin/user-role/edit.gohtml", map[string]interface{}{
		"menu": "user",
		"title": user.Name,
		"user": user,
		"roles": roles,
		"roleIds": roleIds,
	})
}

// 删除用户角色关系
func (this *UserRoleController) Delete(c *gin.Context) {
	id := utils.ToInt64(this.GetPostForm(c)["id"])
	userId := utils.ToInt64(this.GetPostForm(c)["userId"])
	if id <= 0 || userId <= 0 {
		exception.Throw(exception.ParamsError)
	}
	referer := this.GetReferer(c, fmt.Sprintf("/admin/user-role/list?userId=%d", userId), false)
	re := app.Db.Delete(&models.UserRole{}, id)
	app.CheckError(re.Error)
	cache.SetUserRoles(userId)
	this.Redirect(c, referer)
}