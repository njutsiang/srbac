package exception

// 抛出异常
func Throw(code Code, msg... string) {
	message := ""
	if len(msg) >=1 && len(msg) >= 1 {
		message = msg[0]
	} else {
		message = GetMessage(code)
	}
	panic(&Exception{
		code: code,
		message: message,
	})
}

// 获取错误提示语
func GetMessage(code Code) string {
	return messages[code]
}