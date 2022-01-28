package routers

import (
	"github.com/gin-gonic/gin"
	"srbac/controllers"
	"srbac/modules/admin/controllers"
)

// 设置路由
func SetRouters(engine *gin.Engine) {
	// 静态文件
	engine.Static("/assets", "./assets")
	engine.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// 首页
	engine.GET("/", (&controllers.DefaultController{}).Index)

	// 后台首页
	engine.GET("/admin", (&admin.DefaultController{}).Index)

	// 服务管理
	engine.GET("/admin/service/list", (&admin.ServiceController{}).List)
	engine.GET("/admin/service/add", (&admin.ServiceController{}).Add)
	engine.POST("/admin/service/add", (&admin.ServiceController{}).Add)
	engine.GET("/admin/service/edit", (&admin.ServiceController{}).Edit)
	engine.POST("/admin/service/edit", (&admin.ServiceController{}).Edit)
	engine.GET("/admin/service/delete", (&admin.ServiceController{}).Delete)

	// 用户管理
	engine.GET("/admin/user/list", (&admin.UserController{}).List)
	engine.GET("/admin/user/add", (&admin.UserController{}).Add)
	engine.POST("/admin/user/add", (&admin.UserController{}).Add)
	engine.GET("/admin/user/edit", (&admin.UserController{}).Edit)
	engine.POST("/admin/user/edit", (&admin.UserController{}).Edit)
	engine.GET("/admin/user/delete", (&admin.UserController{}).Delete)

	// 角色管理
	engine.GET("/admin/role/list", (&admin.RoleController{}).List)
	engine.GET("/admin/role/add", (&admin.RoleController{}).Add)
	engine.POST("/admin/role/add", (&admin.RoleController{}).Add)
	engine.GET("/admin/role/edit", (&admin.RoleController{}).Edit)
	engine.POST("/admin/role/edit", (&admin.RoleController{}).Edit)
	engine.GET("/admin/role/delete", (&admin.RoleController{}).Delete)

	// 接口节点管理
	engine.GET("/admin/api-item/list", (&admin.ApiItemController{}).List)
	engine.GET("/admin/api-item/add", (&admin.ApiItemController{}).Add)
	engine.POST("/admin/api-item/add", (&admin.ApiItemController{}).Add)
	engine.GET("/admin/api-item/edit", (&admin.ApiItemController{}).Edit)
	engine.POST("/admin/api-item/edit", (&admin.ApiItemController{}).Edit)
	engine.GET("/admin/api-item/delete", (&admin.ApiItemController{}).Delete)

	// 数据节点管理
	engine.GET("/admin/data-item/list", (&admin.DataItemController{}).List)
	engine.GET("/admin/data-item/add", (&admin.DataItemController{}).Add)
	engine.POST("/admin/data-item/add", (&admin.DataItemController{}).Add)
	engine.GET("/admin/data-item/edit", (&admin.DataItemController{}).Edit)
	engine.POST("/admin/data-item/edit", (&admin.DataItemController{}).Edit)
	engine.GET("/admin/data-item/delete", (&admin.DataItemController{}).Delete)

	// 菜单节点管理
	engine.GET("/admin/menu-item/list", (&admin.MenuItemController{}).List)
	engine.GET("/admin/menu-item/add", (&admin.MenuItemController{}).Add)
	engine.POST("/admin/menu-item/add", (&admin.MenuItemController{}).Add)
	engine.GET("/admin/menu-item/edit", (&admin.MenuItemController{}).Edit)
	engine.POST("/admin/menu-item/edit", (&admin.MenuItemController{}).Edit)
	engine.GET("/admin/menu-item/delete", (&admin.MenuItemController{}).Delete)
}
