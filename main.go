package main

import (
	"fmt"
	"srbac/app"
	"srbac/check"
	"srbac/middlewares"
	"srbac/routers"
)

func main() {

	// 初始化
	app.InitEngine()
	app.InitMySQL()
	app.InitRedis()
	app.InitSession()

	// 注册中间件
	app.Engine.Use(middlewares.ErrorHandle)
	app.Engine.Use(middlewares.SessionHandle)
	app.Engine.Use(middlewares.CsrfHandle)
	app.Engine.NoMethod(middlewares.NotFoundHandle)
	app.Engine.NoRoute(middlewares.NotFoundHandle)

	// 配置路由
	routers.SetRouters(app.Engine)

	// 检查初始化数据
	check.InitSrbacData()

	// 启动端口
	err := app.Engine.Run(fmt.Sprintf(":%s", app.Config.Listen.Port))
	app.CheckError(err)
}
