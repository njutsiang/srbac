package main

import (
	"srbac/check"
	"srbac/middlewares"
	"srbac/routers"
	"srbac/srbac"
)

func main() {

	// 初始化
	srbac.InitEngine()
	srbac.InitMySQL()
	srbac.InitRedis()
	srbac.InitSession()

	// 注册中间件
	srbac.Engine.Use(middlewares.ErrorHandle)
	srbac.Engine.Use(middlewares.SessionHandle)
	srbac.Engine.Use(middlewares.CsrfHandle)
	srbac.Engine.NoMethod(middlewares.NotFoundHandle)
	srbac.Engine.NoRoute(middlewares.NotFoundHandle)

	// 配置路由
	routers.SetRouters(srbac.Engine)

	// 检查初始化数据
	check.InitSrbacData()

	// 启动端口
	err := srbac.Engine.Run(":8102")
	srbac.CheckError(err)
}
