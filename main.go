package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"srbac/middlewares"
	"srbac/routers"
	"srbac/srbac"
	"time"
)

func main() {

	var err error
	srbac.Engine = gin.Default()

	// 设置 multipart form 内存限制
	srbac.Engine.MaxMultipartMemory = 8 << 20

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		srbac.Config.Mysql.User,
		srbac.Config.Mysql.Password,
		srbac.Config.Mysql.Host,
		srbac.Config.Mysql.Port,
		srbac.Config.Mysql.Db,
		srbac.Config.Mysql.Charset)

	// 彻底禁用 CreatedAt、UpdatedAt 自动更新
	srbac.Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			now, err := time.Parse(srbac.TimeYmdhis, "1970-01-01 00:00:00")
			srbac.CheckError(err)
			return now
		},
	})
	srbac.CheckError(err)

	// 设置连接池
	sqlDb, err := srbac.Db.DB()
	srbac.CheckError(err)

	// 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxIdleConns(20)
	// 设置打开数据库连接的最大数量
	sqlDb.SetMaxOpenConns(200)
	// 设置连接可复用的最大时间
	sqlDb.SetConnMaxLifetime(time.Hour)

	// 启用 Redis 用于 Session 存储
	store, err := redis.NewStoreWithDB(20, "tcp",
		srbac.Config.Redis.Host+":"+srbac.Config.Redis.Port,
		srbac.Config.Redis.Password,
		srbac.Config.Redis.Db,
		[]byte(srbac.Config.Session.Key))
	srbac.CheckError(err)
	srbac.Engine.Use(sessions.Sessions("srbac_session_id", store))

	// 注册中间件
	srbac.Engine.Use(middlewares.ErrorHandle)

	// 配置路由
	routers.SetRouters(srbac.Engine)

	// 启动端口
	err = srbac.Engine.Run(":8102")
	srbac.CheckError(err)
}
