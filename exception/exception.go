package exception

// 定义异常
type Exception struct {
	code Code
	message string
}

// 异常错误码
func (this *Exception) GetCode() Code {
	return this.code
}

// 异常提示语
func (this *Exception) GetMessage() string {
	return this.message
}

// 实现 error
func (this *Exception) Error() string {
	return this.message
}