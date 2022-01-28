package srbac

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io"
	"os"
)

// 引擎
var Engine *gin.Engine

// 数据库连接
var Db *gorm.DB

// 年月日时分秒常量
var TimeYmdhis = "2006-01-02 15:04:05"

// 全局配置
var Config = getConfig()

// 解析配置文件
func getConfig() ConfigYaml {
	file, err := os.Open("./config.yaml")
	CheckError(err)

	content, err := io.ReadAll(file)
	CheckError(err)

	config := ConfigYaml{}
	err = yaml.Unmarshal(content, &config)
	CheckError(err)

	return config
}

// 判断是否有错误，有则记录错误日志，并抛出 panic
func CheckError(err interface{}) {
	if err == nil {
		return
	}
	panic(err)
}