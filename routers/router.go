package routers

import (
	"github.com/gin-gonic/gin"
	"srbac/app"
	"srbac/controllers"
	"srbac/modules/admin/controllers"
)

// 设置路由
func SetRouters(engine *gin.Engine) {
	// 静态文件
	engine.Static("/assets", "./assets")
	engine.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// 首页
	app.GET("/", (&controllers.DefaultController{}).Index, "首页")

	// 后台首页
	app.GET("/admin", (&admin.DefaultController{}).Index, "后台首页")

	// 登录
	app.GET("/admin/login", (&admin.LoginController{}).Login, "登录")
	app.POST("/admin/login", (&admin.LoginController{}).Login, "登录")
	app.GET("/admin/logout", (&admin.LoginController{}).Logout, "退出登录")

	// 服务管理
	app.GET("/admin/service/list", (&admin.ServiceController{}).List, "服务列表")
	app.GET("/admin/service/add", (&admin.ServiceController{}).Add, "添加服务")
	app.POST("/admin/service/add", (&admin.ServiceController{}).Add, "添加服务")
	app.GET("/admin/service/edit", (&admin.ServiceController{}).Edit, "编辑服务")
	app.POST("/admin/service/edit", (&admin.ServiceController{}).Edit, "编辑服务")
	app.POST("/admin/service/delete", (&admin.ServiceController{}).Delete, "删除服务")

	// 用户管理
	app.GET("/admin/user/list", (&admin.UserController{}).List, "用户列表")
	app.GET("/admin/user/add", (&admin.UserController{}).Add, "添加用户")
	app.POST("/admin/user/add", (&admin.UserController{}).Add, "添加用户")
	app.GET("/admin/user/edit", (&admin.UserController{}).Edit, "编辑用户")
	app.POST("/admin/user/edit", (&admin.UserController{}).Edit, "编辑用户")
	app.POST("/admin/user/delete", (&admin.UserController{}).Delete, "删除用户")

	// 角色管理
	app.GET("/admin/role/list", (&admin.RoleController{}).List, "角色列表")
	app.GET("/admin/role/add", (&admin.RoleController{}).Add, "添加角色")
	app.POST("/admin/role/add", (&admin.RoleController{}).Add, "添加角色")
	app.GET("/admin/role/edit", (&admin.RoleController{}).Edit, "编辑角色")
	app.POST("/admin/role/edit", (&admin.RoleController{}).Edit, "编辑角色")
	app.POST("/admin/role/delete", (&admin.RoleController{}).Delete, "删除角色")

	// 接口节点管理
	app.GET("/admin/api-item/list", (&admin.ApiItemController{}).List, "接口节点列表")
	app.GET("/admin/api-item/add", (&admin.ApiItemController{}).Add, "添加接口节点")
	app.POST("/admin/api-item/add", (&admin.ApiItemController{}).Add, "添加接口节点")
	app.GET("/admin/api-item/edit", (&admin.ApiItemController{}).Edit, "编辑接口节点")
	app.POST("/admin/api-item/edit", (&admin.ApiItemController{}).Edit, "编辑接口节点")
	app.POST("/admin/api-item/delete", (&admin.ApiItemController{}).Delete, "删除接口节点")

	// 数据节点管理
	app.GET("/admin/data-item/list", (&admin.DataItemController{}).List, "数据节点列表")
	app.GET("/admin/data-item/add", (&admin.DataItemController{}).Add, "添加数据节点")
	app.POST("/admin/data-item/add", (&admin.DataItemController{}).Add, "添加数据节点")
	app.GET("/admin/data-item/edit", (&admin.DataItemController{}).Edit, "编辑数据节点")
	app.POST("/admin/data-item/edit", (&admin.DataItemController{}).Edit, "编辑数据节点")
	app.POST("/admin/data-item/delete", (&admin.DataItemController{}).Delete, "删除数据节点")

	// 菜单节点管理
	app.GET("/admin/menu-item/list", (&admin.MenuItemController{}).List, "菜单节点列表")
	app.GET("/admin/menu-item/add", (&admin.MenuItemController{}).Add, "添加菜单节点")
	app.POST("/admin/menu-item/add", (&admin.MenuItemController{}).Add, "添加菜单节点")
	app.GET("/admin/menu-item/edit", (&admin.MenuItemController{}).Edit, "编辑菜单节点")
	app.POST("/admin/menu-item/edit", (&admin.MenuItemController{}).Edit, "编辑菜单节点")
	app.POST("/admin/menu-item/delete", (&admin.MenuItemController{}).Delete, "删除菜单节点")

	// 角色服务分配
	app.GET("/admin/role-service/list", (&admin.RoleServiceController{}).List, "角色服务关系列表")
	app.GET("/admin/role-service/edit", (&admin.RoleServiceController{}).Edit, "给角色分配服务")
	app.POST("/admin/role-service/edit", (&admin.RoleServiceController{}).Edit, "给角色分配服务")
	app.POST("/admin/role-service/delete", (&admin.RoleServiceController{}).Delete, "解除角色服务关系")

	// 角色服务接口节点分配
	app.GET("/admin/role-api-item/edit", (&admin.RoleApiItemController{}).Edit, "给角色分配接口节点")
	app.POST("/admin/role-api-item/edit", (&admin.RoleApiItemController{}).Edit, "给角色分配接口节点")

	// 角色服务数据节点分配
	app.GET("/admin/role-data-item/edit", (&admin.RoleDataItemController{}).Edit, "给角色分配数据节点")
	app.POST("/admin/role-data-item/edit", (&admin.RoleDataItemController{}).Edit, "给角色分配数据节点")

	// 角色服务菜单节点分配
	app.GET("/admin/role-menu-item/edit", (&admin.RoleMenuItemController{}).Edit, "给角色分配菜单节点")
	app.POST("/admin/role-menu-item/edit", (&admin.RoleMenuItemController{}).Edit, "给角色分配菜单节点")

	// 用户角色分配
	app.GET("/admin/user-role/list", (&admin.UserRoleController{}).List, "用户角色关系列表")
	app.GET("/admin/user-role/edit", (&admin.UserRoleController{}).Edit, "给用户分配角色")
	app.POST("/admin/user-role/edit", (&admin.UserRoleController{}).Edit, "给用户分配角色")
	app.POST("/admin/user-role/delete", (&admin.UserRoleController{}).Delete, "解除用户角色关系")

	// 用户服务分配
	app.GET("/admin/user-service/list", (&admin.UserServiceController{}).List, "用户服务关系列表")
	app.GET("/admin/user-service/edit", (&admin.UserServiceController{}).Edit, "给用户分配服务")
	app.POST("/admin/user-service/edit", (&admin.UserServiceController{}).Edit, "给用户分配服务")
	app.POST("/admin/user-service/delete", (&admin.UserServiceController{}).Delete, "解除用户服务关系")

	// 用户服务接口节点分配
	app.GET("/admin/user-api-item/edit", (&admin.UserApiItemController{}).Edit, "给用户分配接口节点")
	app.POST("/admin/user-api-item/edit", (&admin.UserApiItemController{}).Edit, "给用户分配接口节点")

	// 用户服务数据节点分配
	app.GET("/admin/user-data-item/edit", (&admin.UserDataItemController{}).Edit, "给用户分配数据节点")
	app.POST("/admin/user-data-item/edit", (&admin.UserDataItemController{}).Edit, "给用户分配数据节点")

	// 用户服务菜单节点分配
	app.GET("/admin/user-menu-item/edit", (&admin.UserMenuItemController{}).Edit, "给用户分配菜单节点")
	app.POST("/admin/user-menu-item/edit", (&admin.UserMenuItemController{}).Edit, "给用户分配菜单节点")

	// 系统管理
	app.POST("/admin/system/rebuild-cache", (&admin.SystemController{}).RebuildCache, "重建所有缓存")
}
