package app

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 数据库连接
var Db *gorm.DB

// 初始化 MySQL
func InitMySQL() {
	var err error

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		Config.Mysql.User,
		Config.Mysql.Password,
		Config.Mysql.Host,
		Config.Mysql.Port,
		Config.Mysql.Db,
		Config.Mysql.Charset)

	// 彻底禁用 CreatedAt、UpdatedAt 自动更新
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			now, err := time.Parse(TimeYmdhis, "1970-01-01 00:00:00")
			CheckError(err)
			return now
		},
	})
	CheckError(err)

	// 设置连接池
	sqlDb, err := Db.DB()
	CheckError(err)

	// 设置空闲连接池中连接的最大数量
	sqlDb.SetMaxIdleConns(20)

	// 设置打开数据库连接的最大数量
	sqlDb.SetMaxOpenConns(200)

	// 设置连接可复用的最大时间
	sqlDb.SetConnMaxLifetime(time.Hour)
}