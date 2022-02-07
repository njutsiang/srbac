package exception

import "srbac/code"

// 定义异常
type Exception struct {
	code code.Code
	message string
}

// 异常错误码
func (this *Exception) GetCode() code.Code {
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

// 抛出异常
func NewException(err_code code.Code, err_msg... string) {
	msg := ""
	if len(err_msg) >=1 && len(err_msg) >= 1 {
		msg = err_msg[0]
	} else {
		msg = code.GetMessage(err_code)
	}
	panic(&Exception{
		code: err_code,
		message: msg,
	})
}