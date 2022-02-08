package srbac

import "github.com/gin-gonic/gin"

// 引擎
var Engine *gin.Engine

// 初始化 Engine
func InitEngine() {
	Engine = gin.Default()

	// 设置 multipart form 内存限制
	Engine.MaxMultipartMemory = 8 << 20
}