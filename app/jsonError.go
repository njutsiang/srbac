package app

import (
	"srbac/exception"
)

// 错误
// HTTP 响应为 JSON
type JsonError struct {
	Code       int
	StatusCode int
	Message    string
}

// 实例化
// 实现 error interface，因为 interface 是指针，所有要取 JsonError 的指针
// NewJsonError(400000)
// NewJsonError(400000, error)
// NewJsonError(400000, "错误提示语")
// NewJsonError(error)
// NewJsonError("错误提示语")
func NewJsonError(params ...interface{}) *JsonError {
	err := 500000
	statusCode := 500
	message := ""
	if len(params) == 1 {
		switch params[0].(type) {
		case int:
			err = params[0].(int)
		case exception.Code:
			c := params[0].(exception.Code)
			err = int(c)
			message = exception.GetMessage(c)
			statusCode = err / 1000
		case string:
			message = params[0].(string)
		case error:
			message = params[0].(error).Error()
		}
	} else if len(params) == 2 {
		switch params[0].(type) {
		case int:
			err = params[0].(int)
		case exception.Code:
			err = int(params[0].(exception.Code))
		}
		switch params[1].(type) {
		case string:
			message = params[1].(string)
		case error:
			message = params[1].(error).Error()
		}
	}
	if message == "" {
		message = "未知错误"
	}
	return &JsonError{
		Code:       err,
		StatusCode: statusCode,
		Message:    message,
	}
}

// 实现接口 error
func (this *JsonError) Error() string {
	return this.Message
}
