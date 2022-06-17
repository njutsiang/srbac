package app

import "github.com/go-redis/redis/v8"

// Redis 连接
var Rdb *redis.Client

// 初始化 Redis
func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: Config.Redis.Host + ":" + Config.Redis.Port,
		Password: Config.Redis.Password,
		DB: Config.Redis.Db,
	})
}