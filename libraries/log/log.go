package log

import (
	"fmt"
	"runtime"
	"time"
)

// 错误日志
func Error(err interface{}) {
	// 错误信息
	message := ""
	switch err.(type) {
	case string:
		message = err.(string)
	case error:
		message = err.(error).Error()
	}

	// 调用堆栈
	buf := make([]byte, 64 << 10)
	runtime.Stack(buf, false)

	// 打印和保存日志
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] [error] " + message)
	fmt.Println(string(buf))
}

// 记录错误日志，并抛出 panic
func Panic(err interface{}) {
	Error(err)
	panic(err)
}