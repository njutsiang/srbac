package srbac

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
)

// 初始化 Session
// 启用 Redis 用于 Session 存储
func InitSession() {
	store, err := redis.NewStoreWithDB(20, "tcp",
		Config.Redis.Host + ":" + Config.Redis.Port,
		Config.Redis.Password,
		fmt.Sprintf("%d", Config.Redis.Db),
		[]byte(Config.Cookie.Key))
	CheckError(err)
	Engine.Use(sessions.Sessions("srbac_session_id", store))
}