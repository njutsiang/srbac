package srbac

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// 全局配置
var Config = getConfig()

// 配置文件
type ConfigYaml struct {
	Listen struct{
		Port string
	}
	Mysql struct{
		Host string
		Port string
		User string
		Password string
		Db string
		Charset string
	}
	Redis struct{
		Host string
		Port string
		Password string
		Db int
	}
	Cookie struct{
		Key string
	}
}

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